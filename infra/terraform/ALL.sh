#!/bin/bash

# Define the files you want to display
files=(
  "coredns-patch.yaml"
  "eks.tf"
  "vpc.tf"
  "alb-controller.tf"
  "coredns.tf"
  "iam.tf"
  "backend-serviceaccount.yaml"
  "debug-pod.yaml"
  "outputs.tf"
  "backend.tf"
  "dynamo.db"
  "provider.tf"
)

# Loop through each file and display the name and content
for file in "${files[@]}"; do
  if [[ -f "$file" ]]; then
    echo "=== $file ==="
    cat "$file"
    echo ""
  else
    echo "File $file does not exist."
  fi
done

