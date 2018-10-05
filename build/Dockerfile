FROM golang:alpine as gota-build

# build gota
RUN apk add git && go get -u github.com/skanehira/gota && go get -u github.com/samfoo/mdcat

# build gota docker image
FROM alpine:latest
COPY --from=gota-build /go/bin/* /usr/local/bin/
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["gota"]
