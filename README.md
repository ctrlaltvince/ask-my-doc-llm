# ask-my-doc-llm

Secure web app that lets users log in with Google, upload a document, and ask questions about it using OpenAI.

## 💻 Tech Stack

- **Go (Gin)** – REST API backend
- **React + Vite** – Frontend
- **AWS Cognito (Google login)** – Authentication
- **AWS S3** – Secure file storage (with encryption at rest)
- **AWS EKS + Fargate** – Kubernetes-based deployment
- **Terraform** – Infrastructure as Code (EKS, Cognito, S3, IAM, etc.)
- **Helm** – Kubernetes app packaging and deployment

## 🚀 Setup

```bash
# Run backend
cd backend
go run main.go

# Run frontend
cd frontend
npm install
npm run dev
