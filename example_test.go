
package debgo_test

import (
	"github.com/laher/debgo-v0.2"
	"github.com/laher/debgo-v0.2/deb"
	"log"
	"os"
	"path/filepath"
)

func Example_BinaryPackage() {

	pkg := deb.NewPackage("testpkg", "0.0.2", "me")
	pkg.Description = "hiya"

	prep() //prepares files for packaging

	build := deb.NewBuildParams()

	bpkg := deb.NewBinaryPackage(pkg, debgo.BuildBinaryArtifactDefault)
	platform64 := bpkg.InitBinaryArtifact(deb.Arch_amd64, build)
	platform64.Executables = []string{"_test/a.amd64"}
	platform386 := bpkg.InitBinaryArtifact(deb.Arch_i386, build)
	platform386.Executables = []string{"_test/a.i386"}
	platformArm := bpkg.InitBinaryArtifact(deb.Arch_armel, build)
	platformArm.Executables = []string{"_test/a.armel"}

	err := bpkg.Build(build)
	//err = pkg.Build("amd64", exesMap["amd64"])
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Output:
	//
}

func Example_SourcePackage() {

	pkg := deb.NewPackage("testpkg", "0.0.2", "me")
	pkg.Description = "hiya"

	spkg := deb.NewSourcePackage(pkg, debgo.BuildSourcePackageDefault)
	build := deb.NewBuildParams()

	err := spkg.Build(build)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Output:
	//

}

func Example_DevPackage() {

	pkg := deb.NewPackage("testpkg", "0.0.2", "me")
	pkg.Description = "hiya"

	ddpkg := deb.NewDevPackage(pkg, debgo.BuildDevPackageDefault)
	build := deb.NewBuildParams()

	err := ddpkg.Build(build)
	//err = pkg.Build("amd64", exesMap["amd64"])
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Output:
	//
}


func prep() error {
	exesMap := map[string][]string{
		"amd64": []string{"_test/a.amd64"},
		"i386":  []string{"_test/a.i386"},
		"armel": []string{"_test/a.armel"}}
	err := createExes(exesMap)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return err
}

func createExes(exesMap map[string][]string) error {
	for _, exes := range exesMap {
		for _, exe := range exes {
			err := os.MkdirAll(filepath.Dir(exe), 0777)
			if err != nil {
				return err
			}
			fi, err := os.Create(exe)
			if err != nil {
				return err
			}
			_, err = fi.Write([]byte("echo 1"))
			if err != nil {
				return err
			}
			err = fi.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
