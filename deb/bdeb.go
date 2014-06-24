package deb

import (
	"github.com/laher/argo/ar"
	"io"
	"log"
	"os"
	"path/filepath"
)

type BinaryDeb struct {
	Filename string
	TmpDir string
	DebianBinaryVersion string
	ControlArchFile string
	DataArchFile string
}


func (bdeb *BinaryDeb) GetReader() (*ar.Reader, error) {
	fi, err := os.Open(bdeb.Filename)
	if err != nil {
		return nil, err
	}
	arr, err := ar.NewReader(fi)
	if err != nil {
		return nil, err
	}
	return arr, err
}

// ExtractAll extracts all contents from the Ar archive.
// It returns a slice of all filenames.
// In case of any error, it returns the error immediately
func (bdeb *BinaryDeb) ExtractAll() ([]string, error) {
	arr, err := bdeb.GetReader()
	if err != nil {
		return nil, err
	}
	filenames := []string{}
	for {
		hdr, err := arr.Next()
		if err == io.EOF {
			// end of ar archive
			break
		}
		if err != nil {
			return nil, err
		}
		outFilename := filepath.Join(bdeb.TmpDir, hdr.Name)
		//fmt.Printf("Contents of %s:\n", hdr.Name)
		fi, err := os.Create(outFilename)
		if err != nil {
			return filenames, err
		}
		if _, err := io.Copy(fi, arr); err != nil {
			return filenames, err
		}
		err = fi.Close()
		if err != nil {
			return filenames, err
		}
		filenames = append(filenames, outFilename)
		//fmt.Println()
	}
	return filenames, nil
}

func NewBinaryDeb(filename string, tmpDir string) *BinaryDeb {
	bdeb := &BinaryDeb{}
	bdeb.Filename = filename
	bdeb.TmpDir = tmpDir
	return bdeb
}

func (bdeb *BinaryDeb) SetDefaults() {
	bdeb.DebianBinaryVersion = "2.0"
	bdeb.ControlArchFile =filepath.Join(bdeb.TmpDir,"control.tar.gz")
	bdeb.DataArchFile =filepath.Join(bdeb.TmpDir,"data.tar.gz")
}

func (bdeb *BinaryDeb) WriteBytes(aw *ar.Writer, filename string, bytes []byte) error {
	hdr := &ar.Header {
		Name: filename,
		Size: int64(len(bytes))}
	if err := aw.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := aw.Write(bytes); err != nil {
		return err
	}
	return nil
}

func (bdeb *BinaryDeb) WriteFromFile(aw *ar.Writer, filename string) error {
	finf, err := os.Stat(filename)
	if err != nil {
		return err
	}
	log.Printf("Finf size: %d", finf.Size())
	hdr, err := ar.FileInfoHeader(finf)
	if err != nil {
		return err
	}
	if err := aw.WriteHeader(hdr); err != nil {
		return err
	}
	log.Printf("Header Size: %d", hdr.Size)
	fi, err := os.Open(filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(aw, fi); err != nil {
		return err
	}

	err = fi.Close()
	if err != nil {
		return err
	}
	return nil

}

func (bdeb *BinaryDeb) WriteAll() error {
	log.Printf("Building deb %s", bdeb.Filename)
	wtr, err := os.Create(bdeb.Filename)
	if err != nil {
		return err
	}

	aw := ar.NewWriter(wtr)

	log.Printf("Writing debian-binary")
	err = bdeb.WriteBytes(aw, "debian-binary", []byte(bdeb.DebianBinaryVersion + "\n"))
	if err != nil {
		return err
	}
	log.Printf("Writing control file %s", bdeb.ControlArchFile)
	err = bdeb.WriteFromFile(aw, bdeb.ControlArchFile)
	if err != nil {
		return err
	}
	log.Printf("Writing data file %s", bdeb.DataArchFile)
	err = bdeb.WriteFromFile(aw, bdeb.DataArchFile)
	if err != nil {
		return err
	}
	return nil
}
