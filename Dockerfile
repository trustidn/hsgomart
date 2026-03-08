# --- Stage 1: Build frontend ---
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# --- Stage 2: Build Go binary ---
FROM golang:1.23-alpine AS backend
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o hsmart-server ./cmd/server

# --- Stage 3: Minimal runtime ---
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=backend /app/hsmart-server .
COPY --from=backend /app/frontend/dist ./frontend/dist
COPY migrations/ ./migrations/

RUN mkdir -p uploads/logos

EXPOSE 8080

CMD ["./hsmart-server"]
