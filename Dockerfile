FROM golang:1-alpine as builder
RUN apk add --no-cache --virtual .build-deps \
    gcc libc-dev git
COPY ./ /src
WORKDIR /src
RUN go build -o /bin/proxet

FROM alpine:3
COPY --from=builder /bin/proxet /bin/proxet
ENTRYPOINT [ "/bin/proxet" ]