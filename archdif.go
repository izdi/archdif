package archdif

import (
	"archive/tar"
	"archive/zip"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Archive struct {
	cdir string
	name string
}

type Zip struct {
	Archive
	writer *zip.Writer
}

type Tar struct {
	Archive
	writer *tar.Writer
}

type Archiver interface {
	AddFile(filename string)
	Compress(files []string)
}

func ArchiverFactory(format *string) Archiver {
	switch *format {
	case "zip":
		return &Zip{}
	case "tar":
		return &Tar{}
	default:
		return nil
	}
}

func (t *Tar) AddFile(file string) {
	info, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		log.Fatal(err)
	}
	// Set path via name explicitly
	header.Name = file

	if err := t.writer.WriteHeader(header); err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(t.cdir, file)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := t.writer.Write(b); err != nil {
		log.Fatal(err)
	}
}

func (t *Tar) Compress(files []string) {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	t.cdir = root

	name := filepath.Base(root)
	t.name = name + ".tar"

	tarfile, err := os.Create(t.name)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tarfile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	t.writer = tar.NewWriter(tarfile)
	defer func() {
		if err := t.writer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	addFiles(t, files)
}

func (z *Zip) AddFile(file string) {
	ziping, err := z.writer.Create(file)
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(z.cdir, file)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := ziping.Write(b); err != nil {
		log.Fatal(err)
	}
}

func (z *Zip) Compress(files []string) {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	z.cdir = root

	name := filepath.Base(root)
	z.name = name + ".zip"

	zipfile, err := os.Create(z.name)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := zipfile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	z.writer = zip.NewWriter(zipfile)
	defer func() {
		if err := z.writer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	addFiles(z, files)
}

func addFiles(a Archiver, files []string)  {
	for _, f := range files {
		a.AddFile(f)
	}
}