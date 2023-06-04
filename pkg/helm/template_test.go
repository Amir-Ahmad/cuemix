package helm

import (
	"strings"
	"testing"

	tu "github.com/amir-ahmad/cuemix/internal/testutils"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// Additional things to test:
// KubeVersion works

type values = map[string]interface{}
type checks func(t *testing.T, h HelmChart, templates map[string][]*yaml.RNode)

func TestTemplate(t *testing.T) {
	tests := []struct {
		name        string
		txtar       string
		chart       HelmChart
		kubeVersion string
		wantErr     bool
		checks      checks
	}{
		{
			name:  "templating value",
			txtar: "chart.txtar",
			chart: HelmChart{Release: "test-chart", Namespace: "default",
				Values: values{"cmkey": "cmvalue"}},
			wantErr: false,
			checks: func(t *testing.T, h HelmChart, templates map[string][]*yaml.RNode) {
				cm := templates["templates/test1.yaml"][0]

				lookup := "data.key"
				expectedValue := h.Values["cmkey"]
				value, err := cm.GetString(lookup)
				tu.NotErr(t, err, "Failed to lookup "+lookup)
				if value != expectedValue {
					t.Fatalf("Expected %s, got %s for %s", expectedValue, value, lookup)
				}
			},
		},
		{
			name:    "templating with blank values",
			txtar:   "chart.txtar",
			chart:   HelmChart{Release: "test-chart", Namespace: "default"},
			wantErr: false,
		},
		{
			name:    "release namespace is substituted",
			txtar:   "chart.txtar",
			chart:   HelmChart{Release: "test-chart", Namespace: "test-namespace"},
			wantErr: false,
			checks: func(t *testing.T, h HelmChart, templates map[string][]*yaml.RNode) {
				obj := templates["templates/test1.yaml"][0]

				if ns := obj.GetNamespace(); ns != h.Namespace {
					t.Fatalf("Expected namespace %s, got %s", h.Namespace, ns)
				}
			},
		},
		{
			name:    "CRDs are omitted",
			txtar:   "chart.txtar",
			chart:   HelmChart{Release: "test-chart", Namespace: "default", IncludeCRDs: false},
			wantErr: false,
			checks: func(t *testing.T, h HelmChart, templates map[string][]*yaml.RNode) {
				for key := range templates {
					if strings.HasPrefix(key, "crds/") {
						t.Fatalf("Expected no CRDs, but got %s", key)
					}
				}
			},
		},
		{
			name:    "CRDs are outputted",
			txtar:   "chart.txtar",
			chart:   HelmChart{Release: "test-chart", Namespace: "default", IncludeCRDs: true},
			wantErr: false,
			checks: func(t *testing.T, h HelmChart, templates map[string][]*yaml.RNode) {
				key := "crds/dummycrd.yaml"
				if _, ok := templates[key]; !ok {
					t.Fatalf("Output didn't contain CRD %s", key)
				}
			},
		},
		{
			name:        "kubeversion",
			txtar:       "kubeversion.txtar",
			chart:       HelmChart{Release: "test-chart", Namespace: "default"},
			wantErr:     false,
			kubeVersion: "1.10",
			checks: func(t *testing.T, h HelmChart, templates map[string][]*yaml.RNode) {
				// Test that a template specifying the Kubeversion is in the output
				key := "templates/expected.yaml"
				if _, ok := templates[key]; !ok {
					t.Fatalf("Output didn't contain %s", key)
				}

				// Check that a template expecting a newer version is not in the output
				key = "templates/not_expected.yaml"
				if _, ok := templates[key]; ok {
					t.Fatalf("Output contains %s, which is not expected", key)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Write txtar file to directory
			chartPath, cleanup := tu.WriteTxtarToTmp(t, "testdata/"+test.txtar)
			defer cleanup()

			templates, err := test.chart.Template(chartPath, test.kubeVersion)
			tu.NotErr(t, err, "error rendering template")

			output := make(map[string][]*yaml.RNode)

			for key, value := range templates {
				// helm stores the chart name in the key,
				// which is not relevant
				parts := strings.SplitN(key, "/", 2)
				newKey := parts[1]

				// Convert to Rnode for running checks
				nodes, err := kio.FromBytes([]byte(value))
				tu.NotErr(t, err, "error parsing yaml")
				output[newKey] = nodes
			}

			if test.checks != nil {
				test.checks(t, test.chart, output)
			}
		})
	}
}
