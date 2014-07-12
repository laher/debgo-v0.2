package main

import (
	"github.com/laher/debgo-v0.2/cmd"
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/debgen"
	"log"
)

func main() {
	name := "debgen-source"
	log.SetPrefix("[" + name + "] ")
	//set to empty strings because they're being overridden
	pkg := deb.NewPackage("", "", "", "")
	debgen.ApplyGoDefaults(pkg)
	build := deb.NewBuildParams()
	err := build.Init()
	if err != nil {
		log.Fatalf("%v", err)
	}
	fs := cmdutils.InitFlags(name, pkg, build)
	fs.StringVar(&pkg.Architecture, "arch", "all", "Architectures [any,386,armhf,amd64,all]")

	var sourceDir string
	var glob string
	var sourcesRelativeTo string
	fs.StringVar(&sourceDir, "sources", ".", "source dir")
	fs.StringVar(&glob, "sources-glob", debgen.GlobGoSources, "Glob for inclusion of sources")
	fs.StringVar(&sourcesRelativeTo, "sources-relative-to", "", "Sources relative to (it will assume relevant gopath element, unless you specify this)")
	err = cmdutils.ParseFlags(name, pkg, fs)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if sourcesRelativeTo == "" {
		sourcesRelativeTo = debgen.GetGoPathElement(sourceDir)
	}
	spkg := deb.NewSourcePackage(pkg)
	sourcesDestinationDir := pkg.Name + "_" + pkg.Version
	spkg.MappedFiles, err = debgen.GlobForSources(sourcesRelativeTo, sourceDir, glob, sourcesDestinationDir, []string{build.TmpDir, build.DestDir})
	if err != nil {
		log.Fatalf("Error resolving sources: %v", err)
	}
	err = debgen.GenSourceArtifacts(spkg, build) //, sourceDir, sourcesRelativeTo)
	if err != nil {
		log.Fatalf("%v", err)
	}

}
