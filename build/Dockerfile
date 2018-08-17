FROM golang:1.11rc1-alpine as builder
ADD . /go/src/github.com/epimorphics/prometheus-sns-webhook
RUN go build src/github.com/epimorphics/prometheus-sns-webhook/cmd/main.go

FROM alpine:3.8
RUN apk add --update --no-cache ca-certificates  
ADD configs/prometheus-sns-webhook.yaml /etc/prometheusns/
COPY --from=builder /go/main main
CMD ./main
