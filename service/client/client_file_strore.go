package client

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

//FileStore interface to store file
type FileStore interface {
	// save saves a new file
	Save(fileId string, fileType string, fileData bytes.Buffer, fileName string) (string, error)
}

type DiskFileStore struct {
	mutex      sync.RWMutex
	fileFolder string
	// files      map[string]*FileInfo
}

type FileInfo struct {
	Path string
}

// returns a new DiskFileStore
func NewDiskFileStore(fileFolder string) *DiskFileStore {
	return &DiskFileStore{
		fileFolder: fileFolder,
		// files:      make(map[string]*FileInfo),
	}
}

//  Save saves a new file to the store
func (store *DiskFileStore) Save(fileData bytes.Buffer, fileName string) (string, error) {

	// hashId, err := uuid.NewRandom()
	// fileId = hashId.String()

	// if err != nil {
	// 	return "", fmt.Errorf("couldn't generate file id: %v", err)
	// }

	filepath := fmt.Sprintf("%s/%s", store.fileFolder, fileName)
	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("couldn't create file %v: ", err)
	}

	_, err = fileData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("couldn't write file %v: ", err)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	// store.files[fileId] = &FileInfo{
	// 	Path: filepath,
	// }
	fmt.Println(filepath)
	return filepath, nil
}
