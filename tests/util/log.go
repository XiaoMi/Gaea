package util

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type LogModel struct {
	Namespace    string
	User         string
	ClientAddr   string
	BackendAddr  string
	ConnectionID int
	AffectedRows int
	SQL          string
}

func ReadLog(filepath string, searchString string) ([]LogModel, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// 创建一个slice来存储提取的数据
	var result []LogModel

	// 使用正则表达式匹配你需要的部分
	regex := regexp.MustCompile(`ns=(\w+),\s*(\w+)@(\d+\.\d+\.\d+\.\d+:\d+)->(\d+\.\d+\.\d+\.\d+:\d+), mysql_connect_id=(\d+), r=(\d+)\|(.*)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 检查该行是否包含特定的字符串
		if strings.Contains(line, searchString) {
			matches := regex.FindStringSubmatch(line)
			if len(matches) == 8 {
				// 解析ConnectionID为整数
				connectionID, err := strconv.Atoi(matches[5])
				if err != nil {
					return nil, err
				}
				affectedRows, err := strconv.Atoi(matches[6])
				if err != nil {
					return nil, err
				}
				logModel := LogModel{
					Namespace:    matches[1],
					User:         matches[2],
					ClientAddr:   matches[3],
					BackendAddr:  matches[4],
					ConnectionID: connectionID,
					AffectedRows: affectedRows,
					SQL:          matches[6],
				}
				result = append(result, logModel)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func RemoveLog(directory string) error {
	// 检查目录是否存在
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		// 如果目录不存在，则创建目录
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			return err
		}
	}
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			err := os.Remove(directory + "/" + file.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
