FROM golang:alpine as gota-build

# build gota
RUN apk add git && go get -u github.com/skanehira/gota

# build gota docker image
FROM alpine:latest
COPY --from=gota-build /go/bin/gota /usr/local/bin/gota
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["gota"]
