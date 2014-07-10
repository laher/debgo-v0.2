package deb

import (
	"path/filepath"
)

const (
	DebianBinaryVersionDefault = "2.0"         // This is the current version as specified in .deb archives (filename debian-binary)
	DebianCompatDefault        = "9"           // compatibility. Current version
	FormatDefault              = "3.0 (quilt)" // Format as in a .dsc file
	StatusDefault              = "unreleased"  // Status is unreleased by default. Change this once you're happy with it.

	SectionDefault          = "devel"                           //TODO: correct to use this?
	PriorityDefault         = "extra"                           //Most packages should be 'extra'
	DependsDefault          = ""                                //No dependencies
	BuildDependsDefault     = "debhelper (>= 9.1.0), golang-go" //Default build dependencies for Go packages
	StandardsVersionDefault = "3.9.4"                           //Current standards version
	ArchitectureDefault     = "any"                             //Any is the default architecture for packages

	TemplateDirDefault  = "templates"
	ResourcesDirDefault = "resources"
	WorkingDirDefault   = "."

	ExeDirDefault = "/usr/bin" //default directory for exes within the control archive
)

var (
	//	ErrNoBuildFunc     = errors.New("debgo: 'BuildFunc' is nil") // error

	TempDirDefault = filepath.Join("_out", "tmp")
	DistDirDefault = filepath.Join("_out", "dist")

	MaintainerScripts = []string{"postinst", "postrm", "prerm", "preinst"}
)
