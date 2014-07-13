package deb_test

import (
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/targz"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func Example_buildSourceDeb() {
	pkg := deb.NewPackage("testpkg", "0.0.2", "me <me@b.c>", "Nice of all the package\n")
	pkg.Description = "hiya"
	spkg := deb.NewSourcePackage(pkg)
	err := buildOrigArchive(spkg) // it's up to you how to build this
	if err != nil {
		log.Fatalf("Error building source package: %v", err)
	}
	err = buildDebianArchive(spkg) // again - do it yourself
	if err != nil {
		log.Fatalf("Error building source package: %v", err)
	}
	err = buildDscFile(spkg) // yep, same again
	if err != nil {
		log.Fatalf("Error building source package: %v", err)
	}
}

func Test_buildSourceDeb(t *testing.T) {
	pkg := deb.NewPackage("testpkg", "0.0.2", "me <me@b.c>", "Nice of all the package\n")
	pkg.Description = "hiya"
	spkg := deb.NewSourcePackage(pkg)
	err := buildOrigArchive(spkg) // it's up to you how to build this
	if err != nil {
		t.Fatalf("Error building source package: %v", err)
	}
	err = buildDebianArchive(spkg) // again - do it yourself
	if err != nil {
		t.Fatalf("Error building source package: %v", err)
	}
	err = buildDscFile(spkg) // yep, same again
	if err != nil {
		t.Fatalf("Error building source package: %v", err)
	}
}

func buildOrigArchive(spkg *deb.SourcePackage) error {
	origFilePath := filepath.Join(deb.DistDirDefault, spkg.OrigFileName)
	tgzw, err := targz.NewWriterFromFile(origFilePath)
	if err != nil {
		return err
	}
	// Add Sources Here !!
	err = tgzw.Close()
	if err != nil {
		return err
	}
	return nil
}

func buildDebianArchive(spkg *deb.SourcePackage) error {
	tgzw, err := targz.NewWriterFromFile(filepath.Join(deb.DistDirDefault, spkg.DebianFileName))
	if err != nil {
		return err
	}
	// Add Control Files here !!
	err = tgzw.Close()
	if err != nil {
		return err
	}
	return nil
}

func buildDscFile(spkg *deb.SourcePackage) error {
	dscData := []byte{} //generate this somehow. DIY (or see 'debgen' package in this repository)!
	dscFilePath := filepath.Join(deb.DistDirDefault, spkg.DscFileName)
	err := ioutil.WriteFile(dscFilePath, dscData, 0644)
	if err != nil {
		return err
	}
	return nil
}
