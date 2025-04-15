FROM golang:1.24.2-alpine as mybuilder
ENV GO111MODULE=on CGO_ENABLED=0
WORKDIR /syncdns
ADD . /syncdns
RUN go build .

FROM busybox:uclibc
COPY --from=mybuilder /syncdns/syncdns /syncdns
ENTRYPOINT [ "/syncdns" ]
CMD [ "-c", "/etc/config.yaml" ]
