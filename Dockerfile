FROM golang:1.11-alpine as builder

ADD ./ /go/src/github.com/vicanso/tiny-site

RUN apk update \
  && apk add git \
  && go get -u github.com/golang/dep/cmd/dep \
  && cd /go/src/github.com/vicanso/tiny-site \
  && dep ensure \
  && GOOS=linux GOARCH=amd64 go build -tags netgo -o tiny-site 

FROM alpine

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/github.com/vicanso/tiny-site/tiny-site  /usr/local/bin/tiny-site 

CMD [ "tiny-site" ]

