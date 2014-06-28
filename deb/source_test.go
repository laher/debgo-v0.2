package deb

import (
	"testing"
)

/*
func TestSdebCopy(t *testing.T) {
	workingDirectory := "."
	err := os.MkdirAll(workingDirectory, 0777)
	if err != nil {
		t.Fatalf("%v", err)
	}
	tmpDir := filepath.Join(workingDirectory, ".")
	destDir := filepath.Join(tmpDir, "src")
	workingDirectory = "."
	pkg := NewPackage("a", "1", "me")
	spkg := NewSourcePackage(pkg)
	err = spkg.CopySourceRecurse(workingDirectory, destDir)
	if err != nil {
		t.Fatalf("%v", err)
	}
	//TODO: find code & copy
	//ioutil.WriteFile(filepath.Join(debianDir, "control"), sdebControlFile, 0666)
	//TODO: targz
}
*/
func TestSdebBuild(t *testing.T) {
	pkg := NewPackage("testpkg", "0.0.2", "me")
	pkg.Description = "hiya"
	spkg := NewSourcePackage(pkg)
	err := spkg.BuildWithDefaults()
	if err != nil {
		t.Fatalf("Error building source package: %v", err)
	}
}
