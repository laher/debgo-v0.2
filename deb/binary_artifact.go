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
	"fmt"
	"github.com/laher/argo/ar"
	"io"
	"os"
	"path/filepath"
)

// Architecture-specific build information
type BinaryArtifact struct {
	*BinaryPackage
	Architecture        Architecture
	Filename            string
	DebianBinaryVersion string
	ControlArchFile     string
	DataArchFile        string
	Binaries            map[string]string
}

// Factory of platform build information
func NewBinaryArtifact(binaryPackage *BinaryPackage, architecture Architecture) *BinaryArtifact {
	bdeb := &BinaryArtifact{BinaryPackage: binaryPackage, Architecture: architecture}
	bdeb.SetDefaults()
	return bdeb
}

// InitControlArchive initialises and returns the 'control.tar.gz' archive
func (pkg *BinaryPackage) InitControlArchive(build *BuildParams) (*TarGzWriter, error) {
	archiveFilename := filepath.Join(build.TmpDir, "control.tar.gz")
	tgzw, err := NewTarGzWriter(archiveFilename)
	if err != nil {
		return nil, err
	}
	return tgzw, err
}

// InitDataArchive initialises and returns the 'data.tar.gz' archive
func (pkg *BinaryPackage) InitDataArchive(build *BuildParams) (*TarGzWriter, error) {
	archiveFilename := filepath.Join(build.TmpDir, "data.tar.gz")
	tgzw, err := NewTarGzWriter(archiveFilename)
	if err != nil {
		return nil, err
	}
	return tgzw, err
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
func (bdeb *BinaryArtifact) ExtractAll(build *BuildParams) ([]string, error) {
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
		outFilename := filepath.Join(build.TmpDir, hdr.Name)
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
	bdeb.Filename = fmt.Sprintf("%s_%s_%s.deb", bdeb.Name, bdeb.Version, bdeb.Architecture) //goxc_0.5.2_i386.deb")
	bdeb.DebianBinaryVersion = DebianBinaryVersionDefault
	bdeb.ControlArchFile = "control.tar.gz"
	bdeb.DataArchFile = "data.tar.gz"
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
	hdr, err := ar.FileInfoHeader(finf)
	if err != nil {
		return err
	}
	if err := aw.WriteHeader(hdr); err != nil {
		return err
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

func (bdeb *BinaryArtifact) Build(build *BuildParams) error {
	wtr, err := os.Create(filepath.Join(build.DestDir, bdeb.Filename))
	if err != nil {
		return err
	}

	aw := ar.NewWriter(wtr)

	err = bdeb.WriteBytes(aw, "debian-binary", []byte(bdeb.DebianBinaryVersion+"\n"))
	if err != nil {
		return fmt.Errorf("Error writing debian-binary into .ar archive: %v", err)
	}
	err = bdeb.WriteFromFile(aw, filepath.Join(build.TmpDir, bdeb.ControlArchFile))
	if err != nil {
		return fmt.Errorf("Error writing control archive into .ar archive: %v", err)
	}
	err = bdeb.WriteFromFile(aw, filepath.Join(build.TmpDir, bdeb.DataArchFile))
	if err != nil {
		return fmt.Errorf("Error writing data archive into .ar archive: %v", err)
	}
	err = aw.Close()
	if err != nil {
		return fmt.Errorf("Error closing .ar archive: %v", err)
	}
	return nil
}
