#!/bin/bash
# bootstrap.sh
# Use locally only if DynamoDB table needs to be re-created
mkdir terraform-bootstrap
cd terraform-bootstrap
terraform init
terraform apply -auto-approve
cd ..
rm -rf terraform-bootstrap
