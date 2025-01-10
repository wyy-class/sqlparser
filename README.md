# SQL Parser Master

SQL Parser Master is a project designed to parse SQL statements and process CSV files. It currently supports basic SQL statements such as `SELECT`, `INSERT INTO`, `UPDATE`, and `DELETE`. 

## Features

- **SQL Parsing**: Parses basic SQL statements.
- **CSV Processing**: Handles CSV file operations based on SQL commands.

## Limitations

- **Batch Commands**: Batch processing commands are not supported.
- **Concurrency**: Concurrent processing is not supported.
- **SQL Expansion**: SQL statement parsing is still being expanded.

## Usage

To use this project, simply run your SQL commands against your CSV files. Ensure that your SQL statements are simple and adhere to the supported commands.

## Example Usage

Here are some example tests to demonstrate how to use the SQL Parser Master with CSV files:

```go
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
```

These examples show how to perform various SQL operations such as `SELECT`, `UPDATE`, `INSERT INTO`, and `DELETE` on a CSV file using the SQL Parser Master.

## Contribution

Contributions are welcome! Please feel free to submit issues or pull requests.

## Acknowledgements

This project's SQL parser is based on: [https://github.com/marianogappa/sqlparser](https://github.com/marianogappa/sqlparser)


## Future Work

We plan to expand support for more SQL processing file formats in the future.

## License

This project is licensed under the MIT License.


