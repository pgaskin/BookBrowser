# multi-stage build dockerfile:
# first build the app binary
# then put it into a light final image to deploy it

FROM golang:alpine as build
RUN apk update && apk add git
WORKDIR /go/src/app
COPY . .
RUN go get -v
RUN go generate
RUN go build

FROM alpine:latest
RUN mkdir /books
COPY --from=build /go/src/app/app /app
ENTRYPOINT ["/app", "--bookdir", "/books"]
