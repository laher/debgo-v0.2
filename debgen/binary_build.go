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
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/targz"
	"log"
	"os"
	"path/filepath"
)

var ()

// This is the default build process for a BuildArtifact
func GenBinaryArtifact(archArtifact *deb.BinaryArtifact, build *deb.BuildParams) error {
	_, err := GenControlArchive(archArtifact, build)
	if err != nil {
		return err
	}
	_, err = GenDataArchive(archArtifact, build)
	if err != nil {
		return err
	}
	if build.IsVerbose {
		log.Printf("trying to write .deb file for %s", archArtifact.Architecture)
	}
	err = archArtifact.Build(build)
	if err != nil {
		return err
	}
	if build.IsVerbose {
		log.Printf("Closed deb")
	}
	return err
}

func GenControlArchive(archArtifact *deb.BinaryArtifact, build *deb.BuildParams) (string, error) {
	controlTgzw, err := archArtifact.InitControlArchive(build)
	if err != nil {
		return "", err
	}
	templateVars := &TemplateData{Package: archArtifact.BinaryPackage.Package, BinaryArtifact: archArtifact}
	//templateVars.BinaryArtifact = archArtifact

	err = GenControlFile(controlTgzw, templateVars, build)
	if err != nil {
		return "", err
	}
	if build.IsVerbose {
		log.Printf("Wrote control file to control archive")
	}
	// This is where you'd include Postrm/Postinst etc
	for _, scriptName := range deb.MaintainerScripts {
		templatePath := filepath.Join(build.TemplateDir, scriptName+".tpl")
		_, err = os.Stat(templatePath)
		//TODO handle non-EOF errors
		if err == nil {
			scriptData, err := ProcessTemplateFile(templatePath, templateVars)
			if err != nil {
				return "", err
			}
			err = controlTgzw.AddBytes(scriptData, scriptName, 0755)
			if err != nil {
				return "", err
			}
		}
	}

	err = controlTgzw.Close()
	if err != nil {
		return "", err
	}
	if build.IsVerbose {
		log.Printf("Closed control archive")
	}
	return controlTgzw.Filename, err
}

func GenDataArchive(archArtifact *deb.BinaryArtifact, build *deb.BuildParams) (string, error) {
	dataTgzw, err := archArtifact.InitDataArchive(build)
	if err != nil {
		return "", err
	}
	err = dataTgzw.AddFiles(archArtifact.Binaries)
	if err != nil {
		return "", err
	}
	if build.IsVerbose {
		log.Printf("Added executables")
	}
	// TODO add README.debian automatically
	err = dataTgzw.AddFiles(build.Resources)
	if err != nil {
		return "", err
	}
	if build.IsVerbose {
		log.Printf("Added resources")
	}
	err = dataTgzw.Close()
	if err != nil {
		return "", err
	}
	if build.IsVerbose {
		log.Printf("Closed data archive")
	}
	return dataTgzw.Filename, err
}

func GenControlFile(tgzw *targz.Writer, templateVars *TemplateData, build *deb.BuildParams) error {
	controlData, err := ProcessTemplateFileOrString(filepath.Join(build.TemplateDir, "control.tpl"), TemplateBinarydebControl, templateVars)
	if err != nil {
		return err
	}
	if build.IsVerbose {
		log.Printf("Control file:\n%s", string(controlData))
	}
	err = tgzw.AddBytes(controlData, "control", 0644)
	if err != nil {
		return err
	}
	return err
}
