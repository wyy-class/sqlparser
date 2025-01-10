package sqlparser

import (
	"fmt"
	"log"
	"strings"
	"testing"

	. "github.com/marianogappa/sqlparser/lex"
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
func TestSQLParser(t *testing.T) {
	sqls := []string{"select a,b from test where id>'1' and name='wyy' or age<'20' and name like '%a%'"}
	querys, err := ParseMany(sqls)
	if err != nil {
		log.Fatalln("parseMany fail:", err)
		return
	}
	fmt.Println(querys[0].ToString())
}

/*
问题分析：
int型的数据不能使用string比较大小
例如 id='4'>id='10' 结果是true
解决方案:
目前只支持int类型'>','<','>=','<='的比较
*/
func TestQuerySelect(t *testing.T) {
	sql := "select * from test"
	result, err := cp.SelectCSV(sql)
	if err != nil {
		log.Fatalln("QueryCSV fail:", err)
		return
	}
	PrintToTable(result)
	cp.Close()
}
func TestQueryLike(t *testing.T) {
	sql := "select id,age,name from test where name like '%a%'"
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

/*
报错：第一次可以正常插入，第二次之后就插入文件不成功，但不报错
*/
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
