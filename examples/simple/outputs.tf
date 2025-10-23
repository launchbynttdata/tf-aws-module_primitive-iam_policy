output "policy_arn" {
  description = "The ARN of the IAM policy."
  value       = module.aws_iam_policy.policy_arn
}

output "policy_name" {
  description = "The name of the IAM policy."
  value       = module.aws_iam_policy.policy_name
}

output "policy_id" {
  description = "The ID of the IAM policy."
  value       = module.aws_iam_policy.policy_id
}

output "policy_document" {
  description = "The policy document in JSON format."
  value       = module.aws_iam_policy.policy_document
}

output "policy_tags" {
  description = "The tags applied to the IAM policy."
  value       = module.aws_iam_policy.policy_tags
}
