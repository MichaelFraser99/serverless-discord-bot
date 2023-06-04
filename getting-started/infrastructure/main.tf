resource "aws_api_gateway_rest_api" "bot_api" {
  body = templatefile("../api/swagger.yaml", {
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
  filename = "../service/main.zip"
  source_code_hash = filebase64sha256("../service/main.zip")
  architectures = ["arm64"]

  environment {
    variables = {
      PUBLIC_KEY = var.public_key
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

resource "discord-application_command" "poke" {
  application_id = var.application_id
  name = "poke"
  description = "a poke to the application"
  type = 1
}
