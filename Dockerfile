FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY .server gen ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./server/cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]