package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Starts deployment deletion validation webhook",
	Long:  "Starts deployment deletion validation webhook, which stops deletion of the deployment if deletion lock label is set for the deployment",
	Run:   startWebhook,
}

var webhookViper = viper.New()

func init() {
	rootCmd.AddCommand(webhookCmd)

	webhookCmd.Flags().String("tls-cert-file", "",
		"Path to the certificate file. Required, unless --no-tls is set.")
	webhookCmd.Flags().Bool("no-tls", false,
		"Do not use TLS.")
	webhookCmd.Flags().String("tls-private-key-file", "",
		"Path to the certificate key file. Required, unless --no-tls is set.")
	webhookCmd.Flags().Int32("listen-port", 443,
		"Port to listen on.")
	webhookCmd.Flags().String("deployment-deletion-lock-label-key", "deleteLock",
		"The key of the deployment delete lock label")
	webhookCmd.Flags().String("deployment-deletion-lock-label-value", "enabled",
		"The value of the deployment delete lock label")
	webhookCmd.Flags().String("deployment-deletion-failed-message", "The deployment cannot be deleted as deletions are locked for this deployment",
		"Message to display after deletion of the particular deployment fails due to being locked")

	if err := webhookViper.BindPFlags(webhookCmd.Flags()); err != nil {
		errorWithUsage(err)
	}

	webhookViper.AutomaticEnv()
	webhookViper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

func startWebhook(cmd *cobra.Command, args []string) {
	config := &config{}
	if err := webhookViper.Unmarshal(config); err != nil {
		errorWithUsage(err)
	}

	if !config.NoTLS && (config.TLSPrivateKeyFile == "" || config.TLSCertFile == "") {
		errorWithUsage(errors.New("Both --tls-cert-file and --tls-private-key-file are required (unless TLS is disabled by setting --no-tls)"))
	}

	log.Debugf("Configuration is: %+v", config)

	//initialize kube client
	kubeClientSet, kubeClientSetErr := KubeClientSet(true)
	if kubeClientSetErr != nil {
		log.Fatal(kubeClientSetErr)
	}

	http.HandleFunc("/validate", admitFunc(validate).serve(config, kubeClientSet))

	addr := fmt.Sprintf(":%v", config.ListenPort)
	var httpErr error
	if config.NoTLS {
		log.Infof("Starting webserver at %v (no TLS)", addr)
		httpErr = http.ListenAndServe(addr, nil)
	} else {
		log.Infof("Starting webserver at %v (TLS)", addr)
		httpErr = http.ListenAndServeTLS(addr, config.TLSCertFile, config.TLSPrivateKeyFile, nil)
	}

	if httpErr != nil {
		log.Fatal(httpErr)
	} else {
		log.Info("Finished")
	}
}

func getResourceLabels(clientSet *kubernetes.Clientset, ar v1beta1.AdmissionReview) (map[string]string, error) {

	type ResourceData struct {
		Metadata struct {
			Labels map[string]string
		}
	}

	var resourceAPIRequest string
	if ar.Request.Kind.Group != "" {
		resourceAPIRequest = fmt.Sprintf("apis/%s/%s/namespaces/%s/%s/%s", ar.Request.Kind.Group, ar.Request.Kind.Version, ar.Request.Namespace, ar.Request.Resource.Resource, ar.Request.Name)
	} else {
		resourceAPIRequest = fmt.Sprintf("api/%s/namespaces/%s/%s/%s", ar.Request.Kind.Version, ar.Request.Namespace, ar.Request.Resource.Resource, ar.Request.Name)
	}

	jsonData, err := clientSet.RESTClient().Get().AbsPath(resourceAPIRequest).DoRaw()
	var resourceData ResourceData

	if err != nil {
		log.Error(err)
		return nil, err
	}

	errUnmarshall := json.Unmarshal(jsonData, &resourceData)

	if errUnmarshall != nil {
		log.Error(errUnmarshall)
		return nil, errUnmarshall
	}

	return resourceData.Metadata.Labels, nil

}

func validate(ar v1beta1.AdmissionReview, config *config, clientSet *kubernetes.Clientset) *v1beta1.AdmissionResponse {
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	var resourceLabels map[string]string
	var labelFetchError error
	var blockResourceDeletion = false
	resourceDeletionFailureMessage := config.DeploymentDeletionFailedMessage
	resourceDeletionLockLabelKey := config.DeploymentDeletionLockLabelKey
	resourceDeletionLockLabelValue := config.DeploymentDeletionLockLabelValue

	resourceLabels, labelFetchError = getResourceLabels(clientSet, ar)
	if labelFetchError != nil {
		log.Error(labelFetchError)
		return toAdmissionResponse(labelFetchError)
	}

	blockResourceDeletion = isDeletionRequestToBeBlocked(resourceLabels, resourceDeletionLockLabelKey, resourceDeletionLockLabelValue)
	if blockResourceDeletion {
		reviewResponse.Allowed = false
		reviewResponse.Result = &metav1.Status{Code: 403, Message: resourceDeletionFailureMessage}
		log.Infof("%q:%q:%q - %q", ar.Request.Kind.Kind, ar.Request.Namespace, ar.Request.Name, resourceDeletionFailureMessage)
	}

	return &reviewResponse
}
