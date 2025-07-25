resource "helm_release" "frontend" {
  name      = "frontend"
  namespace = "default"
  chart     = "../../charts/frontend"

  set {
    name  = "image.repository"
    value = "242650469816.dkr.ecr.us-west-1.amazonaws.com/ask-my-doc-frontend"
  }

  set {
    name  = "image.tag"
    value = "latest"
  }

  depends_on = [
    aws_eks_fargate_profile.default
  ]
}

