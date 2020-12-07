package types

import (
	"context"

	"github.com/jinghzhu/kutils/pod"
)

const (
	True  bool = true
	False bool = false

	LogCtxID           string = "ContextID"
	LogWorkerName      string = "Worker Name"
	LogWorkerNamespace string = "Worker Namespace"

	WorkerStatusWaiting            string = "Waiting"
	WorkerStatusRunning            string = "Running"
	WorkerStatusTerminated         string = "Terminated"
	WorkerStatusUnknown            string = "Unknown"
	ContainerStatusReasonCompleted string = "Completed"

	// MountTypeNFS is the name of NFS mount type.
	MountTypeNFS string = "NFS"
	// MountTypeKubeSecret is the name of kubesecret mount type.
	MountTypeKubeSecret string = "secret"
)

var (
	ContextRoot      = context.Background()
	defaultPodClient *pod.Client
)
