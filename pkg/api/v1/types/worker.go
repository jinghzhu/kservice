package types

import (
	"strconv"
	"time"

	"github.com/jinghzhu/kservice/pkg/config"
	"github.com/jinghzhu/kservice/pkg/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Params is worker parameters.
type Params struct {
	Body WorkerPod
}

// WorkerPod is to represent the Pod spec.
type WorkerPod struct {
	// Cmd is the command to run. If not set dockerfile.Entrypoint is used.
	Cmd  []string `json:"cmd"`
	Name string   `json:"-"`
	// Mounts is the array of NFS mount points.
	Mounts           []Mount           `json:"mounts,omitempty"`
	Env              map[string]string `json:"env,omitempty"`
	ResourceLimits   *Resource         `json:"limit,omitempty"`
	ResourceRequests *Resource         `json:"resource,omitempty"`
	UserInfo         UserInfo          `json:"userinfo,omitempty"`
	ReadOnlyFS       bool              `json:"readonlyfs,omitempty"`
	Namespace        string            `json:"namespace"`
	Image            string            `json:"image,omitempty"`
	ImageVersion     string            `json:"imageversion,omitempty"`
	Prefix           string            `json:"prefix,omitempty"`
	Annotations      map[string]string `json:"annotations,omitempty"`
	Labels           map[string]string `json:"labels,omitempty"`
}

// DefaultWorkerPod created the WorkerPod with default value.
func DefaultWorkerPod(g *config.Config) *WorkerPod {
	return &WorkerPod{
		Namespace:   g.WorkerNamespace,
		Prefix:      g.PodNamePrefix,
		Name:        g.ContainerNamePrefix,
		ReadOnlyFS:  false,
		Annotations: make(map[string]string),
		Labels:      make(map[string]string),
		UserInfo:    UserInfo{},
	}
}

// TranslateEnv translates WorkerPod.Env to the Pod.api.Envvar
func (wp *WorkerPod) TranslateEnv() (envVar []corev1.EnvVar) {
	env := wp.Env
	for key, value := range env {
		envVar = append(envVar, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}
	return
}

// TranslateMounts translate WorkerPod.Mounts to corev1.VolumeMount and corev1.Volume.
func (wp *WorkerPod) TranslateMounts() ([]corev1.VolumeMount, []corev1.Volume) {
	mounts := wp.Mounts
	vols := []corev1.Volume{}
	volMounts := []corev1.VolumeMount{}
	for _, mount := range mounts {
		name := strconv.FormatInt(time.Now().UnixNano(), 10)
		vols = append(vols, corev1.Volume{
			Name: name,
			VolumeSource: corev1.VolumeSource{
				NFS: &corev1.NFSVolumeSource{
					Server: mount.Server,
					Path:   mount.Share,
				},
			},
		})
		volMounts = append(volMounts, corev1.VolumeMount{
			Name:      name,
			MountPath: mount.MountPath,
		})
	}
	return volMounts, vols
}

// TranslateResourceRequests Translate WorkerPod.ResourceRequests to corev1.ResourceList.
func (wp *WorkerPod) TranslateResourceRequests() (corev1.ResourceList, error) {
	resourceRequests := wp.ResourceRequests
	if resourceRequests == nil {
		resourceRequests = DefaultResourceRequests
	}
	return translateResource(*resourceRequests)
}

// TranslateResourceLimits Translate WorkerPod.ResourceLimits to corev1.ResourceList.
func (wp *WorkerPod) TranslateResourceLimits() (corev1.ResourceList, error) {
	resourceLimits := wp.ResourceLimits
	if resourceLimits == nil {
		resourceLimits = DefaultResourceLimits
	}
	return translateResource(*resourceLimits)
}

func translateResource(workerResource Resource) (corev1.ResourceList, error) {
	kubeResource := corev1.ResourceList{}
	var err error
	if workerResource.Cpu != "" {
		cpu, err := resource.ParseQuantity(workerResource.Cpu)
		if err != nil {
			logger.ErrorFields("Fail to parse CPU info", logger.Fields{
				logger.ERROR: err,
			})
			return kubeResource, err
		}
		kubeResource[corev1.ResourceCPU] = cpu
	}
	if workerResource.Memory != "" {
		memory, err := resource.ParseQuantity(workerResource.Memory)
		if err != nil {
			logger.ErrorFields("Fail to parse memory info", logger.Fields{
				logger.ERROR: err,
			})
			return kubeResource, err
		}
		kubeResource[corev1.ResourceMemory] = memory
	}
	return kubeResource, err
}
