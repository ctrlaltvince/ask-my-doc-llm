#!/bin/bash
terraform apply -target=aws_dynamodb_table.terraform_locks -auto-approve
terraform init
terraform apply -auto-approve

