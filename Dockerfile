FROM golang:1.23 AS server-builder

WORKDIR /code

COPY go.mod go.sum ./
COPY server/ ./server
COPY gen ./gen

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./server/cmd/server/main.go

FROM node:22.9.0-slim AS web-builder

WORKDIR /code

COPY web/ .

RUN npm install
RUN npm run build

FROM gcr.io/distroless/static-debian12:latest AS app

WORKDIR /app

COPY --from=server-builder /server .
COPY --from=web-builder /code/build ./static

# Dev container
#ENV FIRESTORE_EMULATOR_HOST=host.docker.internal:8081

EXPOSE 8080

CMD ["./server"]