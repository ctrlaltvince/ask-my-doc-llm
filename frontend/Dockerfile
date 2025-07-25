# syntax=docker/dockerfile:1

FROM --platform=linux/amd64 node:20 AS builder
WORKDIR /app
COPY . .
RUN npm install && npm run build

FROM --platform=linux/amd64 nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80