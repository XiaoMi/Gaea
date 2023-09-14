package util

import (
	"path/filepath"
	"runtime"
)

var basePath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath = filepath.Dir(currentFile)
}

var defaultLogDirectoryPath = "../cmd/logs"
var defaultLogFilePath = "../cmd/logs/gaea_sql.log"

var defaultE2ECaseDirectory = "../e2e/"
var defaultResultFilePath = "../e2e/noshard/result"

// Path gets the absolute path.
func Path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	return filepath.Join(basePath, rel)
}

func GetTestLogDirectoryAbsPath() string {
	return Path(defaultLogDirectoryPath)
}

func GetTestLogFileAbsPath() string {
	return Path(defaultLogFilePath)
}

func GetTestResultFileAbsPath() string {
	return Path(defaultResultFilePath)
}

func GetCasesFileAbsPath(filename string) string {
	return Path(defaultE2ECaseDirectory + filename)
}
