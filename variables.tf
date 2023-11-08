variable "aws_region" {
  default = "eu-central-1"
}

variable "aws_cli_profile" {
  default = "default"
}

variable "mensa_timezone" {
  default = "Europe/Berlin"
}

variable "project_name" {
  type    = string
  default = "munchy2"
}

variable "table_name" {
  type        = string
  description = "Name of DynamoDB table where food items are stored."
  default     = "munchy"
}

variable "lambda_timeout" {
  type        = number
  description = "Timeout for Lambda execution in seconds."
  default     = 10
}

variable "lambda_memory_size" {
  type        = number
  description = "Memory size for Lambda execution in MB."
  default     = 128
}

variable "webhookurl" {
  description = "Slack webhook url to post to."
  type        = string
}

variable "deepl_target_lang" {
  description = "Slack webhook url to post to."
  type        = string
  default     = "EN"
}

variable "deepl_url" {
  description = "DeepL API URL for translations."
  type        = string
  default     = "https://api-free.deepl.com/v2/translate"
}

variable "deepl_key" {
  description = "DeepL API key."
  type        = string
}
