/*
   Copyright 2013 Am Laher

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package targz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Writer encapsulates the tar, gz and file operations of a TarGz file
type Writer struct {
	Filename string       // Filename
	Fw       io.Writer    // File writer
	Tw       *tar.Writer  // Tar writer (wraps the io.writer)
	Gw       *gzip.Writer // Gzip writer (wraps the tar writer)
}

// Close closes all 3 writers
// Returns the first error
// TODO: close Fw if possible?
func (tgzw *Writer) Close() error {
	err1 := tgzw.Tw.Close()
	err2 := tgzw.Gw.Close()
	//err3 := tgzw.Fw.Close()
	if err1 != nil {
		return fmt.Errorf("Error closing Tar Writer %v", err1)
	}
	if err2 != nil {
		return fmt.Errorf("Error closing Gzip Writer %v", err2)
	}
	//return err3
	return nil
}

// AddFile adds a file from the file system
// This is just a helper function
// TODO: directories
func (tgzw *Writer) AddFile(sourceFile, destName string) error {
	fi, err := os.Open(sourceFile)
	defer fi.Close()
	if err != nil {
		return err
	}
	finf, err := fi.Stat()
	if err != nil {
		return err
	}

	//recurse as necessary
	if finf.IsDir() {
		return fmt.Errorf("Can't add a directory using AddFile. See AddFileOrDir")
	}
	err = tgzw.Tw.WriteHeader(NewTarHeader(destName, finf.Size(), int64(finf.Mode())))
	if err != nil {
		return err
	}
	_, err = io.Copy(tgzw.Tw, fi)
	if err != nil {
		return err
	}
	return nil
}

func (tgzw *Writer) AddFileOrDir(sourceFile, destName string) error {
	finf, err := os.Stat(sourceFile)
	if err != nil {
		return err
	}
	//recurse as necessary
	if finf.IsDir() {
		err = filepath.Walk(sourceFile, func(path string, info os.FileInfo, err2 error) error {
			if info != nil && !info.IsDir() {
				rel, err := filepath.Rel(sourceFile, path)
				if err == nil {
					return tgzw.AddFile(rel, path)
				}
				return err
			}
			return nil
		})
		// return now
		return err
	}

	return tgzw.AddFile(sourceFile, destName)
}

// AddFiles adds resources from file system.
// The key should be the destination filename. Value is the local filesystem path
func (tgzw *Writer) AddFiles(resources map[string]string) error {
	if resources != nil {
		for name, localPath := range resources {
			err := tgzw.AddFile(localPath, name)
			if err != nil {
				tgzw.Close()
				return err
			}
		}
	}
	return nil
}

// AddBytes adds a file by bytes
func (tgzw *Writer) AddBytes(bytes []byte, destName string, mode int64) error {
	err := tgzw.Tw.WriteHeader(NewTarHeader(destName, int64(len(bytes)), mode))
	if err != nil {
		return err
	}
	_, err = tgzw.Tw.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// NewWriter is a factory for Writer
func NewWriterFromFile(archiveFilename string) (*Writer, error) {
	fw, err := os.Create(archiveFilename)
	if err != nil {
		return nil, err
	}
	tgzw := NewWriter(fw)
	tgzw.Filename = archiveFilename
	return tgzw, err
}

// Create creates the file on disk, and
// wraps the io.Writer with a Tar writer and Gzip writer
func NewWriter(w io.Writer) *Writer {

	tgzw := &Writer{Fw: w}
	// gzip writer
	tgzw.Gw = gzip.NewWriter(tgzw.Fw)
	// tar writer
	tgzw.Tw = tar.NewWriter(tgzw.Gw)
	return tgzw
}
