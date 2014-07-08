package deb

import (
	"testing"
)

func TestDevDebBuild(t *testing.T) {

	pkg := NewPackage("testpkg", "0.0.2", "me")
	pkg.Description = "hiya"
	bp := NewBuildParams()
	fn := func(dpkg *DevPackage, bp *BuildParams) error {
		return nil
	}
	ddpkg := NewDevPackage(pkg, fn)
	err := ddpkg.Build(bp)
	//err = pkg.Build("amd64", exesMap["amd64"])
	if err != nil {
		t.Fatalf("%v", err)
	}
}
