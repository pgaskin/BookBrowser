FROM golang:1.13-alpine AS build
ADD . /src
WORKDIR /src
RUN apk add --no-cache git
RUN go mod download
RUN go build .

FROM alpine:latest
RUN mkdir /books
COPY --from=build /src/BookBrowser /
EXPOSE 8090
ENTRYPOINT ["/BookBrowser", "--bookdir", "/books"]