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

const (
	errMsgGetPolicyVersion = "Failed to get IAM policy version"
	errMsgDecodePolicyDoc  = "Failed to URL decode policy document"
	errMsgParsePolicyDoc   = "Failed to parse policy document"
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
		require.NoError(t, err, errMsgGetPolicyVersion)

		// AWS returns URL-encoded policy documents, so we need to decode them
		decodedDocument, err := url.QueryUnescape(*policyVersion.PolicyVersion.Document)
		require.NoError(t, err, errMsgDecodePolicyDoc)

		// Parse the policy documents to compare them
		var expectedPolicy, actualPolicy map[string]interface{}
		err = json.Unmarshal([]byte(policyDocument), &expectedPolicy)
		require.NoError(t, err, errMsgParsePolicyDoc)

		err = json.Unmarshal([]byte(decodedDocument), &actualPolicy)
		require.NoError(t, err, errMsgParsePolicyDoc)

		assert.Equal(t, expectedPolicy, actualPolicy, "Policy documents do not match!")
	})

	t.Run("TestIAMPolicyStatements", func(t *testing.T) {
		policyVersion, err := iamClient.GetPolicyVersion(context.TODO(), &iam.GetPolicyVersionInput{
			PolicyArn: &policyArn,
			VersionId: aws.String("v1"),
		})
		require.NoError(t, err, errMsgGetPolicyVersion)

		// AWS returns URL-encoded policy documents, so we need to decode them
		decodedDocument, err := url.QueryUnescape(*policyVersion.PolicyVersion.Document)
		require.NoError(t, err, errMsgDecodePolicyDoc)

		var policyDoc map[string]interface{}
		err = json.Unmarshal([]byte(decodedDocument), &policyDoc)
		require.NoError(t, err, errMsgParsePolicyDoc)

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

	t.Run("TestIAMPolicyStatementConditions", func(t *testing.T) {
		policyVersion, err := iamClient.GetPolicyVersion(context.TODO(), &iam.GetPolicyVersionInput{
			PolicyArn: &policyArn,
			VersionId: aws.String("v1"),
		})
		require.NoError(t, err, errMsgGetPolicyVersion)

		decodedDocument, err := url.QueryUnescape(*policyVersion.PolicyVersion.Document)
		require.NoError(t, err, errMsgDecodePolicyDoc)

		var policyDoc map[string]interface{}
		err = json.Unmarshal([]byte(decodedDocument), &policyDoc)
		require.NoError(t, err, errMsgParsePolicyDoc)

		statements, ok := policyDoc["Statement"].([]interface{})
		require.True(t, ok, "Policy should contain Statement array")

		var conditionStatement map[string]interface{}
		for _, stmt := range statements {
			statement := stmt.(map[string]interface{})
			if sid, hasSid := statement["Sid"].(string); hasSid && sid == "Stmt2" {
				conditionStatement = statement
				break
			}
		}
		require.NotNil(t, conditionStatement, "Expected to find statement Stmt2 with condition")

		conditionBlock, ok := conditionStatement["Condition"].(map[string]interface{})
		require.True(t, ok, "Statement Stmt2 should include Condition block")

		stringLike, ok := conditionBlock["StringLike"].(map[string]interface{})
		require.True(t, ok, "Condition block should include StringLike test")

		rawPrefixes, ok := stringLike["s3:prefix"]
		require.True(t, ok, "StringLike condition should define s3:prefix values")

		var prefixes []string
		switch v := rawPrefixes.(type) {
		case []interface{}:
			for _, entry := range v {
				prefix, castOk := entry.(string)
				require.True(t, castOk, "Condition s3:prefix entries must be strings")
				prefixes = append(prefixes, prefix)
			}
		case string:
			prefixes = append(prefixes, v)
		default:
			require.Failf(t, "unexpected condition type", "s3:prefix condition expected string or list, got %T", rawPrefixes)
		}

		assert.Contains(t, prefixes, "home/*", "Condition values should include home/* prefix restriction")
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
