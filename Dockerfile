FROM arm64v8/alpine:latest

ADD synology-videostation-reindexer /

ENTRYPOINT ["/synology-videostation-reindexer"]
