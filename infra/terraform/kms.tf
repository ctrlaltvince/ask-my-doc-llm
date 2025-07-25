resource "aws_kms_key" "s3_upload_key" {
  description             = "KMS key for encrypting uploaded files in S3"
  enable_key_rotation     = true
  deletion_window_in_days = 7

  policy = jsonencode({
    Version = "2012-10-17",
    Id      = "key-default-1",
    Statement : [
      {
        Sid : "AllowKeyAdminAccess",
        Effect : "Allow",
        Principal : {
          AWS : "arn:aws:iam::242650469816:root"
        },
        Action : "kms:*",
        Resource : "*"
      },
      {
        Sid : "AllowS3Encryption",
        Effect : "Allow",
        Principal : {
          Service : "s3.amazonaws.com"
        },
        Action : [
          "kms:Encrypt",
          "kms:Decrypt",
          "kms:ReEncrypt*",
          "kms:GenerateDataKey*",
          "kms:DescribeKey"
        ],
        Resource : "*"
      },
      {
        Sid : "AllowAskMyDocUser",
        Effect : "Allow",
        Principal : {
          AWS : "arn:aws:iam::242650469816:user/ask-my-doc-cli"
        },
        Action : [
          "kms:Decrypt",
          "kms:Encrypt",
          "kms:ReEncrypt*",
          "kms:GenerateDataKey*",
          "kms:DescribeKey"
        ],
        Resource : "*"
      },
      {
        Sid : "AllowBackendServicePod",
        Effect : "Allow",
        Principal : {
          AWS : "arn:aws:iam::242650469816:role/backend-service-role"
        },
        Action : [
          "kms:Decrypt",
          "kms:Encrypt",
          "kms:ReEncrypt*",
          "kms:GenerateDataKey*",
          "kms:DescribeKey"
        ],
        Resource : "*"
      }
    ]
  })

  tags = {
    Name        = "ask-my-doc-s3-kms-key"
    Environment = "dev"
  }
}

resource "aws_kms_alias" "s3_upload_alias" {
  name          = "alias/ask-my-doc-s3"
  target_key_id = aws_kms_key.s3_upload_key.key_id
}

