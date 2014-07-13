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

// Deb is an architecture-specific deb
type Deb struct {
	Package             *Package
	Architecture        Architecture
	Filename            string
	DebianBinaryVersion string
	DebianArchive     string
	DataArchive        string
	MappedFiles         map[string]string
}

// NewDebs gets and returns an artifact for each architecture.
// Returns an error if the package's architecture is un-parseable
func NewDebs(pkg *Package) (map[Architecture]*Deb, error) {
	arches, err := pkg.GetArches()
	if err != nil {
		return nil, err
	}
	ret := map[Architecture]*Deb{}
	for _, arch := range arches {
		archArtifact := NewDeb(pkg, arch)
		ret[arch] = archArtifact
	}
	return ret, nil
}

// Factory of platform build information
func NewDeb(pkg *Package, architecture Architecture) *Deb {
	bdeb := &Deb{Package: pkg, Architecture: architecture}
	bdeb.SetDefaults()
	return bdeb
}

// GetReader opens up a new .ar reader
func (bdeb *Deb) GetReader() (*ar.Reader, error) {
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
func (bdeb *Deb) ExtractAll(destDir string) ([]string, error) {
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
		outFilename := filepath.Join(destDir, hdr.Name)
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

func (bdeb *Deb) SetDefaults() {
	bdeb.Filename = fmt.Sprintf("%s_%s_%s.deb", bdeb.Package.Name, bdeb.Package.Version, bdeb.Architecture) //goxc_0.5.2_i386.deb")
	bdeb.DebianBinaryVersion = DebianBinaryVersionDefault
	bdeb.DebianArchive = BinaryControlArchiveNameDefault
	bdeb.DataArchive = BinaryDataArchiveNameDefault
}

func (bdeb *Deb) WriteBytes(aw *ar.Writer, filename string, bytes []byte) error {
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

func (bdeb *Deb) WriteFromFile(aw *ar.Writer, filename string) error {
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

func (bdeb *Deb) Build(tempDir, destDir string) error {
	wtr, err := os.Create(filepath.Join(destDir, bdeb.Filename))
	if err != nil {
		return err
	}
	defer wtr.Close()

	aw := ar.NewWriter(wtr)

	err = bdeb.WriteBytes(aw, "debian-binary", []byte(bdeb.DebianBinaryVersion+"\n"))
	if err != nil {
		return fmt.Errorf("Error writing debian-binary into .ar archive: %v", err)
	}
	err = bdeb.WriteFromFile(aw, filepath.Join(tempDir, bdeb.DebianArchive))
	if err != nil {
		return fmt.Errorf("Error writing control archive into .ar archive: %v", err)
	}
	err = bdeb.WriteFromFile(aw, filepath.Join(tempDir, bdeb.DataArchive))
	if err != nil {
		return fmt.Errorf("Error writing data archive into .ar archive: %v", err)
	}
	err = aw.Close()
	if err != nil {
		return fmt.Errorf("Error closing .ar archive: %v", err)
	}
	return nil
}
