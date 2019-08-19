FROM golang:1.12-alpine as builder
ADD . /go/src/github.com/epimorphics/prometheus-sns-webhook
ENV GO111MODULE=on
WORKDIR /go/src/github.com/epimorphics/prometheus-sns-webhook
RUN apk add --update --no-cache git alpine-sdk && go mod download && go build cmd/main.go

FROM alpine:3.8
RUN apk add --update --no-cache ca-certificates  
ADD configs/prometheus-sns-webhook.yaml /etc/prometheusns/
COPY --from=builder /go/src/github.com/epimorphics/prometheus-sns-webhook/main main
CMD ./main
