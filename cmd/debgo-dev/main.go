package main

import (
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/cmd"
	"log"
)

func main() {
	name := "debgo-dev"
	log.SetPrefix("["+name+"] ")
	//set to empty strings because they're being overridden
	pkg := deb.NewGoPackage("","","")

	fs := cmdutils.InitFlags(name, pkg)
	err := cmdutils.ParseFlags(name, pkg, fs)
	if err != nil {
		log.Fatalf("%v", err)
	}

	ddpkg := deb.NewDevPackage(pkg)
	err = ddpkg.BuildWithDefaults()
	if err != nil {
		log.Fatalf("%v", err)
	}

}
