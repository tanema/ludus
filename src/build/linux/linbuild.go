package linux

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xor-gate/debpkg"

	"github.com/tanema/ludus/src/build/config"
)

// Build creates a linux build
func Build(loveFilePath string, cfg *config.Config, debug bool) error {
	requirments := []string{cfg.Title, cfg.Version, cfg.Author, cfg.Email, cfg.Homepage, cfg.Description}
	for _, req := range requirments {
		if req == "" {
			return fmt.Errorf("Deb build requires Title, Version, Author, Email, Homepage and Description")
		}
	}

	deb := debpkg.New()
	defer deb.Close()
	deb.SetName(cfg.Title)
	deb.SetVersion(cfg.Version)
	deb.SetArchitecture("all")
	deb.SetMaintainer(cfg.Author)
	deb.SetMaintainerEmail(cfg.Email)
	deb.SetHomepage(cfg.Homepage)
	deb.SetDescription(cfg.Description)
	deb.SetDepends("love (>= " + cfg.LoveVersion + " ) ")
	deb.SetPriority("optional")

	if err := deb.AddFile(loveFilePath); err != nil {
		return err
	}

	outFilePath := filepath.Join(cfg.BuildDirectory, "linux", cfg.Title+".deb")
	if err := os.MkdirAll(filepath.Dir(outFilePath), 0755); err != nil {
		return err
	}

	return deb.Write(outFilePath)
}
