package deb

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"archive/tar"
)

type Readable interface {
	GetReader() (io.Reader, error)
}

type StdReadable struct {
	Reader io.Reader
}

func (fr *StdReadable) GetReader() (io.Reader, error) {
	if fr.Reader == nil {
		return nil, errors.New("Reader not set")
	}
	return fr.Reader, nil
}


type TarEntry interface {
	Readable
	GetHeader() (*tar.Header, error)
}

type FileTarEntry struct {
	Filename string
	Header *tar.Header
}

func (fr *FileTarEntry) GetReader() (io.Reader, error) {
	f, err := os.Open(fr.Filename)
	return f, err
}

func (fr *FileTarEntry) GetHeader() (*tar.Header, error) {
	if fr.Header == nil {
		return nil, errors.New("Header not set")
	}
	return fr.Header, nil
}

func TarEntryExecutable(name string, content io.Reader) TarEntry {
	return TarEntryFromReader(TarHeaderExecutable(name), content)
}

func TarHeaderExecutable(name string) *tar.Header {
	return &tar.Header{
		Name: name,
		Mode: 0544}
}

func TarEntryFromReader(header *tar.Header, content io.Reader) TarEntry {
	return &StdTarEntry {
		content,
		header}
}

type StdTarEntry struct {
	Reader io.Reader
	Header *tar.Header
}

func (fr *StdTarEntry) GetReader() (io.Reader, error) {
	if fr.Reader == nil {
		return nil, errors.New("Reader not set")
	}
	return fr.Reader, nil
}

func (fr *StdTarEntry) GetHeader() (*tar.Header, error) {
	if fr.Header == nil {
		return nil, errors.New("Header not set")
	}
	return fr.Header, nil
}


func toBytes(ra Readable) ([]byte, error) {
	if ra == nil {
		return nil, nil
	}
	r, err := ra.GetReader()
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(r)
	println("all: ", string(b))
	return b, err
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


func NewPackage(name, version, maintainer string, executables []string) *DebPackage {
	pkg := new(DebPackage)
	pkg.Name = name
	pkg.Version = version
	pkg.Maintainer = maintainer
	pkg.ExecutablePaths = executables

	pkg.TmpDir = "_test/tmp"
	pkg.DestDir = "_test/dist"
	pkg.IsRmtemp = true
	pkg.WorkingDir = "."

	return pkg
}
