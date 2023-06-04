package helm

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/repo"
)

// GetChartUrl gets the tgz url of a chart and version
func (h HelmChart) GetChartUrl() (string, error) {
	chartRepo, err := repo.NewChartRepository(&repo.Entry{URL: h.Repo}, getter.All(&cli.EnvSettings{}))
	if err != nil {
		return "", fmt.Errorf("failed to initialise ChartRepository '%s': %w", h.Repo, err)
	}

	idxContents, err := chartRepo.DownloadIndexFile()
	if err != nil {
		return "", fmt.Errorf("downloading repository '%s' index file: %w", h.Repo, err)
	}

	index, err := repo.LoadIndexFile(idxContents)
	if err != nil {
		return "", fmt.Errorf("loading repository '%s' index file: %w", h.Repo, err)
	}

	chartInfo, err := index.Get(h.Chart, h.Version)
	if err != nil {
		return "",
			fmt.Errorf("getting chart '%s' version '%s': %w", h.Repo, h.Version, err)
	}

	return chartInfo.URLs[0], nil
}

// DownloadChart downloads a chart to a local directory
// If untar = false, the output is a tgz file
// If untar = true, the chart is extracted and the tgz removed
func (h HelmChart) DownloadChart(destDir string, untar bool) (string, error) {
	// Create dir if it doesn't exist
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// If untar = false, tarDest is the directory the tgz will be downloaded to
	var tarDest string

	// If untar = true:
	// untarDir is the subdirectory of destDir to extract the chart to
	// chartExtractedDir is the final extracted location of the templates,
	// by combining destDir, untarDir, and Chart name
	untarDir := h.Chart + "-" + h.Version
	var chartExtractedDir string

	if registry.IsOCI(h.Repo) {
		name := filepath.Base(h.Repo)
		tarDest = filepath.Join(destDir, name+"-"+h.Version+".tgz")
		chartExtractedDir = filepath.Join(destDir, untarDir, name)
	} else {
		tarDest = filepath.Join(destDir, h.Chart+"-"+h.Version+".tgz")
		chartExtractedDir = filepath.Join(destDir, untarDir, h.Chart)
	}

	// If file already exists, don't redownload
	if _, err := os.Stat(tarDest); err == nil && !untar {
		return tarDest, nil
	}

	// If extracted directory already exists, don't redownload
	if _, err := os.Stat(chartExtractedDir); err == nil && untar {
		return chartExtractedDir, nil
	}

	// Initialise helm action config
	config := new(action.Configuration)
	config.Init(nil, "", "secret", log.Printf)

	var chartUrl string
	if registry.IsOCI(h.Repo) {
		chartUrl = h.Repo
		config.RegistryClient, _ = registry.NewClient()
	} else {
		// For traditional helm repos, get the actual chart url by parsing index.yaml
		chartUrl, err = h.GetChartUrl()
		if err != nil {
			return "", fmt.Errorf("when getting chart url: %w", err)
		}
	}

	// Initialise pull with config
	pull := action.NewPullWithOpts(action.WithConfig(config))
	pull.Settings = cli.New()
	pull.DestDir = destDir
	pull.Untar = untar
	pull.Version = h.Version
	pull.UntarDir = untarDir

	// Download chart
	_, err = pull.Run(chartUrl)
	if err != nil {
		return "", fmt.Errorf("when pulling chart: %w", err)
	}

	if untar {
		return chartExtractedDir, nil
	}

	return tarDest, nil
}
