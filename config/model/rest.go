package model

type Rest struct {
	Port       int    `json:"port"`
	PprofToken string `json:"pprof_token"`
}
