package mysql

import "testing"

func TestGetClientHandshakeCollationID(t *testing.T) {
	// 假设 CharsetIds 和 DefaultCollationID 已经定义
	// CharsetIds = map[string]CollationID{
	// 	"utf8mb4": 45,
	// 	"latin1":  47,
	// }
	// DefaultCollationID = 46

	// 测试正常情况
	t.Run("Valid Charset", func(t *testing.T) {
		charset := "utf8mb4"
		expectedCollationID := CollationID(45)
		result := GetClientHandshakeCollationID(charset)
		if result != expectedCollationID {
			t.Errorf("Expected %d, but got %d", expectedCollationID, result)
		}
	})

	// 测试无效情况
	t.Run("Invalid Charset", func(t *testing.T) {
		charset := "unknown"
		expectedCollationID := CollationID(DefaultCollationID)
		result := GetClientHandshakeCollationID(charset)
		if result != expectedCollationID {
			t.Errorf("Expected %d, but got %d", expectedCollationID, result)
		}
	})

	// 测试边界情况
	t.Run("Empty Charset", func(t *testing.T) {
		charset := ""
		expectedCollationID := CollationID(DefaultCollationID)
		result := GetClientHandshakeCollationID(charset)
		if result != expectedCollationID {
			t.Errorf("Expected %d, but got %d", expectedCollationID, result)
		}
	})

	// 测试特殊字符串
	t.Run("Special Characters", func(t *testing.T) {
		charset := "!@#$%^&*()"
		expectedCollationID := CollationID(DefaultCollationID)
		result := GetClientHandshakeCollationID(charset)
		if result != expectedCollationID {
			t.Errorf("Expected %d, but got %d", expectedCollationID, result)
		}
	})
}
