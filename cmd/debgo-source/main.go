package main

import (
	"github.com/laher/debgo-v0.2/cmd"
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/debgen"
	"log"
)

func main() {
	name := "debgo-dsc"
	log.SetPrefix("[" + name + "] ")
	//set to empty strings because they're being overridden
	pkg := debgen.NewGoPackage("", "", "")
	build := deb.NewBuildParams()
	fs := cmdutils.InitFlags(name, pkg, build)
	fs.StringVar(&pkg.Architecture, "arch", "all", "Architectures [any,386,armhf,amd64,all]")

	var sourceDir string
	var glob string
	var sourcesRelativeTo string
	var sourcesDestinationDir string
	fs.StringVar(&sourceDir, "sources", ".", "source dir")
	fs.StringVar(&glob, "sources-glob", debgen.GLOB_GO_SOURCES, "Glob for inclusion of sources")
	fs.StringVar(&sourcesRelativeTo, "sources-relative-to", "", "Sources relative to (it will assume relevant gopath element, unless you specify this)")
	fs.StringVar(&sourcesDestinationDir, "source-destination", debgen.DEVDEB_GO_PATH_DEFAULT, "Destination dir for sources to be installed")
	err := cmdutils.ParseFlags(name, pkg, fs)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if sourcesRelativeTo == "" {
		sourcesRelativeTo = debgen.GetGoPathElement(sourceDir)
	}
	build.Resources, err = debgen.GlobForSources(sourcesRelativeTo, sourceDir, glob, sourcesDestinationDir, []string{build.TmpDir, build.DestDir})
	if err != nil {
		log.Fatalf("Error resolving sources: %v", err)
	}

	spkg := deb.NewSourcePackage(pkg)
	err = debgen.GenSourceArtifacts(spkg, build) //, sourceDir, sourcesRelativeTo)
	if err != nil {
		log.Fatalf("%v", err)
	}

}
