FROM docker.io/library/golang:1.10.1-alpine AS k8s-term-delay-builder
RUN  apk add --no-cache upx
WORKDIR /go/src/github.com/steigr/k8s-term-delay
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go get -v "${PWD#*/src/}"
ARG UPX_ARGS=-6
ENV UPX_ARGS=$UPX_ARGS
RUN upx ${UPX_ARGS} /go/bin/k8s-term-delay

FROM alpine:3.7 AS k8s-term-delay
COPY --from=k8s-term-delay-builder /go/bin/k8s-term-delay /bin/k8s-term-delay