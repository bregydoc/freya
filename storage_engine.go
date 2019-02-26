package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type StorageEngine struct {
	StoreFolder         string
	AbsoluteStoreFolder string
	ready               bool
	metadata            *os.File
}

type StorageEngineOptions struct {
	ContentType string
}

func NewStorageEngine(folderRoot string) (*StorageEngine, error) {
	if filepath.IsAbs(folderRoot) {
		return &StorageEngine{
			AbsoluteStoreFolder: folderRoot,
			StoreFolder:         folderRoot,
		}, nil
	}

	absRoot, err := filepath.Abs(folderRoot)
	if err != nil {
		return nil, err
	}

	return &StorageEngine{
		StoreFolder:         folderRoot,
		AbsoluteStoreFolder: absRoot,
	}, nil
}

func (s *StorageEngine) Init() error {
	if s.ready && s.metadata != nil {
		return nil
	}

	metadataFileName := path.Join(s.AbsoluteStoreFolder, "metadata.fse")
	file, err := os.Open(metadataFileName)
	if os.IsExist(err) {
		s.metadata = file
	} else {
		f, err := os.Create(metadataFileName)
		if err != nil {
			return err
		}
		s.metadata = f
	}

	s.ready = true
	return nil

}

func (s *StorageEngine) SaveObject(bucketName string, filename string, data []byte, opts *StorageEngineOptions) error {
	if !s.ready {
		return errors.New("storage are not ready, please call to .Init() first")
	}
	bucketDir := path.Join(s.AbsoluteStoreFolder, bucketName)

	err := os.MkdirAll(bucketDir, os.ModeDir)
	if err != nil {
		return err
	}

	nestedDirs, _ := path.Split(filename)
	if nestedDirs != "" {

		err := os.MkdirAll(path.Join(bucketDir, nestedDirs), os.ModeDir)
		if err != nil {
			return err
		}
	}
	finalPath := path.Join(bucketDir, filename)

	err = ioutil.WriteFile(finalPath, data, 0644)
	if err != nil {
		return err
	}

	_, err = s.metadata.WriteString(fmt.Sprintf("[SAVED] %s %s %d\n", bucketName, filename, len(data)))
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageEngine) GetObject(bucketName string, filename string) ([]byte, error) {
	if !s.ready {
		return nil, errors.New("storage are not ready, please call to .Init() first")
	}

	f := path.Join(s.AbsoluteStoreFolder, bucketName, filename)

	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (s *StorageEngine) ListObjects(bucketName string) ([]string, error) {
	bucket := path.Join(s.AbsoluteStoreFolder, bucketName)
	files, err := ioutil.ReadDir(bucket)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, 0)
	for _, f := range files {
		if !f.IsDir() {
			filenames = append(filenames, f.Name())
		}
	}
	return filenames, nil
}
