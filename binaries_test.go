package deb

import (
	"testing"
	"strings"
)

func TestDebBuild(t *testing.T) {
	exes := []string {"a.b"}
	pkg := NewPackage("testpkg", "0.0.2", "me", exes)
	pkg.Description = "hiya"
	pkg.Postinst = &StdReadable{Reader: strings.NewReader("#!/bin/bash\necho 11111")}
	err := pkg.Build("armel")
	if err != nil {
		t.Fatalf("%v", err)
	}
}
