package util

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp      string
	Namespace      string
	User           string
	ClientAddr     string
	BackendAddr    string
	Database       string
	ConnectionID   int
	Query          string
	ResponseTimeMs float64
}

// CompareTimeStrings 比较两个时间字符串的大小
// 返回值为-1，0或1。-1表示time Before time2，0表示time1 = time2，1表示time1 After time2
func CompareTimeStrings(currentTime string, time2 string) (int, error) {
	// 解析时间字符串
	t1, err1 := time.Parse("2006-01-02 15:04:05.999", currentTime)
	t2, err2 := time.Parse("2006-01-02 15:04:05.999", time2)
	if err1 != nil || err2 != nil {
		return 0, fmt.Errorf("解析错误：%v %v", err1, err2)
	}

	// 比较时间
	if t1.Before(t2) {
		return -1, nil
	}
	if t1.After(t2) {
		return 1, nil
	}
	return 0, nil
}

func ParseLogEntries(file *os.File, re *regexp.Regexp, currentTime time.Time, searchString string) ([]LogEntry, error) {
	startTime := currentTime.Format("2006-01-02 15:04:05.999")

	scanner := bufio.NewScanner(file)

	var logEntryRes []LogEntry
	for scanner.Scan() {
		line := scanner.Text()
		// 使用正则表达式匹配日志行
		matches := re.FindStringSubmatch(line)
		if len(matches) != 11 {
			continue
		}
		// 解析并填充结构体
		logEntry := LogEntry{}
		logEntry.Timestamp = matches[1]
		// 检查时间是否在startTime之后
		res, err := CompareTimeStrings(startTime, logEntry.Timestamp)
		if err != nil {
			return []LogEntry{}, nil
		}
		if res != -1 {
			continue
		}
		fmt.Sscanf(matches[3], "%f", &logEntry.ResponseTimeMs)
		logEntry.Namespace = matches[4]
		logEntry.User = matches[5]
		logEntry.ClientAddr = matches[6]
		logEntry.BackendAddr = matches[7]
		logEntry.Database = matches[8]
		fmt.Sscanf(matches[9], "%d", &logEntry.ConnectionID)
		logEntry.Query = matches[10]

		if strings.Compare(searchString, logEntry.Query) != 0 {
			continue
		}
		logEntryRes = append(logEntryRes, logEntry)
	}

	if err := scanner.Err(); err != nil {
		return logEntryRes, fmt.Errorf("error during file scanning:%v", err)
	}
	return logEntryRes, nil
}
