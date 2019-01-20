// +build linux

package build

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/tanema/ludus/src/build/config"
	"github.com/tanema/ludus/src/build/love"
	"github.com/tanema/ludus/src/build/windows"
)

// Run will create a build for the current OS and run it
func Run() error {
	cfg, err := config.Parse()
	if err != nil {
		return err
	}

	loveFilePath, err := love.Build(cfg)
	if err != nil {
		return err
	}

	cmd := exec.Command("love", filepath.Join(cfg.BuildDirectory, cfg.Title+".love"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
