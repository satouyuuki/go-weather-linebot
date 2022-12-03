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
  region = "ap-northeast-1"
  profile = "cli_only_user"
}

resource "aws_ecr_repository" "weather_bot_repo" {
  name = "weather_bot_repo"
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
  # required
  function_name = "weather_bot_lambda"
  role          = aws_iam_role.iam_for_lambda.arn

  # optional
  image_uri = "${aws_ecr_repository.weather_bot_repo.repository_url}:latest"
  package_type = "Image"
}