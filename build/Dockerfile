FROM golang:alpine as gota-build

# build gota
RUN apk add git && \
    git clone https://github.com/skanehira/gota && \
    cd gota && go install

# build gota docker image
FROM alpine:latest
COPY --from=gota-build /go/bin/gota /usr/local/bin/
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["gota"]
