package k8s

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/resid"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	"sigs.k8s.io/kustomize/kyaml/yaml/merge2"
)

// StrategicPatch applies a strategic patch to selected Objects in the store.
func (s *ObjectStore) StrategicPatch(p Object, t Target) error {
	patchApplied := false

	var selector types.Selector
	if t == (Target{}) {
		selector = SelectorFromObject(p)
	} else {
		selector = SelectorFromTarget(t)
	}

	for i, obj := range *s {
		// Check if object matches the selector
		matches, err := MatchesObject(obj, selector)
		if err != nil {
			return fmt.Errorf("error when matching resource: %w", err)
		}
		if !matches {
			continue
		}

		patchApplied = true

		// Apply patch
		mergedNode, err := merge2.Merge(&p, &obj, yaml.MergeOptions{
			ListIncreaseDirection: yaml.MergeOptionsListPrepend,
		})
		if err != nil {
			return fmt.Errorf("failed to apply patch: %w", err)
		}

		(*s)[i] = *mergedNode
	}

	if !patchApplied {
		return fmt.Errorf("patch didn't match any object")
	}

	return nil
}

// JsonPatch applies a json patch to selected Objects in the store.
func (s *ObjectStore) JsonPatch(p []byte, t Target) error {
	patchApplied := false

	selector := SelectorFromTarget(t)

	patch, err := jsonpatch.DecodePatch(p)
	if err != nil {
		return fmt.Errorf("failed to decode patch: %w", err)
	}

	for i, obj := range *s {
		// Check if object matches the selector
		matches, err := MatchesObject(obj, selector)
		if err != nil {
			return fmt.Errorf("error when matching resource: %w", err)
		}
		if !matches {
			continue
		}

		patchApplied = true

		// Convert object to JSON
		objJson, err := obj.MarshalJSON()
		if err != nil {
			return fmt.Errorf("failed to convert object to JSON: %w", err)
		}

		patchedJson, err := patch.Apply(objJson)
		if err != nil {
			return fmt.Errorf("failed to apply patch: %w", err)
		}

		err = obj.UnmarshalJSON(patchedJson)
		if err != nil {
			return fmt.Errorf("failed to convert patched JSON to object: %w", err)
		}

		(*s)[i] = obj
	}

	if !patchApplied {
		return fmt.Errorf("patch didn't match any object")
	}

	return nil
}

func MatchesObject(obj Object, selector types.Selector) (bool, error) {
	// Compile selector to Regex
	sr, err := types.NewSelectorRegex(&selector)
	if err != nil {
		return false, err
	}

	// Skip if mismatch on Namespace
	if !sr.MatchNamespace(obj.GetNamespace()) {
		return false, nil
	}

	// Skip if mismatch on Name
	if !sr.MatchName(obj.GetName()) {
		return false, nil
	}

	// Skip if GVK is mismatch
	gvk := resid.GvkFromNode(&obj)
	if !sr.MatchGvk(gvk) {
		return false, nil
	}

	// matches the label selector
	matched, err := obj.MatchesLabelSelector(selector.LabelSelector)
	if err != nil {
		return false, err
	}
	if !matched {
		return false, nil
	}

	// matches the annotation selector
	matched, err = obj.MatchesAnnotationSelector(selector.AnnotationSelector)
	if err != nil {
		return false, err
	}
	if !matched {
		return false, nil
	}

	return true, nil
}
