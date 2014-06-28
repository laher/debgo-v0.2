package deb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDebBuild(t *testing.T) {

	pkg := NewPackage("testpkg", "0.0.2", "me")
	pkg.Description = "hiya"
	pkg.IsRmtemp = false
	pkg.IsVerbose = true

	exesMap := map[string][]string{
		"amd64": []string{"_test/a.amd64"},
		"i386":  []string{"_test/a.386"},
		"armel": []string{"_test/a.arm"}}
	err := createExes(exesMap)
	if err != nil {
		t.Fatalf("%v", err)
	}

	bpkg := NewBinaryPackage(pkg, exesMap)
	err = bpkg.BuildAllWithDefaults()
	//err = pkg.Build("amd64", exesMap["amd64"])
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
