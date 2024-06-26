package dupfinder

import (
	"crypto/rand"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestFindDuplicateFilesInSingleFolder(t *testing.T) {
	temp := t.TempDir()

	r := getRandomByteSlices()
	h := Sha256Checksum(r)

	fileInfo1, err := writeFile(r, temp)
	if err != nil {
		t.Fail()
	}
	fileInfo2, err := writeFile(r, temp)
	if err != nil {
		t.Fail()
	}
	fileInfo3, err := writeFile(r, temp)
	if err != nil {
		t.Fail()
	}
	want := make([]DuplicateFiles, 1)
	want[0] = DuplicateFiles{
		fileHash: h,
		files: []fs.FileInfo{
			fileInfo1,
			fileInfo2,
			fileInfo3,
		},
	}

	got := FindDuplicateFiles(temp)

	if len(want) != len(got) || len(want[0].files) != len(got[0].files) {
		t.Errorf("Number of elements in `want` do not match `got`")
	}

}

func TestDuplicatesInMultipleFolders(t *testing.T) {
	temp1 := t.TempDir()
	temp2 := t.TempDir()

	r := getRandomByteSlices()
	h := Sha256Checksum(r)

	fileInfo1, err := writeFile(r, temp1)
	if err != nil {
		t.Fail()
	}
	fileInfo2, err := writeFile(r, temp1)
	if err != nil {
		t.Fail()
	}
	fileInfo3, err := writeFile(r, temp2)
	if err != nil {
		t.Fail()
	}

	fileInfo4, err := writeFile(r, temp2)
	if err != nil {
		t.Fail()
	}

	want := make([]DuplicateFiles, 1)
	want[0] = DuplicateFiles{
		fileHash: h,
		files: []fs.FileInfo{
			fileInfo1,
			fileInfo2,
			fileInfo3,
			fileInfo4,
		},
	}

	got := FindDuplicateFiles(temp1, temp2)

	if len(want) != len(got) || len(want[0].files) != len(got[0].files) {
		t.Errorf("Number of elements in `want` do not match `got`")
	}
}

func getRandomByteSlices() []byte {
	a := make([]byte, 10)
	rand.Read(a)
	return a
}

func writeFile(r []byte, dir string) (info fs.FileInfo, err error) {
	path := filepath.Join(dir, uuid.New().String())
	err = os.WriteFile(path, r, 0644)
	if err != nil {
		return
	}
	info, err = os.Stat(path)
	return
}
