package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// IsImage returns true if mediaType is valid mime type and mime type represents image.
// If mediaType is invalid mime type, mime type is defined by filename.
func IsImage(mediaType, filename string) bool {
	mt, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		ext := filepath.Ext(filename)
		if ext != "" {
			mt, _, _ = mime.ParseMediaType(mime.TypeByExtension(ext))
		}
	}

	if mt == "" {
		return false
	}

	if mtp := strings.Split(mt, "/"); len(mtp) > 0 {
		return mtp[0] == "image"
	}

	return false
}

// GetHashPath returns hash path to store file.
func GetHashPath() string {
	return string('a' + rune(rand.Intn(26))) //nolint:gosec
}

// MoveFile moves file from src to dest.
func MoveFile(src, dest string) error {
	// Try to use rename operation.
	err := os.Rename(src, dest)
	if err == nil {
		return nil
	} else {
		var linkError *os.LinkError
		if errors.As(err, &linkError) {
			srcStat, err := os.Stat(src)
			if err != nil {
				return fmt.Errorf("link error %s %w; %s",
					linkError.Op, linkError.Err, dest)
			}

			return fmt.Errorf("link error %s %w; %s; %#v",
				linkError.Op, linkError.Err, dest, srcStat)
		}
	}

	// Copy file content manually if rename operation fails.
	err = CopyFile(src, dest)
	if err != nil {
		return err
	}

	err = os.Remove(src)
	if err != nil {
		return err
	}

	return nil
}

// CopyFile copies file from src to dest.
func CopyFile(src, dest string) (resErr error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer func() {
		err := srcFile.Close()
		if err != nil && resErr == nil {
			resErr = err
		}
	}()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer func() {
		err := destFile.Close()
		if err != nil && resErr == nil {
			resErr = err
		}
	}()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

// DirHasReadableFiles returns false when input directories don't have any file that can be open for read
func DirHasReadableFiles(ctx context.Context, paths []string) bool {
	result := false
	for _, p := range paths {
		fi, err := os.Stat(p)
		if err != nil || !fi.IsDir() {
			continue
		}
		files, err := os.ReadDir(p)
		if err != nil {
			continue
		}
		for _, dirEntry := range files {
			if dirEntry.IsDir() {
				continue
			}
			file, err := os.Open(filepath.Join(p, dirEntry.Name()))
			if err != nil {
				continue
			}
			err = file.Close()
			if err == nil {
				result = true
				break
			}
			select {
			case <-ctx.Done():
				return result
			default:
			}
		}
		if result {
			break
		}
		select {
		case <-ctx.Done():
			return result
		default:
		}
	}
	return result
}

func Sanitize(filename string) string {
	re := regexp.MustCompile(`[/\\?%*:|"<>]`)
	return re.ReplaceAllString(filename, "-")
}
