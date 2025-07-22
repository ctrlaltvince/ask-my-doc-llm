# Will fail unless the cluster and Fargate profile exist. Comment out until those resources are created.
resource "null_resource" "patch_coredns_toleration" {
  provisioner "local-exec" {
    command = <<EOT
      echo "Patching CoreDNS with Fargate toleration..."
      kubectl -n kube-system patch deployment coredns --type merge --patch '{
        "spec": {
          "template": {
            "spec": {
              "tolerations": [
                {
                  "key": "eks.amazonaws.com/compute-type",
                  "operator": "Equal",
                  "value": "fargate",
                  "effect": "NoSchedule"
                }
              ]
            }
          }
        }
      }'
    EOT
  }

  depends_on = [
    aws_eks_cluster.this,
    aws_eks_fargate_profile.kube_system
  ]
}

