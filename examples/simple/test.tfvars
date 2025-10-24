policy_name = "tt-policy-0"
tags = {
  Environment = "dev"
  Project     = "tt-project"
}
policy_statement = {
  "Stmt1" = {
    sid       = "Stmt1"
    actions   = ["s3:ListBucket"]
    resources = ["arn:aws:s3:::example-bucket"]
  }
  "Stmt2" = {
    sid       = "Stmt2"
    actions   = ["s3:GetObject", "s3:PutObject"]
    resources = ["arn:aws:s3:::example-bucket/*"]
  }
}
