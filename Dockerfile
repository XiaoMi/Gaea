# Build step
FROM micr.cloud.mioffice.cn/dockerhub/library/golang:1.16.15 AS builder
ENV LC_ALL en_US.UTF-8
ENV LANG en_US.UTF-8
RUN mkdir -p $GOPATH/src/gaea
ADD . $GOPATH/src/gaea
WORKDIR $GOPATH/src/gaea
RUN pwd && make gaea && mkdir /build &&  cp -r  $GOPATH/src/gaea/bin/gaea /build/gaea

# Final step
FROM micr.cloud.mioffice.cn/devdba/alpine:3.17
ENV LC_ALL en_US.UTF-8
ENV LANG en_US.UTF-8
RUN pwd && mkdir -p /home/work/gaea/bin && mkdir -p /home/work/gaea/etc
COPY --from=builder /build/gaea /home/work/gaea/bin/gaea
WORKDIR /home/work/gaea/
ENTRYPOINT /home/work/gaea/bin/gaea -config /home/work/gaea/etc/gaea.ini
