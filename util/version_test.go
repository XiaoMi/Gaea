package util

import (
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompactServerVersion(t *testing.T) {
	type tCase struct {
		oVersion string
		dVersion string
	}

	tts := []tCase{
		{
			"",
			mysql.ServerVersion,
		},
		{
			"5.7-gaea",
			mysql.ServerVersion,
		},
		{
			"5.7.11-gaea",
			"5.7.11-gaea",
		},
		{
			"5.7.11",
			"5.7.11",
		},
		{
			"8.0.0-gaea",
			"8.0.0-gaea",
		},
	}
	for _, tt := range tts {
		assert.Equal(t, CompactServerVersion(tt.oVersion), tt.dVersion)
	}
}

func TestCheckMySQLVersion(t *testing.T) {
	type tCase struct {
		version  string
		lessThan string
		isTrue   bool
	}

	tts := []tCase{{
		"5.7.20-gaea",
		"< 5.7.25",
		true,
	},
		{
			"5.7.20-gaea",
			">= 5.7.10",
			true,
		},
		{
			"5.7.11-gaea",
			"< 5.7.25",
			true,
		},
		{
			"5.6.20-gaea",
			"< 5.7.25",
			true,
		},
		{
			"5.6.20-gaea",
			"< 8.0.0",
			true,
		},
		{
			"5.6.20",
			"< 8.0.0",
			true,
		},
	}
	for _, tt := range tts {
		assert.Equal(t, CheckMySQLVersion(tt.version, tt.lessThan), tt.isTrue)
	}
}
