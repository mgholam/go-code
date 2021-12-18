package zipfile

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"time"
)

// Add an existing file to a zip file
func AddFile(zipfilename string, filename string) error {

	fni, err := os.Stat(filename)
	if err != nil {
		return err
	}

	if _, err := os.Stat(zipfilename); err != nil {
		return zipCreateAddFile(zipfilename, filename, fni)
	}

	return zipAppendFile(zipfilename, filename, fni)
}

// Add a stream to a zip file
func Add(zipfilename string, filename string, r io.Reader) error {

	if _, err := os.Stat(zipfilename); err != nil {
		return zipCreateAdd(zipfilename, filename, r)
	}

	return zipAppend(zipfilename, filename, r)
}

func zipAppend(zipfilename, filename string, r io.Reader) error {
	zipReader, err := zip.OpenReader(zipfilename)
	if err != nil {
		return err
	}
	targetFile, err := os.Create(zipfilename + ".tmp")
	if err != nil {
		return err
	}
	targetZipWriter := zip.NewWriter(targetFile)

	for _, zipItem := range zipReader.File {
		zipItemReader, _ := zipItem.Open()
		header := zipItem.FileHeader
		targetItem, _ := targetZipWriter.CreateHeader(&header)
		io.Copy(targetItem, zipItemReader)
	}

	z, _ := targetZipWriter.CreateHeader(&zip.FileHeader{
		Name:     filename,
		Modified: time.Now(),
		Method:   zip.Deflate,
		Comment:  time.Now().String(),
	})

	io.Copy(z, r)
	zipReader.Close()
	targetZipWriter.Close()
	targetFile.Close()

	// rename output zipfile
	os.Remove(zipfilename)
	os.Rename(zipfilename+".tmp", zipfilename)
	return nil
}

func zipCreateAdd(zipfilename, filename string, r io.Reader) error {
	archive, err := os.Create(zipfilename)
	if err != nil {
		return err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()
	z, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:     filename,
		Modified: time.Now(),
		Comment:  time.Now().String(),
		Method:   zip.Deflate,
	})
	if err != nil {
		return err
	}
	if _, err := io.Copy(z, r); err != nil {
		return err
	}
	return nil
}

func zipCreateAddFile(zfn string, fn string, fni fs.FileInfo) error {
	archive, err := os.Create(zfn)
	if err != nil {
		return err
	}
	defer archive.Close()

	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()
	z, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:     fn,
		Modified: fni.ModTime(),
		Comment:  fni.ModTime().String(),
		Method:   zip.Deflate,
	})
	if err != nil {
		return err
	}
	if _, err := io.Copy(z, f); err != nil {
		return err
	}
	return nil
}

func zipAppendFile(zfn string, fn string, fni fs.FileInfo) error {
	zipReader, err := zip.OpenReader(zfn)
	if err != nil {
		return err
	}
	targetFile, err := os.Create(zfn + ".tmp")
	if err != nil {
		return err
	}
	targetZipWriter := zip.NewWriter(targetFile)

	for _, zipItem := range zipReader.File {
		zipItemReader, _ := zipItem.Open()
		header := zipItem.FileHeader
		targetItem, _ := targetZipWriter.CreateHeader(&header)
		io.Copy(targetItem, zipItemReader)
	}

	z, _ := targetZipWriter.CreateHeader(&zip.FileHeader{
		Name:     fn,
		Modified: fni.ModTime(),
		Comment:  fni.ModTime().String(),
		Method:   zip.Deflate,
	})
	f, _ := os.Open(fn)
	io.Copy(z, f)
	f.Close()
	zipReader.Close()
	targetZipWriter.Close()
	targetFile.Close()

	// rename output zipfile
	os.Remove(zfn)
	os.Rename(zfn+".tmp", zfn)
	return nil
}
