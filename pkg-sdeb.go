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

/* Nearing completion */
//TODO: various refinements ...
import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	//Tip for Forkers: please 'clone' from my url and then 'pull' from your url. That way you wont need to change the import path.
	//see https://groups.google.com/forum/?fromgroups=#!starred/golang-nuts/CY7o2aVNGZY
	"github.com/laher/goxc/archive"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	//"strings"
	"text/template"
	"time"
)

type Checksum struct {
	Checksum string
	Size     int64
	File     string
}
type TemplateVars struct {
	PackageName      string
	PackageVersion   string
	BuildDepends     string
	Priority         string
	Maintainer       string
	MaintainerEmail  string
	StandardsVersion string
	Architecture     string
	Section          string
	Depends          string
	Description      string
	Other            string
	Status           string
	EntryDate        string
	Format           string
	Files            []Checksum
	ChecksumsSha1    []Checksum
	ChecksumsSha256  []Checksum
}
/*
//runs automatically
func init() {
	Register(Task{
		TASK_PKG_SOURCE,
		"Build a source package. Currently only supports 'source deb' format for Debian/Ubuntu Linux.",
		runTaskPkgSource,
		map[string]interface{}{"metadata": map[string]interface{}{"maintainer": "unknown", "maintainerEmail": "unknown@example.com"}, "metadata-deb": map[string]interface{}{"Depends": "", "Build-Depends": "debhelper (>=4.0.0), golang-go, gcc"}, "templateDir": "debian-templates", "rmtemp": true}})
}
*/
func checksums(path, name string) (*Checksum, *Checksum, *Checksum, error) {
	//checksums
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}

	hashMd5 := md5.New()
	size, err := io.Copy(hashMd5, f)
	if err != nil {
		return nil, nil, nil, err
	}
	checksumMd5 := Checksum{hex.EncodeToString(hashMd5.Sum(nil)), size, name}

	f.Seek(int64(0), 0)
	hash256 := sha256.New()
	size, err = io.Copy(hash256, f)
	if err != nil {
		return nil, nil, nil, err
	}
	checksumSha256 := Checksum{hex.EncodeToString(hash256.Sum(nil)), size, name}

	f.Seek(int64(0), 0)
	hash1 := sha1.New()
	size, err = io.Copy(hash1, f)
	if err != nil {
		return nil, nil, nil, err
	}
	checksumSha1 := Checksum{hex.EncodeToString(hash1.Sum(nil)), size, name}

	err = f.Close()
	if err != nil {
		return nil, nil, nil, err
	}

	return &checksumMd5, &checksumSha1, &checksumSha256, nil

}

func (pkg *DebPackage) SourceBuild() error {
	//default to just package <packagename>
	arches := "all"
	items, err := SdebGetSourcesAsArchiveItems(pkg.WorkingDir, pkg.Name+"-"+pkg.Version)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	//build
	//1. generate orig.tar.gz
	err = os.MkdirAll(pkg.DestDir, 0777)
	if err != nil {
		return err
	}

	//set up template
	templateVars := getTemplateVars(pkg.Name, pkg.Version, pkg.Maintainer, pkg.MaintainerEmail, pkg.Version, arches, pkg.Description, pkg.Metadata)

	//TODO add/exclude resources to /usr/share
	origTgzName := pkg.Name + "_" + pkg.Version + ".orig.tar.gz"
	origTgzPath := filepath.Join(pkg.DestDir, origTgzName)
	err = archive.TarGz(origTgzPath, items)
	if err != nil {
		return err
	}
	log.Printf("Created %s", origTgzPath)

	checksumMd5, checksumSha1, checksumSha256, err := checksums(origTgzPath, origTgzName)
	if err != nil {
		return err
	}
	templateVars.Files = append(templateVars.Files, *checksumMd5)
	templateVars.ChecksumsSha1 = append(templateVars.ChecksumsSha1, *checksumSha1)
	templateVars.ChecksumsSha256 = append(templateVars.ChecksumsSha256, *checksumSha256)

	//2. generate .debian.tar.gz (just containing debian/ directory)

	//debian/control
	controlData, err := getDebMetadataFileContent(filepath.Join(pkg.TemplateDir, "control.tpl"), TEMPLATE_SOURCEDEB_CONTROL, templateVars)
	if err != nil {
		return err
	}

	//compat
	compatData, err := getDebMetadataFileContent(filepath.Join(pkg.TemplateDir, "compat.tpl"), TEMPLATE_DEBIAN_COMPAT, templateVars)
	if err != nil {
		return err
	}

	//debian/rules
	rulesData, err := getDebMetadataFileContent(filepath.Join(pkg.TemplateDir, "rules.tpl"), TEMPLATE_DEBIAN_RULES, templateVars)
	if err != nil {
		return err
	}

	//debian/source/format
	sourceFormatData, err := getDebMetadataFileContent(filepath.Join(pkg.TemplateDir, "source_format.tpl"), TEMPLATE_DEBIAN_SOURCE_FORMAT, templateVars)
	if err != nil {
		return err
	}

	//debian/source/options
	sourceOptionsData, err := getDebMetadataFileContent(filepath.Join(pkg.TemplateDir, "source_options.tpl"), TEMPLATE_DEBIAN_SOURCE_OPTIONS, templateVars)
	if err != nil {
		return err
	}

	//debian/rules
	copyrightData, err := getDebMetadataFileContent(filepath.Join(pkg.TemplateDir, "copyright.tpl"), TEMPLATE_DEBIAN_COPYRIGHT, templateVars)
	if err != nil {
		return err
	}

	//debian/changelog (slightly different)
	var changelogData []byte
	rdr, err := pkg.Changelog.GetReader()
	if err != nil {
		return err
	}
	changelogData, err = ioutil.ReadAll(rdr)
	if err != nil {
		return err
	}
	/*
	_, err = os.Stat(changelogFilename)
	if os.IsNotExist(err) {
		initialChangelogTemplate := TEMPLATE_CHANGELOG_HEADER + "\n\n" + TEMPLATE_CHANGELOG_INITIAL_ENTRY + "\n\n" + TEMPLATE_CHANGELOG_FOOTER
		changelogData, err = getDebMetadataFileContent(filepath.Join(pkg.TemplateDir, "initial-changelog.tpl"), initialChangelogTemplate, templateVars)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		changelogData, err = ioutil.ReadFile(changelogFilename)
		if err != nil {
			return err
		}
	}
	*/
	if len(changelogData) == 0 {
		initialChangelogTemplate := TEMPLATE_CHANGELOG_HEADER + "\n\n" + TEMPLATE_CHANGELOG_INITIAL_ENTRY + "\n\n" + TEMPLATE_CHANGELOG_FOOTER
		changelogData, err = getDebMetadataFileContent(filepath.Join(pkg.TemplateDir, "initial-changelog.tpl"), initialChangelogTemplate, templateVars)
		if err != nil {
			return err
		}
	}
	//generate debian/README.Debian
	//TODO: try pulling in README.md etc
	//debian/README.Debian
	readmeData, err := getDebMetadataFileContent(filepath.Join(pkg.TemplateDir, "readme.tpl"), TEMPLATE_DEBIAN_README, templateVars)
	if err != nil {
		return err
	}
	debTgzName := pkg.Name + "_" + pkg.Version + ".debian.tar.gz"
	debTgzPath := filepath.Join(pkg.DestDir, debTgzName)
	err = archive.TarGz(debTgzPath,
		[]archive.ArchiveItem{
			archive.ArchiveItemFromBytes(changelogData, "debian/changelog"),
			archive.ArchiveItemFromBytes(copyrightData, "debian/copyright"),
			archive.ArchiveItemFromBytes(compatData, "debian/compat"),
			archive.ArchiveItemFromBytes(controlData, "debian/control"),
			archive.ArchiveItemFromBytes(readmeData, "debian/README.Debian"),
			archive.ArchiveItemFromBytes(rulesData, "debian/rules"),
			archive.ArchiveItemFromBytes(sourceFormatData, "debian/source/format"),
			archive.ArchiveItemFromBytes(sourceOptionsData, "debian/source/options")})
	if err != nil {
		return err
	}
	log.Printf("Created %s", debTgzPath)

	checksumMd5, checksumSha1, checksumSha256, err = checksums(debTgzPath, debTgzName)
	if err != nil {
		return err
	}
	templateVars.Files = append(templateVars.Files, *checksumMd5)
	templateVars.ChecksumsSha1 = append(templateVars.ChecksumsSha1, *checksumSha1)
	templateVars.ChecksumsSha256 = append(templateVars.ChecksumsSha256, *checksumSha256)

	dscData, err := getDebMetadataFileContent(filepath.Join(pkg.TemplateDir, "dsc.tpl"), TEMPLATE_DEBIAN_DSC, templateVars)
	if err != nil {
		return err
	}
	//3. generate .dsc file
	dscPath := filepath.Join(pkg.DestDir, pkg.Name+"_"+pkg.Version+".dsc")
	err = ioutil.WriteFile(dscPath, dscData, 0644)
	if err == nil {
		log.Printf("Wrote %s", dscPath)
	}
	return err
}

func getTemplateVars(appName, appVersion, maintainer, maintainerEmail, version, arch, description string, metadataDeb map[string]interface{}) TemplateVars {
	vars := TemplateVars{
		appName,
		appVersion,
		BUILD_DEPENDS_DEFAULT,
		PRIORITY_DEFAULT,
		maintainer,
		maintainerEmail,
		STANDARDS_VERSION_DEFAULT,
		ARCHITECTURE_DEFAULT,
		SECTION_DEFAULT,
		"",
		description,
		"",
		"unreleased",
		time.Now().Format("Mon, 2 Jan 2006 15:04:05 -0700"),
		FORMAT_DEFAULT,
		[]Checksum{},
		[]Checksum{},
		[]Checksum{},
	}
	return vars
}

func getDebMetadataFileContent(templateFile string, templateDefault string, vars interface{}) ([]byte, error) {
	_, err := os.Stat(templateFile)
	var tplText string
	if os.IsNotExist(err) {
		tplText = templateDefault

	} else if err != nil {
		return nil, err
	} else {
		tplBytes, err := ioutil.ReadFile(templateFile)
		if err != nil {
			return nil, err
		}
		tplText = string(tplBytes)
	}
	tpl, err := template.New(templateFile).Parse(tplText)
	if err != nil {
		return nil, err
	}
	var dest bytes.Buffer

	err = tpl.Execute(&dest, vars)
	if err != nil {
		return nil, err
	}
	return dest.Bytes(), nil
}

/*
func getSourceDebControlFileContent(appName, maintainer, version, arch, description string, metadataDeb map[string]interface{}) []byte {
	control := fmt.Sprintf("Source: %s\nPriority: optional\n", appName)
	if maintainer != "" {
		control = fmt.Sprintf("%sMaintainer: %s\n", control, maintainer)
	}
	//mandatory
	control = fmt.Sprintf("%sStandards-Version: %s\n", control, version)

	control = fmt.Sprintf("%s\nPackage: %s\nArchitecture: any\n", control, appName)
	control = fmt.Sprintf("%sArchitecture: %s\n", control, arch)
	//must include Depends and Build-Depends
	for k, v := range metadataDeb {
		control = fmt.Sprintf("%s%s: %s\n", control, k, v)
	}
	control = fmt.Sprintf("%sDescription: %s\n", control, description)
	return []byte(control)
}
*/
