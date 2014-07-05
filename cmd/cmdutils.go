package cmdutils


import (
	"flag"
	"fmt"
	"github.com/laher/debgo-v0.2/deb"
	"os"
)

func InitFlags(name string, pkg *deb.Package) *flag.FlagSet {

	fs := flag.NewFlagSet(name, flag.ContinueOnError)

	fs.StringVar(&pkg.Name, "name", "", "Package name")
	fs.StringVar(&pkg.Version, "version", "", "Package version")
	fs.StringVar(&pkg.Maintainer, "maintainer", "", "Package maintainer")
	fs.StringVar(&pkg.Description, "description", "", "Description")

	pkg.IsRmtemp = false
	pkg.IsVerbose = true

	return fs
}

func ParseFlags(name string, pkg *deb.Package, fs *flag.FlagSet) error {
	err := fs.Parse(os.Args[1:])
	if err == nil {
		err = pkg.Validate()
	}
	if err != nil {
		println("")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", name)
		fs.PrintDefaults()
		println("")
	}
	return err
}


