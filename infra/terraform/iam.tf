data "aws_caller_identity" "current" {}

# Optional: Admin Role (not directly used by Fargate)
resource "aws_iam_role" "eks_admin" {
  name = "eksAdminRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Effect = "Allow",
      Principal = {
        AWS = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
      }
    }]
  })
}

# Required: Fargate Pod Execution Role
#resource "aws_iam_role" "fargate_pod_execution" {
#  name = "FargatePodExecutionRole"

#  assume_role_policy = jsonencode({
#    Version = "2012-10-17",
#    Statement = [
#      {
#        Effect = "Allow",
#        Principal = {
#          Service = "eks-fargate-pods.amazonaws.com"
#        },
#        Action = "sts:AssumeRole"
#      }
#    ]
#  })
#}

#resource "aws_iam_role_policy_attachment" "fargate_pod_execution" {
#  role       = aws_iam_role.fargate_pod_execution.name
#  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSFargatePodExecutionRolePolicy"
#}

resource "aws_iam_role" "eks_node" {
  name = "eksNodeRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Principal = {
        Service = "ec2.amazonaws.com"
      },
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "eks_node_AmazonEKSWorkerNodePolicy" {
  role       = aws_iam_role.eks_node.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
}

resource "aws_iam_role_policy_attachment" "eks_node_AmazonEC2ContainerRegistryReadOnly" {
  role       = aws_iam_role.eks_node.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_iam_role_policy_attachment" "eks_node_AmazonEKS_CNI_Policy" {
  role       = aws_iam_role.eks_node.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
}


