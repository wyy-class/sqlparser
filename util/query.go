package util

import "fmt"

// Query represents a parsed query
type Query struct {
	Type       Type
	TableName  string
	Conditions []Condition
	Updates    map[string]string
	Inserts    [][]string
	Fields     []string // Used for SELECT (i.e. SELECTed field names) and INSERT (INSERTEDed field names)
	Aliases    map[string]string
}

func (q Query) ToString() string {
	ss := make([]byte, 0)
	ss = append(ss, fmt.Sprintf("Type: %v\n", q.Type)...)
	ss = append(ss, fmt.Sprintf("TableName: %v\n", q.TableName)...)
	ss = append(ss, fmt.Sprintf("Conditions: %v\n", q.Conditions)...)
	ss = append(ss, fmt.Sprintf("Updates: %v\n", q.Updates)...)
	ss = append(ss, fmt.Sprintf("Inserts: %v\n", q.Inserts)...)
	ss = append(ss, fmt.Sprintf("Fields: %v\n", q.Fields)...)
	ss = append(ss, fmt.Sprintf("Aliases: %v\n", q.Aliases)...)
	return string(ss)
}

// Type is the type of SQL query, e.g. SELECT/UPDATE
type Type int

const (
	// UnknownType is the zero value for a Type
	UnknownType Type = iota
	// Select represents a SELECT query
	Select
	// Update represents an UPDATE query
	Update
	// Insert represents an INSERT query
	Insert
	// Delete represents a DELETE query
	Delete
)

// TypeString is a string slice with the names of all types in order
var TypeString = []string{
	"UnknownType",
	"Select",
	"Update",
	"Insert",
	"Delete",
}

// Operator is between operands in a condition
type Operator int

const (
	// UnknownOperator is the zero value for an Operator
	UnknownOperator Operator = iota
	// Eq -> "="
	Eq
	// Ne -> "!="
	Ne
	// Gt -> ">"
	Gt
	// Lt -> "<"
	Lt
	// Gte -> ">="
	Gte
	// Lte -> "<="
	Lte
	// Lk -> "LIKE"
	Lk
	// AND -> "AND"
	AND
	// OR -> "OR"
	OR
)

// OperatorString is a string slice with the names of all operators in order
var OperatorString = []string{
	"UnknownOperator",
	"Eq",
	"Ne",
	"Gt",
	"Lt",
	"Gte",
	"Lte",
	"AND",
	"OR",
}

// Condition is a single boolean condition in a WHERE clause
type Condition struct {
	// Operand1 is the left hand side operand
	Operand1 string
	// Operand1IsField determines if Operand1 is a literal or a field name
	Operand1IsField bool
	// Operator is e.g. "=", ">"
	Operator Operator
	// Operand1 is the right hand side operand
	Operand2 string
	// Operand2IsField determines if Operand2 is a literal or a field name
	Operand2IsField bool
	LogicalOperator Operator
}
