package deb

import (
	"os"
	"path/filepath"
	"testing"
)

func Example_buildBinaryDeb(t *testing.T) {

	pkg := NewPackage("testpkg", "0.0.2", "me")
	pkg.Description = "hiya"
	bp := NewBuildParams()
	err := bp.Init()
	if err != nil {
		t.Fatalf("%v", err)
	}

	exesMap := map[string][]string{
		"amd64": []string{"_test/a.amd64"},
		"i386":  []string{"_test/a.i386"},
		"armhf": []string{"_test/a.armhf"}}
	err = createExes(exesMap)
	if err != nil {
		t.Fatalf("%v", err)
	}
	bpkg := NewBinaryPackage(pkg)
	platform64 := bpkg.InitBinaryArtifact(Arch_amd64, bp)
	platform64.Executables = []string{"_test/a.amd64"}
	platform386 := bpkg.InitBinaryArtifact(Arch_i386, bp)
	platform386.Executables = []string{"_test/a.i386"}
	platformArm := bpkg.InitBinaryArtifact(Arch_armhf, bp)
	platformArm.Executables = []string{"_test/a.armhf"}

	buildFunc := func(pkg *BinaryPackage, arch *BinaryArtifact, build *BuildParams) error {
		//Per-platform build logic goes here
		return nil
	}
	err = bpkg.Build(bp, buildFunc)
	if err != nil {
		t.Fatalf("%v", err)
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
