terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region  = var.REGION
  profile = var.AWS_PROFILE
}

data "aws_caller_identity" "current" {}

locals {
  account_id          = data.aws_caller_identity.current.account_id
  codebase_root_path  = abspath("${path.module}/..")
  ecr_repository_name = "weather_bot_repo"
  ecr_image_tag       = "latest"
}

resource "aws_ecr_repository" "weather_bot_repo" {
  name = local.ecr_repository_name
}

resource "null_resource" "ecr_image" {
  triggers = {
    "lambda_file" = md5(file("${local.codebase_root_path}/lambda/main.go"))
    "docker_file" = md5(file("${local.codebase_root_path}/lambda/Dockerfile"))
  }

  provisioner "local-exec" {
    command = <<EOF
      aws ecr get-login-password --region ${var.REGION} --profile ${var.AWS_PROFILE} | docker login --username AWS --password-stdin ${local.account_id}.dkr.ecr.${var.REGION}.amazonaws.com
      cd ${local.codebase_root_path}/lambda
      docker build -t ${aws_ecr_repository.weather_bot_repo.repository_url}:${local.ecr_image_tag} . --platform=linux/amd64
      docker push ${aws_ecr_repository.weather_bot_repo.repository_url}:${local.ecr_image_tag}
    EOF
  }
}

data "aws_ecr_image" "lambda_image" {
  depends_on = [
    null_resource.ecr_image
  ]
  repository_name = local.ecr_repository_name
  image_tag       = local.ecr_image_tag
}

output "ecr_image_id" {
  value = data.aws_ecr_image.lambda_image.id
}

output "aws_profile" {
  value = var.AWS_PROFILE
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "weather_bot_lambda_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_lambda_function" "weather_bot_lambda" {
  depends_on = [
    null_resource.ecr_image
  ]
  # required
  function_name = "weather_bot_lambda"
  role          = aws_iam_role.iam_for_lambda.arn

  # optional
  image_uri    = "${aws_ecr_repository.weather_bot_repo.repository_url}@${data.aws_ecr_image.lambda_image.id}"
  package_type = "Image"

  environment {
    variables = {
      CHANNNE_SECRET        = var.CHANNNE_SECRET
      CHANNNE_TOKEN         = var.CHANNNE_TOKEN
      OPENWEATHER_API_TOKEN = var.OPENWEATHER_API_TOKEN
    }
  }
}

resource "aws_lambda_function_url" "lambda_endpoint" {
  function_name      = aws_lambda_function.weather_bot_lambda.function_name
  authorization_type = "NONE"
}

resource "aws_lambda_permission" "lambda_permission" {
  statement_id  = "LambdaInvokePermission"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.weather_bot_lambda.function_name
  principal     = "*"
}

###
# lambdaの自動起動トリガー
###
resource "aws_cloudwatch_event_rule" "cron_event" {
  name                = "WeatherBotScheduledRule"
  # for production
  schedule_expression = "cron(0 22 ? * SUN-THU *)"
  # for testing
  # schedule_expression = "rate(5 minutes)"
  # 自動起動ON
  is_enabled = true
}

resource "aws_lambda_permission" "cron_permission" {
  statement_id  = "ScheduledEvent"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.weather_bot_lambda.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.cron_event.arn
}

resource "aws_cloudwatch_event_target" "event_target" {
  arn  = aws_lambda_function.weather_bot_lambda.arn
  rule = aws_cloudwatch_event_rule.cron_event.id
}

###
# cloudwatch logging and permissions
###
# See also the following AWS managed policy: AWSLambdaBasicExecutionRole
resource "aws_iam_policy" "lambda_logging" {
  name = "lambda_logging"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.lambda_logging.arn
}
