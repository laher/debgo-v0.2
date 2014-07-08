package deb

import (
	"testing"
)

func TestSdebBuild(t *testing.T) {
	pkg := NewPackage("testpkg", "0.0.2", "me")
	pkg.Description = "hiya"
	bp := NewBuildParams()
	fn := func(*SourcePackage, *BuildParams) error {
		return nil
	}
	spkg := NewSourcePackage(pkg, fn)
	err := spkg.Build(bp)
	if err != nil {
		t.Fatalf("Error building source package: %v", err)
	}
}
