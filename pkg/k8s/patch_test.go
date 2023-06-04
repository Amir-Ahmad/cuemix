package k8s

import (
	"encoding/json"
	"testing"

	tu "github.com/amir-ahmad/cuemix/internal/testutils"
	"github.com/rogpeppe/go-internal/txtar"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func TestStrategicPatch(t *testing.T) {
	// each test is a file in testdata/
	tests := []string{
		"patch_update_value",
		"patch_remove_annotation",
		"patch_multiple",
		"patch_no_match",
		"patch_one_of_multiple",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			// Parse the txtar file
			archive, err := txtar.ParseFile("testdata/" + test + ".txtar")
			tu.NotErr(t, err)

			files := tu.ArchiveToMap(archive)

			// Create a new store and add the input object(s)
			var store ObjectStore
			err = store.AddObject(files["input.yaml"])
			tu.NotErr(t, err, "failed to parse input.yaml")

			// If a error.yaml is found, an error is expected
			_, wantErr := files["error.yaml"]

			// Parse and apply the patch
			patch, err := yaml.Parse(string(files["patch.yaml"]))
			tu.NotErr(t, err, "failed to parse patch.yaml")

			var target Target
			// Parse target.yaml if it exists
			if _, ok := files["target.yaml"]; ok {
				err := yaml.Unmarshal(files["target.yaml"], &target)
				tu.NotErr(t, err, "failed to parse target.yaml")
			}

			err = store.StrategicPatch(*patch, target)
			if err != nil && wantErr == false {
				t.Fatalf("did not expect error when patching, but got: %v", err)
			} else if err == nil && wantErr == true {
				t.Fatalf("expected error when patching but got none")
			} else if err == nil {
				// Create a new store for the expected output and add the expected object(s)
				var expectedOutput ObjectStore
				err = expectedOutput.AddObject(files["output.yaml"])
				tu.NotErr(t, err, "failed to parse output.yaml")

				// verify the patch was applied
				for i, obj := range store {
					objStr, err := obj.String()
					tu.NotErr(t, err, "failed to convert input object to str")
					expectedObjStr, err := expectedOutput[i].String()
					tu.NotErr(t, err, "failed to convert output object to str")
					if objStr != expectedObjStr {
						t.Errorf("Mismatch for %s:\nExpected:\n%sGot:\n%s", test, expectedObjStr, objStr)
					}
				}
			}
		})
	}
}

func TestJsonPatch(t *testing.T) {
	// each test is a file in testdata/
	tests := []string{
		"jsonpatch_remove_label",
		"jsonpatch_no_match",
		"jsonpatch_multiple",
		"jsonpatch_one_of_multiple",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			// Parse the txtar file
			archive, err := txtar.ParseFile("testdata/" + test + ".txtar")
			tu.NotErr(t, err)

			files := tu.ArchiveToMap(archive)

			// Create a new store and add the input object(s)
			var store ObjectStore
			err = store.AddObject(files["input.yaml"])
			tu.NotErr(t, err, "failed to parse input.yaml")

			// If a error.yaml is found, an error is expected
			_, wantErr := files["error.yaml"]

			var jsonpatch JsonPatch
			err = yaml.Unmarshal(files["patch.yaml"], &jsonpatch)
			tu.NotErr(t, err, "failed to parse patch.yaml")

			jsonpatch.Patch, err = json.Marshal(jsonpatch.YamlPatch)
			tu.NotErr(t, err, "failed to parse jsonpatch")

			err = store.JsonPatch(jsonpatch.Patch, jsonpatch.Target)
			if err != nil && wantErr == false {
				t.Fatalf("did not expect error when patching, but got: %v", err)
			} else if err == nil && wantErr == true {
				t.Fatalf("expected error when patching but got none")
			} else if err == nil {
				// Create a new store for the expected output and add the expected object(s)
				var expectedOutput ObjectStore
				err = expectedOutput.AddObject(files["output.yaml"])
				tu.NotErr(t, err, "failed to parse output.yaml")

				// verify the patch was applied
				for i, obj := range store {
					objStr, err := obj.String()
					tu.NotErr(t, err, "failed to convert input object to str")
					expectedObjStr, err := expectedOutput[i].String()
					tu.NotErr(t, err, "failed to convert output object to str")
					if objStr != expectedObjStr {
						t.Errorf("Mismatch for %s:\nExpected:\n%sGot:\n%s", test, expectedObjStr, objStr)
					}
				}
			}
		})
	}
}

func TestMatchesObject(t *testing.T) {
	// each test is a file in testdata/
	tests := []string{
		"match_on_annotation",
		"match_on_group",
		"match_on_label",
		"match_on_name",
		"match_on_namespace",
		"match_on_version",
		"no_match_on_annotation",
		"no_match_on_group_and_version",
		"no_match_on_kind",
		"no_match_on_label",
		"no_match_on_name",
		"no_match_on_name_and_ns",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			// Parse the txtar file
			archive, err := txtar.ParseFile("testdata/" + test + ".txtar")
			tu.NotErr(t, err)

			files := tu.ArchiveToMap(archive)

			obj, err := yaml.Parse(string(files["object.yaml"]))
			tu.NotErr(t, err, "failed to parse object.yaml")

			// By default a match is expected. But if a no_match.yaml is found,
			// the test asserts that the object won't match the target
			_, expectMatch := files["no_match.yaml"]
			expectMatch = !expectMatch

			var target Target
			err = yaml.Unmarshal(files["target.yaml"], &target)
			tu.NotErr(t, err, "failed to parse target.yaml")
			selector := SelectorFromTarget(target)

			match, err := MatchesObject(*obj, selector)
			tu.NotErr(t, err, "when matching object")
			if match && !expectMatch {
				t.Fatalf("Expected no match, but target matched object")
			}
			if !match && expectMatch {
				t.Fatalf("Expected match, but target didn't match object")
			}
		})
	}
}
