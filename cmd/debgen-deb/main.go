package main

import (
	"github.com/laher/debgo-v0.2/cmd"
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/debgen"
	"log"
	"os"
	"path/filepath"
)

func main() {
	name := "debgen-deb"
	log.SetPrefix("[" + name + "] ")
	//set to empty strings because they're being overridden
	pkg := debgen.NewGoPackage("", "", "", "")
	build := deb.NewBuildParams()
	fs := cmdutils.InitFlags(name, pkg, build)

	var binDir string
	var resourcesDir string
	fs.StringVar(&binDir, "binaries", "", "directory containing binaries for each architecture. Directory names should end with the architecture")
	fs.StringVar(&pkg.Architecture, "arch", "any", "Architectures [any,386,armhf,amd64,all]")
	fs.StringVar(&resourcesDir, "resources", "", "directory containing resources for this platform")
	err := cmdutils.ParseFlags(name, pkg, fs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = build.Init()
	if err != nil {
		log.Fatalf("%v", err)
	}

	pkg.MappedFiles = map[string]string{}
	err = filepath.Walk(resourcesDir, func(path string, info os.FileInfo, err2 error) error {
		if info != nil && !info.IsDir() {
			rel, err := filepath.Rel(resourcesDir, path)
			if err == nil {
				pkg.MappedFiles[rel] = path
			}
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	//log.Printf("Resources: %v", build.Resources)
	// TODO determine this platform
	//err = bpkg.Build(build, debgen.GenBinaryArtifact)
	artifacts, err := deb.GetArtifacts(pkg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	for arch, artifact := range artifacts {
		if artifact.MappedFiles == nil {
			artifact.MappedFiles = map[string]string{}
		}
		archBinDir := filepath.Join(binDir, string(arch))
		err = filepath.Walk(archBinDir, func(path string, info os.FileInfo, err2 error) error {
			if info != nil && !info.IsDir() {
				rel, err := filepath.Rel(binDir, path)
				if err == nil {
					artifact.MappedFiles[rel] = path
				}
				return err
			}
			return nil
		})

		err = debgen.GenBinaryArtifact(artifact, build)
		if err != nil {
			log.Fatalf("Error building for '%s': %v", arch, err)
		}
	}
}
