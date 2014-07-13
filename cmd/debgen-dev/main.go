package main

import (
	"github.com/laher/debgo-v0.2/cmd"
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/debgen"
	"log"
)

func main() {
	name := "debgen-dev"
	log.SetPrefix("[" + name + "] ")
	//set to empty strings because they're being overridden
	pkg := deb.NewPackage("", "", "", "")
	debgen.ApplyGoDefaults(pkg)
	build := debgen.NewBuildParams()
	fs := cmdutils.InitFlags(name, pkg, build)
	fs.StringVar(&pkg.Architecture, "arch", "all", "Architectures [any,386,armhf,amd64,all]")
	ddpkg := deb.NewDevPackage(pkg)

	var sourceDir string
	var glob string
	var sourcesRelativeTo string
	var sourcesDestinationDir string
	fs.StringVar(&sourceDir, "sources", build.WorkingDir, "source dir")
	fs.StringVar(&glob, "sources-glob", debgen.GlobGoSources, "Glob for inclusion of sources")
	fs.StringVar(&sourcesRelativeTo, "sources-relative-to", "", "Sources relative to (it will assume relevant gopath element, unless you specify this)")
	fs.StringVar(&sourcesDestinationDir, "sources-destination", debgen.DevGoPathDefault, "Destination dir for sources to be installed")
	err := cmdutils.ParseFlags(name, pkg, fs)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if sourcesRelativeTo == "" {
		sourcesRelativeTo = debgen.GetGoPathElement(sourceDir)
	}
	ddpkg.MappedFiles, err = debgen.GlobForSources(sourcesRelativeTo, sourceDir, glob, sourcesDestinationDir, []string{build.TmpDir, build.DestDir})
	if err != nil {
		log.Fatalf("Error resolving sources: %v", err)
	}

	err = debgen.GenDevArtifact(ddpkg, build)
	if err != nil {
		log.Fatalf("Error building -dev: %v", err)
	}

}
