FROM arm32v7/alpine:latest

ADD TorrentFetcher /

ENTRYPOINT ["/TorrentFetcher"]
