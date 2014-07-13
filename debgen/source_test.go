package debgen_test

import (
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/debgen"
	"log"
)

func Example_genSourcePackage() {

	pkg := deb.NewPackage("testpkg", "0.0.2", "me <a@me.org>", "Dummy package for doing nothing\n")
	build := debgen.NewBuildParams()
	build.IsRmtemp = false
	debgen.ApplyGoDefaults(pkg, build)
	spkg := deb.NewSourcePackage(pkg)
	err := build.Init()
	if err != nil {
		log.Fatalf("Error initializing dirs: %v", err)
	}
	sourcesDestinationDir := pkg.Name + "_" + pkg.Version
	sourceDir := ".."
	sourcesRelativeTo := debgen.GetGoPathElement(sourceDir)
	spkg.MappedFiles, err = debgen.GlobForSources(sourcesRelativeTo, sourceDir, debgen.GlobGoSources, sourcesDestinationDir, []string{build.TmpDir, build.DestDir})
	if err != nil {
		log.Fatalf("Error resolving sources: %v", err)
	}

	err = debgen.GenSourceArtifacts(spkg, build)

	if err != nil {
		log.Fatalf("Error building source: %v", err)
	}

	// Output:
	//
}
