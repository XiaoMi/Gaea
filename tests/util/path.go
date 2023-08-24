package util

import (
	"path/filepath"
	"runtime"
)

var basepath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

var defaulTrelLogDirectoryPath = "../cmd/logs"
var defaulTrelLogFilePath = "../cmd/logs/gaea_sql.log"

var defaulSqlCaseDirectory = "../sql/"

// Path gets the absolute path.
func Path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	return filepath.Join(basepath, rel)
}

func GetTestLogDirectoryAbsPath() string {
	return Path(defaulTrelLogDirectoryPath)
}

func GetTestLogFileAbsPath() string {
	return Path(defaulTrelLogFilePath)
}

func GetCasesFileAbsPath(filename string) string {
	return Path(defaulSqlCaseDirectory + filename)
}
