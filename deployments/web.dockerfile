FROM golang:1.21.5-bullseye
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./configs ./configs
CMD go run cmd/goFipe/main.go