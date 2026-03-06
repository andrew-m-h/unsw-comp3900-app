package main

// Shared configuration constants and CloudFormation export names for cross-stack references.

// --- App / ECR ---
const (
	DefaultECRRepositoryName = "unsw-comp3900-app"
	AppRunnerServicePort     = "8080"
)

// --- App Runner instance ---
const (
	AppRunnerCPU        = "1024"
	AppRunnerMemory     = "2048"
	AppRunnerAutoDeploy = true
)

// --- GitHub OIDC (https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-amazon-web-services) ---
const (
	GitHubOIDCURL        = "https://token.actions.githubusercontent.com"
	GitHubOIDCAudience   = "sts.amazonaws.com"
	GitHubOIDCThumbprint = "6938fd4d98bab03faadb97b34396831e3780aea1"
)

// --- CloudFormation export names (base stack exports; app stack imports via Fn.importValue) ---
const (
	ExportECRRepositoryUri       = "UnswComp3900App-ECRRepositoryUri"
	ExportAppRunnerEcrAccessRole = "UnswComp3900App-AppRunnerEcrAccessRoleArn"
	ExportStaticAssetsBucketName = "UnswComp3900App-StaticAssetsBucketName"
	ExportGitHubECRPushRoleArn   = "UnswComp3900App-GitHubECRPushRoleArn"
)
