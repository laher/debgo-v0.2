package main

import (
	"flag"
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/debgen"
	"log"
	"os"
)

func main() {
	name := "debgo-deb"
	log.SetPrefix("[" + name + "] ")
	build := debgen.NewBuildParams()

	err := build.Init()
	if err != nil {
		log.Fatalf("%v", err)
	}
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.BoolVar(&build.IsRmtemp, "rmtemp", false, "Remove 'temp' dirs")
	fs.BoolVar(&build.IsVerbose, "verbose", false, "Show log messages")

	var isControl, isContents, isDebianContents bool
	fs.BoolVar(&isControl, "control", false, "Show control")
	fs.BoolVar(&isContents, "contents", false, "Show contents of data archive")
	fs.BoolVar(&isDebianContents, "debian-contents", false, "Show contents of 'debian' archive (metadata and scripts)")

	//var debFile string
	//fs.StringVar(&debFile, "file", "", ".deb file")
	err = fs.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("%v", err)
	}
	args := fs.Args()
	if len(args) < 1 {
		log.Fatalf("File not specified")
	}
	if isControl {
		for _, debFile := range args {
			rdr, err := os.Open(debFile)
			if err != nil {
				log.Fatalf("%v", err)
			}
			log.Printf("File: %+v", debFile)
			err = deb.ExtractFileL2(rdr, "control.tar.gz", "control", os.Stdout)
			if err != nil {
				log.Fatalf("%v", err)
			}
		}

	} else if isContents {
		for _, debFile := range args {
			rdr, err := os.Open(debFile)
			if err != nil {
				log.Fatalf("%v", err)
			}
			log.Printf("File: %+v", debFile)
			files, err := deb.Contents(rdr, "data.tar.gz")
			if err != nil {
				log.Fatalf("%v", err)
			}
			for _, file := range files {
				log.Printf("%s", file)
			}
		}
	} else if isDebianContents {
		for _, debFile := range args {
			rdr, err := os.Open(debFile)
			if err != nil {
				log.Fatalf("%v", err)
			}
			log.Printf("File: %+v", debFile)
			files, err := deb.Contents(rdr, "control.tar.gz")
			if err != nil {
				log.Fatalf("%v", err)
			}
			for _, file := range files {
				log.Printf("%s", file)
			}
		}
	} else {
		log.Fatalf("No command specified")
	}

}
