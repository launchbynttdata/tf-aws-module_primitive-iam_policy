module "aws_iam_policy" {
  source = "../../"

  policy_name      = var.policy_name
  policy_statement = var.policy_statement
  tags             = var.tags
}
