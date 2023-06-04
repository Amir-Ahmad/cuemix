package k8s

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/kio"
)

// Validate that the Object has a few required fields set
func ValidateObject(o Object) error {

	if o.GetKind() == "" {
		return fmt.Errorf("validating object, kind not set")
	}
	if o.GetApiVersion() == "" {
		return fmt.Errorf("validating object, apiVersion not set")
	}

	return nil
}

// Print all Objects in store
func (s *ObjectStore) Print(sort bool) error {

	var output ObjectStore

	// If sort = true, sort objects
	if sort {
		output = sortManifestsByKind(*s)
	} else {
		output = *s
	}

	for i, o := range output {
		out, err := o.String()
		if err != nil {
			return fmt.Errorf("printing object at index %d: %w", i, err)
		}

		fmt.Printf("%s---\n", out)
	}
	return nil
}

// Add object to store
func (s *ObjectStore) AddObject(yamlBytes []byte) error {

	nodes, err := kio.FromBytes(yamlBytes)
	if err != nil {
		return fmt.Errorf("parsing YAML: %w", err)
	}

	if len(nodes) == 0 {
		return fmt.Errorf("addObject: no objects parsed")
	}

	for _, obj := range nodes {
		err = ValidateObject(*obj)
		if err != nil {
			return err
		}

		*s = append(*s, *obj)
	}

	return nil
}

// Read Object from a URL
func (s *ObjectStore) AddObjectFromURL(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("downloading file from URL: %w", err)
	}
	defer resp.Body.Close()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading downloaded file: %w", err)
	}

	return s.AddObject(contents)
}

// Read Object from Path
func (s *ObjectStore) AddObjectFromPath(path string) error {

	// Read the file from the local filesystem
	contents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("adding object from path: %w", err)
	}

	return s.AddObject(contents)
}

// Read Object/s from a Directory
func (s *ObjectStore) AddObjectsFromDir(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("reading directory: %w", err)
	}

	for _, entry := range entries {
		// Skip directories
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(path, entry.Name())

		// skip non yaml files
		ext := filepath.Ext(filePath)
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		// Read file and add to Store
		err = s.AddObjectFromPath(filePath)
		if err != nil {
			return fmt.Errorf("adding object from path %q: %w", filePath, err)
		}
	}

	return nil
}

// Filter ObjectStore by kind
func (s *ObjectStore) FilterByKind(kind string) (ObjectStore, error) {
	filteredObjects := ObjectStore{}
	for _, obj := range *s {
		objKind := obj.GetKind()
		if strings.EqualFold(objKind, kind) {
			filteredObjects = append(filteredObjects, obj)
		}
	}
	return filteredObjects, nil
}
