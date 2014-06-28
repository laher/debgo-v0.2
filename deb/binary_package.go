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

type BinaryPackage struct {
	*Package
	ExecutablePaths map[string][]string
	DebFiles map[string]*DebFile
}

func NewBinaryPackage(pkg *Package, executablePaths map[string][]string) *BinaryPackage {
	return &BinaryPackage{Package: pkg, ExecutablePaths: executablePaths}
}

// Generates the control file content for a binary deb, for a given architecture
func (pkg *BinaryPackage) GenerateControlFileContentForArch(arch string) []byte {
	control := fmt.Sprintf("Package: %s\nPriority: Extra\n", pkg.Name)
	if pkg.Maintainer != "" {
		control = fmt.Sprintf("%sMaintainer: %s\n", control, pkg.Maintainer)
	}
	//mandatory
	control = fmt.Sprintf("%sVersion: %s\n", control, pkg.Version)

	control = fmt.Sprintf("%sArchitecture: %s\n", control, arch)

	// extra metadata ..
	for k, v := range pkg.Metadata {
		control = fmt.Sprintf("%s%s: %s\n", control, k, v)
	}
	control = fmt.Sprintf("%sDescription: %s\n", control, pkg.Description)
	return []byte(control)
}

/*
func getDebArch(destArch string, armArchName string) string {
	architecture := "all"
	switch destArch {
	case "386":
		architecture = "i386"
	case "arm":
		architecture = armArchName
	case "amd64":
		architecture = "amd64"
	}
	return architecture
}

func getArmArchName(settings *config.Settings) string {
	armArchName := settings.GetTaskSettingString(TASK_PKG_BUILD, "armarch")
	if armArchName == "" {
		//derive it from GOARM version:
		goArm := settings.GetTaskSettingString(TASK_XC, "GOARM")
		if goArm == "5" {
			armArchName = "armel"
		} else {
			armArchName = "armhf"
		}
	}
	return armArchName
}

func debBuild(dest platforms.Platform, tp TaskParams) (err error) {
	metadata := tp.Settings.GetTaskSettingMap(TASK_PKG_BUILD, "metadata")
	armArchName := getArmArchName(tp.Settings)
	metadataDeb := tp.Settings.GetTaskSettingMap(TASK_PKG_BUILD, "metadata-deb")
	rmtemp := tp.Settings.GetTaskSettingBool(TASK_PKG_BUILD, "rmtemp")
	debDir := filepath.Join(tp.OutDestRoot, tp.Settings.GetFullVersionName()) //v0.8.1 dont use platform dir
	tmpDir := filepath.Join(debDir, ".goxc-temp")
}
*/



//Builds debs for all arches. 
func (pkg *BinaryPackage) BuildAllWithDefaults() error {
	arches, err := pkg.GetArches()
	if err != nil {
		return err
	}

	for _, arch := range arches {
		err := pkg.BuildWithDefaults(arch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pkg *BinaryPackage) BuildWithDefaults(arch string) error {
	log.Printf("Building for arch %s", arch)
	//defer removal ...
	if pkg.IsRmtemp {
		defer os.RemoveAll(pkg.TmpDir)
	}
	//make tmpDir
	err := os.MkdirAll(pkg.TmpDir, 0755)
	if err != nil {
		return err
	}
	controlTgzw, err := pkg.InitControlArchive()
	if err != nil {
		return err
	}
	err = pkg.AddDefaultControlFile(arch, controlTgzw)
	if err != nil {
		return err
	}
	if pkg.IsVerbose {
		log.Printf("Wrote control file")
	}
	err = controlTgzw.Close()
	if err != nil {
		return err
	}
	if pkg.IsVerbose {
		log.Printf("Closed control archive")
	}

	dataTgzw, err := pkg.InitDataArchive()
	if err != nil {
		return err
	}
	err = pkg.AddExecutablesByFilepath(pkg.ExecutablePaths[arch], dataTgzw)
	if err != nil {
		return err
	}
	if pkg.IsVerbose {
		log.Printf("Added executables")
	}
	err = dataTgzw.Close()
	if err != nil {
		return err
	}
	if pkg.IsVerbose {
		log.Printf("Closed data archive")
	}
	err = pkg.BuildDebFile(arch, controlTgzw.Filename, dataTgzw.Filename)
	if err != nil {
		return err
	}
	if pkg.IsVerbose {
		log.Printf("Closed deb")
	}
	return err
}

func (pkg *BinaryPackage) AddDefaultControlFile(arch string, tgzw *TarGzWriter) error {
	controlContent := pkg.GenerateControlFileContentForArch(arch)
	if pkg.IsVerbose {
		log.Printf("Control file:\n%s", string(controlContent))
	}
	err := tgzw.AddBytes(controlContent, "control", 0644)
	if err != nil {
		return err
	}
	return err
}

func (pkg *BinaryPackage) InitControlArchive() (*TarGzWriter, error) {
	archiveFilename := filepath.Join(pkg.TmpDir, "control.tar.gz")
	tgzw, err := NewTarGzWriter(archiveFilename)
	if err != nil {
		return nil, err
	}
	return tgzw, err
}

func (pkg *BinaryPackage) InitDataArchive() (*TarGzWriter, error) {
	archiveFilename := filepath.Join(pkg.TmpDir, "data.tar.gz")
	tgzw, err := NewTarGzWriter(archiveFilename)
	if err != nil {
		return nil, err
	}
	return tgzw, err
}

func (pkg *BinaryPackage) AddExecutablesByFilepath(executablePaths []string, tgzw *TarGzWriter) error {
	for _, executable := range executablePaths {
		exeName := "/usr/bin/" + filepath.Base(executable)
		err := tgzw.AddFile(executable, exeName)
		if err != nil {
			tgzw.Close()
			return err
		}
	}
	return nil
	//TODO add resources to /usr/share/appName/
}

func (pkg *BinaryPackage) InitDebFile(arch, controlArchFile, dataArchFile string) {
	targetFile := filepath.Join(pkg.DestDir, fmt.Sprintf("%s_%s_%s.deb", pkg.Name, pkg.Version, arch)) //goxc_0.5.2_i386.deb")
	if pkg.DebFiles == nil {
		pkg.DebFiles = map[string]*DebFile{}
	}
	log.Printf("prepared for arch %s", arch)
	pkg.DebFiles[arch] = NewDebFile(targetFile, pkg.TmpDir)
	pkg.DebFiles[arch].ControlArchFile = controlArchFile
	pkg.DebFiles[arch].DataArchFile = dataArchFile

}

func (pkg *BinaryPackage) BuildDebFile(arch, controlArchFile, dataArchFile string) error {
	if pkg.DebFiles == nil || pkg.DebFiles[arch] == nil {
		pkg.InitDebFile(arch, controlArchFile, dataArchFile)
	}
	err := os.MkdirAll(pkg.DestDir, 0755)
	if err != nil {
		return err
	}
	//err = ioutil.WriteFile(filepath.Join(pkg.TmpDir, "debian-binary"), []byte("2.0\n"), 0644)

	log.Printf("trying to write .deb file for %s", arch)
	err = pkg.DebFiles[arch].WriteAll()
	return err
}
