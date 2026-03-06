package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// GitHubOIDCConfig configures which GitHub repo can assume the ECR-push role via OIDC.
type GitHubOIDCConfig struct {
	Owner  string // GitHub org or username, e.g. "myorg"
	Repo   string // Repository name, e.g. "unsw-comp3900-app"
	Branch string // Optional: restrict to a branch, e.g. "main". Leave empty to allow any ref (e.g. "repo:owner/repo:*").
}

// getContextString reads a string from CDK context (e.g. cdk.json "context" section). Returns defaultVal if missing or not a string.
func getContextString(scope constructs.Construct, key string, defaultVal string) string {
	v := scope.Node().TryGetContext(jsii.String(key))
	if v == nil {
		return defaultVal
	}
	switch s := v.(type) {
	case string:
		return s
	case *string:
		if s != nil {
			return *s
		}
		return defaultVal
	default:
		return defaultVal
	}
}

// newGitHubOIDCRoleForECR creates an IAM OIDC provider for GitHub and a role that GitHub Actions can assume to push to ECR (no long-lived credentials).
// Uses the L1 AWS::IAM::OIDCProvider so no CDK bootstrap bucket is required.
func newGitHubOIDCRoleForECR(stack awscdk.Stack, ecrRepo awsecr.IRepository, config *GitHubOIDCConfig) awsiam.IRole {
	githubOIDC := awsiam.NewCfnOIDCProvider(stack, jsii.String("GitHubOIDC"), &awsiam.CfnOIDCProviderProps{
		Url:           jsii.String(GitHubOIDCURL),
		ClientIdList:  &[]*string{jsii.String(GitHubOIDCAudience)},
		ThumbprintList: &[]*string{jsii.String(GitHubOIDCThumbprint)},
	})

	subClaim := "repo:" + config.Owner + "/" + config.Repo + ":*"
	if config.Branch != "" {
		subClaim = "repo:" + config.Owner + "/" + config.Repo + ":ref:refs/heads/" + config.Branch
	}
	conditions := map[string]interface{}{
		"StringEquals": map[string]interface{}{
			"token.actions.githubusercontent.com:aud": GitHubOIDCAudience,
		},
		"StringLike": map[string]interface{}{
			"token.actions.githubusercontent.com:sub": subClaim,
		},
	}

	principal := awsiam.NewOpenIdConnectPrincipal(githubOIDC, &conditions)
	role := awsiam.NewRole(stack, jsii.String("GitHubECRPushRole"), &awsiam.RoleProps{
		RoleName:    jsii.String("github-ecr-push-" + config.Repo),
		AssumedBy:   principal,
		Description: jsii.String("Allows GitHub Actions to push images to ECR via OIDC"),
	})
	ecrRepo.GrantPush(role)
	return role
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewBaseStack(app, "BaseStack", &BaseStackProps{
		StackProps: awscdk.StackProps{Env: env()},
	})

	NewAppStack(app, "AppStack", &AppStackProps{
		StackProps: awscdk.StackProps{Env: env()},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to be deployed.
func env() *awscdk.Environment {
	return nil
}
