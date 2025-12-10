FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# Runtime stage
FROM gcr.io/distroless/static-debian12

COPY --from=builder /out/api /api
COPY templates ./templates

EXPOSE 8000
ENTRYPOINT ["/api"]
