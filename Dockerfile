ARG DB
FROM golang:alpine
ARG DB
ENV DB=${DB}
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base
RUN mkdir /short-url
WORKDIR /short-url
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build ./cmd/shorturl
EXPOSE 9000
ENTRYPOINT shorturl $DB