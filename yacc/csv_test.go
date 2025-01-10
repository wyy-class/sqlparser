package yacc

import (
	"log"
	"os"
	"testing"
)

func TestMatchPattern(t *testing.T) {
	pattern := ".bc%"
	data := "abcdefg"
	log.Println(matchPattern(pattern, data))
	pattern = "%bc%"
	log.Println(matchPattern(pattern, data))
}
func TestCSVDB_Open(t *testing.T) {
	db := NewCSVDB("test.csv", "csv")
	err := db.Open()
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer db.file.Close()
}

func TestCSVDB_Write(t *testing.T) {
	db := NewCSVDB("test.csv", "csv")
	err := db.Open()
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer func() {
		db.file.Close()
	}()

	data := []byte("name,age\nJohn,30\nDoe,40")
	n, err := db.Write(data, os.O_TRUNC)
	if err != nil {
		t.Fatalf("Failed to write to CSV file: %v", err)
	}
	if n != len(data) {
		t.Fatalf("Expected to write %d bytes, wrote %d", len(data), n)
	}
}

func TestCSVDB_Read(t *testing.T) {
	db := NewCSVDB("test.csv", "csv")
	err := db.Open()
	if err != nil {
		t.Fatalf("Failed to open CSV file: %v", err)
	}
	defer func() {
		db.file.Close()
	}()

	data := []byte("name,age,city\nAlice,30,New York\nBob,25,Los Angeles\nCharlie,35,Chicago")
	_, err = db.Write(data, os.O_TRUNC)
	if err != nil {
		t.Fatalf("Failed to write to CSV file: %v", err)
	}

	records, err := db.Read()
	if err != nil {
		t.Fatalf("Failed to read from CSV file: %v", err)
	}

	expected := []string{
		"name,age,city",
		"Alice,30,New York",
		"Bob,25,Los Angeles",
		"Charlie,35,Chicago",
	}

	for i, record := range records {
		if record != expected[i] {
			t.Fatalf("Expected to read %s, read %s", expected[i], record)
		}
	}
}
