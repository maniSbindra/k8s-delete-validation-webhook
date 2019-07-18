FROM golang:1.10-stretch AS builder

RUN curl https://glide.sh/get | sh

WORKDIR /go/src/github.com/manisbindra/k8s-delete-validation-webhook

COPY internal/k8s-delete-validation-webhook/glide.* /go/src/github.com/manisbindra/k8s-delete-validation-webhook/
RUN glide install -v

COPY internal/k8s-delete-validation-webhook/*.* ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

FROM alpine:latest

COPY --from=builder /go/src/github.com/manisbindra/k8s-delete-validation-webhook/k8s-delete-validation-webhook /k8s-delete-validation-webhook/internal/k8s-delete-validation-webhook

ENTRYPOINT ["/k8s-delete-validation-webhook/internal/k8s-delete-validation-webhook", "webhook"]
