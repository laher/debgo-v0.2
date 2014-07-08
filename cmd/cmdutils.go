package cmdutils

import (
	"flag"
	"github.com/laher/debgo-v0.2/deb"
	"os"
)

func InitFlags(name string, pkg *deb.Package, build *deb.BuildParams) *flag.FlagSet {

	fs := flag.NewFlagSet(name, flag.ContinueOnError)

	fs.StringVar(&pkg.Name, "name", "", "Package name")
	fs.StringVar(&pkg.Version, "version", "", "Package version")
	fs.StringVar(&pkg.Maintainer, "maintainer", "", "Package maintainer")
	fs.StringVar(&pkg.Description, "description", "", "Description")
	fs.BoolVar(&build.IsRmtemp, "rmtemp", false, "Remove 'temp' dirs")
	fs.BoolVar(&build.IsVerbose, "verbose", false, "Show log messages")

	return fs
}

func ParseFlags(name string, pkg *deb.Package, fs *flag.FlagSet) error {
	err := fs.Parse(os.Args[1:])
	if err == nil {
		err = pkg.Validate()
	}
//	if err != nil {
//		println("")
//		fmt.Fprintf(os.Stderr, "Usage of %s:\n", name)
//		fs.PrintDefaults()
//		println("")
//	}
	return err
}
