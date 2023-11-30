package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

const dirPerm os.FileMode = 0770

func NewStorage(filesRoot string, etagEncoder EtagEncoder) Storage {
	return &storage{
		filesRoot:   filesRoot,
		etagEncoder: etagEncoder,
	}
}

// storage implements file operations. Filepath in local storage is
// filesRoot + hash path + file id.
type storage struct {
	filesRoot   string
	etagEncoder EtagEncoder
}

func (s *storage) Move(id, src string) (string, error) {
	path, hashPath, err := s.createDestDir()
	if err != nil {
		return "", err
	}

	dest := filepath.Join(path, id)
	err = MoveFile(src, dest)
	if err != nil {
		return "", err
	}

	return hashPath, nil
}

func (s *storage) Copy(id, src string) (string, error) {
	path, hashPath, err := s.createDestDir()
	if err != nil {
		return "", err
	}

	dest := filepath.Join(path, id)
	err = CopyFile(src, dest)
	if err != nil {
		return "", err
	}

	return hashPath, nil
}

func (s *storage) CopyReader(id string, src io.Reader) (hashPath string, err error) {
	path, hashPath, err := s.createDestDir()
	if err != nil {
		return "", err
	}

	dest := filepath.Join(path, id)
	destFile, err := os.Create(dest)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(destFile, src)
	if err != nil {
		return "", err
	}

	err = destFile.Close()
	if err != nil {
		return "", err
	}

	return hashPath, nil
}

func (s *storage) Delete(id, hashPath string) error {
	return os.Remove(s.GetFilepath(id, hashPath))
}

func (s *storage) DeleteByName(name string) (bool, error) {
	err := os.Remove(name)
	if err != nil {
		if errors.Is(err, syscall.ENOTEMPTY) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *storage) GetEtag(id, hashPath string) (etag string, resErr error) {
	path := s.GetFilepath(id, hashPath)
	r, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer func() {
		err = r.Close()
		if err != nil && resErr == nil {
			resErr = err
		}
	}()

	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return s.etagEncoder.Encode(b)
}

func (s *storage) createDestDir() (string, string, error) {
	hashPath := GetHashPath()
	path := filepath.Join(s.filesRoot, hashPath)

	err := os.MkdirAll(path, os.ModeDir|dirPerm)
	if err != nil {
		var pathError *os.PathError
		if errors.As(err, &pathError) {
			err = fmt.Errorf("permission error %s %s", pathError.Op, hashPath)
		}

		return "", "", err
	}

	return path, hashPath, nil
}

func (s *storage) GetFilepath(id, hashPath string) string {
	return filepath.Join(s.filesRoot, hashPath, id)
}
