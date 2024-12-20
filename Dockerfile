FROM golang:alpine AS builder
ARG buildtags=""
ARG version="none given"

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN apk add git
RUN go build -v -tags "$buildtags" -ldflags="-X main.version=$version" -o synology-videostation-reindexer2 .

FROM alpine
COPY --from=builder /app/synology-videostation-reindexer .
ENTRYPOINT ["/synology-videostation-reindexer"]
