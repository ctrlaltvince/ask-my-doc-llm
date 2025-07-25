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

# Copy built assets from builder
COPY --from=builder /app/dist /usr/share/nginx/html
COPY --from=builder /app/dist/background.png /usr/share/nginx/html/background.png


# Custom nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

# Optional: set non-root user for security (if nginx supports it)
# USER nginx
