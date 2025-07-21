# Development stage with hot-reloading
FROM node:18-alpine AS development

WORKDIR /app

# Copy package files first for better caching
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install

# Copy all files
COPY . .

# Expose Vite's default port
EXPOSE 5173

# Start development server
CMD ["npm", "run", "dev"]

# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .
RUN npm run build

# Production stage
FROM nginx:alpine AS production

# Copy nginx config
COPY --from=builder /app/nginx.conf /etc/nginx/conf.d/default.conf

# Copy built assets
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]