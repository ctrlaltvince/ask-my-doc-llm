# ask-my-doc-llm

Secure web app that lets users log in with Google, upload a document, and ask questions about it using OpenAI.

## Tech Stack

- **Go (Gin)** – REST API backend
- **React + Vite** – Frontend
- **AWS Cognito (Google login)** – Authentication
- **AWS S3** – Secure file storage, with encryption at rest using **KMS**
- **AWS Secrets Manager** – Secrets management for sensitive credentials (OpenAI key, Cognito client info)
- **AWS EKS + Fargate** – Kubernetes-based deployment
- **Terraform** – Infrastructure as Code (EKS, Cognito, S3, IAM, etc.) using an **S3 remote backend with state locking via DynamoDB**
- **Helm** – Kubernetes app packaging and deployment

---

## 🚀 Setup Installation

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

Then update your kubeconfig to talk to the EKS cluster:
```bash
aws eks update-kubeconfig --region us-west-1 --name ask-my-doc-cluster
```

Finally, deploy the backend and frontend via Helm:
```bash
helm upgrade --install backend ./backend
helm upgrade --install frontend ./frontend

```

## 🚨 Setup: Manual and Recovery Instructions
Some AWS resources must be created or managed manually, especially if you're recovering from a partial destroy:

1. If terraform apply fails due to existing resources:
	- You may need to manually delete IAM resources (role, policy, etc.)
	- Or import the existing ones into Terraform
2. Common Terraform fixes:
	- Re-import OIDC provider if missing:
	```bash
	terraform import aws_iam_openid_connect_provider.oidc arn:aws:iam::<ACCOUNT_ID>:oidc-provider/oidc.eks.us-west-1.amazonaws.com/id/<OIDC_ID>
	```
	- If S3 bucket exists:
	```bash
	terraform import aws_s3_bucket.uploads ask-my-doc-llm-files
	```
	- If IAM role already exists:
	```bash
	terraform import aws_iam_role.alb_controller AmazonEKSLoadBalancerControllerRole
	```

3. DNS fix: If your domain (e.g., askmydoc.dev) stops resolving:
	- Delete the existing A record in Route 53 and let ALB recreate it.

4. Secrets must be created manually via AWS CLI:
	```bash
	aws secretsmanager create-secret \
  	--name askmydoc/backend \
  	--secret-string '{
    	"CLIENT_ID": "YOUR_CLIENT_ID",
    	"CLIENT_SECRET": "YOUR_CLIENT_SECRET",
    	"OPENAI_API_KEY": "YOUR_OPENAI_API_KEY"
  	}'
	```

## ⚠️ Destroy Notes
Terraform does not delete the following critical AWS resources:
	- DynamoDB table (Terraform backend locking)
	- AWS Secrets
	- IAM roles, policies, or OIDC provider
	- Route53 hosted zone / A record
	- S3 bucket

| Resource                       | Terraform Managed? | Manual Action Needed? | Estimated Monthly Cost |
|---------------------------------|-----------------------------|----------------------------------|---------------------------------|
| EKS Cluster                   |             ✅                 |                  ❌                  |                ~$72               |
| Fargate Profiles             |             ✅                 |                  ❌                  |     ~$0 (pay-per-use)      |
| Load Balancer (ALB)     |             ✅                 |                  ❌                  |                ~$18                |
| Route 53 Record           |             ✅                 |      ✅ (sometimes)          |                ~$0.50            |
| ACM Certificate             |             ✅                 |                  ❌                  |                $0 (free)           |
| S3 Bucket                      |             ✅                 |                  ✅                  |      Depends on usage    |
| DynamoDB Table          |             ✅                 |                  ✅                   |                ~$0.25             |
| Cognito User Pool         |             ✅                 |                  ❌                   |       ~$0.005 per MAU     |
| AWS Secrets Manager  |       ❌ (manual)        |                  ✅                   |   ~$0.40 per secret/mo   |
| IAM Roles & Policies     |             ✅                 |                  ✅ (if errors)    |                $0                    |