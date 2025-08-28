package fileUtil

import (
	"os"
	"path/filepath"

	"radius-server/src/common/logger"
)

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func CreateDir(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func CreateFile(path string) (*os.File, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func OpenFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func WriteFile(path string, content string) error {
	err := os.WriteFile(path, []byte(content), 0o777)
	if err != nil {
		return err
	}
	return nil
}

func CloseFile(file *os.File, throwError bool) error {
	if err := file.Close(); err != nil {
		if throwError {
			return err
		} else {
			logger.Logger.Error().Msgf("Close file error. %s", err.Error())
			return nil
		}
	}
	return nil
}

func GetAbsolutePathFromRoot(relativePath string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		dir = filepath.Dir(dir)
	}
	return filepath.Join(dir, relativePath), nil
}
