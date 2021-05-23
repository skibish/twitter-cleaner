FROM alpine:latest as builder

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY twitter-cleaner /

ENTRYPOINT ["/twitter-cleaner"]
