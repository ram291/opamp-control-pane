# Stage 1: Build React frontend
FROM node:20-alpine AS ui-builder
WORKDIR /app
COPY ui/package*.json ./
RUN npm ci
COPY ui/ .
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.22-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy the built React frontend
COPY --from=ui-builder /app/dist ./internal/embedfs/ui/dist/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o supervisor ./cmd/supervisor/

# Stage 3: Minimal runtime
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=go-builder /app/supervisor /usr/local/bin/supervisor
COPY configs/supervisor.example.yaml /etc/opamp/supervisor.yaml
EXPOSE 8080
ENTRYPOINT ["supervisor"]
CMD ["--config", "/etc/opamp/supervisor.yaml"]