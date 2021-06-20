FROM alpine:latest

ADD k8s-labels-validation-webhook /k8s-labels-validation-webhook
ENTRYPOINT ["./k8s-labels-validation-webhook"]
