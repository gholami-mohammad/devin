FROM golang:1.9.3-alpine3.7
RUN apk update && apk add git

WORKDIR /go/src/gogit
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 80
CMD ["gogit"]
