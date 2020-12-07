package config

import (
	"os"
)

func init() {
	initConfig()
}

func initConfig() {
	config = &Config{}
	config.PodNamePrefix = defaultPodNamePrefix
	config.CRDNamePrefix = defaultCRDNamePrefix
	config.ContainerNamePrefix = defaultContainerNamePrefix

	config.CRDNamespace = os.Getenv("KSERVICE_CRD_NAMESPACE")
	if config.CRDNamespace == "" {
		config.CRDNamespace = defaultCRDNamespace
	}

	config.CRDNamespace = os.Getenv("KSERVICE_WORKER_NAMESPACE")
	if config.WorkerNamespace == "" {
		config.WorkerNamespace = defaultWorkerNamespace
	}

	config.KubeContext = os.Getenv("KSERVICE_KUBECONTEXT")
	if config.KubeContext == "" {
		config.KubeContext = defaultKubeContext
	}

	config.ListenPort = os.Getenv("KSERVICE_PORT")
	if config.ListenPort == "" {
		config.ListenPort = defaultPort
	}
	config.ListenAddress = "0.0.0.0:" + config.ListenPort

	config.LoggingLocalOptions = []*LoggingLocalOpts{defaultLoggingLocalOpts()}

	config.Kubeconfig = os.Getenv("KSERVICE_KUBECONFIG")
	if config.Kubeconfig == "" {
		config.Kubeconfig = defaultKubeconfig
	}
}

// GetConfig returns a pointer to the current config.
func GetConfig() *Config {
	return config
}
