# update together with ci.yml
FROM golang:1.19.1 AS go

# update together with ci.yml
FROM golangci/golangci-lint:v1.49.0 AS linter


FROM go AS dev
ENV INSIDE_DEV_CONTAINER 1
COPY --from=linter /usr/bin/golangci-lint /usr/bin/
WORKDIR /app
