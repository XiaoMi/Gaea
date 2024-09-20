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

type VersionCompareStatus struct {
	LessThanMySQLVersion80  bool
	LessThanMySQLVersion803 bool
}

func NewVersionCompareStatus(version string) *VersionCompareStatus {
	v := &VersionCompareStatus{}
	return v.SetVersion(version)
}

// 用于需保留版本状态时
func (v *VersionCompareStatus) SetVersion(srcVersion string) *VersionCompareStatus {
	if strings.HasSuffix(srcVersion, "-gaea") {
		srcVersion = srcVersion[:len(srcVersion)-5]
	} else {
		srcVersion = strings.Split(srcVersion, "-")[0]
	}
	src, err := version.NewVersion(srcVersion)
	if err != nil {
		return v
	}
	for _, compareVersion := range []string{LessThanMySQLVersion80, LessThanMySQLVersion803} {
		constraints, err := version.NewConstraint(compareVersion)
		if err != nil || !constraints.Check(src) {
			continue
		}
		switch compareVersion {
		case LessThanMySQLVersion80:
			v.LessThanMySQLVersion80 = true
		case LessThanMySQLVersion803:
			v.LessThanMySQLVersion803 = true
		}
	}
	return v
}

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
