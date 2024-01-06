terraform {
  required_version = ">= 1.5.7"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.31.0"
    }
    discord-application = {
      source = "MichaelFraser99/discord-application"
      version = "0.2.2"
    }
    terracurl = {
      source = "devops-rob/terracurl"
      version = "1.1.0"
    }
  }
  backend "s3" {
    bucket         = "harrymooredev-tf-state"
    key            = "tf-discord-bot-deployment/terraform.tfstate"
    region         = "eu-west-1"
    encrypt        = true
    dynamodb_table = "tf-state"
  }
}

# Configure the AWS Provider
provider "aws" {
  region = "eu-west-2"
}

# Configure the Discord Application Provider
provider "discord-application" {
  token = var.application_secret
}

provider "terracurl" {}  
