terraform {
  backend "s3" {
    bucket         = "ask-my-doc-llm-terraform-state"
    key            = "eks/terraform.tfstate"
    region         = "us-west-1"
    encrypt        = true
    dynamodb_table = "terraform-locks" # Optional: for state locking
  }
}

