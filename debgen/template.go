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
	"bytes"
	"github.com/laher/debgo-v0.2/deb"
	"io/ioutil"
	"os"
	"text/template"
)

// initialize "template data" object
func NewTemplateData(pkg *deb.Package) *TemplateData {
	templateVars := TemplateData{Package: pkg, EntryDate: "", Checksums: nil}
	return &templateVars
}

//Data for templates
type TemplateData struct {
	Package        *deb.Package
	BinaryArtifact *deb.BinaryArtifact
	EntryDate      string
	Checksums      *deb.Checksums
}

func ProcessTemplateFileOrString(templateFile string, templateDefault string, vars interface{}) ([]byte, error) {
	_, err := os.Stat(templateFile)
	var tplText string
	if os.IsNotExist(err) {
		tplText = templateDefault
		return ProcessTemplateString(tplText, vars)
	} else if err != nil {
		return nil, err
	} else {
		return ProcessTemplateFile(tplText, vars)
	}
}

func ProcessTemplateFile(templateFile string, vars interface{}) ([]byte, error) {
	tplBytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}
	tplText := string(tplBytes)
	return ProcessTemplateString(tplText, vars)
}

func ProcessTemplateString(tplText string, vars interface{}) ([]byte, error) {
	tpl, err := template.New("template").Parse(tplText)
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
func newTemplateData(appName, appVersion, maintainer, maintainerEmail, version, arch, description, depends, buildDepends, priority, status, standardsVersion, section, format string, extraData map[string]interface{}, metadataDeb map[string]string) TemplateData {
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
		extraData,
		nil}
	return vars
}
*/
