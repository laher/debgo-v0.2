package debgen_test

import (
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/debgen"
	"log"
)

func Example_sourcePackage() {

	pkg := deb.NewPackage("testpkg", "0.0.2", "me", "Dummy package for doing nothing")

	spkg := deb.NewSourcePackage(pkg)
	build := deb.NewBuildParams()
	build.IsRmtemp = false
	err := build.Init()
	if err != nil {
		log.Fatalf("Error initializing dirs: %v", err)

	}
	err = debgen.GenSourceArtifacts(spkg, build)

	if err != nil {
		log.Fatalf("Error building source: %v", err)
	}

	// Output:
	//
}
