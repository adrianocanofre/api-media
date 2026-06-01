package storage

import (
	"os"
)

type Storage struct {
	DownloadsDir string
}

func New(downloadsDir string) *Storage {
	return &Storage{
		DownloadsDir: downloadsDir,
	}
}

func (s *Storage) StatFile(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
