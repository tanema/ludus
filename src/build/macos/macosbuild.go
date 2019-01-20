package macos

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tanema/ludus/src/build/bitbucket"
	"github.com/tanema/ludus/src/build/config"
	"github.com/tanema/ludus/src/build/zipbuild"
)

var (
	plistBundleIdent  = regexp.MustCompile("(<key>CFBundleIdentifier</key>\n.*<string>)(.*)(</string>)")
	plistBundleName   = regexp.MustCompile("(<key>CFBundleName</key>\n.*<string>)(.*)(</string>)")
	plistExportedType = regexp.MustCompile("<key>UTExportedTypeDeclarations</key>.*\n.*<array>(?s:.*)</array>")
)

// Build creates a macos build
func Build(loveFilePath string, cfg *config.Config, debug bool) error {
	macOSZip, err := bitbucket.Get(cfg.LoveVersion, "macos")
	if err != nil {
		return err
	}

	builder, err := zipbuild.New(filepath.Join(cfg.BuildDirectory, "darwin", cfg.Title+".zip"), debug)
	if err != nil {
		return err
	}
	defer builder.Close()

	for _, zipFile := range macOSZip.File {
		if info := zipFile.FileInfo(); info.IsDir() {
			continue
		}

		f, err := zipFile.Open()
		if err != nil {
			return err
		}

		if zipFile.Name == "love.app/Contents/Info.plist" {
			f = fixPlist(cfg, f)
		}

		if err := builder.Write(strings.Replace(zipFile.Name, "love.app", cfg.Title+".app", 1), f); err != nil {
			return err
		}
	}

	file, err := os.Open(loveFilePath)
	if err != nil {
		return err
	}

	return builder.Write(cfg.Title+".app/Contents/Resources/game.love", file)
}

func fixPlist(cfg *config.Config, original io.ReadCloser) io.ReadCloser {
	buf := new(bytes.Buffer)
	buf.ReadFrom(original)
	original.Close()
	plist := plistBundleIdent.ReplaceAllString(buf.String(), "$1 "+cfg.Identifier+" $3")
	plist = plistBundleName.ReplaceAllString(plist, "$1 "+cfg.Title+" $3")
	plist = plistExportedType.ReplaceAllString(plist, "")
	return ioutil.NopCloser(strings.NewReader(plist))
}
