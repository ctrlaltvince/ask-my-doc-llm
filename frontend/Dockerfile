# syntax=docker/dockerfile:1

FROM --platform=linux/amd64 node:20 AS builder
WORKDIR /app

# Copy only package files first to leverage Docker cache
COPY package*.json ./

# Install deps first (cached if unchanged)
RUN npm install

# Now copy the rest of the app
COPY . .

# Build the project
RUN npm run build

# Final lightweight image
FROM --platform=linux/amd64 nginx:alpine

# ✅ Copy everything Vite built, including favicon.ico
COPY --from=builder /app/dist /usr/share/nginx/html

# ✅ Only keep this if background.png is not in dist
# (But it should be — Vite copies from /public)
# COPY public/background.png /usr/share/nginx/html/background.png

# nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80