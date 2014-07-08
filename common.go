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

package debgo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func absAndResolveSymlinks(path string) (string, error) {
	pathResolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	pathAbs, err := filepath.Abs(pathResolved)
	if err != nil {
		return "", err
	}
	return pathAbs, nil
}

// Tries to find the most relevant GOPATH element.
// First, tries to find an element which is a parent of the current directory.
// If not, it uses the first one.
func getGoPathElement(workingDirectory string) string {
	var gopath string
	gopathVar := os.Getenv("GOPATH")
	if gopathVar == "" {
		log.Printf("GOPATH env variable not set! Using '.'")
		gopath = "."
	} else {
		gopaths := filepath.SplitList(gopathVar)
		validGopaths := []string{}
		workingDirectoryAbs, err := filepath.Abs(workingDirectory)
		if err != nil {
			//strange. TODO: investigate
			workingDirectoryAbs = workingDirectory
		}
		//see if you can match the workingDirectory
		for _, gopathi := range gopaths {
			//if empty or GOROOT, continue
			//logic taken from http://tip.golang.org/src/pkg/go/build/build.go
			if gopathi == "" || gopathi == runtime.GOROOT() || strings.HasPrefix(gopathi, "~") {
				continue
			} else {
				validGopaths = append(validGopaths, gopathi)
				gopathAbs, err := filepath.Abs(gopathi)
				if err != nil {
					//strange. TODO: investigate
					gopathAbs = gopathi
				}
				//working directory is inside this path element. Use it!
				if strings.HasPrefix(workingDirectoryAbs, gopathAbs) {
					return gopathi
				}
			}
		}
		if len(validGopaths) > 0 {
			gopath = validGopaths[0]

		} else {
			log.Printf("GOPATH env variable not valid! Using '.'")
			gopath = "."
		}
	}
	return gopath
}

// Glob for Go sources.
// This function looks for source files, and prepares their paths for
func globForSources(goPathRoot, codeDir, destinationPrefix string, ignoreFiles []string) (map[string]string, error) {
	goPathRootAbs, err := absAndResolveSymlinks(goPathRoot)
	if err != nil {
		return nil, err
	}
	//	log.Printf("Globbing %s into %s relative to %s", codeDir, destinationPrefix, goPathRootAbs)

	//1. Glob for files in this dir
	//log.Printf("Globbing %s", codeDir)
	matches, err := filepath.Glob(filepath.Join(codeDir, "*.go"))
	if err != nil {
		return nil, err
	}
	sources := map[string]string{}
	for _, match := range matches {
		ignore := false
		for _, ignoreFile := range ignoreFiles {
			if match == ignoreFile {
				ignore = true
			}
		}
		if !ignore {
			absMatch, err := absAndResolveSymlinks(match)
			if err != nil {
				return nil, fmt.Errorf("Error finding go sources (match %s): %v,", match, err)
			}
			relativeMatch, err := filepath.Rel(goPathRootAbs, absMatch)
			if err != nil {
				return nil, fmt.Errorf("Error finding go sources (match %s): %v,", match, err)
			}
			//log.Printf("Adding %s into %s / %s relative to %s", match, destinationPrefix, relativeMatch, goPathRootAbs)

			destName := filepath.Join(destinationPrefix, relativeMatch)
			sources[destName] = match
		}
	}

	//2. Recurse into subdirs
	fis, err := ioutil.ReadDir(codeDir)
	for _, fi := range fis {
		if fi.IsDir() {
			ignore := false
			for _, ignoreFile := range ignoreFiles {
				if fi.Name() == ignoreFile {
					ignore = true
				}
			}
			if !ignore {
				moreSources, err := globForSources(goPathRoot, filepath.Join(codeDir, fi.Name()), destinationPrefix, ignoreFiles)
				if err != nil {
					return nil, err
				}
				for k, v := range moreSources {
					sources[k] = v
				}
			}
		}
	}
	return sources, err

}
