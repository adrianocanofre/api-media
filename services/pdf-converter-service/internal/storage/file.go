package storage

import (
	"io"
	"os"
	"path/filepath"
)

type Storage struct {
	DownloadsDir string
}

func New(downloadsDir string) *Storage {
	return &Storage{
		DownloadsDir: downloadsDir,
	}
}

func (s *Storage) EnsureDir() error {
	return os.MkdirAll(s.DownloadsDir, 0755)
}

func (s *Storage) SaveUploadedFile(file io.Reader, filename string) (string, error) {
	path := filepath.Join(s.DownloadsDir, filename)

	outFile, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	return path, err
}

func (s *Storage) StatFile(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
