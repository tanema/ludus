package zipbuild

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Builder wraps the operations for creating and writing to a zip file
type Builder struct {
	debug     bool
	dest      string
	outFile   *os.File
	zipWriter *zip.Writer
}

// New creates a new zipbuilder
func New(dest string, debug bool) (*Builder, error) {
	err := os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return nil, err
	}

	outFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &Builder{
		dest:      filepath.Dir(dest),
		zipWriter: zip.NewWriter(outFile),
		outFile:   outFile,
		debug:     debug,
	}, nil
}

func (builder *Builder) Write(dest string, sources ...io.ReadCloser) error {
	zipFileWriter, err := builder.zipWriter.Create(dest)
	if err != nil {
		return err
	}
	writers := []io.Writer{zipFileWriter}

	if builder.debug {
		path := filepath.Join(builder.dest, dest)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		outFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
		defer outFile.Close()
		writers = append(writers, outFile)
	}

	for _, source := range sources {
		_, err = io.Copy(io.MultiWriter(writers...), source)
		defer source.Close()
	}
	return err
}

// Close closes all the builder resources
func (builder *Builder) Close() error {
	if err := builder.zipWriter.Close(); err != nil {
		return err
	}
	return builder.outFile.Close()
}
