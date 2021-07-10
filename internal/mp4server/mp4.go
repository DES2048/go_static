package mp4server

import (
	"os"
	"path/filepath"
	"time"

	"github.com/alfg/mp4"
)

type FileInfo struct {
	os.FileInfo
	Duration time.Duration
}

func ListMp4(path string) ([]FileInfo, error) {
	absPath, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	globPath := filepath.Join(absPath, "*.mp4")

	paths, err := filepath.Glob(globPath)

	if err != nil {
		return nil, err
	}

	infos := make([]FileInfo, 0, len(paths))

	// gather file info's
	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			return nil, err
		}

		// get mp4 duration
		mp4, err := mp4.Open(path)

		if err != nil {
			return nil, err
		}

		durSeconds := mp4.Moov.Mvhd.Duration / mp4.Moov.Mvhd.Timescale
		dur := time.Duration(durSeconds) * time.Second

		infos = append(infos, FileInfo{
			FileInfo: info,
			Duration: dur,
		})
	}

	return infos, nil
}
