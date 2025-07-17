# ask-my-doc-llm

Secure web app that lets users log in with Google, upload a document, and ask questions about it using OpenAI.

## Tech Stack

- Go (Gin) backend
- React + Vite frontend
- AWS Cognito (Google login)
- AWS S3 (file storage)
- AWS EKS + Fargate (deployment)
- Helm (infrastructure as code)

## Setup

```bash
# Run backend
cd backend
go run main.go

# Run frontend
cd frontend
npm install
npm run dev
