# multi-stage build dockerfile:
# first build the app binary
# then put it into a light final image to deploy it

FROM golang:alpine as build
RUN apk update && apk add git
ENV GO111MODULE=on
WORKDIR /go/src/app
COPY . .
RUN go generate
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM alpine:latest
RUN mkdir /books
COPY --from=build ./app /app
ENTRYPOINT ["/app", "--bookdir", "/books"]
