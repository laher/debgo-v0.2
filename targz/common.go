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
	"strings"
	"time"
)

// NewTarHeader is a factory for a tar header. Fixes slashes, populates ModTime
func NewTarHeader(path string, datalen int64, mode int64) *tar.Header {
	h := new(tar.Header)
	//slash-only paths
	h.Name = strings.Replace(path, "\\", "/", -1)
	h.Size = datalen
	h.Mode = mode
	h.ModTime = time.Now()
	return h
}


