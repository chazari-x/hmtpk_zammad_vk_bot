FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.10

RUN adduser -DH zammad-vk-bot

WORKDIR /app

COPY --from=builder /app/main /app/

RUN chown zammad-vk-bot:zammad-vk-bot /app
RUN chmod +x /app

USER zammad-vk-bot

CMD ["/app/main"]