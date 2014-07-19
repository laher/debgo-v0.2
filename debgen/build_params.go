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
	"errors"
	"github.com/laher/debgo-v0.2/deb"
	"os"
)

// BuildParams provides information about a particular build
type BuildParams struct {
	// Properties below are mainly for build-related properties rather than metadata

	IsVerbose  bool   // Whether to log debug information
	TmpDir     string // Directory in-which to generate intermediate files & archives
	IsRmtemp   bool   // Delete tmp dir after execution?
	DestDir    string // Where to generate .deb files and source debs (.dsc files etc)
	WorkingDir string // This is the root from which to find .go files, templates, resources, etc

	TemplateDir  string // Optional. Only required if you're using templates
	ResourcesDir string // Optional. Only if debgo packages your resources automatically.

	//TemplateStringsSource map[string]string //Populate this to fulfil templates for the different control files.
}

//Factory for BuildParams. Populates defaults.
func NewBuildParams() *BuildParams {
	bp := &BuildParams{IsVerbose: false}
	bp.TmpDir = deb.TempDirDefault
	bp.IsRmtemp = true
	bp.DestDir = deb.DistDirDefault
	bp.WorkingDir = deb.WorkingDirDefault
	bp.TemplateDir = deb.TemplateDirDefault
	bp.ResourcesDir = deb.ResourcesDirDefault
	return bp
}


//Initialise build directories (make Temp and Dest directories)
func (bp *BuildParams) Init() error {
	//make tmpDir
	if bp.TmpDir == "" {
		return errors.New("Temp directory not specified")
	}
	err := os.MkdirAll(bp.TmpDir, 0755)
	if err != nil {
		return err
	}
	//make destDir
	if bp.DestDir == "" {
		return errors.New("Destination directory not specified")
	}
	err = os.MkdirAll(bp.DestDir, 0755)
	if err != nil {
		return err
	}
	return err
}
