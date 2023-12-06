// Package file contains operations with files.
package file

import (
	"io"
)

//go:generate mockgen -destination=../../mocks/lib/file/file.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/file Storage

// Storage is used to implement file modification operations.
type Storage interface {
	// Move moves file from src to local storage.
	Move(id, src string) (hashPath string, err error)
	// Copy copies file from src to local storage.
	Copy(id, src string) (hashPath string, err error)
	CopyReader(id string, src io.Reader) (hashPath string, err error)
	// GetFilepath returns absolute file path.
	GetFilepath(id, hashPath string) string
	// Delete removes file from local storage.
	Delete(id, hashPath string) error
	// DeleteByName removes file or empty dir. It returns (false, nil) if dir is not empty.
	DeleteByName(name string) (bool, error)
	// GetEtag returns etag of file content.
	GetEtag(id, hashPath string) (string, error)
}

// EtagEncoder is used to generate etag based on file content.
type EtagEncoder interface {
	Encode([]byte) (string, error)
}
