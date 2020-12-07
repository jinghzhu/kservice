package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/jinghzhu/kservice/pkg/config"
	"github.com/jinghzhu/kservice/pkg/logger"

	apitypes "github.com/jinghzhu/kservice/pkg/api/types"
	"github.com/jinghzhu/kservice/pkg/api/v1/types"
	v1types "github.com/jinghzhu/kservice/pkg/api/v1/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InitWorkerPod inits a default WorkPod spec.
func InitWorkerPod(ctx context.Context, jsonBody io.ReadCloser) (*v1types.WorkerPod, error) {
	g := config.GetConfig()
	wp := v1types.DefaultWorkerPod(g)
	id := ctx.Value(apitypes.LogCtxID).(uuid.UUID)
	idArr := strings.Split(id.String(), "-")
	wp.Name = wp.Name + idArr[0]
	logger.InfoFields("This default", logger.Fields{
		apitypes.LogCtxID:  ctx.Value(apitypes.LogCtxID),
		"Container-prefix": wp.Name,
	})

	err := json.NewDecoder(jsonBody).Decode(&wp)
	wp.Prefix = wp.Prefix + "-" + idArr[0] + "-"
	if err != nil {
		return wp, err
	}
	if wp.Cmd == nil {
		return wp, errors.New("cmd is a mandatory parameter")
	}

	return wp, err
}

func setPodContainer(wp *types.WorkerPod, mounts []corev1.VolumeMount) (container corev1.Container, err error) {
	image := bytes.NewBufferString("")
	fmt.Fprint(image, wp.Image, ":", wp.ImageVersion)
	env := wp.TranslateEnv()

	limits, err := wp.TranslateResourceLimits()
	if err != nil {
		return container, err
	}
	requests, err := wp.TranslateResourceRequests()
	if err != nil {
		return container, err
	}

	container = corev1.Container{
		Command: wp.Cmd,
		Env:     env,
		Image:   image.String(),
		Name:    wp.Name,
		SecurityContext: &corev1.SecurityContext{
			ReadOnlyRootFilesystem: &wp.ReadOnlyFS,
		},
		VolumeMounts:    mounts,
		ImagePullPolicy: "Always",
	}
	if len(limits) > 0 {
		container.Resources.Limits = limits
	}
	if len(requests) > 0 {
		container.Resources.Requests = requests
	}

	return container, err
}

// TranslateWorkerPodToPod translates the WorkerPod to Kubernetes Pod.
func TranslateWorkerPodToPod(ctx context.Context, wp *types.WorkerPod) (*corev1.Pod, error) {
	g := config.GetConfig()
	mounts, vols := wp.TranslateMounts()
	container, err := setPodContainer(wp, mounts)
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: wp.Prefix,
			Namespace:    wp.Namespace,
			Labels:       wp.Labels,
			Annotations:  wp.Annotations,
		},
		Spec: corev1.PodSpec{
			SecurityContext: &corev1.PodSecurityContext{},
			Containers:      []corev1.Container{container},
			HostNetwork:     false,
			RestartPolicy:   "Never",
			Volumes:         vols,
		},
	}
	logger.InfoFields("Output kubeconfig name", logger.Fields{
		apitypes.LogCtxID:           ctx.Value(apitypes.LogCtxID),
		"kubeconfig":                g.Kubeconfig,
		apitypes.LogWorkerNamespace: pod.Namespace,
	})
	if wp.UserInfo.UserID != nil {
		pod.Spec.SecurityContext.RunAsUser = wp.UserInfo.UserID
		pod.Spec.SecurityContext.SupplementalGroups = wp.UserInfo.GroupID
		pod.Spec.Containers[0].Env = append(pod.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  v1types.EnvUser,
			Value: wp.UserInfo.UserName,
		})
	}

	return pod, err
}

func TranslatePodToWorkerPod(ctx context.Context, pod *corev1.Pod) *v1types.WorkerPod {
	wp := &v1types.WorkerPod{
		Env:              TranslatePodEnv(ctx, pod),
		ResourceLimits:   TranslatePodResourceLimits(pod),
		ResourceRequests: TranslateResourceRequests(pod),
		UserInfo:         TranslateUserInfo(ctx, pod),
		Namespace:        pod.GetNamespace(),
		Labels:           pod.GetLabels(),
		Annotations:      pod.GetAnnotations(),
		Cmd:              TranslatePodCmd(pod),
		Mounts:           TranslatePodMounts(ctx, pod),
	}
	image := strings.SplitN(pod.Spec.Containers[0].Image, ":", 2)
	wp.Image = image[0]
	if len(image) > 1 {
		wp.ImageVersion = image[1]
	}

	return wp
}

func TranslateUserInfo(ctx context.Context, pod *corev1.Pod) v1types.UserInfo {
	userinfo := v1types.UserInfo{}
	userinfo.UserID = pod.Spec.SecurityContext.RunAsUser
	userinfo.GroupID = pod.Spec.SecurityContext.SupplementalGroups
	env := TranslatePodEnv(ctx, pod)
	userinfo.UserName = env[v1types.EnvUser]

	return userinfo
}

func TranslateResourceRequests(pod *corev1.Pod) *v1types.Resource {
	return translateResource(pod.Spec.Containers[0].Resources.Requests)
}

func TranslatePodResourceLimits(pod *corev1.Pod) *v1types.Resource {
	return translateResource(pod.Spec.Containers[0].Resources.Limits)
}

func translateResource(resourceReq corev1.ResourceList) *v1types.Resource {
	resource := types.Resource{}
	resource.Cpu = resourceReq.Cpu().String()
	resource.Memory = resourceReq.Memory().String()

	return &resource
}

func TranslatePodEnv(ctx context.Context, pod *corev1.Pod) map[string]string {
	envVars := pod.Spec.Containers[0].Env
	env := map[string]string{}
	logger.InfoFields("Env info", logger.Fields{
		apitypes.LogCtxID:      ctx.Value(apitypes.LogCtxID),
		apitypes.LogWorkerName: pod.GetName(),
		"Env":                  envVars,
	})

	for _, envVar := range envVars {
		env[envVar.Name] = envVar.Value
	}

	return env
}

func TranslatePodCmd(pod *corev1.Pod) []string {
	return pod.Spec.Containers[0].Command
}

func TranslatePodMounts(ctx context.Context, pod *corev1.Pod) []v1types.Mount {
	volumes := pod.Spec.Volumes
	volumeMounts := pod.Spec.Containers[0].VolumeMounts
	mounts := []v1types.Mount{}
	logger.InfoFields("Volumes info", logger.Fields{
		apitypes.LogCtxID:      ctx.Value(apitypes.LogCtxID),
		apitypes.LogWorkerName: pod.GetName(),
		"Volumes Mounts":       volumeMounts,
		"Volumes":              volumes,
	})
	for _, volMount := range volumeMounts {
		mount := &v1types.Mount{}
		mount.Name = volMount.Name
		mount.MountPath = volMount.MountPath
		var volume corev1.Volume
		for _, vol := range volumes {
			logger.InfoFields("Volume info", logger.Fields{
				apitypes.LogCtxID:      ctx.Value(apitypes.LogCtxID),
				apitypes.LogWorkerName: pod.GetName(),
				"Volume Name":          vol.Name,
				"Mount Name":           mount.Name,
			})
			if vol.Name == mount.Name {
				volume = vol
				break
			}
		}
		if volume.VolumeSource.NFS != nil {
			mount.Server = volume.VolumeSource.NFS.Server
			mount.Share = volume.VolumeSource.NFS.Path
			mount.Type = apitypes.MountTypeNFS
			mounts = append(mounts, *mount)
		} else {
			continue
		}
	}
	return mounts
}
