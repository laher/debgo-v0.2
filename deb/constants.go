package deb

import (
	"errors"
)

const (
	DEBIAN_BINARY_VERSION_DEFAULT = "2.0"         // This is the current version as specified in .deb archives (filename debian-binary)
	DEBIAN_COMPAT_DEFAULT         = "9"           // compatibility. Current version
	FORMAT_DEFAULT                = "3.0 (quilt)" // Format as in a .dsc file
	STATUS_DEFAULT                = "unreleased"  //

	SECTION_DEFAULT           = "devel" //TODO: correct to use this?
	PRIORITY_DEFAULT          = "extra"
	DEPENDS_DEFAULT           = ""
	BUILD_DEPENDS_DEFAULT     = "debhelper (>= 9.1.0), golang-go"
	STANDARDS_VERSION_DEFAULT = "3.9.4"
	ARCHITECTURE_DEFAULT      = "any"

	TEMP_DIR_DEFAULT    = "_test/tmp"
	DIST_DIR_DEFAULT    = "_test/dist"
	WORKING_DIR_DEFAULT = "."
)

var (
	ErrNoBuildFunc = errors.New("debgo: 'BuildFunc' is nil")
)
