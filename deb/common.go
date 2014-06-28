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
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"text/template"
	"time"
)

type Checksum struct {
	Checksum string
	Size     int64
	File     string
}

type Checksums struct {
	ChecksumsMd5     []Checksum
	ChecksumsSha1    []Checksum
	ChecksumsSha256  []Checksum
}

type TemplateData struct {
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
	AdditionalControlData    map[string]string
	Checksums	 *Checksums
}

func (cs *Checksums) Add(filepath, basename string) error {
	checksumMd5, checksumSha1, checksumSha256, err := checksums(filepath, basename)
	if err != nil {
		return err
	}
	cs.ChecksumsMd5 = append(cs.ChecksumsMd5, *checksumMd5)
	cs.ChecksumsSha1 = append(cs.ChecksumsSha1, *checksumSha1)
	cs.ChecksumsSha256 = append(cs.ChecksumsSha256, *checksumSha256)
	return nil
}

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
/*
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
*/
func ProcessTemplateFileOrString(templateFile string, templateDefault string, vars interface{}) ([]byte, error) {
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

func newTemplateData(appName, appVersion, maintainer, maintainerEmail, version, arch, description, depends, buildDepends, priority, status, standardsVersion, section, format string, metadataDeb map[string]string) TemplateData {
	vars := TemplateData{
		appName,
		appVersion,
		buildDepends,
		priority,
		maintainer,
		maintainerEmail,
		standardsVersion,
		arch,
		section,
		depends,
		description,
		"",
		status,
		time.Now().Format("Mon, 2 Jan 2006 15:04:05 -0700"),
		format,
		metadataDeb,
		nil}
	return vars
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
