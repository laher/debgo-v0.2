package deb

import (
	"testing"
)

func TestDevDebBuild(t *testing.T) {

	pkg := NewPackage("testpkg", "0.0.2", "me")
	pkg.Description = "hiya"
	pkg.IsRmtemp = false
	pkg.IsVerbose = true
	ddpkg := NewDevPackage(pkg)
	err := ddpkg.Build()
	//err = pkg.Build("amd64", exesMap["amd64"])
	if err != nil {
		t.Fatalf("%v", err)
	}
}
