FROM golang:1.18 AS builder

RUN mkdir /app
WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build addCron.go

FROM alpine:latest

RUN apk add --no-cache tzdata

RUN mkdir /app
WORKDIR /app

ENV TZ=America/Sao_Paulo

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=builder /app/addCron ./
COPY --from=builder /app/runner ./

CMD ["sh", "./runner"]

