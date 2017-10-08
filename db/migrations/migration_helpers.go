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

type ColumnType int

const (
	String ColumnType = iota
	Text
	Int
	Float
	Bool
	Timestamp
)

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

type Column struct {
	Name          string
	Type          ColumnType
	Default       *string
	NotNull       bool
	PrimaryKey    bool
	AutoIncrement bool
}

type createTable struct {
	idHelper
	name    string
	columns []Column
}

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
