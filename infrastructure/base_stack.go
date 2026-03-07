package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// BaseStackProps configures the base stack (S3, ECR, IAM roles).
type BaseStackProps struct {
	awscdk.StackProps
	GitHub *GitHubOIDCConfig
}

// NewBaseStack creates the base stack: S3 bucket for static assets, ECR repository, IAM role for App Runner to pull from ECR, and optionally the GitHub OIDC role for ECR push.
// Outputs are exported for the app stack to import.
func NewBaseStack(scope constructs.Construct, id string, props *BaseStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	ecrRepositoryName := getContextString(scope, "ecrRepositoryName", DefaultECRRepositoryName)
	githubOwner := getContextString(scope, "githubOwner", "")
	githubRepo := getContextString(scope, "githubRepo", "")
	githubBranch := getContextString(scope, "githubBranch", "main")

	// S3 bucket for static assets (CloudFront will use via OAC from app stack)
	staticBucket := awss3.NewBucket(stack, jsii.String("StaticAssets"), &awss3.BucketProps{
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		EnforceSSL:        jsii.Bool(true),
		Versioned:         jsii.Bool(false),
	})
	// Allow CloudFront (OAC) to read objects; app stack uses this bucket as /static/* origin
	awss3.NewCfnBucketPolicy(stack, jsii.String("StaticAssetsBucketPolicy"), &awss3.CfnBucketPolicyProps{
		Bucket: staticBucket.BucketName(),
		PolicyDocument: map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []interface{}{
				map[string]interface{}{
					"Sid":       "AllowCloudFrontServicePrincipal",
					"Effect":    "Allow",
					"Principal": map[string]interface{}{"Service": "cloudfront.amazonaws.com"},
					"Action":    "s3:GetObject",
					"Resource":  awscdk.Fn_Join(jsii.String(""), &[]*string{staticBucket.BucketArn(), jsii.String("/*")}),
				},
			},
		},
	})

	// ECR repository
	ecrRepo := awsecr.NewRepository(stack, jsii.String("AppRepo"), &awsecr.RepositoryProps{
		RepositoryName:     jsii.String(ecrRepositoryName),
		RemovalPolicy:      awscdk.RemovalPolicy_DESTROY,
		ImageScanOnPush:    jsii.Bool(true),
		ImageTagMutability: awsecr.TagMutability_MUTABLE,
	})

	// IAM role for App Runner to pull from ECR
	appRunnerAccessRole := awsiam.NewRole(stack, jsii.String("AppRunnerEcrAccessRole"), &awsiam.RoleProps{
		AssumedBy:   awsiam.NewServicePrincipal(jsii.String("build.apprunner.amazonaws.com"), nil),
		Description: jsii.String("Allows App Runner to pull images from ECR (used by build and runtime)"),
	})
	ecrRepo.GrantPull(appRunnerAccessRole)

	// GitHub OIDC role for ECR push (optional)
	var githubConfig *GitHubOIDCConfig
	if props != nil && props.GitHub != nil {
		githubConfig = props.GitHub
	} else if githubOwner != "" && githubRepo != "" {
		githubConfig = &GitHubOIDCConfig{Owner: githubOwner, Repo: githubRepo, Branch: githubBranch}
	}
	if githubConfig != nil {
		provider := newGitHubOIDCProvider(stack)
		principal := githubOIDCPrincipal(provider, githubConfig)
		githubRole := newGitHubOIDCRoleForECR(stack, ecrRepo, principal, githubConfig)
		cdkDeployRole := newGitHubCDKDeployRole(stack, principal, githubConfig)
		awscdk.NewCfnOutput(stack, jsii.String("GitHubECRPushRoleArn"), &awscdk.CfnOutputProps{
			Value:       githubRole.RoleArn(),
			Description: jsii.String("ARN of the IAM role for GitHub Actions to push to ECR; set as AWS_ROLE_ARN in build-and-push workflow."),
			ExportName:  jsii.String(ExportGitHubECRPushRoleArn),
		})
		awscdk.NewCfnOutput(stack, jsii.String("GitHubCDKDeployRoleArn"), &awscdk.CfnOutputProps{
			Value:       cdkDeployRole.RoleArn(),
			Description: jsii.String("ARN of the IAM role for GitHub Actions to deploy CDK; set as AWS_CDK_DEPLOY_ROLE_ARN secret."),
			ExportName:  jsii.String(ExportGitHubCDKDeployRoleArn),
		})
	}

	// Exports for app stack
	awscdk.NewCfnOutput(stack, jsii.String("ECRRepositoryUri"), &awscdk.CfnOutputProps{
		Value:       ecrRepo.RepositoryUri(),
		Description: jsii.String("ECR repository URI; push your Docker image here and deploy manually in App Runner"),
		ExportName:  jsii.String(ExportECRRepositoryUri),
	})
	awscdk.NewCfnOutput(stack, jsii.String("AppRunnerEcrAccessRoleArn"), &awscdk.CfnOutputProps{
		Value:       appRunnerAccessRole.RoleArn(),
		Description: jsii.String("ARN of the IAM role for App Runner to pull from ECR"),
		ExportName:  jsii.String(ExportAppRunnerEcrAccessRole),
	})
	awscdk.NewCfnOutput(stack, jsii.String("StaticAssetsBucketName"), &awscdk.CfnOutputProps{
		Value:       staticBucket.BucketName(),
		Description: jsii.String("S3 bucket for static assets; served under /static/* via CloudFront"),
		ExportName:  jsii.String(ExportStaticAssetsBucketName),
	})

	return stack
}
