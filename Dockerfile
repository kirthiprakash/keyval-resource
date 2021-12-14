FROM golang:alpine as ginkgo
RUN go get -u "github.com/onsi/ginkgo/ginkgo"



FROM ginkgo as builder
ENV CGO_ENABLED 0

WORKDIR /keyval-resource

COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o /assets/out ./out \
 && go build -o /assets/in ./in \
 && go build -o /assets/check ./check
# RUN set -e; for pkg in $(go list ./...); do \
# 		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
# 	done
RUN ACK_GINKGO_RC="true" ginkgo -r -progress .



FROM alpine:latest AS resource
RUN apk add --update bash tzdata
COPY --from=builder /assets /opt/resource



FROM resource AS tests
# COPY --from=builder /tests /tests
# RUN set -e; for test in /tests/*.test; do \
# 		$test; \
# 	done



FROM resource
