FROM node:12-alpine as webbuilder

ADD . /tiny-site
RUN cd /tiny-site/web \
  && yarn \
  && yarn build \
  && rm -rf node_module

FROM golang:1.14-alpine as builder

COPY --from=webbuilder /tiny-site /tiny-site

RUN apk update \
  && apk add git make \
  && go get -u github.com/gobuffalo/packr/v2/packr2 \
  && cd /tiny-site \
  && make build

FROM alpine 

EXPOSE 7001

RUN addgroup -g 1000 go \
  && adduser -u 1000 -G go -s /bin/sh -D go \
  && apk add --no-cache ca-certificates

COPY --from=builder /tiny-site/tiny-site /usr/local/bin/tiny-site

USER go

WORKDIR /tiny-site

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 CMD [ "wget", "http://127.0.0.1:7001/ping", "-q", "-O", "-"]

CMD ["tiny-site"]
