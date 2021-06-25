# k8s-labels-validation-webhook
Kubernetes Admission Validation Webhook enforcing recommended set of labels for pods

This is a code described in [this blogpost](https://bojan.tech/kubernetes/devops/automation/golang/2021/06/25/enforcing-best-practices-with-k8s-admission-controller-webhooks.html)

## Deploy

1. Install cert-manager `make install-cert-manager`
2. Deploy issuer and certificate `make deploy-cert`
3. Deploy app and service: `make deploy`
