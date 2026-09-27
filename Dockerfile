FROM golang:1.21 AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /scheduler .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /scheduler /app/scheduler
COPY web /app/web
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
VOLUME ["/data"]
CMD ["/app/scheduler"]
