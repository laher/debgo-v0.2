package main

import (
	"github.com/laher/debgo-v0.2/cmd"
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/debgen"
	"log"
	"path/filepath"
	"os"
	"text/template"
)

func main() {
	name := "debgen-changelog"
	log.SetPrefix("[" + name + "] ")
	//set to empty strings because they're being overridden
	pkg := deb.NewPackage("", "", "", "")
	build := debgen.NewBuildParams()
	debgen.ApplyGoDefaults(pkg)
	fs := cmdutils.InitFlags(name, pkg, build)
	fs.StringVar(&pkg.Architecture, "arch", "all", "Architectures [any,386,armhf,amd64,all]")
	var entry string
	fs.StringVar(&entry, "entry", "", "Changelog entry data")

	err := cmdutils.ParseFlags(name, pkg, fs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if entry == "" {
		log.Fatalf("Error: --entry is a required flag")

	}
	filename := filepath.Join(build.ResourcesDir, "debian", "changelog")
	templateVars := debgen.NewTemplateData(pkg)
	templateVars.ChangelogEntry = entry
	err = os.MkdirAll(filepath.Join(build.ResourcesDir, "debian"), 0777)
	if err != nil {
		log.Fatalf("Error making dirs: %v", err)
	}

	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		tpl, err := template.New("template").Parse(debgen.TemplateChangelogInitial)
		if err != nil {
			log.Fatalf("Error parsing template: %v", err)
		}
		//create ..
		f, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Error creating file: %v", err)
		}
		defer f.Close()
		err = tpl.Execute(f, templateVars)
		if err != nil {
			log.Fatalf("Error executing template: %v", err)
		}
		err = f.Close()
		if err != nil {
			log.Fatalf("Error closing written file: %v", err)
		}
	} else if err != nil {
		log.Fatalf("Error reading existing changelog: %v", err)
	} else {
		tpl, err := template.New("template").Parse(debgen.TemplateChangelogAdditionalEntry)
		if err != nil {
			log.Fatalf("Error parsing template: %v", err)
		}
		//append..
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer f.Close()
		err = tpl.Execute(f, templateVars)
		if err != nil {
			log.Fatalf("Error executing template: %v", err)
		}
		err = f.Close()
		if err != nil {
			log.Fatalf("Error closing written file: %v", err)
		}
	}

}
