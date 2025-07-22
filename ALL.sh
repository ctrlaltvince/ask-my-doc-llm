#!/bin/bash

# Define the directories and files to search
root_dir="."

# Use find to list all files recursively in the root directory and subdirectories
find "$root_dir" -type f | while read file; do
  echo "=== $file ==="
  cat "$file"
  echo ""
done

