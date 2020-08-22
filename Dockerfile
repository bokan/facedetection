# Build
FROM golang:1.15-alpine as gobuild

RUN apk --update upgrade \
    && apk --no-cache --no-progress add ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

COPY . /go/src/github.com/bokan/facedetection

WORKDIR /go/src/github.com/bokan/facedetection

RUN GO111MODULE=on GOPROXY=https://proxy.golang.org go mod download
RUN GO111MODULE=on CGO_ENABLED=0 go build -o ./dist/facedetection cmd/facedetection/facedetection.go

# Image
FROM scratch

COPY --from=gobuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=gobuild /go/src/github.com/bokan/facedetection/dist/facedetection /
COPY --from=gobuild /go/src/github.com/bokan/facedetection/pkg/facedetect/pigofacedetect/cascades /cascades

EXPOSE 8000
VOLUME ["/tmp"]

ENTRYPOINT ["/facedetection", "-c", "/cascades"]
