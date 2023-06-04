package k8s

import (
	"path/filepath"
	"testing"

	tu "github.com/amir-ahmad/cuemix/internal/testutils"
	"github.com/rogpeppe/go-internal/txtar"
	"sigs.k8s.io/kustomize/kyaml/kio"
)

func TestAddObject(t *testing.T) {
	testcases := []struct {
		name       string
		yamlBytes  []byte
		wantErr    bool
		numObjects *int
	}{
		{
			name: "Valid Object",
			yamlBytes: []byte(`
apiVersion: v1
kind: Pod
metadata:
  name: my-pod`),
			wantErr:    false,
			numObjects: tu.PtrInt(1),
		},
		{
			name: "Two Objects",
			yamlBytes: []byte(`
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
---
apiVersion: v1
kind: Pod
metadata:
  name: second-pod`),
			wantErr:    false,
			numObjects: tu.PtrInt(2),
		},
		{
			name: "Invalid Object",
			yamlBytes: []byte(`
apiVersion: v1`),
			wantErr: true,
		},
		{
			name: "Not YAML",
			yamlBytes: []byte(`
{"apiVersion": "v1"}`),
			wantErr: true,
		},
		{
			name:       "Empty YAML",
			yamlBytes:  []byte{},
			wantErr:    true,
			numObjects: tu.PtrInt(0),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var store ObjectStore
			err := store.AddObject(tc.yamlBytes)
			hasError := err != nil
			if hasError != tc.wantErr {
				t.Errorf("Expected %t, but got %t: %v", tc.wantErr, hasError, err)
			}
			if tc.numObjects != nil && len(store) != *tc.numObjects {
				t.Errorf("Expected %d object/s, but got %d", *tc.numObjects, len(store))
			}
		})
	}
}

func TestAddObjectFromPath(t *testing.T) {
	tmpDir, cleanup := tu.WriteTxtarToTmp(t, "testdata/manifests.txtar")
	defer cleanup()

	tests := []string{
		"test1.yaml",
		"test2.yaml",
		"test3.yaml",
	}

	for _, file := range tests {
		var store ObjectStore
		t.Run(file, func(t *testing.T) {
			err := store.AddObjectFromPath(filepath.Join(tmpDir, file))
			tu.NotErr(t, err, "failed to parse objects from "+file)
		})
	}
}

func TestAddObjectsFromDir(t *testing.T) {
	tmpDir, cleanup := tu.WriteTxtarToTmp(t, "testdata/manifests.txtar")
	defer cleanup()

	var store ObjectStore
	err := store.AddObjectsFromDir(tmpDir)
	if err != nil || len(store) != 3 {
		t.Errorf("failed to read objects from directory: %v", err)
	}
}

func TestFilterByKind(t *testing.T) {
	archive, err := txtar.ParseFile("testdata/manifests.txtar")
	tu.NotErr(t, err)

	var store ObjectStore
	for _, file := range archive.Files {
		err = store.AddObject(file.Data)
		tu.NotErr(t, err, "failed to add objects for "+file.Name)
	}
	if len(store) != 3 {
		t.Fatalf("failed to read objects from testdata: %v", err)
	}

	var filtered ObjectStore
	kind := "ConfigMap"
	filtered, err = store.FilterByKind(kind)
	tu.NotErr(t, err, "when filtering by kind "+kind)

	if len(filtered) == 0 {
		t.Errorf("No objects in filtered store")
	}

	for _, obj := range filtered {
		if obj.GetKind() != kind {
			t.Errorf("Filtering on %s failed, %s found", kind, obj.GetKind())
		}
	}
}

func TestValidateObject(t *testing.T) {
	tests := []struct {
		name      string
		yamlBytes []byte
		wantErr   bool
	}{
		{
			name: "ValidObject",
			yamlBytes: []byte(`
apiVersion: v1
kind: Pod
metadata:
  name: my-pod`),
			wantErr: false,
		},
		{
			name: "MissingKind",
			yamlBytes: []byte(`
apiVersion: v1
metadata:
  name: my-pod`),
			wantErr: true,
		},
		{
			name: "MissingApiVersion",
			yamlBytes: []byte(`
kind: Pod
metadata:
  name: my-second-pod`),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nodes, err := kio.FromBytes(tc.yamlBytes)
			tu.NotErr(t, err, "parsing YAML")

			err = ValidateObject(*nodes[0])
			hasError := err != nil
			if hasError != tc.wantErr {
				t.Errorf("Expected %t, but got %t: %v", tc.wantErr, hasError, err)
			}
		})
	}
}
