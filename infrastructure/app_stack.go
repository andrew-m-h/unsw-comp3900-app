package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapprunner"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// AppStackProps configures the app stack (App Runner, CloudFront). The stack imports ECR URI, App Runner access role ARN, and S3 bucket name from the base stack outputs.
type AppStackProps struct {
	awscdk.StackProps
}

// NewAppStack creates the app stack: App Runner service (using ECR and role from base stack) and CloudFront distribution (App Runner + S3 from base).
// Requires the base stack to be deployed first so that Fn.importValue can resolve the exports.
func NewAppStack(scope constructs.Construct, id string, props *AppStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	ecrRepositoryName := getContextString(scope, "ecrRepositoryName", DefaultECRRepositoryName)

	// Import values from base stack
	ecrRepositoryUri := awscdk.Fn_ImportValue(jsii.String(ExportECRRepositoryUri))
	appRunnerEcrAccessRoleArn := awscdk.Fn_ImportValue(jsii.String(ExportAppRunnerEcrAccessRole))
	staticAssetsBucketName := awscdk.Fn_ImportValue(jsii.String(ExportStaticAssetsBucketName))

	// Image identifier for App Runner: <ecr-uri>:latest
	imageIdentifier := awscdk.Fn_Join(jsii.String(":"), &[]*string{
		ecrRepositoryUri,
		jsii.String("latest"),
	})

	// Guestbook DynamoDB table
	guestbookTable := NewGuestbookTable(stack, jsii.String("GuestbookTable"))

	// IAM role for App Runner instance (runtime) — allows app to call DynamoDB
	appRunnerInstanceRole := awsiam.NewRole(stack, jsii.String("AppRunnerInstanceRole"), &awsiam.RoleProps{
		AssumedBy:   awsiam.NewServicePrincipal(jsii.String("tasks.apprunner.amazonaws.com"), nil),
		Description: jsii.String("Allows App Runner service instance to access DynamoDB (e.g. guestbook table)"),
	})
	guestbookTable.GrantReadWriteData(appRunnerInstanceRole)

	// App Runner service
	appRunnerService := awsapprunner.NewCfnService(stack, jsii.String("AppRunnerService"), &awsapprunner.CfnServiceProps{
		ServiceName: jsii.String(ecrRepositoryName),
		SourceConfiguration: &awsapprunner.CfnService_SourceConfigurationProperty{
			AuthenticationConfiguration: &awsapprunner.CfnService_AuthenticationConfigurationProperty{
				AccessRoleArn: appRunnerEcrAccessRoleArn,
			},
			AutoDeploymentsEnabled: jsii.Bool(AppRunnerAutoDeploy),
			ImageRepository: &awsapprunner.CfnService_ImageRepositoryProperty{
				ImageIdentifier:     imageIdentifier,
				ImageRepositoryType: jsii.String("ECR"),
				ImageConfiguration: &awsapprunner.CfnService_ImageConfigurationProperty{
					Port: jsii.String(AppRunnerServicePort),
					RuntimeEnvironmentVariables: &[]*awsapprunner.CfnService_KeyValuePairProperty{
						{Name: jsii.String("GUESTBOOK_TABLE_NAME"), Value: guestbookTable.TableName()},
					},
				},
			},
		},
		InstanceConfiguration: &awsapprunner.CfnService_InstanceConfigurationProperty{
			Cpu:             jsii.String(AppRunnerCPU),
			Memory:          jsii.String(AppRunnerMemory),
			InstanceRoleArn: appRunnerInstanceRole.RoleArn(),
		},
		HealthCheckConfiguration: &awsapprunner.CfnService_HealthCheckConfigurationProperty{
			Path:     jsii.String(AppRunnerHealthCheckPath),
			Protocol: jsii.String("HTTP"),
			Timeout:  jsii.Number(10),
			Interval: jsii.Number(10),
		},
	})

	// S3 bucket reference (from base stack) for CloudFront origin
	staticBucket := awss3.Bucket_FromBucketName(stack, jsii.String("StaticAssets"), staticAssetsBucketName)

	// CloudFront: default to S3 (static/index.html for "/" and SPA fallback); API paths to App Runner
	appRunnerOrigin := awscloudfrontorigins.NewHttpOrigin(appRunnerService.AttrServiceUrl(), &awscloudfrontorigins.HttpOriginProps{
		ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
		CustomHeaders:  &map[string]*string{},
	})
	s3Origin := awscloudfrontorigins.S3BucketOrigin_WithOriginAccessControl(staticBucket, &awscloudfrontorigins.S3BucketOriginWithOACProps{})

	distribution := awscloudfront.NewDistribution(stack, jsii.String("Distribution"), &awscloudfront.DistributionProps{
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			Origin:               s3Origin,
			ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
			AllowedMethods:       awscloudfront.AllowedMethods_ALLOW_GET_HEAD_OPTIONS(),
			CachedMethods:        awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
		},
		AdditionalBehaviors: &map[string]*awscloudfront.BehaviorOptions{
			// API / backend routes → App Runner (path patterns must include leading slash to match request path)
			"/health": {
				Origin:               appRunnerOrigin,
				ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
				AllowedMethods:       awscloudfront.AllowedMethods_ALLOW_GET_HEAD_OPTIONS(),
				CachedMethods:        awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
			},
			"/api/*": {
				Origin:               appRunnerOrigin,
				ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
				AllowedMethods:       awscloudfront.AllowedMethods_ALLOW_ALL(),
				CachedMethods:        awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
			},
			"/static/*": {
				Origin:               s3Origin,
				ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
				AllowedMethods:       awscloudfront.AllowedMethods_ALLOW_GET_HEAD_OPTIONS(),
				CachedMethods:        awscloudfront.CachedMethods_CACHE_GET_HEAD_OPTIONS(),
			},
		},
		DefaultRootObject: jsii.String("static/index.html"),
		// Error responses for SPA Apps - not helpful for debugging
		/*
			ErrorResponses: &[]*awscloudfront.ErrorResponse{
				{
					HttpStatus:         jsii.Number(403),
					ResponseHttpStatus: jsii.Number(200),
					ResponsePagePath:   jsii.String("/static/index.html"),
					Ttl:                awscdk.Duration_Seconds(jsii.Number(0)),
				},
				{
					HttpStatus:         jsii.Number(404),
					ResponseHttpStatus: jsii.Number(200),
					ResponsePagePath:   jsii.String("/static/index.html"),
					Ttl:                awscdk.Duration_Seconds(jsii.Number(0)),
				},
			},
		*/
		Comment:                jsii.String("CDN for App Runner app and static assets"),
		PriceClass:             awscloudfront.PriceClass_PRICE_CLASS_100,
		MinimumProtocolVersion: awscloudfront.SecurityPolicyProtocol_TLS_V1_2_2021,
	})

	// Outputs (no export needed unless other stacks depend on these)
	awscdk.NewCfnOutput(stack, jsii.String("GuestbookTableName"), &awscdk.CfnOutputProps{
		Value:       guestbookTable.TableName(),
		Description: jsii.String("DynamoDB guestbook table name (set as GUESTBOOK_TABLE_NAME in App Runner)"),
	})
	awscdk.NewCfnOutput(stack, jsii.String("AppRunnerServiceArn"), &awscdk.CfnOutputProps{
		Value:       appRunnerService.AttrServiceArn(),
		Description: jsii.String("App Runner service ARN (for start-deployment after CDK deploy)"),
	})
	awscdk.NewCfnOutput(stack, jsii.String("AppRunnerServiceUrl"), &awscdk.CfnOutputProps{
		Value:       appRunnerService.AttrServiceUrl(),
		Description: jsii.String("App Runner service URL (also available behind CloudFront)"),
	})
	awscdk.NewCfnOutput(stack, jsii.String("CloudFrontDistributionUrl"), &awscdk.CfnOutputProps{
		Value:       awscdk.Fn_Join(jsii.String(""), &[]*string{jsii.String("https://"), distribution.DomainName(), jsii.String("/")}),
		Description: jsii.String("CloudFront URL; use this as the public entry point (app + static assets)"),
	})

	return stack
}
