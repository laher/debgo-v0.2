package deb_test

import (
	"github.com/laher/debgo-v0.2/deb"
	"log"
	"os"
	"path/filepath"
)

func Example_buildBinaryDeb() {

	pkg := deb.NewPackage("testpkg", "0.0.2", "me", "lovely package")
	pkg.Description = "hiya"
	build := deb.NewBuildParams()
	err := build.Init()
	if err != nil {
		log.Fatalf("%v", err)
	}

	exesMap := map[string][]string{
		"amd64": []string{"_test/a.amd64"},
		"i386":  []string{"_test/a.i386"},
		"armhf": []string{"_test/a.armhf"}}
	err = createExes(exesMap)
	if err != nil {
		log.Fatalf("%v", err)
	}
	bpkg := deb.NewBinaryPackage(pkg)
	artifacts, err := bpkg.GetArtifacts()
	if err != nil {
		log.Fatalf("Error building binary: %v", err)
	}
	artifacts[deb.Arch_amd64].Binaries = map[string]string{"/usr/bin/a": "_test/a.amd64"}
	artifacts[deb.Arch_i386].Binaries = map[string]string{"/usr/bin/a": "_test/a.i386"}
	artifacts[deb.Arch_armhf].Binaries = map[string]string{"/usr/bin/a": "_test/a.armhf"}
	buildBinaryArtifact := func(art *deb.BinaryArtifact, build *deb.BuildParams) error {
		//generate artifact here ...
		return nil
	}
	for arch, artifact := range artifacts {
		//build binary deb here ...
		err = buildBinaryArtifact(artifact, build)
		if err != nil {
			log.Fatalf("Error building for '%s': %v", arch, err)
		}
	}
}

func createExes(exesMap map[string][]string) error {
	for _, exes := range exesMap {
		for _, exe := range exes {
			err := os.MkdirAll(filepath.Dir(exe), 0777)
			if err != nil {
				return err
			}
			fi, err := os.Create(exe)
			if err != nil {
				return err
			}
			_, err = fi.Write([]byte("echo 1"))
			if err != nil {
				return err
			}
			err = fi.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
