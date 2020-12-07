package types

import (
	"time"
)

const (
	// TimeoutConnection is the maximum length of time to wait before giving up on a server request when calling
	// Kubernetes API. A value of 0 means no timeout.
	TimeoutConnection time.Duration = 15 * time.Second

	// EnvUser is used to set the system variable for container.
	EnvUser string = "USER"
)

type Resource struct {
	Cpu    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

var DefaultResourceRequests = &Resource{
	Cpu:    "100m",
	Memory: "300Mi",
}

var DefaultResourceLimits = &Resource{
	Cpu:    "200m",
	Memory: "400Mi",
}

// WorkerDetails is the ID of target Kubernetes resource.
type WorkerDetails struct {
	Id string `json:"id"`
}

type UserInfo struct {
	UserName string  `json:"username"`
	UserID   *int64  `json:"uid"`
	GroupID  []int64 `json:"gid"`
}

type Logs struct {
	Log string
}

type Mount struct {
	Name      string `string:"name,omitempty"`
	Type      string `string:"type,omitempty"`
	Share     string `string:"share"`
	Server    string `string:"server"`
	MountPath string `string:"mountpath"`
}

type ExitCode *int32

type WorkerStatus struct {
	Status   string   `json:"status"`
	Reason   string   `json:"reason"`
	State    string   `json:"state"`
	Message  string   `json:"msg"`
	ExitCode ExitCode `json:"exitCode"`
}
