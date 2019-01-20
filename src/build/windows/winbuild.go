package windows

import (
	"os"
	"path/filepath"

	"github.com/tanema/ludus/src/build/bitbucket"
	"github.com/tanema/ludus/src/build/config"
	"github.com/tanema/ludus/src/build/zipbuild"
)

var requiredFiles = []string{
	"SDL2.dll",
	"OpenAL32.dll",
	"license.txt",
	"love.exe",
	"love.dll",
	"lua51.dll",
	"mpg123.dll",
	"msvcp120.dll",
	"msvcr120.dll",
}

// Build will generate a windows build
func Build(loveFilePath string, cfg *config.Config, debug bool) error {
	if err := buildArch("win64", loveFilePath, cfg, debug); err != nil {
		return err
	}
	return buildArch("win32", loveFilePath, cfg, debug)
}

func buildArch(arch string, loveFilePath string, cfg *config.Config, debug bool) error {
	winZip, err := bitbucket.Get(cfg.LoveVersion, arch)
	if err != nil {
		return err
	}

	builder, err := zipbuild.New(filepath.Join(cfg.BuildDirectory, "win", cfg.Title+"-"+arch+".zip"), debug)
	if err != nil {
		return err
	}
	defer builder.Close()

	for _, zipFile := range winZip.File {
		baseName := filepath.Base(zipFile.Name)
		info := zipFile.FileInfo()

		if info.IsDir() || !contains(requiredFiles, baseName) {
			continue
		}

		f, err := zipFile.Open()
		if err != nil {
			return err
		}

		if baseName == "love.exe" {
			lovefile, err := os.Open(loveFilePath)
			if err != nil {
				return err
			}
			err = builder.Write(cfg.Title+".exe", f, lovefile)
		} else {
			err = builder.Write(baseName, f)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
