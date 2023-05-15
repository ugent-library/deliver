# build stage
FROM golang:alpine AS build
WORKDIR /build
COPY . .
RUN go get -d -v ./...
RUN go build -o app -v

# final stage
ARG SOURCE_BRANCH
ARG SOURCE_COMMIT
ENV SOURCE_BRANCH $SOURCE_BRANCH
ENV SOURCE_COMMIT $SOURCE_COMMIT
FROM alpine:latest
WORKDIR /dist
COPY --from=build /build .
EXPOSE 3000
CMD ["/dist/app", "app"]
