package deb

import (
	"testing"
	"os"
	"path/filepath"
)

func TestDebBuild(t *testing.T) {
	exes := []string{"_test/a.b"}
	err := createExes(exes)
	if err != nil {
		t.Fatalf("%v", err)
	}
	pkg := NewPackage("testpkg", "0.0.2", "me", exes)
	pkg.Description = "hiya"
	pkg.IsRmtemp = false
	pkg.IsVerbose = true
	err = pkg.Build("amd64")
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func createExes(exes []string) error {
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
	return nil
}
