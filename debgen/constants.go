package debgen

import (
	"github.com/laher/debgo-v0.2/deb"
)

const (
	GlobGoSources               = "*.go"
	TemplateDebianSourceFormat  = deb.FormatDefault                                      // Debian source formaat
	TemplateDebianSourceOptions = `tar-ignore = .hg
tar-ignore = .git
tar-ignore = .bzr` //specifies files to ignore while building.

	// The debian rules file describes how to build a 'source deb' into a binary deb. The default template here invokes debhelper scripts to automate this process for simple cases.
	TemplateDebianRules = `#!/usr/bin/make -f
# -*- makefile -*-

# Uncomment this to turn on verbose mode.
#export DH_VERBOSE=1

export GOPATH=$(CURDIR){{range $i, $gpe := .Package.ExtraData.GoPathExtra }}:{{$gpe}}{{end}}

PKGDIR=debian/{{.Package.Name}}

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

	mkdir -p $(PKGDIR)/usr/bin $(CURDIR)/bin/
	mkdir -p $(PKGDIR)/usr/share/gopkg/ $(CURDIR)/pkg/

	BINFILES=$(wildcard $(CURDIR)/bin/*)

	for x in$(BINFILES); do \
		cp $$x $(PKGDIR)/usr/bin/; \
	done;

	PKGFILES=$(wildcard $(CURDIR)/pkg/*.a)
	for x in$(PKGFILES); do \
		cp $$x $(PKGDIR)/usr/share/gopkg/; \
	done;

	dh_strip
	dh_compress
	dh_fixperms
	dh_installdeb
	dh_gencontrol
	dh_md5sums
	dh_builddeb

binary: binary-arch`

	// The debian control file (binary debs) defines package metadata
	TemplateBinarydebControl = `Package: {{.Package.Name}}
Priority: {{.Package.Priority}}
{{if .Package.Maintainer}}Maintainer: {{.Package.Maintainer}}
{{end}}Section: {{.Package.Section}}
Version: {{.Package.Version}}
Architecture: {{.Deb.Architecture}}
{{if .Package.Depends}}Depends: {{.Package.Depends}}
{{end}}{{range $key, $value := .Package.AdditionalControlData}}{{$key}}: {{$value}}
{{end}}Description: {{.Package.Description}}
`

	// The debian control file (source debs) defines build metadata AND package metadata
	TemplateSourcedebControl = `Source: {{.Package.Name}}
Build-Depends: {{.Package.BuildDepends}}
Priority: {{.Package.Priority}}
Maintainer: {{.Package.Maintainer}}
Standards-Version: {{.Package.StandardsVersion}}
Section: {{.Package.Section}}

Package: {{.Package.Name}}
Architecture: {{.Package.Architecture}}
Depends: ${misc:Depends}{{.Package.Depends}}
Description: {{.Package.Description}}
{{.Package.Other}}`

	// The dsc file defines package metadata AND checksums
	TemplateDebianDsc = `Format: {{.Package.Format}}
Source: {{.Package.Name}}
Binary: {{.Package.Name}}
Architecture: {{.Package.Architecture}}
Version: {{.Package.Version}}
Maintainer: {{.Package.Maintainer}}
Standards-Version: {{.Package.StandardsVersion}}
Build-Depends: {{.Package.BuildDepends}}
Priority: {{.Package.Priority}}
Section: {{.Package.Section}}
Checksums-Sha1:{{range .Checksums.ChecksumsSha1}}
 {{.Checksum}} {{.Size}} {{.File}}{{end}}
Checksums-Sha256:{{range .Checksums.ChecksumsSha256}}
 {{.Checksum}} {{.Size}} {{.File}}{{end}}
Files:{{range .Checksums.ChecksumsMd5}}
 {{.Checksum}} {{.Size}} {{.File}}{{end}}
{{.Package.Other}}`

	TemplateChangelogHeader       = `{{.Package.Name}} ({{.Package.Version}}) {{.Package.Status}}; urgency=low`
	TemplateChangelogInitialEntry = `  * Initial import`
	TemplateChangelogFooter       = ` -- {{.Package.Maintainer}}  {{.EntryDate}}`
	TemplateChangelogInitial      = TemplateChangelogHeader + "\n\n" + TemplateChangelogInitialEntry + "\n\n" + TemplateChangelogFooter // + "\n\n"
	TemplateDebianCopyright       = `Copyright 2014 {{.Package.Name}}`
	TemplateDebianReadme          = `{{.Package.Name}}
==========

`
	DevGoPathDefault   = "/usr/share/gocode" // This is used by existing -dev.deb packages e.g. golang-doozer-dev and golang-protobuf-dev
	GoPathExtraDefault = ":" + DevGoPathDefault

	DebianDir    = "debian"
	TplExtension = ".tpl"

	ChangelogDateLayout = "Mon, 02 Jan 2006 15:04:05 -0700"
)

var (
	SourceDebianFiles = map[string]string{
		"control":        TemplateSourcedebControl,
		"compat":         deb.DebianCompatDefault,
		"rules":          TemplateDebianRules,
		"source/format":  TemplateDebianSourceFormat,
		"source/options": TemplateDebianSourceOptions,
		"copyright":      TemplateDebianCopyright,
		"changelog":      TemplateChangelogInitial,
		"README.debian":  TemplateDebianReadme}
)
