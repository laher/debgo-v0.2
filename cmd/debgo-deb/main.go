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
	name := "debgo-deb"
	log.SetPrefix("[" + name + "] ")
	//set to empty strings because they're being overridden
	pkg := debgen.NewGoPackage("", "", "")
	build := deb.NewBuildParams()
	fs := cmdutils.InitFlags(name, pkg, build)
	bpkg := deb.NewBinaryPackage(pkg)

	var binDir string
	var resourcesDir string
	fs.StringVar(&binDir, "binaries", "", "directory containing binaries for each architecture. Directory names should end with the architecture")
	fs.StringVar(&pkg.Architecture, "arch", "any", "Architectures [any,386,armel,amd64,all]")
	fs.StringVar(&resourcesDir, "resources", "", "directory containing resources for this platform")
	err := cmdutils.ParseFlags(name, pkg, fs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = build.Init()
	if err != nil {
		log.Fatalf("%v", err)
	}

	build.Resources = map[string]string{}
	err = filepath.Walk(resourcesDir, func(path string, info os.FileInfo, err2 error) error {
		if info!=nil && !info.IsDir() {
			rel, err := filepath.Rel(resourcesDir,path)
			if err == nil {
				build.Resources[rel] = path
			}
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("Resources: %v", build.Resources)
	// TODO determine this platform
	//err = bpkg.Build(build, debgen.GenBinaryArtifact)
	artifacts, err := bpkg.GetArtifacts()
	if err != nil {
		log.Fatalf("%v", err)
	}
	for arch, artifact := range artifacts {
		artifact.Binaries = map[string]string{}
		archBinDir := filepath.Join(binDir, string(arch))
		err = filepath.Walk(archBinDir, func(path string, info os.FileInfo, err2 error) error {
			if info!=nil && !info.IsDir() {
				rel, err := filepath.Rel(binDir,path)
				if err == nil {
					artifact.Binaries[rel] = path
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
