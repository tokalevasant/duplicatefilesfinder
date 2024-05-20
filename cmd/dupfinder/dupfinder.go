package dupfinder

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

type DuplicateFiles struct {
	fileHash string
	files    []os.FileInfo
}

type fileDetails struct {
	fileHash string
	filePath string
	fileInfo os.FileInfo
}

func FindDuplicateFiles(path string) []DuplicateFiles {

	var files []fileDetails
	listFilesInDir(path, files)

	fileMap := make(map[string][]fileDetails)

	for _, fi := range files {
		val, ok := fileMap[fi.fileHash]
		if !ok {
			val = make([]fileDetails, 1)
		}
		val = append(val, fi)
		fileMap[fi.fileHash] = val
	}

	return computeFileDuplicates(fileMap)
}

func computeFileDuplicates(fileMap map[string][]fileDetails) []DuplicateFiles {
	duplicateFiles := []DuplicateFiles{}
	for sum, fd := range fileMap{
		if len(fd) > 1{
			fis := []fs.FileInfo{}
			for _, f := range fd{
				fis = append(fis, f.fileInfo)
			}
			dup := DuplicateFiles{
				sum,
				fis,
			}
			duplicateFiles = append(duplicateFiles, dup)
		}
	}
	return duplicateFiles
}

func listFilesInDir(path string, files []fileDetails) {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, e := range entries {
		fp, err := filepath.Abs(e.Name())
		if err != nil {
			panic(err)
		}

		if e.IsDir() {
			listFilesInDir(fp, files)
		}
		if e.Type().IsRegular() {
			fi, err := os.Stat(fp)
			if err != nil {
				panic(err)
			}
			sum := FileSha256Checksum(fp)
			fd := fileDetails{sum, fp, fi}
			files = append(files, fd)
		}
	}

}

func FileSha256Checksum(file string) string {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return Sha256Checksum(data)
}

func Sha256Checksum(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}
