package app

import (
	"github.com/amir-ahmad/cuemix/pkg/helm"
	"github.com/amir-ahmad/cuemix/pkg/k8s"
)

type HelmConfig struct {
	DestDir     string `json:"destDir"`
	Untar       bool   `json:"untar"`
	KubeVersion string
}

type AppConfig struct {
	Helm HelmConfig `json:"helm"`
}

type StrategicPatch struct {
	Object map[string]interface{} `json:"patch"`
	Target k8s.Target             `json:"target"`
}
type Config struct {
	Config           AppConfig                 `json:"cuemix"`
	HelmCharts       map[string]helm.HelmChart `json:"helm"`
	Manifests        []string                  `json:"manifests"`
	StrategicPatches map[string]StrategicPatch `json:"strategicpatch"`
	JsonPatches      map[string]k8s.JsonPatch  `json:"jsonpatch"`
}
