# Multi-stage build for the tge CLI. The functional tests build this image with
# testcontainers and exec commands against the resulting binary.
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Pure-Go SQLite driver: CGO is not required, so we ship a static binary.
RUN CGO_ENABLED=0 go build -o /tge ./cmd/tge

FROM alpine:3.20
COPY --from=build /tge /usr/local/bin/tge
ENV TGE_DB=/data/tge.db
RUN mkdir -p /data
# Keep the container alive so the tests can exec multiple CLI commands that
# share the same on-disk SQLite database.
ENTRYPOINT ["tail", "-f", "/dev/null"]
