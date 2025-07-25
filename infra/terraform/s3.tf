resource "aws_s3_bucket" "uploads" {
  bucket = "ask-my-doc-llm-files"

  tags = {
    Name        = "AskMyDoc File Uploads"
    Environment = "dev"
  }
}

resource "aws_s3_bucket_versioning" "uploads_versioning" {
  bucket = aws_s3_bucket.uploads.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "uploads_encryption" {
  bucket = aws_s3_bucket.uploads.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = "aws:kms"
      kms_master_key_id = aws_kms_key.s3_upload_key.arn
    }
  }
}
