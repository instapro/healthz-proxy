FROM golang:1.16-alpine AS builder
WORKDIR /usr/src/app
COPY go.mod proxy.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o proxy .

FROM scratch
USER nobody:nogroup
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /usr/src/app/proxy /proxy
ENTRYPOINT ["/proxy"]
