package main

import (
	"github.com/laher/debgo/deb"
	"log"
	"os"
)

func main() {
	log.SetPrefix("[debgo] ")
	args := os.Args()
	name := args[0]
	version := args[1]
	maintainer := args[2]
	pkg := deb.NewPackage(name, version, maintainer)
	pkg.Description = description
	pkg.IsRmtemp = false
	pkg.IsVerbose = true
	// TODO determine this platform
	// TODO find executables for this platform
}
