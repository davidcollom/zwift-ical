FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY . .
RUN go build -o zwiftcal ./cmd/zwiftcal/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/zwiftcal /app/zwiftcal
EXPOSE 3000
CMD ["/app/zwiftcal"]
