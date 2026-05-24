FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/web ./web

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db

CMD ["./server"]
