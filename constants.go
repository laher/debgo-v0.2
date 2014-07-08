package debgo

import(
	"github.com/laher/debgo-v0.2/deb"
)

const (
	TEMPLATE_DEBIAN_SOURCE_FORMAT  = deb.FORMAT_DEFAULT                                         // Debian source formaat
	TEMPLATE_DEBIAN_SOURCE_OPTIONS = `tar-ignore = .hg
tar-ignore = .git
tar-ignore = .bzr` //specifies files to ignore while building.

	// The debian rules file describes how to build a 'source deb' into a binary deb. The default template here invokes debhelper scripts to automate this process for simple cases.
	TEMPLATE_DEBIAN_RULES = `#!/usr/bin/make -f
# -*- makefile -*-

# Uncomment this to turn on verbose mode.
#export DH_VERBOSE=1

export GOPATH=$(CURDIR){{range $i, $gpe := .ExtraData.GoPathExtra }}:{{$gpe}}{{end}}

PKGDIR=debian/{{.PackageName}}

%:
	dh $@

clean:
	dh_clean
	rm -rf $(CURDIR)/bin/* $(CURDIR)/pkg/*
	#cd $(CURDIR)/src && find * -name '*.go' -exec dirname {} \; | xargs -n1 go clean
	rm -f $(CURDIR)/goinstall.log

binary-arch: clean
	dh_prep
	dh_installdirs
	cd $(CURDIR)/src && find * -name '*.go' -exec dirname {} \; | xargs -n1 go install
	mkdir -p $(PKGDIR)/usr/bin
	cp $(CURDIR)/bin/* $(PKGDIR)/usr/bin/
	dh_strip
	dh_compress
	dh_fixperms
	dh_installdeb
	dh_gencontrol
	dh_md5sums
	dh_builddeb

binary: binary-arch`

	// The debian control file (binary debs) defines package metadata
	TEMPLATE_BINARYDEB_CONTROL = `Package: {{.PackageName}}
Priority: {{.Priority}}
{{if .Maintainer}}Maintainer: {{.Maintainer}}
{{end}}Section: {{.Section}}
Version: {{.PackageVersion}}
Architecture: {{.Architecture}}
{{if .Depends}}Depends: {{.Depends}}
{{end}}{{range $key, $value := .AdditionalControlData}}{{$key}}: {{$value}}
{{end}}Description: {{.Description}}
`

	// The debian control file (source debs) defines build metadata AND package metadata
	TEMPLATE_SOURCEDEB_CONTROL = `Source: {{.PackageName}}
Build-Depends: {{.BuildDepends}}
Priority: {{.Priority}}
Maintainer: {{.Maintainer}}
Standards-Version: {{.StandardsVersion}}
Section: {{.Section}}

Package: {{.PackageName}}
Architecture: {{.Architecture}}
Depends: ${misc:Depends}{{.Depends}}
Description: {{.Description}}
{{.Other}}`

	// The dsc file defines package metadata AND checksums
	TEMPLATE_DEBIAN_DSC = `Format: {{.Format}}
Source: {{.PackageName}}
Binary: {{.PackageName}}
Architecture: {{.Architecture}}
Version: {{.PackageVersion}}
Maintainer: {{.Maintainer}}
Standards-Version: {{.StandardsVersion}}
Build-Depends: {{.BuildDepends}}
Priority: {{.Priority}}
Section: {{.Section}}
Checksums-Sha1:{{range .Checksums.ChecksumsSha1}}
 {{.Checksum}} {{.Size}} {{.File}}{{end}}
Checksums-Sha256:{{range .Checksums.ChecksumsSha256}}
 {{.Checksum}} {{.Size}} {{.File}}{{end}}
Files:{{range .Checksums.ChecksumsMd5}}
 {{.Checksum}} {{.Size}} {{.File}}{{end}}
{{.Other}}`

	TEMPLATE_CHANGELOG_HEADER        = `{{.PackageName}} ({{.PackageVersion}}) {{.Status}}; urgency=low`
	TEMPLATE_CHANGELOG_INITIAL_ENTRY = `  * Initial import`
	TEMPLATE_CHANGELOG_FOOTER        = ` -- {{.Maintainer}} <{{.MaintainerEmail}}>  {{.EntryDate}}`
	TEMPLATE_DEBIAN_COPYRIGHT        = `Copyright 2013 {{.PackageName}}`
	TEMPLATE_DEBIAN_README           = `{{.PackageName}}
==========

`
	DEVDEB_GO_PATH_DEFAULT    = "/usr/share/gocode" // This is used by existing -dev.deb packages e.g. golang-doozer-dev and golang-protobuf-dev
	GO_PATH_EXTRA_DEFAULT     = ":" + DEVDEB_GO_PATH_DEFAULT
	
)
