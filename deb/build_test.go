package deb

import (
	"testing"
	"strings"
)

func TestDebBuild(t *testing.T) {
	exes := []string {"a.b"}
	pkg := NewPackage("testpkg", "0.0.2", "me", exes)
	pkg.Description = "hiya"
	pkg.IsRmtemp = false
	pkg.Preinst = TarEntryExecutable("preinst", strings.NewReader("#!/bin/bash\necho 11111"))
	pkg.Postinst = TarEntryExecutable("postinst", strings.NewReader("#!/bin/bash\necho 22222"))
	pkg.Prerm = TarEntryExecutable("prerm", strings.NewReader("#!/bin/bash\necho 33333"))
	pkg.Postrm = TarEntryExecutable("postrm", strings.NewReader("#!/bin/bash\necho 44444"))
	err := pkg.BuildAll()
	if err != nil {
		t.Fatalf("%v", err)
	}
}
