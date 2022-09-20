# update together with ci.yml
FROM golang:1.19.1 AS go

# update together with ci.yml
FROM golangci/golangci-lint:v1.49.0 AS linter

FROM kjconroy/sqlc:1.15.0 AS sqlc

FROM go as sql-migrate
RUN go install github.com/rubenv/sql-migrate/sql-migrate@v1.1.1

FROM go AS protoc
ARG TARGETARCH
RUN VERSION=21.6 && \
    apt-get update && \
    apt-get install unzip && \
    if [ "$TARGETARCH" = "arm64" ]; then ARCH="aarch_64"; else ARCH="x86_64"; fi; \
    wget -O x.zip https://github.com/protocolbuffers/protobuf/releases/download/v$VERSION/protoc-$VERSION-linux-$ARCH.zip && \
    unzip x.zip -d protoc

FROM go AS protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

FROM go AS twirp
RUN go install github.com/twitchtv/twirp/protoc-gen-twirp@v8.1.2


FROM go AS dev
ENV INSIDE_DEV_CONTAINER 1
RUN apt-get update && apt-get install -y netcat && rm -rf /var/lib/apt/lists/*
COPY --from=linter      /usr/bin/golangci-lint   /usr/bin/
COPY --from=sqlc        /workspace/sqlc          /usr/bin/
COPY --from=sql-migrate /go/bin/sql-migrate      /usr/bin/
COPY --from=protoc      /go/protoc/              /usr/
COPY --from=protobuf    /go/bin/protoc-gen-go    /usr/bin/
COPY --from=twirp       /go/bin/protoc-gen-twirp /usr/bin/
WORKDIR /app
