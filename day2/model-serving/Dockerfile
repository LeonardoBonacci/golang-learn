FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod tidy
COPY . .
RUN go build -o model-serving ./cmd/app

FROM golang:1.24
WORKDIR /app
COPY --from=builder /app/model-serving .
EXPOSE 8080
CMD ["./model-serving"]