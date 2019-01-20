package love

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/tanema/ludus/src/build/config"
	"github.com/tanema/ludus/src/build/zipbuild"
)

// Build compiles the lua code into a .love file
func Build(cfg *config.Config) (string, error) {
	outFilePath := filepath.Join(cfg.BuildDirectory, cfg.Title+".love")

	builder, err := zipbuild.New(outFilePath, false)
	if err != nil {
		return "", err
	}
	defer builder.Close()

	basePath := filepath.Dir(cfg.SourceDirectory)
	return outFilePath, filepath.Walk(cfg.SourceDirectory, func(currentPath string, fileInfo os.FileInfo, err error) error {
		if err != nil || fileInfo.IsDir() {
			return err
		}

		for _, p := range cfg.ExcludeFileList {
			if matched, _ := regexp.MatchString(p, currentPath); matched {
				return nil
			}
		}

		file, err := os.Open(currentPath)
		if err != nil {
			return err
		}

		relativeFilePath, err := filepath.Rel(basePath, currentPath)
		if err != nil {
			return err
		}

		return builder.Write(relativeFilePath, file)
	})
}
