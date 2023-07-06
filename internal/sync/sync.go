package sync

import (
	"os"
	"path/filepath"
	"reflect"
	"time"
)

func FilesSync(localPath, remotePath, fileType string) error {

	local, remote, err := fileScan(localPath, remotePath, fileType)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(local, remote) {
		return nil
	}

	for name, tt := range local {
		t, ok := remote[name]
		if !ok || tt.After(t) {
			err := copyFile(filepath.Clean(localPath+"/"+name), filepath.Clean(remotePath+"/"+name))
			if err != nil {
				return err
			}
		}
		delete(remote, name)
	}

	if len(remote) != 0 {
		for name := range remote {
			err := os.Remove(filepath.Clean(remotePath + "/" + name))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func fileScan(localPath, remotePath, fileType string) (local map[string]time.Time, remote map[string]time.Time, err error) {

	local = make(map[string]time.Time)
	remote = make(map[string]time.Time)

	localFiles, err := os.ReadDir(localPath)
	if err != nil {
		return nil, nil, err
	}

	for _, file := range localFiles {
		if filepath.Ext(file.Name()) == fileType {
			info, err := file.Info()
			if err != nil {
				return nil, nil, err
			}
			local[file.Name()] = info.ModTime()
		}
	}

	remoteFiles, err := os.ReadDir(remotePath)
	if err != nil {
		return nil, nil, err
	}

	for _, file := range remoteFiles {
		if filepath.Ext(file.Name()) == fileType {
			info, err := file.Info()
			if err != nil {
				return nil, nil, err
			}
			remote[file.Name()] = info.ModTime()
		}
	}

	return local, remote, nil
}

func copyFile(src, dest string) error {

	file, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dest, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
