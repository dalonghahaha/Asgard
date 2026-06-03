 # 多阶段构建 web-admin 静态资源（基于 nginx 托管 dist/）
 FROM node:20-alpine AS builder
 WORKDIR /app
 COPY package.json package-lock.json* ./
 RUN npm install
 COPY . .
 RUN npm run build

 FROM nginx:1.27-alpine
 COPY --from=builder /app/dist /usr/share/nginx/html
 COPY deploy/nginx.conf /etc/nginx/conf.d/default.conf
 EXPOSE 80
 CMD ["nginx", "-g", "daemon off;"]
