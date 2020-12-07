package types

import (
	"github.com/jinghzhu/kservice/pkg/config"
	"github.com/jinghzhu/kutils/pod"
)

func init() {
	initDefaultPodClient()
}

func initDefaultPodClient() {
	c, err := pod.New(config.ContextRoot, "", config.GetConfig().Kubeconfig)
	if err != nil {
		panic(err)
	}
	defaultPodClient = c
}

// DefaultPodClient returns the default Pod client.
func DefaultPodClient() *pod.Client {
	return defaultPodClient
}
