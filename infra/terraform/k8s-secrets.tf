data "aws_secretsmanager_secret_version" "backend" {
  secret_id = "askmydoc/backend"
}

locals {
  backend_secrets = jsondecode(data.aws_secretsmanager_secret_version.backend.secret_string)
}

resource "kubernetes_secret" "backend_secrets" {
  metadata {
    name      = "backend-secrets"
    namespace = "default"
  }

  data = {
    clientId     = base64encode(local.backend_secrets["CLIENT_ID"])
    clientSecret = base64encode(local.backend_secrets["CLIENT_SECRET"])
  }

  type = "Opaque"
}

resource "kubernetes_secret" "openai_secret" {
  metadata {
    name      = "openai-secret"
    namespace = "default"
  }

  data = {
    OPENAI_API_KEY = base64encode(local.backend_secrets["OPENAI_API_KEY"])
  }

  type = "Opaque"
}

