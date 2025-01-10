package sqlparser

import (
	"fmt"
	"log"

	. "github.com/marianogappa/sqlparser/lex"
	. "github.com/marianogappa/sqlparser/yacc"
)

type CSVParser struct {
	sqlParser *SQLParser
}

// inupt sql and csv path return csv data
func (cp *CSVParser) SelectCSV(sql string) (CSVTable, error) {
	query, err := Parse(sql)
	if err != nil {
		log.Fatalln("parseMany fail:", err)
	}
	return cp.sqlParser.SQLSelect(query)
}
func (cp *CSVParser) UpdateCSV(sql string) (int, error) {
	query, err := Parse(sql)
	if err != nil {
		log.Fatalln("parseMany fail:", err)
	}
	return cp.sqlParser.SQLUpdate(query)
}
func (cp *CSVParser) InsertCSV(sql string) (bool, error) {
	query, err := Parse(sql)
	if err != nil {
		log.Fatalln("parse fail:", err)
	}
	return cp.sqlParser.SQLInsert(query)
}
func (cp *CSVParser) DeleteCSV(sql string) (int, error) {
	query, err := Parse(sql)
	if err != nil {
		log.Fatalln("parse fail:", err)
	}
	return cp.sqlParser.SQLDelete(query)
}
func NewCSVParse(enginer *Enginer) *CSVParser {
	err := enginer.DB.(*CSVDB).Open()
	if err != nil {
		log.Fatalln("csvDB open fail.")
	}
	return &CSVParser{sqlParser: NewSQLParser(enginer)}
}
func (cp *CSVParser) Close() {
	cp.sqlParser.Close()
}
func PrintToTable(table CSVTable) {
	for _, v := range table.Headers {
		fmt.Printf("%v\t", v)
	}
	fmt.Println()
	for _, row := range table.Rows {
		for _, v := range row {
			fmt.Printf("%v\t", v)
		}
		fmt.Println()
	}
}
