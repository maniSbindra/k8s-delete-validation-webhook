package main

type config struct {
	NoTLS                            bool   `mapstructure:"no-tls"`
	TLSCertFile                      string `mapstructure:"tls-cert-file"`
	TLSPrivateKeyFile                string `mapstructure:"tls-private-key-file"`
	ListenPort                       int    `mapstructure:"listen-port"`
	Namespace                        string `mapstructure:"namespace"`
	DeploymentDeletionFailedMessage  string `mapstructure:"deployment-deletion-failed-message"`
	DeploymentDeletionLockLabelKey   string `mapstructure:"deployment-deletion-lock-label-key"`
	DeploymentDeletionLockLabelValue string `mapstructure:"deployment-deletion-lock-label-value"`
}
