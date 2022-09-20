# update together with ci.yml
FROM golang:1.19.1 AS go

# update together with ci.yml
FROM golangci/golangci-lint:v1.49.0 AS linter

FROM go as sql-migrate
RUN go install github.com/rubenv/sql-migrate/sql-migrate@v1.1.1


FROM go AS dev
ENV INSIDE_DEV_CONTAINER 1
RUN apt-get update && apt-get install -y netcat && rm -rf /var/lib/apt/lists/*
COPY --from=linter      /usr/bin/golangci-lint /usr/bin/
COPY --from=sql-migrate /go/bin/sql-migrate    /usr/bin/
WORKDIR /app
