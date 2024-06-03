package dupfinder

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type DuplicateFiles struct {
	fileHash string
	files    []os.FileInfo
}

type fileDetails struct {
	fileHash string
	filePath string
	fileInfo os.FileInfo
}

func FindDuplicateFiles(paths ...string) []DuplicateFiles {

	var files []fileDetails
	for _, p := range paths {
		files = append(files, listFilesInDir(p)...)
	}
	fileMap := make(map[string][]fileDetails)

	for _, fi := range files {
		val, ok := fileMap[fi.fileHash]
		if !ok {
			val = make([]fileDetails, 0)
		}
		val = append(val, fi)
		fileMap[fi.fileHash] = val
	}

	return computeFileDuplicates(fileMap)
}

func computeFileDuplicates(fileMap map[string][]fileDetails) []DuplicateFiles {
	duplicateFiles := []DuplicateFiles{}
	for sum, fd := range fileMap {
		if len(fd) > 1 {
			fis := []fs.FileInfo{}
			for _, f := range fd {
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

func listFilesInDir(path string) []fileDetails {
	files := []fileDetails{}
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, e := range entries {
		fp := filepath.Join(path, e.Name())

		if e.IsDir() {
			f := listFilesInDir(fp)
			files = append(files, f...)
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
	return files
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
