package deb

import (
	"log"
)

func Example_buildDevPackage() {

	pkg := NewPackage("testpkg", "0.0.2", "me")
	pkg.Description = "hiya"
	bp := NewBuildParams()
	buildFunc := func(dpkg *DevPackage, bp *BuildParams) error {
		// Generate files here.
		return nil
	}
	dpkg := NewDevPackage(pkg)
	err := buildFunc(dpkg, bp)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
