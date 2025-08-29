package site

import (
	"encoding/json"
	"os"

	"github.com/tailscale/hujson"
)

type Config struct {
	Title        string `json:"title"`
	Author       string `json:"author"`
	Domain       string `json:"domain"`
	DirBuild     string `json:"dir_build"`
	DirAssets    string `json:"dir_assets"`
	EmailList    string `json:"email_list"`
	EnableDrafts bool
}

func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	b, err = hujson.Standardize(b)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
