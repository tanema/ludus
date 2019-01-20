package build

import (
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/tanema/ludus/src/build/config"
	"github.com/tanema/ludus/src/build/linux"
	"github.com/tanema/ludus/src/build/love"
	"github.com/tanema/ludus/src/build/macos"
	"github.com/tanema/ludus/src/build/windows"
)

const downloadsURL = "https://bitbucket.org/rude/love/downloads/"

type buildPipepipeline func(string, *config.Config, bool) error

// Build will create all the configured builds
func Build() error {
	cfg, err := config.Parse()
	if err != nil {
		return err
	}

	loveFilePath, err := love.Build(cfg)
	if err != nil {
		return err
	}

	var g errgroup.Group
	pipelines := []buildPipepipeline{macos.Build, windows.Build, linux.Build}
	for _, p := range pipelines {
		build := p
		g.Go(func() error {
			return build(loveFilePath, cfg, false)
		})
	}
	return g.Wait()
}

// Clean will clean up the build directory
func Clean() error {
	cfg, err := config.Parse()
	if err != nil {
		return err
	}
	return os.RemoveAll(cfg.BuildDirectory)
}
