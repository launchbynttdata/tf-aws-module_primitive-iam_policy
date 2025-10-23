data "aws_iam_policy_document" "this" {
  dynamic "statement" {
    for_each = var.policy_statement
    content {
      sid       = statement.value.sid
      actions   = statement.value.actions
      resources = statement.value.resources
    }
  }
}

resource "aws_iam_policy" "this" {
  name   = var.policy_name
  policy = data.aws_iam_policy_document.this.json
  tags   = local.tags
}
