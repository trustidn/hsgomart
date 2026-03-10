# ---------- Stage 1: Build frontend ----------
  FROM node:20-alpine AS frontend

  WORKDIR /app/frontend
  
  # copy dependency terlebih dahulu untuk docker cache
  COPY frontend/package*.json ./
  
  # gunakan mirror agar lebih cepat dari Asia
  RUN npm config set registry https://registry.npmmirror.com \
      && npm ci
  
  # copy source
  COPY frontend/ ./
  
  # build production
  RUN npm run build
  
  
  # ---------- Stage 2: Build Go binary ----------
  FROM golang:1.26-alpine AS backend
  
  RUN apk add --no-cache git
  
  WORKDIR /app
  
  # cache dependency
  COPY go.mod go.sum ./
  RUN go mod download
  
  # copy source
  COPY . .
  
  # copy hasil build frontend
  COPY --from=frontend /app/frontend/dist ./frontend/dist
  
  # build binary server
  RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
      go build -ldflags="-s -w" -o hsmart-server ./cmd/server
  
  # build seed binary
  RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
      go build -ldflags="-s -w" -o hsmart-seed ./cmd/seed
  
  
  # ---------- Stage 3: Runtime ----------
  FROM alpine:3.20
  
  RUN apk add --no-cache \
      ca-certificates \
      tzdata \
      curl
  
  # buat user non root
  RUN addgroup -S appgroup \
      && adduser -S appuser -G appgroup
  
  WORKDIR /app
  
  # copy binary
  COPY --from=backend /app/hsmart-server .
  COPY --from=backend /app/hsmart-seed .
  
  # copy frontend static
  COPY --from=backend /app/frontend/dist ./frontend/dist
  
  # copy migration
  COPY migrations ./migrations
  
  # upload folder
  RUN mkdir -p uploads/logos \
      && chown -R appuser:appgroup /app
  
  USER appuser
  
  EXPOSE 8080
  
  HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1
  
  CMD ["./hsmart-server"]