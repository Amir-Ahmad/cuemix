package k8s

import (
	"encoding/json"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// Alias Object to kyaml.RNode
type Object = yaml.RNode

type ObjectStore []Object

type JsonPatch struct {
	YamlPatch interface{}     `yaml:"patch"`
	Patch     json.RawMessage `json:"patch" yaml:"-"`
	Target    Target          `json:"target" yaml:"target"`
}

type Target struct {
	Group              string `json:"group" yaml:"group"`
	Version            string `json:"version" yaml:"version"`
	Kind               string `json:"kind" yaml:"kind"`
	Name               string `json:"name" yaml:"name"`
	Namespace          string `json:"namespace" yaml:"namespace"`
	AnnotationSelector string `json:"annotationSelector" yaml:"annotationSelector"`
	LabelSelector      string `json:"labelSelector" yaml:"labelSelector"`
}
