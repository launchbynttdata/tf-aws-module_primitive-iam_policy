variable "policy_name" {
  description = "The name of the IAM policy."
  type        = string
}

variable "tags" {
  description = "A map of tags to assign to the IAM policy."
  type        = map(string)
  default     = {}
}

variable "policy_statement" {
  description = "The policy statement in JSON format."
  type = map(object({
    sid       = string
    actions   = list(string)
    resources = list(string)
  }))
}
