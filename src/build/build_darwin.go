// +build darwin

package build

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tanema/ludus/src/build/config"
	"github.com/tanema/ludus/src/build/love"
	"github.com/tanema/ludus/src/build/macos"
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

	if err := macos.Build(loveFilePath, cfg, true); err != nil {
		return err
	}

	executablePath := filepath.Join(
		cfg.BuildDirectory,
		"darwin",
		cfg.Title+".app",
		"Contents",
		"MacOS",
		"love",
	)

	cmd := exec.Command(executablePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
