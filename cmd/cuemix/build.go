package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/amir-ahmad/cuemix/internal/app"
	"github.com/amir-ahmad/cuemix/pkg/k8s"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [PATH]",
	Args:  cobra.ExactArgs(1),
	Short: "Build and output Kubernetes resources",
	RunE:  build,
}

var kindFilter string

func init() {
	buildCmd.Flags().StringVarP(&kindFilter, "kind", "k", "", "Filter output by kind")
	rootCmd.AddCommand(buildCmd)
}

func build(cmd *cobra.Command, args []string) error {
	path := args[0]

	config, err := app.ParseConfig(path)
	if err != nil {
		return fmt.Errorf("parsing cue config: %w", err)
	}

	// Determine basePath, if path is a file, basePath is the dir containing it
	fileInfo, err := os.Stat(path)
	var basePath string
	if err != nil {
		return fmt.Errorf("checking path: %w", err)
	}
	if fileInfo.IsDir() {
		basePath = path
	} else {
		basePath = filepath.Dir(path)
	}

	// Load all objects found in the config
	store, err := config.LoadObjects(basePath)
	if err != nil {
		return fmt.Errorf("loading objects: %w", err)
	}

	// Apply patches
	for name, patch := range config.StrategicPatches {
		patchObject, err := k8s.FromMap(patch.Object, false)
		if err != nil {
			return fmt.Errorf("converting patch %s to Object: %w", name, err)
		}
		err = store.StrategicPatch(patchObject, patch.Target)
		if err != nil {
			return fmt.Errorf("when applying patch %s: %w", name, err)
		}
	}

	// Apply json patches
	for name, v := range config.JsonPatches {
		err = store.JsonPatch(v.Patch, v.Target)
		if err != nil {
			return fmt.Errorf("when applying json patch %s: %w", name, err)
		}
	}

	// Filter by kind if provided
	var output k8s.ObjectStore
	if kindFilter != "" {
		output, err = store.FilterByKind(kindFilter)
		if err != nil {
			return fmt.Errorf("filtering by kind: %w", err)
		}
	} else {
		output = store
	}

	// Output objects
	err = output.Print(true)
	if err != nil {
		return err
	}

	return nil
}
