package archiveutil

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	ZIP = "zip"
	TAR = "tar"
)

type Archive interface {
	AddFile(filePath string) error
	AddFolder(folderPath string) error
	Close() error
}

type WriteFunc func(info os.FileInfo, file io.Reader, entryName string) (err error)

// Create Archiver by archiveType.
// Supported archiveTypes: "zip"
func CreateArchive(archiveType string, w io.Writer) Archive {
	switch archiveType {
	case ZIP:
		return CreateZipArchive(w)
	}
	return nil
}

func getSubDir(dir string, rootDir string) (subDir string) {
	subDir = strings.Replace(dir, rootDir, "", 1)
	parts := strings.Split(rootDir, string(os.PathSeparator))
	subDir = path.Join(parts[len(parts)-1], subDir)
	return
}

func addAllEntry(dir string, baseDir string, writerFunc WriteFunc) error {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, info := range fileInfos {
		full := path.Join(dir, info.Name())
		var file io.Reader
		if !info.IsDir() {
			file, err = os.Open(full)
			if err != nil {
				return err
			}
		}
		subDir := getSubDir(dir, baseDir)
		entryName := path.Join(subDir, info.Name())
		if err := writerFunc(info, file, entryName); err != nil {
			return err
		}
		if info.IsDir() {
			addAllEntry(full, baseDir, writerFunc)
		}
	}
	return nil
}
