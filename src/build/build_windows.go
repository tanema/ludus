// +build windows

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

	if err := windows.Build(loveFilePath, cfg, true); err != nil {
		return err
	}

	arch := "32"
	if runtime.GOARCH == "amd64" {
		arch := "64"
	}

	executablePath := filepath.Join(
		cfg.BuildDirectory,
		"win",
		cfg.Title+"-win"+arch,
		cfg.Title+".exe",
	)

	cmd := exec.Command(executablePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
