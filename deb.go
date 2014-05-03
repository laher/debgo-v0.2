package deb

import (
	"io"
	"io/ioutil"
	"os"
)

type Readable interface {
	GetReader() (io.Reader, error)
}

type FileReadable struct {
	Filename string
}

func (fr *FileReadable) GetReader() (io.Reader, error) {
	f, err := os.Open(fr.Filename)
	return f, err
}

type StdReadable struct {
	Reader io.Reader
}

func (dr *StdReadable) GetReader() (io.Reader, error) {
	return dr.Reader, nil
}

func toBytes(ra Readable) ([]byte, error) {
	if ra == nil {
		return nil, nil
	}
	r, err := ra.GetReader()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

//package
type DebPackage struct {
	Name string
	Version string
	Description string
	Maintainer string
	MaintainerEmail string
	Metadata map[string]interface{}

	Architecture string

	Preinst Readable
	Postinst Readable
	Prerm Readable
	Postrm Readable
	Changelog Readable

	ExecutablePaths []string
	OtherFiles map[string]string

	IsVerbose bool

	//only required for sourcedebs
	Depends string
	BuildDepends string
	TemplateDir string

	IsRmtemp bool
	TmpDir string
	DestDir string
	WorkingDir string
}
