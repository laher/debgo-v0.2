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
	"errors"
	"os"
)

// Build information ..
type BuildParams struct {
	// Properties below are mainly for build-related properties rather than metadata

	IsVerbose  bool   // Whether to log debug information
	TmpDir     string // Directory in-which to generate intermediate files & archives
	IsRmtemp   bool   // Delete tmp dir after execution?
	DestDir    string // Where to generate .deb files and source debs (.dsc files etc)
	WorkingDir string // This is the root from which to find .go files, templates, resources, etc

	TemplateDir string            // Optional. Only required if you're using templates
	Resources   map[string]string // Optional. Only if debgo packages your resources automatically. Key is the destination file. Value is the local file
}

func NewBuildParams() *BuildParams {
	pb := &BuildParams{IsVerbose: false}

	pb.TmpDir = TempDirDefault
	pb.IsRmtemp = true
	pb.DestDir = DistDirDefault
	pb.WorkingDir = WorkingDirDefault
	pb.TemplateDir = TemplateDirDefault
	pb.Resources = map[string]string{}
	return pb
}

//Initialise build process (make Temp and Dest directories)
func (pkg *BuildParams) Init() error {
	//make tmpDir
	if pkg.TmpDir == "" {
		return errors.New("Temp directory not specified")
	}
	err := os.MkdirAll(pkg.TmpDir, 0755)
	if err != nil {
		return err
	}
	//make destDir
	if pkg.DestDir == "" {
		return errors.New("Destination directory not specified")
	}
	err = os.MkdirAll(pkg.DestDir, 0755)
	if err != nil {
		return err
	}
	return err
}
