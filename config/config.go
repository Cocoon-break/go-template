package config

import (
	"fmt"
	"os"

	"go-template/config/model"
	"go-template/pkg/env"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

var (
	_default        *Config
	appBaseDir      = "/usr/local/go-template"
	etcPath         = appBaseDir + "/etc"
	runtime         = etcPath + "/runtime.yml"
	dumpPath        = appBaseDir + "/dump"
	DefaultFilename = etcPath + "/go-template.yml"
)

func init() {
	appBaseDir = "/usr/local/" + env.AppName()
	etcPath = appBaseDir + "/etc"
	dumpPath = appBaseDir + "/dump"
	DefaultFilename = etcPath + fmt.Sprintf("/%s.yml", env.AppName())
	runtime = etcPath + "/runtime.yml"
}

func NewDefault() *Config {
	return &Config{
		Rest: &model.Rest{
			Port:       24407,
			PprofToken: "go-template",
		},
	}
}

// Parse parse from local file
func (c *Config) Parse(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}

// DumpAndSetGlobal dump and set global variables
func DumpAndSetGlobal(c *Config) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "dump and set global marshal")
	}
	_default = c

	if err := os.WriteFile(runtime, data, os.ModePerm); err != nil {
		return errors.Wrap(err, "dump and set global write file")
	}
	return nil
}

// Reload reload config from file
func Reload() error {
	if _default != nil {
		if err := _default.Parse(DefaultFilename); err != nil {
			return errors.Wrap(err, "reload")
		}
	} else {
		return errors.New("not init default config")
	}
	return nil
}
