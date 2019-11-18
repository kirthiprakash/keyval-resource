FROM golang:alpine as builder
ENV CGO_ENABLED 0

WORKDIR /keyval-resource

COPY go.sum .
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /assets/out ./out
RUN go build -o /assets/in ./in
RUN go build -o /assets/check ./check
# RUN set -e; for pkg in $(go list ./...); do \
# 		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
# 	done

FROM alpine:edge AS resource
RUN apk add --update bash tzdata
COPY --from=builder /assets /opt/resource

FROM resource AS tests
# COPY --from=builder /tests /tests
# RUN set -e; for test in /tests/*.test; do \
# 		$test; \
# 	done

FROM resource
