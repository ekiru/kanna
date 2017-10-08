package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

func FromRow(row *sql.Rows, cols map[string]interface{}) error {
	found := make(map[string]bool, len(cols))
	dsts := make([]interface{}, 0, len(cols))
	rowCols, err := row.Columns()
	if err != nil {
		return err
	}
	for _, col := range rowCols {
		if found[col] {
			return fmt.Errorf("duplicate column %q", col)
		}
		if dst, ok := cols[col]; ok {
			dsts = append(dsts, dst)
			found[col] = true
		} else {
			return fmt.Errorf("unexpected column %q", col)
		}
	}
	if len(dsts) != len(cols) {
		var missing []string
		for col, _ := range cols {
			if !found[col] {
				missing = append(missing, strconv.Quote(col))
			}
		}
		return fmt.Errorf("missing columns %s", strings.Join(missing, ", "))
	}
	return row.Scan(dsts...)
}
