// Copyright 2013 Am Laher.
// This code is adapted from code within the Go tree.
// See Go's licence information below:
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package targz_test

import (
	"bytes"
	"fmt"
	"github.com/laher/debgo-v0.2/targz"
	"archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"
)

// change this to true for generating an archive on the Filesystem
var (
 filename = filepath.Join("_test", "tmp.tar.gz")
 isFs = true
)
func Example() {
	// Create a buffer to write our archive to.
	wtr := writer(isFs)

	// Create a new ar archive.
	tgzw := targz.NewWriter(wtr)

	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive contains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling licence."},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Size: int64(len(file.Body)),
		}
		if err := tgzw.Tw.WriteHeader(hdr); err != nil {
			log.Fatalln(err)
		}
		if _, err := tgzw.Tw.Write([]byte(file.Body)); err != nil {
			log.Fatalln(err)
		}
	}
	// Make sure to check the error on Close.
	if err := tgzw.Close(); err != nil {
		log.Fatalln(err)
	}
	rdr := reader(isFs, wtr)
	tgzr, err := targz.NewReader(rdr)
	if err != nil {
		log.Fatalln(err)
	}

	// Iterate through the files in the archive.
	for {
		hdr, err := tgzr.Next()
		if err == io.EOF {
			// end of ar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Contents of %s:\n", hdr.Name)
		if _, err := io.Copy(os.Stdout, tgzr); err != nil {
			log.Fatalln(err)
		}

		fmt.Println()
	}

	// Output:
	// Contents of readme.txt:
	// This archive contains some text files.
	// Contents of gopher.txt:
	// Gopher names:
	// George
	// Geoffrey
	// Gonzo
	// Contents of todo.txt:
	// Get animal handling licence.
}

func reader(isFs bool, w io.Writer) io.Reader {
	if isFs {
		fi := w.(*os.File)
		err := fi.Close()
		if err != nil {
			log.Fatalln(err)
		}

		r, err := os.Open(filename)
		if err != nil {
			log.Fatalln(err)
		}
		return r
	} else {
		buf := w.(*bytes.Buffer)
		// Open the ar archive for reading.
		r := bytes.NewReader(buf.Bytes())
		return r
	}

}

func writer(isFs bool) io.Writer {
	if isFs {
		fi, err := os.Create(filename)
		if err != nil {
			log.Fatalln(err)
		}
		return fi
	} else {
		return new(bytes.Buffer)
	}

}
