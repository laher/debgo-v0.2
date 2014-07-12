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
