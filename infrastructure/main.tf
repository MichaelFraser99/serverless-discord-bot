resource "aws_api_gateway_rest_api" "bot_api" {
  body = templatefile("${path.module}/api/swagger.yaml", {
    botLambdaInvocationArn = aws_lambda_function.bot_lambda.invoke_arn
    name = "${var.name}-api"
  })

  name = "${var.name}-api"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_deployment" "bot_api_deployment" {
  rest_api_id = aws_api_gateway_rest_api.bot_api.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.bot_api.body))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "bot_api_stage" {
  deployment_id = aws_api_gateway_deployment.bot_api_deployment .id
  rest_api_id   = aws_api_gateway_rest_api.bot_api.id
  stage_name    = var.stage_name
}

resource "aws_lambda_function" "bot_lambda" {
  function_name = "${var.name}-lambda"
  role          = aws_iam_role.bot_lambda_role.arn
  handler = "main"
  runtime = "provided.al2"
  memory_size = 128
  filename = "${path.module}/service/main.zip"
  source_code_hash = filebase64sha256("${path.module}/service/main.zip")
  architectures = ["arm64"]

  environment {
    variables = {
      PUBLIC_KEY = var.public_key
      GITHUB_TOKEN = var.github_token
    }
  }

  depends_on = [aws_cloudwatch_log_group.bot_log_group]
}

resource "aws_lambda_permission" "lambda_permission" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.bot_lambda.arn
  principal     = "apigateway.amazonaws.com"

  # The /* part allows invocation from any stage, method and resource path
  # within API Gateway.
  source_arn = "${aws_api_gateway_rest_api.bot_api.execution_arn}/*"
}

resource "aws_cloudwatch_log_group" "bot_log_group" {
  name              = "/aws/lambda/${var.name}-lambda"
  retention_in_days = 14
}

resource "aws_iam_role" "bot_lambda_role" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

data "aws_iam_policy_document" "lambda_logging" {
  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = ["arn:aws:logs:*:*:*"]
  }
}

resource "aws_iam_policy" "lambda_logging" {
  name        = "lambda_logging"
  path        = "/"
  description = "IAM policy for logging from a lambda"
  policy      = data.aws_iam_policy_document.lambda_logging.json
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.bot_lambda_role.name
  policy_arn = aws_iam_policy.lambda_logging.arn
}

resource "discord-application_command" "create_server" {
  application_id = var.application_id
  name = "create"
  description = "Create a new server"
  type = 1
  options = [
    {
      name = "map"
      description = "select the Map"
      type = 3
      required = false
    },
    {
      name = "modsenabled"
      description = "Select if mods are enabled"
      type = 5
      required = false
    },
    {
      name = "maxplayers"
      description = "Select the max players"
      type = 4
      required = false
    },
    {
      name = "maxcars"
      description = "Select the max cars per player"
      type = 4
      required = false
    },
    {
      name = "private"
      description = "Select if the server is private"
      type = 5
      required = false
    }
  ]
}

resource "discord-application_command" "destroy_server" {
  application_id = var.application_id
  name = "destroy"
  description = "Destroy the server"
  type = 1
}

resource "terracurl_request" "discord_application_interaction_url" {
  name = "discord_application_interaction_url"
  url = "https://discord.com/api/v10/applications/${var.application_id}"
  method = "PATCH"
  request_body = <<EOF
{
  "interactions_endpoint_url": "${aws_api_gateway_stage.bot_api_stage.invoke_url}/interaction"
}
  EOF
  headers = {
    Authorization = "Bot ${var.application_secret}"
    Content-Type = "application/json"
    Accept = "application/json"
  }
  response_codes = [ 200, 204 ]

  destroy_url = "https://discord.com/api/v10/applications/${var.application_id}"
  destroy_method = "PATCH"
  destroy_request_body = <<EOF
{
  "interactions_endpoint_url": ""
}
  EOF
  destroy_headers = {
    Authorization = "Bot ${var.application_secret}"
    Content-Type = "application/json"
    Accept = "application/json"
  }
  destroy_response_codes = [ 200, 204 ]
}
