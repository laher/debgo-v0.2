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
	"io"
	"strings"
	"os"
	"time"
	"log"
)


func newTarHeader(path string, datalen int64, mode int64) *tar.Header {
	h := new(tar.Header)
	//backslash-only paths
	h.Name = strings.Replace(path, "\\", "/", -1)
	h.Size = datalen
	h.Mode = mode
	h.ModTime = time.Now()
	return h
}

type TarGzWriter struct {
	Fw io.WriteCloser
	Tw *tar.Writer
	Gw *gzip.Writer
}

func (tgzw *TarGzWriter) Open(archiveFilename string) error {
	var err error
	tgzw.Fw, err = os.Create(archiveFilename)
	if err != nil {
		return err
	}
	// gzip write
	tgzw.Gw = gzip.NewWriter(tgzw.Fw)
	// tar write
	tgzw.Tw = tar.NewWriter(tgzw.Gw)
	return nil
}

func (tgzw *TarGzWriter) Close() error {
	err1 := tgzw.Tw.Close()
	err2 := tgzw.Gw.Close()
	err3 := tgzw.Fw.Close()
	if err1 != nil {
		log.Printf("Error closing Tar Writer {}", err1)
		return err1
	}
	if err2 != nil {
		log.Printf("Error closing Gzip Writer {}", err2)
		return err2
	}
	return err3
}

func (tgzw *TarGzWriter) AddFile(sourceFile, destName string) error {
	fi, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	finf, err := fi.Stat()
	if err != nil {
		return err
	}
	err = tgzw.Tw.WriteHeader(newTarHeader(destName, finf.Size(), int64(finf.Mode())))
	if err != nil {
		return err
	}
	_, err = io.Copy(tgzw.Tw, fi)
	if err != nil {
		return err
	}
	return nil
}

func (tgzw *TarGzWriter) AddBytes(bytes []byte, destName string, mode int64) error {
	err := tgzw.Tw.WriteHeader(newTarHeader(destName, int64(len(bytes)), mode))
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
	tgzw := &TarGzWriter{}
	err := tgzw.Open(archiveFilename)
	return tgzw, err
}

