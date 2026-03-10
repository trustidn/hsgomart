# --- Stage 1: Build frontend ---
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install --include=dev
COPY frontend/ ./
RUN npx vite build

# --- Stage 2: Build Go binary ---
FROM golang:1.26-alpine AS backend
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o hsmart-server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o hsmart-seed ./cmd/seed

# --- Stage 3: Minimal runtime ---
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata curl \
    && addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app

COPY --from=backend /app/hsmart-server .
COPY --from=backend /app/hsmart-seed .
COPY --from=backend /app/frontend/dist ./frontend/dist
COPY migrations/ ./migrations/

RUN mkdir -p uploads/logos && chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./hsmart-server"]
