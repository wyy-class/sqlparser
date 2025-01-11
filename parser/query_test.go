package sqlparser

import (
	"log"
	"strings"
	"testing"

	. "github.com/marianogappa/sqlparser/yacc"
)

var cp *CSVParser

func init() {
	cp = initObject()
}
func initObject() *CSVParser {
	db := NewCSVDB("test.csv", "csv")
	parser := NewCSVParse(&Enginer{DB: db})
	return parser
}

/*
问题分析：
int型的数据不能使用string比较大小
例如 id='4'>id='10' 结果是true
解决方案:
目前只支持int类型'>','<','>=','<='的比较
*/
// 报错：输出跟debug不一样，文件内容修改了，但输出依然不变
// go test自带的缓存机制，导致文件内容修改后，读取的还是缓存的内容
// 解决方案：go test -count=1
func TestQuerySelect(t *testing.T) {
	sql := "select distinct * from test"
	result, err := cp.SelectCSV(sql)
	if err != nil {
		log.Fatalln("QueryCSV fail:", err)
		return
	}
	PrintToTable(result)
	cp.Close()
}
func TestQueryLike(t *testing.T) {
	sql := "select id as ID,age as AGE,name as NAME from test where name like '%a%'"
	result, err := cp.SelectCSV(sql)
	if err != nil {
		log.Fatalln("QueryCSV fail:", err)
		return
	}
	PrintToTable(result)
	cp.Close()
}
func TestQueryUpdate(t *testing.T) {
	sql := "update test set id='6' where id>='4'"
	resultNum, err := cp.UpdateCSV(sql)
	if err != nil {
		log.Fatalln("QueryCSV fail:", err)
		return
	}
	log.Print(resultNum)
	TestQuerySelect(t)
	cp.Close()
}

func TestQueryInsert(t *testing.T) {
	sql := "INSERT INTO 'test' (id,age,city) VALUES ('0','1','3' ),('4','5','6')"
	flag, err := cp.InsertCSV(sql)
	if err != nil {
		log.Fatalln("QueryCSV fail:", err)
		return
	}
	log.Print(flag)
	TestQuerySelect(t)
	cp.Close()
}
func TestQueryDelete(t *testing.T) {
	sql := "delete from test where id='0'"
	resultNum, err := cp.DeleteCSV(sql)
	if err != nil {
		log.Fatalln("QueryCSV fail:", err)
		return
	}
	log.Print(resultNum)
	TestQuerySelect(t)
	cp.Close()
}
func TestDemo1(t *testing.T) {
	log.Println(strings.ToUpper("like"))
	log.Println(strings.Compare("4", "10"))
}
