package yacc

import (
	"encoding/csv"
	"os"
	"strings"

	. "github.com/marianogappa/sqlparser/util"
)

type Enginer struct {
	DB DB
}
type CSVDB struct {
	path     string
	filetype string
	file     *os.File
}

func NewCSVDB(path string, filetype string) *CSVDB {
	return &CSVDB{path: path, filetype: filetype}
}
func (c *CSVDB) Open() error {
	file, err := os.OpenFile(c.path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	c.file = file
	return nil
}
func (c *CSVDB) Write(p []byte, tp int) (n int, err error) {
	var writer *csv.Writer
	if tp == os.O_APPEND {
		// 重新打开文件，并使用 os.O_APPEND 标志追加数据
		file, err := os.OpenFile(c.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			return 0, err
		}
		defer file.Close()
		writer = csv.NewWriter(file)
	}
	if tp == os.O_TRUNC {
		// 重新打开文件，并使用 os.O_TRUNC 标志截断文件
		file, err := os.OpenFile(c.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return 0, err
		}
		defer file.Close()
		writer = csv.NewWriter(file)
	}
	// Convert byte slice to string and split by new lines
	lines := string(p)
	records := [][]string{}
	for _, line := range strings.Split(lines, "\n") {
		records = append(records, strings.Split(line, ","))
	}
	writer.WriteAll(records)
	writer.Flush()
	return len(p), nil
}
func (c *CSVDB) Read() ([]string, error) {
	// Ensure the file is flushed and seek to the beginning before reading
	if err := c.file.Sync(); err != nil {
		return nil, err
	}
	if _, err := c.file.Seek(0, 0); err != nil {
		return nil, err
	}

	reader := csv.NewReader(c.file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var result []string
	for _, record := range records {
		result = append(result, strings.Join(record, ","))
	}

	return result, nil
}

func (c *CSVDB) Close() error {
	return c.file.Close()
}
