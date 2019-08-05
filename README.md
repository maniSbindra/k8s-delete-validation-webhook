# Kubernetes resource deletion validation webhook
AdmissionControl validating webhook to block deletion of resources based on resource labels.

## What this solution does
For Configured Kubernetes namespaces and resources (deployments, crds etc) deletion requests are first processed by this webhook, and based on labels associated with the resource the resource deletion requests (kubectl delete resource .) are either rejected or allowed. The label keys and values for which resource deletion requests are rejected can also be configured.

## Solution Overview
![alt](./images/deletion-validation-webhook.png)
* In the above diagram the validating webhook has been configured to monitor **deployments** and **myCRDs** in the **default** namespace
* So if deletion request for all 6 resources shown in the diagram above are made to the kubernetes API, the requests for the 2 deployments and myCRD resource in the default namespace will be forwarded to our webhook. 
* Since the webhook has been configured to reject resource deletion requests where resources have a label **deleteLock=enabled** the deletion requests for the deployment and myCRD highlighted in orange will be rejected.
* For the kubernetes core components to communicate with the webhook, TLS certificate and TLS key need to be created and associated with the webhook

## Reference solution
* The [avast/k8s-admission-webhook](https://github.com/avast/k8s-admission-webhook) is a good reference for create and update resource validations. 
* The avast solution has been referenced for the initialization / configuration of this application (using cobra, viper etc). The script to generate the certificate key files has also been referenced

## Key Functions of webhook implementation
* webhook.go->validate : This operation receives the AdmissionReview request object from Kubernetes. It then extracts the name of the resource to be deleted, the namespace, the API group, the API version, and the resource type. These are then passed to the getResourceLabels Operations. The labels are then passed to the validate.go->isDeletionRequestToBeBlocked operation to check if the deletion is to be blocked
* webhook.go->getResourceLabels : This operation uses the REST client to fetch the labels for resource and returns the map or labels
* validate.go->isDeletionRequestToBeBlocked : This checks if the resource labels indicate that the resource is to be deleted or not.

## Installation
* **Building and Pushing the container image** :The **Make target docker-build or docker-build-local** can be used to create the container image. The **make docker-push**  makefile target can be used to push the container image to the container registry. With the **make docker-build-local** makefile target you need dependencies like glide on your machine, the **make docker-build** makefile target uses a multi stage build for building the go binary. Make sure you change the values of **CONTAINER_NAME**, **CONTAINER_VERSION** in the Make file.
* The [ **deployments/webhook-k8s-resources.template.yaml**](https://github.com/maniSbindra/k8s-delete-validation-webhook/blob/master/deployments/webhook-k8s-resources.template.yaml#L7) is the kubernetes manifest template for this solution. The main kubernetes resources to be created are ValidatingWebhookConfiguration, a deployment and a service. The template has place holders for the [TLS Certificate](https://github.com/maniSbindra/k8s-delete-validation-webhook/blob/12d9aacc757b6c6208e47618e7282ad623eb05b8/deployments/webhook-k8s-resources.template.yaml#L16), the [TLS Key](https://github.com/maniSbindra/k8s-delete-validation-webhook/blob/12d9aacc757b6c6208e47618e7282ad623eb05b8/deployments/webhook-k8s-resources.template.yaml#L7), the [CA Bundle](https://github.com/maniSbindra/k8s-delete-validation-webhook/blob/12d9aacc757b6c6208e47618e7282ad623eb05b8/deployments/webhook-k8s-resources.template.yaml#L98) and the [container image](https://github.com/maniSbindra/k8s-delete-validation-webhook/blob/12d9aacc757b6c6208e47618e7282ad623eb05b8/deployments/webhook-k8s-resources.template.yaml#L35).
* The [Steps](https://github.com/avast/k8s-admission-webhook#example-configuration) mentioned of the avast repository explain how to replace values in the yaml template file. Instead of manually doing it you can using the **make gen-k8s-manifests** Makefile target from this solution. This is described in more detail as follows
* The makefile target **make gen-k8s-manifests** in this solution has all steps to replace values in the template, and as an output it generates the deployments/webhook-k8s-resources.yaml which has certificate, key and the ca bundle in the yaml. To execute this make target **you need to have access to the target kubernets** cluster (KUBECONFIG or ./kube/config). Before running the make target, verify that the values of the **CONTAINER_NAME, CONTAINER_VERSION, WEBHOOK_NAMESPACE and WEBHOOK_SERVICE_NAME** in the Makefile are correct.  
* After this applying the generated **deployments/webhook-k8s-resources.yaml** file creates all required kubernetes resources. By default entry for this generated file is in the .gitignore file.

### Commands
* Clone the repo
  ```sh
  git clone https://github.com/maniSbindra/k8s-delete-validation-webhook.git
  cd k8s-delete-validation-webhook
  ```
* Modify Makefile values for **CONTAINER_NAME, CONTAINER_VERSION, WEBHOOK_NAMESPACE and WEBHOOK_SERVICE_NAME**
* Make sure you are logged in to the container registry and have access to the kubernetes cluster
* Build and Push container image
  ```sh
  make docker-build
  make docker-push
  ```
* Generate certs, keys and generate actual yaml from yaml template.
  ```sh
  make gen-k8s-manifests
  ```
* Apply the yaml
  ```sh
  kubectl apply -f deployments/webhook-k8s-resources.yaml
  ```

## Verify installation
* Check the deployment k8s-delete-validation-webhook, for which a pod should be running
  ```sh
  kubectl get deploy k8s-delete-validation-webhook
  ```
* add the label webhook=enabled to the default namespace. Note this is the default value and can be changed in [**the namespace selector**](https://github.com/maniSbindra/k8s-delete-validation-webhook/blob/9f86e415d4365c66f484e5a543935e950f3026a1/deployments/webhook-k8s-resources.template.yaml#L107)
  ```sh
  kubectl label namespace default webhook=enabled
  ```
* Create Deployment with deleteLock=enabled label
  ```sh
  kubectl run nginx --image=nginx --port=80 --labels=deleteLock=enabled
  ```
* Trying to delete this deployment should throw following error
  ```sh
  kubectl delete deploy nginx
  Error from server: admission webhook "k8s-deletion-validation-webhook.test.com" denied the request: The deployment cannot be deleted as deletions are locked for this deployment
  ```



## Configuring the namespaces, resources and labels
* By modifying the [**rules section**](https://github.com/maniSbindra/k8s-delete-validation-webhook/blob/12d9aacc757b6c6208e47618e7282ad623eb05b8/deployments/webhook-k8s-resources.template.yaml#L100-L103) of the validatingwebhookconfiguration resource we can change the namespaces, api groups, api versions and resource types which this webhook handles. We also need to change the namespaces, and resource permissions for the cluster role in the [**Cluster Role and Binding section**](https://github.com/maniSbindra/k8s-delete-validation-webhook/blob/a4fde8f34445fed1cbfd44a5bf157d235cd603e5/deployments/webhook-k8s-resources.template.yaml#L111-L132) 
* By modifying the [**environment settings**](https://github.com/maniSbindra/k8s-delete-validation-webhook/blob/12d9aacc757b6c6208e47618e7282ad623eb05b8/deployments/webhook-k8s-resources.template.yaml#L48-L53) of the webhook deployment, we can control the delete rejection message, and the key and value of the label used to determine if resource deletion requests are rejected
* You need to add the label specified in [**the namespace selector**](https://github.com/maniSbindra/k8s-delete-validation-webhook/blob/9f86e415d4365c66f484e5a543935e950f3026a1/deployments/webhook-k8s-resources.template.yaml#L107) to any namespace where you expect the webhook to perform the validation



