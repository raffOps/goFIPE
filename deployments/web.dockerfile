FROM golang:1.21.5-bullseye
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY cmd ./cmd
CMD go run cmd/gofipe/main.go