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
	"archive/tar"
	"github.com/laher/argo/ar"
	"github.com/laher/debgo-v0.2/targz"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type DebReader struct {
	Reader io.Reader
	ArReader *ar.Reader
	HasDebianVersion bool
}

func NewDebReader(rdr io.Reader) (*DebReader, error) {
	drdr := &DebReader{Reader:rdr, HasDebianVersion: false}
	arr, err := ar.NewReader(rdr)
	if err != nil {
		return nil, err
	}
	drdr.ArReader = arr
	return drdr, err
}

//Gets next tar header for supported types
func (drdr *DebReader) NextTar() (string, *tar.Reader, error) {
	for {
		hdr, err := drdr.ArReader.Next()
		if err != nil {
			return "", nil, err
		}
		if hdr.Name == "debian-binary" {
			drdr.HasDebianVersion = true
			continue
		}
		if strings.HasSuffix(hdr.Name, ".tar.gz") {
			tgzr, err := targz.NewReader(drdr.ArReader)
			return hdr.Name, tgzr.Reader, err
		}
		// else return error
		return hdr.Name, nil, fmt.Errorf("Unsuported file type: %s", hdr.Name)
	}
}


func DebGetContents(rdr io.Reader, topLevelFilename string) ([]string, error) {
	ret := []string{}
	fileNotFound := true
	arr, err := ar.NewReader(rdr)
	if err != nil {
		return nil, err
	}
	for {
		hdr, err := arr.Next()
		if err == io.EOF {
			// end of ar archive
			break
		}
		if err != nil {
			return nil, err
		}
		if hdr.Name == topLevelFilename {
			fileNotFound = false
			tgzr, err := targz.NewReader(arr)
			if err != nil {
				tgzr.Close()
				return nil, err
			}
			for {
				thdr, err := tgzr.Next()
				if err == io.EOF {
					// end of tar.gz archive
					break
				}
				if err != nil {
					tgzr.Close()
					return nil, err
				}
				ret = append(ret, thdr.Name)
			}
		}
	}
	if fileNotFound {
		return nil, fmt.Errorf("File not found")
	}
	return ret, nil
}

func DebExtractFileL2(rdr io.Reader, topLevelFilename string, secondLevelFilename string, destination io.Writer) error {
	arr, err := ar.NewReader(rdr)
	if err != nil {
		return err
	}
	for {
		hdr, err := arr.Next()
		if err == io.EOF {
			// end of ar archive
			break
		}
		if err != nil {
			return err
		}
		if hdr.Name == topLevelFilename {
			tgzr, err := targz.NewReader(arr)
			if err != nil {
				tgzr.Close()
				return err
			}
			for {
				thdr, err := tgzr.Next()
				if err == io.EOF {
					// end of tar.gz archive
					break
				}
				if err != nil {
					tgzr.Close()
					return err
				}
				if thdr.Name == secondLevelFilename {
					_, err = io.Copy(destination, tgzr)
					tgzr.Close()
					return nil
				} else {
					//SKIP
					log.Printf("File %s", thdr.Name)
				}
			}
			tgzr.Close()
		}
	}
	return fmt.Errorf("File not found")
}

// ParseDebMetadata reads an artifact's contents.
func DebParseMetadata(rdr io.Reader) (*Package, error) {

	arr, err := ar.NewReader(rdr)
	if err != nil {
		return nil, err
	}

	var pkg *Package
	pkg = nil
	hasDataArchive := false
	hasControlArchive := false
	hasControlFile := false
	hasDebianBinaryFile := false

	// Iterate through the files in the archive.
	for {
		hdr, err := arr.Next()
		if err == io.EOF {
			// end of ar archive
			break
		}
		if err != nil {
			return nil, err
		}
		//		t.Logf("File %s:\n", hdr.Name)
		if hdr.Name == BinaryDataArchiveNameDefault {
			// SKIP!
			hasDataArchive = true
		} else if hdr.Name == BinaryControlArchiveNameDefault {
			// Find control file
			hasControlArchive = true
			tgzr, err := targz.NewReader(arr)
			if err != nil {
				return nil, err
			}
			for {
				thdr, err := tgzr.Next()
				if err == io.EOF {
					// end of tar.gz archive
					break
				}
				if err != nil {
					return nil, err
				}
				if thdr.Name == "control" {
					hasControlFile = true
					dscr := NewDscReader(tgzr)
					pkg, err = dscr.Parse()
					if err != nil {
						return nil, err
					}
				} else {
					//SKIP
					log.Printf("File %s", thdr.Name)
				}
			}

		} else if hdr.Name == "debian-binary" {
			b, err := ioutil.ReadAll(arr)
			if err != nil {
				return nil, err
			}
			hasDebianBinaryFile = true
			if string(b) != "2.0\n" {
				return nil, fmt.Errorf("Binary version not valid: %s", string(b))
			}
		} else {
			return nil, fmt.Errorf("Unsupported file %s", hdr.Name)
		}
	}

	if !hasDebianBinaryFile {
		return nil, fmt.Errorf("No debian-binary file in .deb archive")
	}
	if !hasDataArchive {
		return nil, fmt.Errorf("No data.tar.gz file in .deb archive")
	}
	if !hasControlArchive {
		return nil, fmt.Errorf("No control.tar.gz file in .deb archive")
	}
	if !hasControlFile {
		return nil, fmt.Errorf("No debian/control file in control.tar.gz")
	}
	return pkg, nil
}
