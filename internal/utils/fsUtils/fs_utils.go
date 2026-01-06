package fsUtils

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

// EnsureDir creates a directory (and parents) if not exists
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// DirExists checks if directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir()
	}

	return false
}

// RemoveDir removes a directory recursively
func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

// EmptyDir removes all contents but keeps the directory
func EmptyDir(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if err := os.RemoveAll(filepath.Join(path, e.Name())); err != nil {
			return err
		}
	}
	return nil
}

// FileExists checks whether a file exists
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}

	return false
}

// CreateFile creates a file (and parent dirs if needed)
func CreateFile(path string) (*os.File, error) {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return nil, err
	}
	return os.Create(path)
}

// WriteFile writes data to a file (overwrite)
func WriteFile(path string, data []byte, perm os.FileMode) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}

// AppendFile appends data to a file
func AppendFile(path string, data []byte, perm os.FileMode) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

// ReadFile reads entire file
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// DeleteFile removes a file
func DeleteFile(path string) error {
	return os.Remove(path)
}

// CopyFile copies file contents
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

// MoveFile moves or renames a file
func MoveFile(src, dst string) error {
	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}
	return os.Rename(src, dst)
}

// Join joins paths safely
func Join(parts ...string) string {
	return filepath.Join(parts...)
}

// Abs returns absolute path
func Abs(path string) (string, error) {
	return filepath.Abs(path)
}

// Base returns last element
func Base(path string) string {
	return filepath.Base(path)
}

// Ext returns file extension
func Ext(path string) string {
	return filepath.Ext(path)
}

// Chmod changes file permissions
func Chmod(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

// FileSize returns size in bytes
func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// ModTime returns file modification time
func ModTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}
