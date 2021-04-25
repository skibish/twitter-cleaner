FROM scratch

COPY twitter-cleaner /

ENTRYPOINT ["/twitter-cleaner"]
