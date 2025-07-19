provider "aws" {
  region = "us-west-1"
}

resource "aws_eks_cluster" "this" {
  name     = "ask-my-doc-cluster"
  role_arn = aws_iam_role.eks_cluster.arn

  vpc_config {
    subnet_ids = data.aws_subnets.default.ids
  }

  version = "1.29"
}

resource "aws_iam_role" "eks_cluster" {
  name = "eksClusterRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Principal = {
        Service = "eks.amazonaws.com"
      },
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "eks_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.eks_cluster.name
}


#resource "aws_eks_fargate_profile" "default" {
#  cluster_name           = aws_eks_cluster.this.name
#  fargate_profile_name   = "default"
#  pod_execution_role_arn = aws_iam_role.fargate_pod_execution.arn
#  subnet_ids             = data.aws_subnets.default.ids

#  selector {
#    namespace = "default"
#  }

#  depends_on = [aws_eks_cluster.this]
#}

resource "aws_eks_node_group" "default" {
  cluster_name    = aws_eks_cluster.this.name
  node_group_name = "ask-my-doc-node-group"
  node_role_arn   = aws_iam_role.eks_node.arn
  subnet_ids      = data.aws_subnets.default.ids

  scaling_config {
    desired_size = 1
    max_size     = 2
    min_size     = 1
  }

  depends_on = [
    aws_eks_cluster.this,
    aws_iam_role_policy_attachment.eks_node_AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.eks_node_AmazonEC2ContainerRegistryReadOnly,
    aws_iam_role_policy_attachment.eks_node_AmazonEKS_CNI_Policy,
  ]
}

