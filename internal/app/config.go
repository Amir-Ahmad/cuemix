package app

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

//go:embed schema.cue
var configSchema string

func ParseConfig(path string) (Config, error) {

	ctx := cuecontext.New()

	fileInfo, err := os.Stat(path)
	if err != nil {
		return Config{}, err
	}

	// Load file or directory depending on what's received
	// load.Instances will always return one element here
	var inst *build.Instance
	if fileInfo.IsDir() {
		inst = load.Instances([]string{"."}, &load.Config{
			Dir: filepath.Join(path),
		})[0]
	} else {
		inst = load.Instances([]string{path}, nil)[0]
	}

	if inst.Err != nil {
		return Config{}, fmt.Errorf("loading cue instance: %w", inst.Err)
	}

	data := ctx.BuildInstance(inst)
	if data.Err() != nil {
		return Config{},
			fmt.Errorf("building cue instance '%s': %w", inst.DisplayPath, data.Err())
	}

	schema := ctx.CompileString(configSchema)
	if schema.Err() != nil {
		return Config{}, fmt.Errorf("compiling schema: %w", schema.Err())
	}

	// Unify schema with cue data, similar to & in cue
	value := schema.Unify(data)
	if value.Err() != nil {
		return Config{}, fmt.Errorf("unifying schema with data: %w", value.Err())
	}

	// Validate cue data and ensure all required fields are set
	err = value.Validate(cue.Concrete(true))
	if err != nil {
		return Config{}, fmt.Errorf("validating cue data: %w", err)
	}

	var config Config
	err = value.Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("decoding cue into struct: %w", err)
	}

	return config, nil
}
