package yacc

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	. "github.com/marianogappa/sqlparser/util"
)

// sql解析器
type SQLParser struct {
	enginer *Enginer
}

// CSVQuery converts a query to a CSV query
func (sp *SQLParser) SQLSelect(query Query) (CSVTable, error) {
	if sp.enginer == nil {
		log.Fatalln("enginer is nil!")
	}
	log.Println(query.ToString())
	switch query.Type {
	case Select:
		return CSVSelect(sp.enginer, query)
	default:
		return CSVTable{}, fmt.Errorf("unknown    V query type: %v", query.Type)
	}
}
func (sp *SQLParser) SQLUpdate(query Query) (int, error) {
	if sp.enginer == nil {
		log.Fatalln("enginer is nil")
	}
	log.Println(query.ToString())
	switch query.Type {
	case Update:
		return CSVUpdate(sp.enginer, query)
	default:
		return -1, fmt.Errorf("unknown query type: %v", query.Type)
	}
}
func (sp *SQLParser) SQLInsert(query Query) (bool, error) {
	if sp.enginer == nil {
		log.Fatalln("enginer is nil")
	}
	log.Println(query.ToString())
	switch query.Type {
	case Insert:
		return CSVInsert(sp.enginer, query)
	default:
		return false, fmt.Errorf("unknown query type: %v", query.Type)
	}
}
func (sp *SQLParser) SQLDelete(query Query) (int, error) {
	if sp.enginer == nil {
		log.Fatalln("enginer is nil")
	}
	log.Println(query.ToString())
	switch query.Type {
	case Delete:
		return CSVDelete(sp.enginer, query)
	default:
		return -1, fmt.Errorf("unknown query type: %v", query.Type)
	}
}
func (sp *SQLParser) Close() {
	sp.enginer.DB.Close()
}
func NewSQLParser(enginer *Enginer) *SQLParser {
	return &SQLParser{enginer: enginer}
}

type CSVTable struct {
	Headers []string
	Rows    [][]string
}

func readData(enginer *Enginer) (CSVTable, error) {
	result, err := enginer.DB.Read()
	if err != nil {
		return CSVTable{}, fmt.Errorf("Failed to read from CSV file: %v", err)
	}
	// 解析 CSV 数据
	var headers []string
	var resultRows [][]string
	for i, line := range result {
		fields := strings.Split(line, ",")
		if i == 0 {
			headers = fields
		} else {
			resultRows = append(resultRows, fields)
		}
	}
	return CSVTable{Headers: headers, Rows: resultRows}, nil
}
func readHeader(db *CSVDB) ([]string, error) {
	reader := csv.NewReader(db.file)
	record, err := reader.Read()
	if err != nil {
		return nil, err
	}
	return record, nil
}
func CSVSelect(enginer *Enginer, query Query) (CSVTable, error) {
	tb, err := readData(enginer)
	if err != nil {
		return CSVTable{}, err
	}
	headers := tb.Headers
	resultRows := tb.Rows

	//根据where条件过滤数据
	log.Println("befor where filter...")
	log.Println("header:", headers)
	log.Println("rows:", resultRows)
	if len(query.Conditions) > 0 {
		_, resultRows = filterRows(headers, resultRows, query.Conditions)
	}
	log.Println("after where filter")
	log.Println("condition:", query.Conditions)
	log.Println("Condition filter:", resultRows)

	// 根据fileds过滤列数据
	log.Println("field filter")
	if !(len(query.Fields) == 1 && query.Fields[0] == "*") {
		resultRows = filterData(headers, resultRows, query.Fields)
		headers = query.Fields
	}
	log.Println("field filter:", resultRows)
	// 处理别名
	if len(query.Aliases) > 0 {
		headers = replaceHeaders(headers, query.Aliases)
	}
	// 处理distinct
	if query.Distinct {
		resultRows = distinct(resultRows)
	}
	//返回结果
	return CSVTable{Headers: headers, Rows: resultRows}, nil
}

// filterData filters the rows based on the specified fields
func filterData(headers []string, rows [][]string, fields []string) [][]string {
	var resultRows [][]string
	for _, row := range rows {
		var resultRow []string
		for _, field := range fields {
			for i, header := range headers {
				if header == field {
					// log.Println("index", i)
					resultRow = append(resultRow, row[i])
				}
			}
		}
		resultRows = append(resultRows, resultRow)
	}
	return resultRows
}

// Condition represents a condition in a query
func filterRows(header []string, rows [][]string, conditions []Condition) ([]int, [][]string) {
	var resultRows [][]string
	var idxs []int
	for idx, row := range rows {
		if evalConditions(header, row, conditions) {
			resultRows = append(resultRows, row)
			idxs = append(idxs, idx)
		}
	}
	return idxs, resultRows
}

// evalConditions evaluates the conditions for a row
func evalConditions(header, row []string, conditions []Condition) bool {
	var result bool
	for _, condition := range conditions {
		tmpflag := false
		// 获取字段索引
		var index int
		// log.Println("header:", header)
		// log.Println("condition:", condition)
		for i, h := range header {
			if condition.Operand1IsField && h == condition.Operand1 {
				index = i
				break
			}
		}
		// 获取字段值
		value := row[index]
		// 判断条件
		switch condition.Operator {
		case Eq:
			if !condition.Operand2IsField && value == condition.Operand2 {
				tmpflag = true
			}
		case Ne:
			if !condition.Operand2IsField && value != condition.Operand2 {
				tmpflag = true
			}
		case Gt:
			value_int, err := strconv.Atoi(value)
			op_int, err := strconv.Atoi(condition.Operand2)
			if err != nil {
				log.Fatalf("strocnv.Atoi: %v %v", value, condition.Operator)
			}
			if !condition.Operand2IsField && value_int > op_int {
				tmpflag = true
			}
		case Lt:
			value_int, err := strconv.Atoi(value)
			op_int, err := strconv.Atoi(condition.Operand2)
			if err != nil {
				log.Fatalf("strocnv.Atoi: %v %v", value, condition.Operator)
			}
			if !condition.Operand2IsField && value_int < op_int {
				tmpflag = true
			}
		case Gte:
			value_int, err := strconv.Atoi(value)
			op_int, err := strconv.Atoi(condition.Operand2)
			if err != nil {
				log.Fatalf("strocnv.Atoi: %v %v", value, condition.Operator)
			}
			if !condition.Operand2IsField && value_int >= op_int {
				tmpflag = true
			}
		case Lte:
			value_int, err := strconv.Atoi(value)
			op_int, err := strconv.Atoi(condition.Operand2)
			if err != nil {
				log.Fatalf("strocnv.Atoi: %v %v", value, condition.Operator)
			}
			if !condition.Operand2IsField && value_int <= op_int {
				tmpflag = true
			}
		case Lk:
			if !condition.Operand2IsField && containString(value, condition.Operand2) {
				tmpflag = true
			}
		default:
			log.Fatalf("Unknown operator: %v", condition.Operator)
		}
		// 判断AND OR
		switch condition.LogicalOperator {
		case AND:
			result = result && tmpflag
		case OR:
			result = result || tmpflag
		default:
			result = tmpflag
		}
	}
	return result
}
func containString(value string, contition string) bool {
	// 判断是否为模糊查询
	if strings.Contains(contition, ".") || strings.Contains(contition, "%") {
		return matchPattern(contition, value)
	} else {
		return value == contition
	}
}

// matchPattern 函数用于判断字符串是否匹配给定的模式
func matchPattern(pattern, str string) bool {
	// 处理模式中的特殊字符
	pattern = strings.ReplaceAll(pattern, ".", ".{1}")
	pattern = strings.ReplaceAll(pattern, "%", ".*")

	// 使用 regexp.QuoteMeta 转义其他特殊字符
	regexPattern := "^" + regexp.QuoteMeta(pattern) + "$"
	regexPattern = strings.ReplaceAll(regexPattern, "\\.\\{1\\}", ".{1}")
	regexPattern = strings.ReplaceAll(regexPattern, "\\.\\*", ".*")

	// log.Println(regexPattern)
	// 编译正则表达式
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return false
	}
	// 使用正则表达式匹配字符串
	return re.MatchString(str)
}

func replaceHeaders(headers []string, aliases map[string]string) []string {
	for i, header := range headers {
		if alias, ok := aliases[header]; ok {
			headers[i] = alias
		}
	}
	return headers
}

func distinct(rows [][]string) [][]string {
	var result [][]string
	m := make(map[string]bool)
	for _, row := range rows {
		key := strings.Join(row, ",")
		if _, ok := m[key]; !ok {
			m[key] = true
			result = append(result, row)
		}
	}
	return result
}

func CSVUpdate(enginer *Enginer, query Query) (int, error) {
	tb, err := readData(enginer)
	if err != nil {
		return -1, err
	}
	headers := tb.Headers
	resultRows := tb.Rows
	var idxs []int
	//根据where条件过滤数据
	log.Println("befor where filter...")
	log.Println("header:", headers)
	log.Println("rows:", resultRows)
	if len(query.Conditions) > 0 {
		idxs, resultRows = filterRows(headers, resultRows, query.Conditions)
	}
	log.Println("after where filter")
	log.Println("condition:", query.Conditions)
	log.Println("Condition filter:", resultRows)

	//根据update字段更新数据
	log.Println("update filter")
	log.Println("field filter:", resultRows)
	UpdateData(headers, resultRows, query.Updates)
	//更新csv文件数据
	log.Println("index:", idxs)
	log.Println("update filter:", resultRows)
	updateCSV(enginer, &tb, idxs, resultRows)
	return len(resultRows), nil
}

func UpdateData(headers []string, rows [][]string, updates map[string]string) {
	for _, row := range rows {
		for field, value := range updates {
			for i, header := range headers {
				if header == field {
					row[i] = value
				}
			}
		}
	}
}
func updateCSV(enginer *Enginer, table *CSVTable, idxs []int, rows [][]string) {
	headers, resultRows := table.Headers, table.Rows
	// 更新数据
	for i, idx := range idxs {
		resultRows[idx] = rows[i]
	}
	log.Println("update data:", resultRows)
	// 写入数据
	writeCSVData(enginer.DB, headers, resultRows, os.O_TRUNC)
}
func writeCSVData(db DB, headers []string, rows [][]string, tp int) {
	var data []string
	if tp == os.O_TRUNC {
		data = append(data, strings.Join(headers, ","))
	}
	for _, row := range rows {
		data = append(data, strings.Join(row, ","))
	}
	_, err := db.Write([]byte(strings.Join(data, "\n")), tp)
	if err != nil {
		log.Fatalf("Failed to write to CSV file: %v", err)
	}
}
func parseCSVData(result []string) ([]string, [][]string, error) {
	var headers []string
	var resultRows [][]string
	for i, line := range result {
		fields := strings.Split(line, ",")
		if i == 0 {
			headers = fields
		} else {
			resultRows = append(resultRows, fields)
		}
	}
	return headers, resultRows, nil
}

// CSVInsert executes an INSERT query on the CSV file
func CSVInsert(enginer *Enginer, query Query) (bool, error) {
	headers, err := readHeader(enginer.DB.(*CSVDB))
	if err != nil {
		return false, fmt.Errorf("Failed to read headers from CSV file: %v", err)
	}
	//插入数据
	resultRows := insertData(headers, query.Fields, query.Inserts)
	//更新csv文件数据
	writeCSVData(enginer.DB, headers, resultRows, os.O_APPEND)
	return true, nil
}
func insertData(headers, fields []string, inserts [][]string) [][]string {
	var resultRows [][]string
	for _, insert := range inserts {
		tmp := make([]string, len(headers))
		for i, field := range fields {
			for j, h := range headers {
				if field == h {
					tmp[j] = insert[i]
					break
				}
			}
		}
		resultRows = append(resultRows, tmp)
	}
	return resultRows
}
func CSVDelete(enginer *Enginer, query Query) (int, error) {
	tb, err := readData(enginer)
	if err != nil {
		return -1, err
	}
	headers := tb.Headers
	resultRows := tb.Rows
	var idxs []int
	//根据where条件过滤数据
	log.Println("befor where filter...")
	log.Println("header:", headers)
	log.Println("rows:", resultRows)
	if len(query.Conditions) > 0 {
		idxs, resultRows = filterRows(headers, resultRows, query.Conditions)
	}
	log.Println("after where filter")
	log.Println("condition:", query.Conditions)
	log.Println("Condition filter:", resultRows)
	//删除数据
	log.Println("delete filter")
	log.Println("field filter:", resultRows)
	tb.Rows = deleteData(tb.Rows, idxs)
	// 删除csv文件数据
	log.Println("index:", idxs)
	log.Println("delete filter:", resultRows)
	writeTable(enginer.DB, tb)
	return len(idxs), nil
}
func deleteData(rows [][]string, idxs []int) [][]string {
	var result [][]string
	for i, row := range rows {
		var flag bool = false
		for _, idx := range idxs {
			if i == idx {
				flag = true
				break
			}
		}
		if !flag {
			result = append(result, row)
		}
	}
	return result
}
func writeTable(db DB, table CSVTable) {
	headers, rows := table.Headers, table.Rows
	writeCSVData(db, headers, rows, os.O_TRUNC)
}
