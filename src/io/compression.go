package io

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(file string, dst string) error {

	r, error := os.Open(file)
	if nil != error {
		return error
	}

	gzr, error := gzip.NewReader(r)
	if error != nil {
		return error
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, error := tr.Next()

		switch {

		// if no more files are found return
		case error == io.EOF:
			return nil

		// return any other error
		case error != nil:
			return error

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, error := os.Stat(target); error != nil {
				if error := os.MkdirAll(target, 0755); error != nil {
					return error
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, error := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if error != nil {
				return error
			}

			// copy over contents
			if _, error := io.Copy(f, tr); error != nil {
				return error
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
