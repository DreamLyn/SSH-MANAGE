FROM node:20-alpine3.19 AS front-builder

WORKDIR /app

COPY . /app/

RUN \
  cd /app/ui && \
  npm install && \
  npm run build:prod


FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ../. /app/

RUN rm -rf /app/ui/dist

COPY --from=front-builder /app/ui/dist /app/ui/dist

RUN go build -o ssh-manage



FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ssh-manage .

ENTRYPOINT ["./ssh-manage", "serve", "--http", "0.0.0.0:8090"]