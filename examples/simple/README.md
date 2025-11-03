# simple

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.0 |

## Providers

No providers.

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_aws_iam_policy"></a> [aws\_iam\_policy](#module\_aws\_iam\_policy) | ../../ | n/a |

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_policy_name"></a> [policy\_name](#input\_policy\_name) | The name of the IAM policy. | `string` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | A map of tags to assign to the IAM policy. | `map(string)` | `{}` | no |
| <a name="input_policy_statement"></a> [policy\_statement](#input\_policy\_statement) | The policy statements, supporting optional IAM conditions. | <pre>map(object({<br/>    sid       = string<br/>    actions   = list(string)<br/>    resources = list(string)<br/>    conditions = optional(list(object({<br/>      test     = string<br/>      variable = string<br/>      values   = list(string)<br/>    })))<br/>  }))</pre> | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_policy_arn"></a> [policy\_arn](#output\_policy\_arn) | The ARN of the IAM policy. |
| <a name="output_policy_name"></a> [policy\_name](#output\_policy\_name) | The name of the IAM policy. |
| <a name="output_policy_id"></a> [policy\_id](#output\_policy\_id) | The ID of the IAM policy. |
| <a name="output_policy_document"></a> [policy\_document](#output\_policy\_document) | The policy document in JSON format. |
| <a name="output_policy_tags"></a> [policy\_tags](#output\_policy\_tags) | The tags applied to the IAM policy. |
<!-- END_TF_DOCS -->
