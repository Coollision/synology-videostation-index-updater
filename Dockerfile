FROM arm32v7/alpine:latest

ADD synology-videostation-reindexer /

ENTRYPOINT ["/synology-videostation-reindexer"]
