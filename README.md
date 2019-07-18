# k8s-delete-validation-webhook
AdmissionControl validating webhook to block deletion of resources based on resource labels.

## Once deployed what this solution does?
For Configured Kubernetes namespaces and resources (deployments, crds etc) deletion requests are first processed by this webhook, and based on labels associated with the resource the resource deletion requests (kubectl delete resource .) are either rejected or allowed. The label keys and values for which resource deletion requests are rejected can also be configured.

## Referenced Respository

## Solution Components

## Key Functions of webhook implementation

## Installation / deployment

## Configuring the namespaces, resources and labels
