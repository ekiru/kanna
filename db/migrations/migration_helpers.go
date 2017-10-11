// The migrations package defines several types implementing the
// db.Migration interface.
package migrations

import (
	"bytes"
	"fmt"

	"github.com/ekiru/kanna/db"
)

type idHelper struct {
	id string
}

func (id idHelper) ID() string {
	return id.id
}

// A ColumnType identifies a column type when creating tables or adding
// or modifying columns.
type ColumnType int

const (
	String ColumnType = iota
	Text
	Int
	Float
	Bool
	Timestamp
)

// String converts a ColumnType to a string.
func (t ColumnType) String() string {
	switch t {
	case String, Text:
		return "text"
	case Int:
		return "int"
	case Float:
		return "real"
	case Bool:
		return "boolean"
	case Timestamp:
		return "int"
	default:
		return "" // TODO maybe panic or something
	}
}

// A Column represents a database column.
type Column struct {
	// Name is the name of the column.
	Name string
	// Type represents the type of the column and may be translated
	// to a database-specific type name when applying a Migration.
	Type ColumnType
	// Default, if non-nil, contains a string that will be used as
	// a SQL expression to define the default value for a column.
	Default *string
	// NotNull, if true, specifies that a column may not contain
	// null values.
	NotNull bool
	// PrimaryKey, if true, specifies that the column is the
	// primary key of the table.
	PrimaryKey bool
	// AutoIncrement, if true, requests that the database
	// automatically supply incrementing values for a primary key
	// column.
	AutoIncrement bool
	// Unique, if true, specifies that no two rows in the table can
	// hold the same value.
	Unique bool
}

type createTable struct {
	idHelper
	name    string
	columns []Column
}

// CreateTable defines a migration that creates a table.
func CreateTable(id string, name string, cols_and_options ...interface{}) db.Migration {
	var mi createTable
	mi.idHelper.id = id
	mi.name = name
	for _, col := range cols_and_options {
		mi.columns = append(mi.columns, col.(Column))
	}
	return &mi
}

func (mi *createTable) Up(tx db.MigrationTx) {
	var stmt bytes.Buffer
	stmt.WriteString("create table ")
	stmt.WriteString(mi.name)
	stmt.WriteString(" (")
	sep := ""
	for _, col := range mi.columns {
		stmt.WriteString(sep)
		stmt.WriteString(col.Name)
		stmt.WriteString(" ")
		stmt.WriteString(col.Type.String())
		if col.Default != nil {
			stmt.WriteString(" default ")
			stmt.WriteString(*col.Default)
		}
		if col.PrimaryKey {
			stmt.WriteString(" primary key")
			if col.AutoIncrement {
				stmt.WriteString(" autoincrement")
			}
		}
		if col.Unique {
			stmt.WriteString(" unique")
		}
		if col.NotNull {
			stmt.WriteString(" not null")
		}
		sep = ", "
	}
	stmt.WriteString(")")
	tx.Exec(stmt.String())
}

func (mi *createTable) Down(tx db.MigrationTx) {
	tx.Exec(fmt.Sprintf("drop table %s", mi.name))
}

// A FreeForm Migration allows specifying arbitrary SQL queries for the
// Up and Down directions.
type FreeForm struct {
	Identifier       string
	Upward, Downward func(tx db.MigrationTx)
}

// ID returns the Identifier field of the FreeForm Migration.
func (mi FreeForm) ID() string {
	return mi.Identifier
}

// Up simply calls the Upward field of the FreeForm Migration.
func (mi FreeForm) Up(tx db.MigrationTx) {
	mi.Upward(tx)
}

// Down simply calls the Downward field of the FreeForm Migration.
func (mi FreeForm) Down(tx db.MigrationTx) {
	mi.Downward(tx)
}
