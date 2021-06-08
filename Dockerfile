FROM golang:1.16-alpine AS builder
WORKDIR /usr/src/app
COPY go.mod proxy.go ./
RUN go build -o proxy .

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
USER nobody
COPY --from=builder /usr/src/app/proxy /proxy
ENTRYPOINT ["/proxy"]
