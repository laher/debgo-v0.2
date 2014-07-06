package main

import (
	"github.com/laher/debgo-v0.2/cmd"
	"github.com/laher/debgo-v0.2/deb"
	"log"
)

func main() {
	name := "debgo-deb"
	log.SetPrefix("[" + name + "] ")
	//set to empty strings because they're being overridden
	pkg := deb.NewGoPackage("", "", "")

	fs := cmdutils.InitFlags(name, pkg)
	bpkg := deb.NewBinaryPackage(pkg)
	var binDir string
	fs.StringVar(&binDir, "binaries", "", "directory containing binaries for this platform")
	fs.StringVar(&pkg.Architecture, "arch", "any", "Architectures [any,386,armel,amd64,all]")
	err := cmdutils.ParseFlags(name, pkg, fs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	// TODO determine this platform
	// TODO find executables for this platform
	err = bpkg.Build()
	if err != nil {
		log.Fatalf("%v", err)
	}

}
