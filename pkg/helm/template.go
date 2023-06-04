package helm

import (
	"fmt"
	"strings"

	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

// Template does the equivalent of a `helm template`
func (h HelmChart) Template(chartPath string, kubeVersion string) (map[string]string, error) {
	chart, err := loader.Load(chartPath)
	if err != nil {
		return map[string]string{}, fmt.Errorf("loading chart: %w", err)
	}

	if kubeVersion == "" {
		// Set a sensible default for kubeVersion
		kubeVersion = "1.25"
	}
	parsedKubeVersion, err := chartutil.ParseKubeVersion(kubeVersion)
	if err != nil {
		return map[string]string{}, fmt.Errorf("error parsing kubeVersion: %w", err)
	}

	capabilities := chartutil.Capabilities{
		KubeVersion: *parsedKubeVersion,
	}

	releaseOptions := chartutil.ReleaseOptions{
		Name:      h.Release,
		Namespace: h.Namespace,
		IsInstall: true,
	}

	// Merge chart values with our provided ones
	renderValues, err := chartutil.ToRenderValues(chart, h.Values, releaseOptions, &capabilities)
	if err != nil {
		return map[string]string{}, fmt.Errorf("preparing render values: %w", err)
	}

	// Helm template
	engine := engine.Engine{}
	renderedTemplates, err := engine.Render(chart, renderValues)
	if err != nil {
		return map[string]string{}, fmt.Errorf("rendering helm templates: %w", err)
	}

	// Iterate through charts CRDs and add to map
	if h.IncludeCRDs {
		for _, crd := range chart.CRDObjects() {
			renderedTemplates["crds/"+crd.File.Name] = string(crd.File.Data)
		}
	}

	for key, val := range renderedTemplates {
		// Remove NOTES.txt which is sometimes in the templated output
		if strings.HasSuffix(key, "NOTES.txt") {
			delete(renderedTemplates, key)
		}

		// Delete all empty templates
		if strings.TrimSpace(val) == "" {
			delete(renderedTemplates, key)
		}
	}

	return renderedTemplates, nil
}
