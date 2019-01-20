package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

var defaultExclude = []string{
	`\.git`, `\.hg`, `\.svn`, "_darcs", "CVS",
	`\.DS_Store`, `.sublime-(project|workspace)`,
	`Thumbs\.db`, `desktop\.ini`, `ludus\.yaml`,
}

// Config descripts a wrp config
type Config struct {
	Title           string   `mapstructure:"title"`
	OSXIcon         string   `mapstructure:"macos_icon"`
	LoveVersion     string   `mapstructure:"love_version"`
	Version         string   `mapstructure:"version"`
	Author          string   `mapstructure:"author"`
	Email           string   `mapstructure:"email"`
	Description     string   `mapstructure:"description"`
	Homepage        string   `mapstructure:"homepage"`
	Identifier      string   `mapstructure:"identifier"`
	ExcludeFileList []string `mapstructure:"exclude"`
	SourceDirectory string   `mapstructure:"source_directory"`
	BuildDirectory  string   `mapstructure:"build_directory"`
}

// Parse will find and parse the config file
func Parse() (*Config, error) {
	cfg := &Config{}
	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return cfg, err
	}
	if cfg.LoveVersion == "" {
		return cfg, fmt.Errorf("A love version is required to build")
	}
	releasePath, _ := filepath.Rel(filepath.Dir(cfg.SourceDirectory), cfg.BuildDirectory)
	cfg.ExcludeFileList = append(cfg.ExcludeFileList, releasePath)
	cfg.ExcludeFileList = append(cfg.ExcludeFileList, defaultExclude...)
	return cfg, nil
}
