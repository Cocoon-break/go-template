package option

import (
	"fmt"
	"os"
	"strings"

	"go-template/config"
	"go-template/pkg/zlog"
)

type Option struct {
	ConfigFile string
}

func New() *Option {
	return &Option{
		ConfigFile: config.DefaultFilename,
	}
}

func (o *Option) Validate() error {
	if !IsFileExist(o.ConfigFile) {
		return fmt.Errorf("%s not exist", o.ConfigFile)
	}
	return nil
}

func (o *Option) Config(filePath string) (*config.Config, error) {
	if strings.TrimSpace(filePath) != "" {
		o.ConfigFile = filePath
	}
	cfg := config.NewDefault()
	if err := o.Validate(); err != nil {
		zlog.Warn("option_config", zlog.String("u_msg", "will use default config"), zlog.Any("err", err.Error()))
	} else {
		if err := cfg.Parse(o.ConfigFile); err != nil {
			zlog.Error("option_config", zlog.String("u_msg", "parse "+o.ConfigFile), zlog.Any("err", err.Error()))
			return nil, err
		}
	}
	return cfg, nil
}

func IsPathExist(path string) bool {
	return isPathExist(path)
}

func IsFileExist(file string) bool {
	return isPathExist(file)
}

func isPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
