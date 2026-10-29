# dev target
FROM node:current-alpine AS dev
WORKDIR /src
COPY --from=golang:1.27.3-alpine /usr/local/go/ /usr/local/go/
ENV GOPATH=/usr/local/go
ENV PATH="${GOPATH}/bin:${PATH}"
RUN go install github.com/jackc/tern/v2@v2.4.3
RUN go install github.com/cespare/reflex@latest
ENV TERN_CONFIG=/src/tern.docker.conf
ENV TERN_MIGRATIONS=/src/db/migrations
CMD ["reflex", "-d",  "none",  "-c", "reflex.docker.conf"]

# build stage
FROM golang:1.27.3-alpine AS build
WORKDIR /build
COPY . .
RUN go get -v ./...
RUN go build -o app -v
RUN GOBIN=/build/ go install github.com/jackc/tern/v2@v2.4.3

# final stage
FROM alpine:3.24
# hooks/build passes a new value on every build, so this layer is never served from cache.
ARG SNAPSHOT_EPOCH
RUN apk upgrade --no-cache

ARG SOURCE_BRANCH
ARG SOURCE_COMMIT
ARG IMAGE_NAME
ENV SOURCE_BRANCH=$SOURCE_BRANCH
ENV SOURCE_COMMIT=$SOURCE_COMMIT
ENV IMAGE_NAME=$IMAGE_NAME

ENV TERN_CONFIG=/dist/tern.docker.conf
ENV TERN_MIGRATIONS=/dist/db/migrations

WORKDIR /dist
COPY --from=build /build .
EXPOSE 3000
CMD ["/dist/app", "server"]
