package config

import (
	"context"
	"log"
	"log/syslog"
)

const (
	defaultLoggingFlag         int             = log.LstdFlags
	defaultLoggingPriority     syslog.Priority = syslog.LOG_INFO
	defaultKubeContext         string          = "default"
	defaultPodNamePrefix       string          = "worker"
	defaultCRDNamePrefix       string          = "crd"
	defaultContainerNamePrefix string          = "worker-"
	defaultPort                string          = "80"
	defaultKubeconfig          string          = "/kube/tempconfig"
	defaultWorkerNamespace     string          = "worker"
	defaultCRDNamespace        string          = "crd"
)

var (
	ContextRoot = context.Background()
	config      *Config
)

// Config consists fields to connect to Kubernetes.
type Config struct {
	// CRDNamespace is the namespace where kservice deals with CRD.
	CRDNamespace string `json:"crdNamespace"`
	// WorkerNamespace is the namespace where kservice deals with Pod.
	WorkerNamespace     string              `json:"workerNamespace"`
	ListenPort          string              `json:"port"`
	ListenAddress       string              `json:"address"`
	KubeContext         string              `json:"kubeContest"`
	PodNamePrefix       string              `json:"podNamePrefix"`
	ContainerNamePrefix string              `json:"containerNamePrefix"`
	CRDNamePrefix       string              `json:"crdNamePrefix"`
	LoggingLocalOptions []*LoggingLocalOpts `json:"localLogOpts"`
	Kubeconfig          string              `json:"kubeconfig"`
}
