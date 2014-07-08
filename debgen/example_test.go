package debgen_test

import (
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/debgen"
	"log"
	"os"
	"path/filepath"
)

func Example_binaryPackage() {

	pkg := deb.NewPackage("testpkg", "0.0.2", "me")
	pkg.Description = "Dummy package for doing nothing"

	build := deb.NewBuildParams()
	build.IsRmtemp = false

	bpkg := deb.NewBinaryPackage(pkg)
	artifacts, err := bpkg.GetArtifacts()
	if err != nil {
		log.Fatalf("Error building binary: %v", err)
	}
	platform64 := artifacts[deb.Arch_amd64]
	platform64.Binaries = map[string]string{"/usr/bin/a": "_test/a.amd64"}
	artifacts[deb.Arch_i386].Binaries = map[string]string{"/usr/bin/a": "_test/a.x86"}
	artifacts[deb.Arch_armel].Binaries = map[string]string{"/usr/bin/a": "_test/a.armel"}

	prep() //prepare files for packaging using some other means.
	err = bpkg.Build(build, debgen.GenBinaryArtifact)
	if err != nil {
		log.Fatalf("Error building binary: %v", err)
	}

	// Output:
	//
}

func Example_devPackage() {
	pkg := deb.NewPackage("testpkg", "0.0.2", "me")
	pkg.Description = "hiya"

	ddpkg := deb.NewDevPackage(pkg)
	build := deb.NewBuildParams()
	build.IsRmtemp = false
	var err error
	build.Resources, err = debgen.GlobForGoSources(".", []string{build.TmpDir, build.DestDir})
	if err != nil {
		log.Fatalf("Error building -dev: %v", err)
	}

	err = debgen.GenDevArtifact(ddpkg, build)
	if err != nil {
		log.Fatalf("Error building -dev: %v", err)
	}

	// Output:
	//
}

func prep() error {
	exesMap := map[string][]string{
		"amd64": []string{"_test/a.amd64"},
		"i386":  []string{"_test/a.i386"},
		"armel": []string{"_test/a.armel"}}
	err := createExes(exesMap)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return err
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
