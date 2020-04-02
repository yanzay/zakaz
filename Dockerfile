FROM golang:latest as builder
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
RUN go build -v .

FROM alpine:latest
RUN apk update && apk add ca-certificates
COPY --from=builder /app/zakaz /zakaz
CMD ["/zakaz"]
