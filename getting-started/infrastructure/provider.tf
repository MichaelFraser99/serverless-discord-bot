terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }

    discord-application = {
      source = "MichaelFraser99/discord-application"
      version = "0.1.2"
    }
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "eu-west-2"
}

# Configure the Discord Application Provider
provider "discord-application" {
  token = "*your token here*"
}