package archiveutil

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type Zip struct {
	Writer *zip.Writer
}

func CreateZipArchive(w io.Writer) *Zip {
	z := new(Zip)
	z.Writer = zip.NewWriter(w)
	return z
}
func (zipp *Zip) AddFile(filePath string) error {
	bytearq, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	filep, err := zipp.Writer.Create(filePath)
	if err != nil {
		return err
	}
	_, err = filep.Write(bytearq)
	if err != nil {
		return err
	}
	return nil
}

func (zipp *Zip) AddFolder(folderPath string) error {
	folderPath = path.Clean(folderPath)
	return addAllEntry(folderPath, folderPath, func(info os.FileInfo, file io.Reader, entryName string) (err error) {
		if file == nil {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			header.Method = zip.Deflate
		}
		header.Name = entryName
		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		}
		writer, err := zipp.Writer.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err := io.Copy(writer, file); err != nil {
			return err
		}
		return nil
	})
}

func (zipp *Zip) Close() error {
	return zipp.Writer.Close()
}
