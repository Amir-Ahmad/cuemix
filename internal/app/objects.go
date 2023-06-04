package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/amir-ahmad/cuemix/pkg/k8s"
)

// Load all objects found in config from helm, manifests, etc
func (c Config) LoadObjects(basePath string) (k8s.ObjectStore, error) {
	var store k8s.ObjectStore

	// Process Helm charts
	for _, chart := range c.HelmCharts {

		chartType := GetRepoType(basePath, chart.Repo)

		var dest string
		var err error
		if chartType == "local" {
			dest = ResolveLocalChartDir(basePath, chart.Repo)
		} else {
			dest, err = chart.DownloadChart(c.Config.Helm.DestDir, c.Config.Helm.Untar)
			if err != nil {
				return nil, fmt.Errorf("error downloading chart '%s' : %w", chart.Chart, err)
			}
		}

		objects, err := chart.Template(dest, c.Config.Helm.KubeVersion)
		if err != nil {
			return nil, fmt.Errorf("error templating chart '%s' : %w", chart.Chart, err)
		}

		for name, obj := range objects {
			err = store.AddObject([]byte(obj))
			if err != nil {
				return nil, fmt.Errorf("error loading '%s' : %w", name, err)
			}
		}
	}

	// Load manifests
	for _, path := range c.Manifests {
		var err error

		manifestType, err := GetPathType(basePath, path)
		if err != nil {
			return nil, err
		}

		switch manifestType {
		case "file":
			err = store.AddObjectFromPath(filepath.Join(basePath, path))
		case "directory":
			err = store.AddObjectsFromDir(filepath.Join(basePath, path))
		case "url":
			err = store.AddObjectFromURL(path)
		}

		if err != nil {
			return nil, fmt.Errorf("error loading '%s' : %w", path, err)
		}
	}
	return store, nil
}

// GetPathType checks if manifest path is a url, file, or directory
func GetPathType(basePath string, path string) (string, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return "url", nil
	}

	fileInfo, err := os.Stat(filepath.Join(basePath, path))
	if err != nil {
		return "", fmt.Errorf("GetPathType: %w", err)
	}

	if fileInfo.IsDir() {
		return "directory", nil
	}

	return "file", nil
}

func ResolveLocalChartDir(basePath string, repo string) string {
	if filepath.IsAbs(repo) {
		return repo
	}

	return filepath.Join(basePath, repo)
}

// GetRepoType checks if Helm Repo is helm, oci, or a local directory
func GetRepoType(basePath string, repo string) string {
	if strings.HasPrefix(repo, "http://") || strings.HasPrefix(repo, "https://") {
		return "helm"
	} else if strings.HasPrefix(repo, "oci://") {
		return "oci"
	}
	return "local"
}
