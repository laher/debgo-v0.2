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
	"bufio"
	"io"
	"log"
	"strings"
)

// DscReader reads a control file.
type DscReader struct {
	Reader io.Reader
}

//NewDscReader is a factory for reading Dsc files.
func NewDscReader(rdr io.Reader) (*DscReader) {
	return &DscReader{rdr}
}

// Parse parses a file into a package.
func (dscr *DscReader) Parse() (*Package, error) {
	pkg := &Package{}
	br := bufio.NewReader(dscr.Reader)
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
			pkg.SetField(res[0], res[1])
		} else {

		}
	}
	return pkg, nil

}

