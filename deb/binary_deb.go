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
	"github.com/laher/argo/ar"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Architecture-specific build information
type BinaryArtifact struct {
	Architecture        Architecture
	Filename            string
	TmpDir              string
	DebianBinaryVersion string
	ControlArchFile     string
	DataArchFile        string
	Executables         []string
	IsVerbose           bool
}

// Factory of platform build information
func NewBinaryArtifact(architecture Architecture, filename string, tmpDir string, isVerbose bool) *BinaryArtifact {
	bdeb := &BinaryArtifact{Architecture: architecture, Filename: filename, TmpDir: tmpDir, IsVerbose: isVerbose}
	bdeb.SetDefaults()
	return bdeb
}

func (bdeb *BinaryArtifact) GetReader() (*ar.Reader, error) {
	fi, err := os.Open(bdeb.Filename)
	if err != nil {
		return nil, err
	}
	arr, err := ar.NewReader(fi)
	if err != nil {
		return nil, err
	}
	return arr, err
}

// ExtractAll extracts all contents from the Ar archive.
// It returns a slice of all filenames.
// In case of any error, it returns the error immediately
func (bdeb *BinaryArtifact) ExtractAll() ([]string, error) {
	arr, err := bdeb.GetReader()
	if err != nil {
		return nil, err
	}
	filenames := []string{}
	for {
		hdr, err := arr.Next()
		if err == io.EOF {
			// end of ar archive
			break
		}
		if err != nil {
			return nil, err
		}
		outFilename := filepath.Join(bdeb.TmpDir, hdr.Name)
		//fmt.Printf("Contents of %s:\n", hdr.Name)
		fi, err := os.Create(outFilename)
		if err != nil {
			return filenames, err
		}
		if _, err := io.Copy(fi, arr); err != nil {
			return filenames, err
		}
		err = fi.Close()
		if err != nil {
			return filenames, err
		}
		filenames = append(filenames, outFilename)
		//fmt.Println()
	}
	return filenames, nil
}

func (bdeb *BinaryArtifact) SetDefaults() {
	bdeb.DebianBinaryVersion = DEBIAN_BINARY_VERSION_DEFAULT
	bdeb.ControlArchFile = filepath.Join(bdeb.TmpDir, "control.tar.gz")
	bdeb.DataArchFile = filepath.Join(bdeb.TmpDir, "data.tar.gz")
}

func (bdeb *BinaryArtifact) WriteBytes(aw *ar.Writer, filename string, bytes []byte) error {
	hdr := &ar.Header{
		Name: filename,
		Size: int64(len(bytes))}
	if err := aw.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := aw.Write(bytes); err != nil {
		return err
	}
	return nil
}

func (bdeb *BinaryArtifact) WriteFromFile(aw *ar.Writer, filename string) error {
	finf, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if bdeb.IsVerbose {
		log.Printf("Finf size: %d", finf.Size())
	}
	hdr, err := ar.FileInfoHeader(finf)
	if err != nil {
		return err
	}
	if err := aw.WriteHeader(hdr); err != nil {
		return err
	}
	if bdeb.IsVerbose {
		log.Printf("Header Size: %d", hdr.Size)
	}
	fi, err := os.Open(filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(aw, fi); err != nil {
		return err
	}

	err = fi.Close()
	if err != nil {
		return err
	}
	return nil

}

func (bdeb *BinaryArtifact) Build() error {
	if bdeb.IsVerbose {
		log.Printf("Building deb %s", bdeb.Filename)
	}
	wtr, err := os.Create(bdeb.Filename)
	if err != nil {
		return err
	}

	aw := ar.NewWriter(wtr)

	if bdeb.IsVerbose {
		log.Printf("Writing debian-binary")
	}
	err = bdeb.WriteBytes(aw, "debian-binary", []byte(bdeb.DebianBinaryVersion+"\n"))
	if err != nil {
		return err
	}
	if bdeb.IsVerbose {
		log.Printf("Writing control file %s", bdeb.ControlArchFile)
	}
	err = bdeb.WriteFromFile(aw, bdeb.ControlArchFile)
	if err != nil {
		return err
	}
	if bdeb.IsVerbose {
		log.Printf("Writing data file %s", bdeb.DataArchFile)
	}
	err = bdeb.WriteFromFile(aw, bdeb.DataArchFile)
	if err != nil {
		return err
	}
	return nil
}
