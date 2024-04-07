package config

import "go-template/config/model"

func GetRest() *model.Rest {
	return _default.Rest
}
