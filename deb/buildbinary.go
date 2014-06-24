package deb

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

import (
	"fmt"
	"archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"compress/gzip"
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
		return []string{ "i386", "armel", "amd64" }, nil
	}
	return []string{arches}, nil
}

func (pkg *DebPackage) BuildAll() error {
	arches, err := resolveArches(pkg.Architecture)
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
	if pkg.IsRmtemp {
		defer os.RemoveAll(pkg.TmpDir)
	}
	err := os.MkdirAll(pkg.TmpDir, 0755)
	if err != nil {
		return err
	}
	err = pkg.BuildControlArchive(arch)
	if err != nil {
		return err
	}
	err = pkg.BuildDataArchive(arch)
	if err != nil {
		return err
	}
	err = pkg.BuildDebFile(arch)
	if err != nil {
		return err
	}
	return err
}

func newTarHeader(path string, datalen int64, mode int64) *tar.Header {
	h := new(tar.Header)
	//backslash-only paths
	h.Name = strings.Replace(path, "\\", "/", -1)
	h.Size = datalen
	h.Mode = mode
	h.ModTime = time.Now()
	return h
}

func (pkg *DebPackage) BuildControlArchive(arch string) error {
	archiveFilename := filepath.Join(pkg.TmpDir, "control.tar.gz")
	fw, err := os.Create(archiveFilename)
	if err != nil {
		return err
	}
	defer fw.Close()

	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	controlContent := pkg.getControlFileContent(arch)
	if pkg.IsVerbose {
		log.Printf("Control file:\n%s", string(controlContent))
	}
	err = tw.WriteHeader(newTarHeader("control", int64(len(controlContent)), 0644))
	if err != nil {
		return err
	}
	_, err = tw.Write(controlContent)
	if err != nil {
		return err
	}
	if pkg.Postinst != nil {
		pi, err := toBytes(pkg.Postinst)
		if err != nil {
			return err
		}
		err = tw.WriteHeader(newTarHeader("postinst", int64(len(pi)), 0755))
		if err != nil {
			return err
		}
		_, err = tw.Write(pi)
		if err != nil {
			return err
		}
	}
	/*
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
	return tw.Close()
}

func (pkg *DebPackage) BuildDataArchive(arch string) error {
	archiveFilename := filepath.Join(pkg.TmpDir, "data.tar.gz")
	fw, err := os.Create(archiveFilename)
	if err != nil {
		return err
	}
	defer fw.Close()

	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, executable := range pkg.ExecutablePaths {
		exeName := "/usr/bin/" + filepath.Base(executable)
		fi, err := os.Open(executable)
		if err != nil {
			return err
		}
		finf, err := fi.Stat()
		if err != nil {
			return err
		}
		err = tw.WriteHeader(newTarHeader(exeName, finf.Size(), int64(finf.Mode())))
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, fi)
		if err != nil {
			return err
		}

	}
	//TODO add resources to /usr/share/appName/
	//err := archive.TarGz(filepath.Join(pkg.TmpDir, "data.tar.gz"), items)
	return tw.Close()
}

func (pkg *DebPackage) BuildDebFile(arch string) error {
	err := os.MkdirAll(pkg.DestDir, 0755)
	if err != nil {
		return err
	}
	//err = ioutil.WriteFile(filepath.Join(pkg.TmpDir, "debian-binary"), []byte("2.0\n"), 0644)

	targetFile := filepath.Join(pkg.DestDir, fmt.Sprintf("%s_%s_%s.deb", pkg.Name, pkg.Version, arch)) //goxc_0.5.2_i386.deb")

	log.Printf("prepared for arch %s", arch)
	bdeb := NewBinaryDeb(targetFile, pkg.TmpDir)
	bdeb.SetDefaults()
	log.Printf("trying to write .deb file %s", arch)
	err = bdeb.WriteAll()
	return err
}
