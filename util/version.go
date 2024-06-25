package util

import (
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/hashicorp/go-version"
	"strings"
)

const (
	LessThanMySQLVersion80  = "< 8.0"
	LessThanMySQLVersion803 = "< 8.0.3"
)

// CompactServerVersion CompactVersion get compact version from origin string
func CompactServerVersion(sv string) string {
	version := strings.Trim(sv, " ")
	if version != "" {
		v := strings.Split(sv, ".")
		if len(v) < 3 {
			return mysql.ServerVersion
		}
		return version
	} else {
		return mysql.ServerVersion
	}
}

// CheckMySQLVersion check if curVersion satisfies constraintVersion, which should be format like ">= 5.7, < 8.0"
func CheckMySQLVersion(curVersion string, constraintVersion string) bool {
	if strings.HasSuffix(curVersion, "-gaea") {
		curVersion = curVersion[:len(curVersion)-5]
	} else {
		curVersion = strings.Split(curVersion, "-")[0]
	}

	cur, err := version.NewVersion(curVersion)
	if err != nil {
		return false
	}

	constraints, err := version.NewConstraint(constraintVersion)
	if err != nil {
		return false
	}

	if constraints.Check(cur) {
		return true
	}
	return false
}
