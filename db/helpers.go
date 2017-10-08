package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// FromRow scans a row returned from a SQL query into destinations
// chosen based on column names. The sql.Rows object must have already
// had its Next method called. The cols map maps column names to the
// pointer into which the column should be scanned. FromRow returns an
// error if either the Columns or Scan methods on the sql.Rows object
// returns an error, if a column name appear twice in the columns
// returned by the query, if a column name that is not present in cols
// was returned by the query, or if a column name present in cols was
// not returned by the query.
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
