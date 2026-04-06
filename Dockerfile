# Stage 1: Build frontend
FROM node:22-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.23.8-alpine AS backend
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
ENV GOTOOLCHAIN=local
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=1 go build -o dashboard .

# Stage 3: Final image
FROM alpine:3.20
RUN apk add --no-cache ca-certificates docker-cli docker-compose
WORKDIR /app

COPY --from=backend /app/dashboard .
COPY --from=frontend /app/frontend/dist ./static

EXPOSE 3000
VOLUME /app/data

CMD ["./dashboard"]
