package k8s

import (
	"fmt"

	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/resid"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// Convert from map[string]interface{} to Object and validate
func FromMap(m map[string]interface{}, validate bool) (Object, error) {
	obj, err := yaml.FromMap(m)
	if err != nil {
		return Object{}, fmt.Errorf("when converting to Object: %w", err)
	}

	if validate {
		err = ValidateObject(*obj)
		if err != nil {
			return Object{}, err
		}
	}

	return *obj, nil
}

func SelectorFromTarget(t Target) types.Selector {
	return types.Selector{
		ResId: resid.ResId{
			Gvk:       resid.Gvk{Group: t.Group, Version: t.Version, Kind: t.Kind},
			Name:      t.Name,
			Namespace: t.Namespace,
		},
		AnnotationSelector: t.AnnotationSelector,
		LabelSelector:      t.LabelSelector,
	}
}

func SelectorFromObject(o Object) types.Selector {
	return types.Selector{
		ResId: resid.ResId{
			Gvk:       resid.GvkFromNode(&o),
			Name:      o.GetName(),
			Namespace: o.GetNamespace(),
		},
	}
}
