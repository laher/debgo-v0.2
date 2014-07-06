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
	"log"
	"os"
	"path/filepath"
)

// BuildDevPackageFunc specifies a function which can build a DevPackage
type BuildBinaryArtifactFunc func(*BinaryPackage, *BinaryArtifact) error

// *BinaryPackage specifies functionality for building binary '.deb' packages.
// This encapsulates a Package plus information about platform-specific debs and executables
type BinaryPackage struct {
	*Package
	ExeDest         string
	BinaryArtifacts []*BinaryArtifact       //BinaryArtifact-specific builds
	BuildFunc       BuildBinaryArtifactFunc //Overrideable build function
}

// Factory for BinaryPackage
func NewBinaryPackage(pkg *Package) *BinaryPackage {
	return &BinaryPackage{Package: pkg, ExeDest: "/usr/bin", BuildFunc: BuildBinaryArtifactDefault}
}

// Builds debs for all arches.
func (pkg *BinaryPackage) Build() error {
	arches, err := pkg.GetArches()
	if err != nil {
		return err
	}
	err = pkg.Init()
	if err != nil {
		return err
	}
	for _, arch := range arches {
		archArtifact := pkg.GetBinaryArtifact(arch)
		if pkg.IsVerbose {
			log.Printf("Building for arch %s", arch)
		}
		err = pkg.BuildFunc(pkg, archArtifact)
		//even with an error, remove temp
		if pkg.IsRmtemp {
			os.RemoveAll(pkg.TmpDir)

		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Builds debs for a single architecture. Assumes default behaviours of a typical Go package.
// This allows for a limited amount of flexibility (e.g. the use of templates for metadata files).
// To get full flexibility, please use the more granular methods to return archives and add manual work
func (pkg *BinaryPackage) GetBinaryArtifact(arch Architecture) *BinaryArtifact {
	platform := pkg.getBinaryArtifact(arch)
	if platform == nil {
		platform = pkg.InitBinaryArtifact(arch)
	}
	return platform
}

// This is the default build process for a BuildArtifact
func BuildBinaryArtifactDefault(pkg *BinaryPackage, archArtifact *BinaryArtifact) error {
	_, err := pkg.BuildDefaultControlArchive(archArtifact)
	if err != nil {
		return err
	}
	_, err = pkg.BuildDefaultDataArchive(archArtifact)
	if err != nil {
		return err
	}
	err = pkg.BuildDebFile(archArtifact)
	if err != nil {
		return err
	}
	if pkg.IsVerbose {
		log.Printf("Closed deb")
	}
	return err
}

func (pkg *BinaryPackage) BuildDefaultControlArchive(archArtifact *BinaryArtifact) (string, error) {
	controlTgzw, err := pkg.InitControlArchive()
	if err != nil {
		return "", err
	}
	err = pkg.AddDefaultControlFile(archArtifact, controlTgzw)
	if err != nil {
		return "", err
	}
	if pkg.IsVerbose {
		log.Printf("Wrote control file to control archive")
	}
	// This is where you'd include Postrm/Postinst etc

	err = controlTgzw.Close()
	if err != nil {
		return "", err
	}
	if pkg.IsVerbose {
		log.Printf("Closed control archive")
	}
	return controlTgzw.Filename, err
}

func (pkg *BinaryPackage) BuildDefaultDataArchive(archArtifact *BinaryArtifact) (string, error) {
	dataTgzw, err := pkg.InitDataArchive()
	if err != nil {
		return "", err
	}
	err = pkg.AddExecutablesByFilepath(archArtifact.Executables, dataTgzw)
	if err != nil {
		return "", err
	}
	if pkg.IsVerbose {
		log.Printf("Added executables")
	}
	err = pkg.AddResources(pkg.Resources, dataTgzw)
	if err != nil {
		return "", err
	}
	if pkg.IsVerbose {
		log.Printf("Added resources")
	}
	err = dataTgzw.Close()
	if err != nil {
		return "", err
	}
	if pkg.IsVerbose {
		log.Printf("Closed data archive")
	}
	return dataTgzw.Filename, err
}

func (pkg *BinaryPackage) AddDefaultControlFile(archArtifact *BinaryArtifact, tgzw *TarGzWriter) error {
	templateVars := pkg.NewTemplateData()
	templateVars.Architecture = string(archArtifact.Architecture)
	controlData, err := ProcessTemplateFileOrString(filepath.Join(pkg.TemplateDir, "control.tpl"), TEMPLATE_BINARYDEB_CONTROL, templateVars)
	if err != nil {
		return err
	}
	if pkg.IsVerbose {
		log.Printf("Control file:\n%s", string(controlData))
	}
	err = tgzw.AddBytes(controlData, "control", 0644)
	if err != nil {
		return err
	}
	return err
}

// Initialise and return the 'control.tar.gz' archive
func (pkg *BinaryPackage) InitControlArchive() (*TarGzWriter, error) {
	archiveFilename := filepath.Join(pkg.TmpDir, "control.tar.gz")
	tgzw, err := NewTarGzWriter(archiveFilename)
	if err != nil {
		return nil, err
	}
	return tgzw, err
}

// Initialise and return the 'data.tar.gz' archive
func (pkg *BinaryPackage) InitDataArchive() (*TarGzWriter, error) {
	archiveFilename := filepath.Join(pkg.TmpDir, "data.tar.gz")
	tgzw, err := NewTarGzWriter(archiveFilename)
	if err != nil {
		return nil, err
	}
	return tgzw, err
}

// Add executables from file system.
// Be careful to make sure these are the relevant executables for the correct architecture
func (pkg *BinaryPackage) AddExecutablesByFilepath(executablePaths []string, tgzw *TarGzWriter) error {
	return pkg.AddFilesByFilepath(pkg.ExeDest, executablePaths, tgzw)
}

func (pkg *BinaryPackage) AddFilesByFilepath(destinationDir string, executablePaths []string, tgzw *TarGzWriter) error {
	if executablePaths != nil {
		for _, executable := range executablePaths {

			exeName := filepath.Join(destinationDir, filepath.Base(executable))
			err := tgzw.AddFile(executable, exeName)
			if err != nil {
				tgzw.Close()
				return err
			}
		}
	}
	return nil
}

// Add resources from file system.
// In this context, resources are simply files to include untouched to every architecture
// TODO add README.debian automatically
func (pkg *BinaryPackage) AddResources(resources map[string]string, tgzw *TarGzWriter) error {
	if resources != nil {
		for name, localPath := range resources {
			err := tgzw.AddFile(localPath, name)
			if err != nil {
				tgzw.Close()
				return err
			}
		}
	}
	return nil
}

//Initialise the generation parameters for a given architecture
func (pkg *BinaryPackage) InitBinaryArtifact(arch Architecture) *BinaryArtifact {
	targetFile := filepath.Join(pkg.DestDir, fmt.Sprintf("%s_%s_%s.deb", pkg.Name, pkg.Version, arch)) //goxc_0.5.2_i386.deb")
	if pkg.BinaryArtifacts == nil {
		pkg.BinaryArtifacts = []*BinaryArtifact{}
	}
	if pkg.IsVerbose {
		log.Printf("prepared for arch %s", arch)
	}
	platform := NewBinaryArtifact(arch, targetFile, pkg.TmpDir, pkg.IsVerbose)
	pkg.BinaryArtifacts = append(pkg.BinaryArtifacts, platform)
	return platform
}

// get the generation parameters for a given architecture
func (pkg *BinaryPackage) getBinaryArtifact(arch Architecture) *BinaryArtifact {
	if pkg.BinaryArtifacts == nil {
		pkg.BinaryArtifacts = []*BinaryArtifact{}
	}
	for _, platform := range pkg.BinaryArtifacts {
		if platform.Architecture == arch {
			return platform
		}
	}
	return nil
}

func (pkg *BinaryPackage) BuildDebFile(archArtifact *BinaryArtifact) error {
	if pkg.IsVerbose {
		log.Printf("trying to write .deb file for %s", archArtifact.Architecture)
	}
	err := archArtifact.WriteAll()
	return err
}
