# dev target
FROM node:current-alpine as dev
WORKDIR /src
COPY --from=golang:alpine /usr/local/go/ /usr/local/go/
ENV GOPATH=/usr/local/go
ENV PATH="${GOPATH}/bin:${PATH}"
RUN go install github.com/jackc/tern/v2@latest
RUN go install github.com/cespare/reflex@latest
ENV TERN_CONFIG /src/tern.docker.conf
ENV TERN_MIGRATIONS /src/db/migrations
CMD ["reflex", "-d",  "none",  "-c", "reflex.docker.conf"]

# build stage
FROM golang:alpine AS build
WORKDIR /build
COPY . .
RUN go get -d -v ./...
RUN go build -o app -v
RUN GOBIN=/build/ go install github.com/jackc/tern/v2@latest

# final stage
FROM alpine:latest

ARG SOURCE_BRANCH
ARG SOURCE_COMMIT
ARG IMAGE_NAME
ENV SOURCE_BRANCH $SOURCE_BRANCH
ENV SOURCE_COMMIT $SOURCE_COMMIT
ENV IMAGE_NAME $IMAGE_NAME

ENV TERN_CONFIG /dist/tern.docker.conf
ENV TERN_MIGRATIONS /dist/db/migrations

WORKDIR /dist
COPY --from=build /build .
EXPOSE 3000
CMD ["/dist/app", "server"]
