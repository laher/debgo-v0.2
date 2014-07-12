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

package debgen

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// NewTarHeader is a factory for a tar header. Fixes slashes, populates ModTime
func NewTarHeader(path string, datalen int64, mode int64) *tar.Header {
	h := new(tar.Header)
	//slash-only paths
	h.Name = strings.Replace(path, "\\", "/", -1)
	if strings.HasPrefix(h.Name, "/") {
		h.Name = h.Name[1:]
	}
	h.Size = datalen
	h.Mode = mode
	h.ModTime = time.Now()
	return h
}

// TarAddFile adds a file from the file system
// This is just a helper function
// TODO: directories
func TarAddFile(tw *tar.Writer, sourceFile, destName string) error {
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
		return fmt.Errorf("Can't add a directory using TarAddFile. See AddFileOrDir")
	}
	err = tw.WriteHeader(NewTarHeader(destName, finf.Size(), int64(finf.Mode())))
	if err != nil {
		return err
	}
	_, err = io.Copy(tw, fi)
	if err != nil {
		return err
	}
	return nil
}

func TarAddFileOrDir(tw *tar.Writer, sourceFile, destName string) error {
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
					return TarAddFile(tw, rel, path)
				}
				return err
			}
			return nil
		})
		// return now
		return err
	}

	return TarAddFile(tw, sourceFile, destName)
}

// TarAddFiles adds resources from file system.
// The key should be the destination filename. Value is the local filesystem path
func TarAddFiles(tw *tar.Writer, resources map[string]string) error {
	if resources != nil {
		for name, localPath := range resources {
			err := TarAddFile(tw, localPath, name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// TarAddBytes adds a file by bytes with a given path
func TarAddBytes(tw *tar.Writer, bytes []byte, destName string, mode int64) error {
	err := tw.WriteHeader(NewTarHeader(destName, int64(len(bytes)), mode))
	if err != nil {
		return err
	}
	_, err = tw.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
