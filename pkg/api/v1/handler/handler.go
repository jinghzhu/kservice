package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinghzhu/kservice/pkg/api/v1/adapter"
	"github.com/jinghzhu/kservice/pkg/api/v1/types"
	"github.com/jinghzhu/kservice/pkg/config"
	"github.com/jinghzhu/kservice/pkg/logger"

	apitypes "github.com/jinghzhu/kservice/pkg/api/types"
	v1types "github.com/jinghzhu/kservice/pkg/api/v1/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreatePod creates a Pod.
func CreatePod(ctx context.Context, r *http.Request) (result []byte, status int, err error) {
	status = http.StatusOK
	logger.InfoFields("Start to create Pod", logger.Fields{
		apitypes.LogCtxID: ctx.Value(apitypes.LogCtxID),
	})
	if r.Body == nil {
		errMsg := "No POST parameters found in the request"
		logger.Error(errMsg)

		return result, 400, errors.New(errMsg)
	}
	// Set Pod spec from input parameters.
	wp, err := adapter.InitWorkerPod(ctx, r.Body)
	if err != nil {
		errMsg := "Fail to parse JSON POST params"
		logger.ErrorFields(errMsg, logger.Fields{
			apitypes.LogCtxID: ctx.Value(apitypes.LogCtxID),
			logger.ERROR:      err,
		})

		return result, 400, fmt.Errorf("%s because of %v", errMsg, err)
	}

	// Translate to pod
	podObj, err := adapter.TranslateWorkerPodToPod(ctx, wp)
	if err != nil {
		errMsg := "Fail to init Pod object"
		logger.ErrorFields(errMsg, logger.Fields{
			apitypes.LogCtxID: ctx.Value(apitypes.LogCtxID),
			logger.ERROR:      err,
		})

		return result, 500, fmt.Errorf("%s because of %v", errMsg, err)
	}

	// Create Pod in Kubernetes.
	pod, err := apitypes.DefaultPodClient().CreatePod(podObj, podObj.GetNamespace(), metav1.CreateOptions{})
	if err != nil {
		errMsg := "Fail to create Pod"
		logger.ErrorFields(errMsg, logger.Fields{
			apitypes.LogCtxID: ctx.Value(apitypes.LogCtxID),
			logger.ERROR:      err,
		})

		return result, 400, fmt.Errorf("%s because of %v", errMsg, err)
	}
	podName, podNamespace := pod.GetName(), pod.GetNamespace()
	podLabel, podAnnotation := pod.GetLabels(), pod.GetAnnotations()
	ld := &v1types.WorkerDetails{
		Id: podName,
	}
	result, err = json.Marshal(ld)
	if err != nil {
		errMsg := "Fail to marshal result into JSON"
		logger.ErrorFields(errMsg, logger.Fields{
			apitypes.LogCtxID:           ctx.Value(apitypes.LogCtxID),
			apitypes.LogWorkerName:      podName,
			apitypes.LogWorkerNamespace: podNamespace,
			"pod.label":                 podLabel,
			"pod.annotation":            podAnnotation,
			logger.ERROR:                err,
		})

		return result, 500, fmt.Errorf("%s for %s because of %v", errMsg, podName, err)
	}

	logger.InfoFields("Successfully create Pod", logger.Fields{
		apitypes.LogCtxID:           ctx.Value(apitypes.LogCtxID),
		apitypes.LogWorkerName:      podName,
		apitypes.LogWorkerNamespace: podNamespace,
		"pod.label":                 podLabel,
		"pod.annotation":            podAnnotation,
	})

	return result, http.StatusOK, err
}

// GetPodStatus Retrieve Pod status.
func GetPodStatus(ctx context.Context, r *http.Request) (result []byte, status int, err error) {
	status = http.StatusOK
	vars := mux.Vars(r)
	podName := vars["key"]
	logger.InfoFields("Calling GetPodStatus", logger.Fields{
		apitypes.LogCtxID:      ctx.Value(apitypes.LogCtxID),
		apitypes.LogWorkerName: podName,
	})
	g := config.GetConfig()

	// Get Pod from Kubernetes.
	pod, err := apitypes.DefaultPodClient().GetPod(g.WorkerNamespace, podName, metav1.GetOptions{})
	if err != nil {
		errMsg := "Fail to get Pod"
		logger.ErrorFields(errMsg, logger.Fields{
			apitypes.LogCtxID:      ctx.Value(apitypes.LogCtxID),
			apitypes.LogWorkerName: podName,
			logger.ERROR:           err,
		})

		return result, 404, fmt.Errorf("%s for %s because of %v", errMsg, podName, err)
	}

	// Parse worker status.
	workerStatus := &v1types.WorkerStatus{
		Reason:  pod.Status.Reason,
		Message: pod.Status.Message,
		State:   apitypes.WorkerStatusUnknown,
		Status:  string(pod.Status.Phase),
	}
	if pod.Status.ContainerStatuses != nil {
		if pod.Status.ContainerStatuses[0].State.Waiting != nil {
			workerStatus.Reason = pod.Status.ContainerStatuses[0].State.Waiting.Reason
			workerStatus.Message = pod.Status.ContainerStatuses[0].State.Waiting.Message
			workerStatus.State = apitypes.WorkerStatusWaiting
		} else if pod.Status.ContainerStatuses[0].State.Running != nil {
			workerStatus.State = apitypes.WorkerStatusRunning
		} else if pod.Status.ContainerStatuses[0].State.Terminated != nil {
			workerStatus.State = apitypes.WorkerStatusTerminated
			workerStatus.Reason = pod.Status.ContainerStatuses[0].State.Terminated.Reason
			workerStatus.Message = pod.Status.ContainerStatuses[0].State.Terminated.Message
			workerStatus.ExitCode = &pod.Status.ContainerStatuses[0].State.Terminated.ExitCode
		}
	}
	result, err = json.Marshal(workerStatus)
	podNamespace := pod.GetNamespace()
	podAnnotation, podLabel := pod.GetAnnotations(), pod.GetLabels()
	if err != nil {
		errMsg := "Fail to marshal Pod status into JSON"
		logger.ErrorFields(errMsg, logger.Fields{
			apitypes.LogCtxID:      ctx.Value(apitypes.LogCtxID),
			apitypes.LogWorkerName: podName,
			"Pod.Namespace":        podNamespace,
			"Pod.Annotation":       podAnnotation,
			"Pod.Label":            podLabel,
			logger.ERROR:           err,
		})

		return result, 500, fmt.Errorf("%s for %s because of %v", errMsg, podName, err)
	}

	logger.InfoFields("Successfully get Pod status", logger.Fields{
		apitypes.LogCtxID:      ctx.Value(apitypes.LogCtxID),
		apitypes.LogWorkerName: podName,
		"pod.namespace":        podNamespace,
		"pod.annotation":       podAnnotation,
		"pod.label":            podLabel,
	})

	return result, status, err
}

// GetPodLog retrieves Pod logs.
func GetPodLog(ctx context.Context, r *http.Request) (result []byte, status int, err error) {
	status = http.StatusOK
	vars := mux.Vars(r)
	podName := vars["key"]
	logFields := logger.Fields{
		apitypes.LogCtxID:      ctx.Value(apitypes.LogCtxID),
		apitypes.LogWorkerName: podName,
	}
	logger.InfoFields("Calling getPodLog", logFields)
	g := config.GetConfig()
	logFields[apitypes.LogWorkerNamespace] = g.WorkerNamespace

	// Get logs.
	podLog, err := apitypes.DefaultPodClient().GetLogString(g.WorkerNamespace, podName, &corev1.PodLogOptions{})
	if err != nil {
		errMsg := "Fail to find log stream"
		logFields[logger.ERROR] = err
		logger.ErrorFields(errMsg, logFields)

		return result, 404, fmt.Errorf("%s for %s because of %v", errMsg, podName, err)
	}

	// Parser logs into responses.
	logs := &types.Logs{
		Log: podLog,
	}
	result, err = json.Marshal(logs)
	if err != nil {
		errMsg := "Fail to marshal logs into JSON"
		logFields[logger.ERROR] = err
		logger.ErrorFields(errMsg, logFields)

		return result, 500, fmt.Errorf("%s for %s because of %v", errMsg, podName, err)
	}

	logger.InfoFields("Successfully get Pod logs", logFields)

	return result, status, err
}

// GetPodInfo gets pod spec.
func GetPodInfo(ctx context.Context, r *http.Request) (result []byte, status int, err error) {
	status = http.StatusOK
	vars := mux.Vars(r)
	podName := vars["key"]
	logger.InfoFields("Calling GetPodInfo", logger.Fields{
		apitypes.LogCtxID:      ctx.Value(apitypes.LogCtxID),
		apitypes.LogWorkerName: podName,
	})
	g := config.GetConfig()

	// Get Pod from Kubernetes.
	pod, err := apitypes.DefaultPodClient().GetPod(g.WorkerNamespace, podName, metav1.GetOptions{})
	if err != nil {
		errMsg := "Fail to get Pod"
		logger.ErrorFields(errMsg, logger.Fields{
			apitypes.LogCtxID:      ctx.Value(apitypes.LogCtxID),
			apitypes.LogWorkerName: podName,
			logger.ERROR:           err,
		})

		return result, 404, fmt.Errorf("%s for %s because of %v", errMsg, podName, err)
	}

	// Translate Pod spec into WorkerPod spec.
	wp := adapter.TranslatePodToWorkerPod(ctx, pod)
	podNamespace := pod.GetNamespace()
	podLabel, podAnnotation := pod.GetLabels(), pod.GetAnnotations()
	result, err = json.Marshal(wp)
	if err != nil {
		errMsg := "Fail to marshal Pod info"
		logger.ErrorFields(errMsg, logger.Fields{
			apitypes.LogCtxID:           ctx.Value(apitypes.LogCtxID),
			apitypes.LogWorkerName:      podName,
			apitypes.LogWorkerNamespace: podNamespace,
			"Pod.annotation":            podAnnotation,
			"Pod.label":                 podLabel,
			logger.ERROR:                err,
		})

		return result, 500, fmt.Errorf("%s for %s because of %v", errMsg, podName, err)
	}

	logger.InfoFields("Successfully get Pod info", logger.Fields{
		apitypes.LogCtxID:           ctx.Value(apitypes.LogCtxID),
		apitypes.LogWorkerName:      podName,
		apitypes.LogWorkerNamespace: podNamespace,
		"pod.annotation":            podAnnotation,
		"pod.label":                 podLabel,
	})

	return result, status, err
}
