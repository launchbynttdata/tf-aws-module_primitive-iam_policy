package testimpl

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComposableComplete(t *testing.T, ctx types.TestContext) {
	iamClient := GetAWSIAMClient(t)

	policyArn := terraform.Output(t, ctx.TerratestTerraformOptions(), "policy_arn")
	policyName := terraform.Output(t, ctx.TerratestTerraformOptions(), "policy_name")
	policyDocument := terraform.Output(t, ctx.TerratestTerraformOptions(), "policy_document")

	t.Run("TestIAMPolicyExists", func(t *testing.T) {
		policy, err := iamClient.GetPolicy(context.TODO(), &iam.GetPolicyInput{
			PolicyArn: &policyArn,
		})
		require.NoError(t, err, "Failed to get IAM policy")
		assert.Equal(t, policyArn, *policy.Policy.Arn, "Expected policy ARN did not match actual ARN!")
		assert.Equal(t, policyName, *policy.Policy.PolicyName, "Expected policy name did not match actual name!")
	})

	t.Run("TestIAMPolicyDocument", func(t *testing.T) {
		policyVersion, err := iamClient.GetPolicyVersion(context.TODO(), &iam.GetPolicyVersionInput{
			PolicyArn: &policyArn,
			VersionId: aws.String("v1"),
		})
		require.NoError(t, err, "Failed to get IAM policy version")

		// AWS returns URL-encoded policy documents, so we need to decode them
		decodedDocument, err := url.QueryUnescape(*policyVersion.PolicyVersion.Document)
		require.NoError(t, err, "Failed to URL decode policy document")

		// Parse the policy documents to compare them
		var expectedPolicy, actualPolicy map[string]interface{}
		err = json.Unmarshal([]byte(policyDocument), &expectedPolicy)
		require.NoError(t, err, "Failed to parse expected policy document")

		err = json.Unmarshal([]byte(decodedDocument), &actualPolicy)
		require.NoError(t, err, "Failed to parse actual policy document")

		assert.Equal(t, expectedPolicy, actualPolicy, "Policy documents do not match!")
	})

	t.Run("TestIAMPolicyStatements", func(t *testing.T) {
		policyVersion, err := iamClient.GetPolicyVersion(context.TODO(), &iam.GetPolicyVersionInput{
			PolicyArn: &policyArn,
			VersionId: aws.String("v1"),
		})
		require.NoError(t, err, "Failed to get IAM policy version")

		// AWS returns URL-encoded policy documents, so we need to decode them
		decodedDocument, err := url.QueryUnescape(*policyVersion.PolicyVersion.Document)
		require.NoError(t, err, "Failed to URL decode policy document")

		var policyDoc map[string]interface{}
		err = json.Unmarshal([]byte(decodedDocument), &policyDoc)
		require.NoError(t, err, "Failed to parse policy document")

		statements, ok := policyDoc["Statement"].([]interface{})
		require.True(t, ok, "Policy should contain Statement array")
		require.Greater(t, len(statements), 0, "Policy should have at least one statement")

		// Verify each statement has required fields
		for i, stmt := range statements {
			statement := stmt.(map[string]interface{})
			assert.Contains(t, statement, "Sid", "Statement %d should have Sid", i)
			assert.Contains(t, statement, "Action", "Statement %d should have Action", i)
			assert.Contains(t, statement, "Resource", "Statement %d should have Resource", i)
		}
	})
}

func GetAWSIAMClient(t *testing.T) *iam.Client {
	awsIAMClient := iam.NewFromConfig(GetAWSConfig(t))
	return awsIAMClient
}

func GetAWSConfig(t *testing.T) (cfg aws.Config) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoErrorf(t, err, "unable to load SDK config, %v", err)
	return cfg
}
