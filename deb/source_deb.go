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
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	//"text/template"
)


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
			}
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
		if len(validGopaths) > 0 {
			gopath = validGopaths[0]

		} else {
			log.Printf("GOPATH env variable not valid! Using '.'")
			gopath = "."
		}
	}
	return gopath
}


// TODO: unfinished: need to discover root dir to determine which dirs to pre-make.
func SdebAddSources(codeDir, prefix string, tgzw *TarGzWriter) error {
	goPathRoot := getGoPathElement(codeDir)
	goPathRootResolved, err := filepath.EvalSymlinks(goPathRoot)
	if err != nil {
		log.Printf("Could not evaluate symlinks for '%s'", goPathRoot)
		goPathRootResolved = goPathRoot
	}
	log.Printf("Code dir '%s' (using goPath element '%s')", codeDir, goPathRootResolved)
	return sdebAddSources(goPathRootResolved, codeDir, prefix, tgzw)
}

// Get sources and append them
func sdebAddSources(goPathRoot, codeDir, prefix string, tgzw *TarGzWriter) error {
	//1. Glob for files in this dir
	//log.Printf("Globbing %s", codeDir)
	matches, err := filepath.Glob(filepath.Join(codeDir, "*.go"))
	if err != nil {
		return err
	}
	for _, match := range matches {
		absMatch, err := filepath.Abs(match)
		if err != nil {
			return fmt.Errorf("Error finding go sources (match %s): %v,", match, err)
		}
		relativeMatch, err := filepath.Rel(goPathRoot, absMatch)
		if err != nil {
			return fmt.Errorf("Error finding go sources (match %s): %v,", match, err)
		}
		destName := filepath.Join(prefix, relativeMatch)
		err = tgzw.AddFile(match, destName)
		if err != nil {
			return fmt.Errorf("Error adding go sources (match %s): %v,", match, err)
		}
		/*
		//log.Printf("Putting file %s in %s", match, destName)
		finf, err := os.Stat(destName)
		if err != nil {
			return fmt.Errorf("Error finding go sources (match %s): %v,", match, err)
		}
		tgzw.Tw.WriteHeader(NewTarHeader(destName, int64(finf.Size()), 0644))
		if err != nil {
			return err
		}
		_, err = tgzw.Tw.Write(controlContent)
		if err != nil {
			return err
		}

		sources = append(sources, archive.ArchiveItemFromFileSystem(match, destName))
		*/
	}

	//2. Recurse into subdirs
	fis, err := ioutil.ReadDir(codeDir)
	for _, fi := range fis {
		if fi.IsDir() && fi.Name() != DIRNAME_TEMP {
			err := sdebAddSources(goPathRoot, filepath.Join(codeDir, fi.Name()), prefix, tgzw)
			//sources = append(sources, additionalItems...)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func SdebCopySourceRecurse(codeDir, destDir string) (err error) {
	log.Printf("Globbing %s", codeDir)
	//get all files and copy into destDir
	matches, err := filepath.Glob(filepath.Join(codeDir, "*.go"))
	if err != nil {
		return err
	}
	if len(matches) > 0 {
		err = os.MkdirAll(destDir, 0777)
		if err != nil {
			return err
		}
	}
	for _, match := range matches {
		//TODO copy files
		log.Printf("copying %s into %s", match, filepath.Join(destDir, filepath.Base(match)))
		r, err := os.Open(match)
		if err != nil {
			return err
		}
		defer func() {
			err := r.Close()
			if err != nil {
				panic(err)
			}
		}()
		w, err := os.Create(filepath.Join(destDir, filepath.Base(match)))
		if err != nil {
			return err
		}
		defer func() {
			err := w.Close()
			if err != nil {
				panic(err)
			}
		}()

		_, err = io.Copy(w, r)
		if err != nil {
			return err
		}
	}
	fis, err := ioutil.ReadDir(codeDir)
	for _, fi := range fis {
		if fi.IsDir() && fi.Name() != DIRNAME_TEMP {
			err = SdebCopySourceRecurse(filepath.Join(codeDir, fi.Name()), filepath.Join(destDir, fi.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

/*
// prepare folders and debian/ files.
// (everything except copying source)
func SdebPrepare(workingDirectory, appName, maintainer, version, arches, description, buildDepends string, metadataDeb map[string]interface{}) (err error) {
	//make temp dir & subfolders
	tmpDir := filepath.Join(workingDirectory, DIRNAME_TEMP)
	debianDir := filepath.Join(tmpDir, "debian")
	err = os.MkdirAll(filepath.Join(debianDir, "source"), 0777)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(tmpDir, "src"), 0777)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(tmpDir, "bin"), 0777)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(tmpDir, "pkg"), 0777)
	if err != nil {
		return err
	}
	//write control file and related files
	tpl, err := template.New("rules").Parse(TEMPLATE_DEBIAN_RULES)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(debianDir, "rules"))
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()
	err = tpl.Execute(file, appName)
	if err != nil {
		return err
	}
	sdebControlFile := getSdebControlFileContent(appName, maintainer, version, arches, description, buildDepends, metadataDeb)
	ioutil.WriteFile(filepath.Join(debianDir, "control"), sdebControlFile, 0666)
	//copy source into folders
	//call dpkg-build, if available
	return err
}
*/
func getSdebControlFileContent(appName, maintainer, version, arches, description, buildDepends string, metadataDeb map[string]interface{}) []byte {
	control := fmt.Sprintf("Source: %s\nPriority: extra\n", appName)
	if maintainer != "" {
		control = fmt.Sprintf("%sMaintainer: %s\n", control, maintainer)
	}
	if buildDepends == "" {
		buildDepends = BUILD_DEPENDS_DEFAULT
	}
	control = fmt.Sprintf("%sBuildDepends: %s\n", control, buildDepends)
	control = fmt.Sprintf("%sStandards-Version: %s\n", control, STANDARDS_VERSION_DEFAULT)

	//TODO - homepage?

	control = fmt.Sprintf("%sVersion: %s\n", control, version)

	control = fmt.Sprintf("%sPackage: %s\n", control, appName)

	//mandatory
	control = fmt.Sprintf("%sArchitecture: %s\n", control, arches)
	for k, v := range metadataDeb {
		control = fmt.Sprintf("%s%s: %s\n", control, k, v)
	}
	control = fmt.Sprintf("%sDescription: %s\n", control, description)
	return []byte(control)
}
