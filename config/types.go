package config

import "go-template/config/model"

type Config struct {
	Rest *model.Rest `json:"rest"`
}
