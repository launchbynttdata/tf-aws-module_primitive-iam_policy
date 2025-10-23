output "policy_arn" {
  description = "The ARN of the IAM policy."
  value       = aws_iam_policy.this.arn
}

output "policy_name" {
  description = "The name of the IAM policy."
  value       = aws_iam_policy.this.name
}

output "policy_id" {
  description = "The ID of the IAM policy."
  value       = aws_iam_policy.this.id
}

output "policy_document" {
  description = "The policy document in JSON format."
  value       = aws_iam_policy.this.policy
}

output "policy_tags" {
  description = "The tags applied to the IAM policy."
  value       = aws_iam_policy.this.tags
}
