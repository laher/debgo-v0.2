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

package deb

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// TarGzWriter encapsulates the tar, gz and file operations of a TarGz file
type TarGzWriter struct {
	Filename string         // Filename
	Fw       io.WriteCloser // File writer
	Tw       *tar.Writer    // Tar writer (wraps the io.writer)
	Gw       *gzip.Writer   // Gzip writer (wraps the tar writer)
}

func NewTarHeader(path string, datalen int64, mode int64) *tar.Header {
	h := new(tar.Header)
	//backslash-only paths
	h.Name = strings.Replace(path, "\\", "/", -1)
	h.Size = datalen
	h.Mode = mode
	h.ModTime = time.Now()
	return h
}

// creates the file on disk, and wraps the os.File with a Tar writer and Gzip writer
func (tgzw *TarGzWriter) Create() error {
	var err error
	tgzw.Fw, err = os.Create(tgzw.Filename)
	if err != nil {
		return err
	}
	// gzip write
	tgzw.Gw = gzip.NewWriter(tgzw.Fw)
	// tar write
	tgzw.Tw = tar.NewWriter(tgzw.Gw)
	return nil
}

// Closes all 3 writers
func (tgzw *TarGzWriter) Close() error {
	err1 := tgzw.Tw.Close()
	err2 := tgzw.Gw.Close()
	err3 := tgzw.Fw.Close()
	if err1 != nil {
		return fmt.Errorf("Error closing Tar Writer %v", err1)
	}
	if err2 != nil {
		return fmt.Errorf("Error closing Gzip Writer %v", err2)
	}
	return err3
}

// add a file from the file system
func (tgzw *TarGzWriter) AddFile(sourceFile, destName string) error {
	fi, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	finf, err := fi.Stat()
	if err != nil {
		return err
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

// Add resources from file system.
// The key should be the destination filename. Value is the local filesystem path
func (tgzw *TarGzWriter) AddFiles(resources map[string]string) error {
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

// add a file by bytes
func (tgzw *TarGzWriter) AddBytes(bytes []byte, destName string, mode int64) error {
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
func NewTarGzWriter(archiveFilename string) (*TarGzWriter, error) {
	tgzw := &TarGzWriter{Filename: archiveFilename}
	err := tgzw.Create()
	return tgzw, err
}
