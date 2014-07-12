package deb_test

import (
	"github.com/laher/debgo-v0.2/deb"
	"log"
	"testing"
)

func Example_buildDevPackage() {

	pkg := deb.NewPackage("testpkg", "0.0.2", "me", "A package\ntestpkg is a lovel package with many wow")
	bp := deb.NewBuildParams()
	buildFunc := func(dpkg *deb.Package, bp *deb.BuildParams) error {
		// Generate files here.
		return nil
	}
	dpkg := deb.NewDevPackage(pkg)
	err := buildFunc(dpkg, bp)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func Test_buildDevPackage(t *testing.T) {

	pkg := deb.NewPackage("testpkg", "0.0.2", "me", "A package\ntestpkg is a lovel package with many wow")
	bp := deb.NewBuildParams()
	bp.Init()
	buildFunc := func(dpkg *deb.Package, bp *deb.BuildParams) error {
		// Generate files here.
		return nil
	}
	dpkg := deb.NewDevPackage(pkg)
	err := buildFunc(dpkg, bp)
	if err != nil {
		t.Fatalf("%v", err)
	}
}
