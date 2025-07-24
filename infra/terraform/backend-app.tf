resource "helm_release" "backend" {
  name       = "backend"
  namespace  = "default"
  chart = "../../charts/backend"

  set = [
    {
      name  = "env.clientId"
      value = "39u7iped9gp9cfnfutjp1ras8b"
    },
    {
      name  = "env.clientSecret"
      value = "22hgbmveqbd36du39hbg43hgs18nm9vtjmqlop13o165b9ea3kj"
    }
  ]

  depends_on = [
    aws_eks_fargate_profile.default
  ]

}

