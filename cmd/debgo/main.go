package main

import (
	"github.com/laher/debgo-v0.2/deb"
	"log"
	"os"
)

func main() {
	log.SetPrefix("[debgo] ")
	args := os.Args
	name := args[1]
	version := args[2]
	maintainer := args[3]
	description := args[4]
	pkg := deb.NewPackage(name, version, maintainer)
	pkg.Description = description
	pkg.IsRmtemp = false
	pkg.IsVerbose = true
	// TODO determine this platform
	// TODO find executables for this platform
	bpkg := deb.NewBinaryPackage(pkg)
	err := bpkg.BuildAllWithDefaults()
	if err != nil {
		log.Fatalf("%v", err)
	}

}
