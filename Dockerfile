FROM golang:1.9-alpine AS builder
WORKDIR /usr/src/app
COPY proxy.go .
RUN go build -o proxy .

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
USER nobody
COPY --from=builder /usr/src/app/proxy /proxy
ENTRYPOINT ["/proxy"]
