package models

type MusicDTO struct {
	Url string `json:"url"`
}

type Config struct {
	App struct {
		Version string `yaml:"version"`
	} `yaml:"app"`
}
