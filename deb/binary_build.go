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

func (pkg *DebPackage) getControlFileContent(arch string) []byte {
	control := fmt.Sprintf("Package: %s\nPriority: Extra\n", pkg.Name)
	if pkg.Maintainer != "" {
		control = fmt.Sprintf("%sMaintainer: %s\n", control, pkg.Maintainer)
	}
	//mandatory
	control = fmt.Sprintf("%sVersion: %s\n", control, pkg.Version)

	control = fmt.Sprintf("%sArchitecture: %s\n", control, arch)
	for k, v := range pkg.Metadata {
		control = fmt.Sprintf("%s%s: %s\n", control, k, v)
	}
	control = fmt.Sprintf("%sDescription: %s\n", control, pkg.Description)
	return []byte(control)
}

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

/*
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

func resolveArches(arches string) ([]string, error) {
	if arches == "any" || arches == "" {
		return []string{"i386", "armel", "amd64"}, nil
	}
	return []string{arches}, nil
}
func (pkg *DebPackage) GetArches() ([]string, error) {
	arches, err := resolveArches(pkg.Architecture)
	return arches, err
}

func (pkg *DebPackage) DefaultBuildAllArches() error {
	arches, err := pkg.GetArches()
	if err != nil {
		return err
	}

	for _, arch := range arches {
		err := pkg.Build(arch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pkg *DebPackage) Build(arch string) error {
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

	controlTgzw, err := pkg.InitControlArchive(arch, true)
	if err != nil {
		return err
	}
	if pkg.IsVerbose {
		log.Printf("Wrote control archive")
	}

	err = controlTgzw.Close()
	if err != nil {
		return err
	}
	if pkg.IsVerbose {
		log.Printf("Closed control archive")
	}

	dataTgzw, err := pkg.InitDataArchive(arch, true)
	if err != nil {
		return err
	}
	err = dataTgzw.Close()
	if err != nil {
		return err
	}
	err = pkg.BuildDebFile(arch, controlTgzw.Filename, dataTgzw.Filename)
	if err != nil {
		return err
	}
	return err
}

func (pkg *DebPackage) InitControlArchive(arch string, generateControlFile bool) (*TarGzWriter, error) {
	archiveFilename := filepath.Join(pkg.TmpDir, "control.tar.gz")
	tgzw, err := NewTarGzWriter(archiveFilename)
	if err != nil {
		return nil, err
	}

	if generateControlFile {
		controlContent := pkg.getControlFileContent(arch)
		if pkg.IsVerbose {
			log.Printf("Control file:\n%s", string(controlContent))
		}
		err = tgzw.AddBytes(controlContent, "control", 0644)
		if err != nil {
			return nil, err
		}
	}

	/*
		//err = tgzw.Tw.WriteHeader(NewTarHeader("control", int64(len(controlContent)), 0644))
	_, err = tgzw.Tw.Write(controlContent)
	if err != nil {
		return nil, err
	}
	
	if pkg.Postinst != nil {
		pi, err := toBytes(pkg.Postinst)
		if err != nil {
			return nil, err
		}
		//err = tgzw.Tw.WriteHeader(NewTarHeader("postinst", int64(len(pi)), 0755))
		err = tgzw.AddBytes(pi, "postinst", 0644)
		if err != nil {
			return nil, err
		}
	}
		err = ioutil.WriteFile(filepath.Join(pkg.TmpDir, "control"), controlContent, 0644)
		if err != nil {
			return err
		}
		controlFiles := []archive.ArchiveItem{archive.ArchiveItem{FileSystemPath: filepath.Join(pkg.TmpDir, "control"), ArchivePath: "control"}}
		barr, err := toBytes(pkg.Postinst)
		if err != nil {
			return err
		}
		if barr != nil {
			controlFiles = append(controlFiles, archive.ArchiveItem{Data: barr, ArchivePath: "postinst"})
		}

		barr2, err := toBytes(pkg.Preinst)
		if err != nil {
			return err
		}
		if barr2 != nil {
			controlFiles = append(controlFiles, archive.ArchiveItem{Data: barr2, ArchivePath: "preinst"})
		}

		barr3, err := toBytes(pkg.Postrm)
		if err != nil {
			return err
		}
		if barr3 != nil {
			controlFiles = append(controlFiles, archive.ArchiveItem{Data: barr3, ArchivePath: "postrm"})
		}

		barr4, err := toBytes(pkg.Prerm)
		if err != nil {
			return err
		}
		if barr4 != nil {
			controlFiles = append(controlFiles, archive.ArchiveItem{Data: barr4, ArchivePath: "prerm"})
		}

		err = archive.TarGz(filepath.Join(pkg.TmpDir, "control.tar.gz"), controlFiles)
		if err != nil {
			return err
		}
	*/
	return tgzw, err
}

func (pkg *DebPackage) InitDataArchive(arch string, addExecutables bool) (*TarGzWriter, error) {
	archiveFilename := filepath.Join(pkg.TmpDir, "data.tar.gz")
	tgzw, err := NewTarGzWriter(archiveFilename)
	if err != nil {
		return nil, err
	}
	if addExecutables {
		for _, executable := range pkg.ExecutablePaths {
			exeName := "/usr/bin/" + filepath.Base(executable)
			err = tgzw.AddFile(executable, exeName)
			if err != nil {
				tgzw.Close()
				return nil, err
			}
		}
	}
	//TODO add resources to /usr/share/appName/
	return tgzw, err
}

func (pkg *DebPackage) BuildDebFile(arch string, controlArchFile, dataArchFile string) error {
	err := os.MkdirAll(pkg.DestDir, 0755)
	if err != nil {
		return err
	}
	//err = ioutil.WriteFile(filepath.Join(pkg.TmpDir, "debian-binary"), []byte("2.0\n"), 0644)

	targetFile := filepath.Join(pkg.DestDir, fmt.Sprintf("%s_%s_%s.deb", pkg.Name, pkg.Version, arch)) //goxc_0.5.2_i386.deb")

	log.Printf("prepared for arch %s", arch)
	bdeb := NewBinaryDeb(targetFile, pkg.TmpDir)
	bdeb.ControlArchFile = controlArchFile
	bdeb.DataArchFile = dataArchFile
	log.Printf("trying to write .deb file %s", arch)
	err = bdeb.WriteAll()
	return err
}
