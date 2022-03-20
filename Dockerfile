FROM node:16-alpine as webbuilder

COPY . /tiny-site
RUN cd /tiny-site/web \
  && npm i \
  && npm run build \
  && rm -rf node_module

FROM golang:1.17-alpine as builder

COPY --from=webbuilder /tiny-site /tiny-site

RUN apk update \
  && apk add git make curl jq \
  && cd /tiny-site \
  && rm -rf asset/dist \
  && cp -rf web/dist asset/ \
  && make install \
  && make generate \
  && ./download-swagger.sh \
  && make build

FROM alpine 

EXPOSE 7001

# tzdata 安装所有时区配置或可根据需要只添加所需时区

RUN addgroup -g 1000 go \
  && adduser -u 1000 -G go -s /bin/sh -D go \
  && apk add --no-cache ca-certificates tzdata

COPY --from=builder /tiny-site/tiny-site /usr/local/bin/tiny-site
COPY --from=builder /tiny-site/entrypoint.sh /entrypoint.sh

USER go

WORKDIR /home/go

HEALTHCHECK --timeout=10s --interval=10s CMD [ "wget", "http://127.0.0.1:7001/ping", "-q", "-O", "-"]

CMD ["tiny-site"]

ENTRYPOINT ["/entrypoint.sh"]
