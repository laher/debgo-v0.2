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


import(
	"io"
	"log"
	"fmt"
	"io/ioutil"
	"github.com/laher/argo/ar"
	"github.com/laher/debgo-v0.2/targz"
	"bufio"
	"strings"
)

// ParseBinaryArtifactMetadata reads an artifact's contents.
func ParseBinaryArtifactMetadata(rdr io.Reader) (*BinaryArtifact, error) {

	arr, err := ar.NewReader(rdr)
	if err != nil {
		return nil, err
	}

	art := &BinaryArtifact{}
	art.BinaryPackage = &BinaryPackage{}
	art.BinaryPackage.Package = &Package{}

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
					br := bufio.NewReader(tgzr)
					for {
						line, err := br.ReadString('\n')
						if err == io.EOF {
							break
						}
						if err != nil {
							return nil, err
						}
						if strings.Contains(line, ":") {
							res := strings.SplitN(line, ":", 2)
							log.Printf("Control File entry: '%s': %s", res[0], res[1])
							art.BinaryPackage.Package.SetField(res[0], res[1])
						}
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
	return art, nil
}

