package db

import "database/sql"

// MigrationTx wraps the database/sql package's Tx type to provide a
// more restricted interface when running migrations.
type MigrationTx struct {
	tx *sql.Tx
}

// Exec executes a query without returning any rows by calling the
// underlying Tx object's Exec method.
func (tx MigrationTx) Exec(q string, vs ...interface{}) {
	_, err := tx.tx.Exec(q, vs...)
	if err != nil {
		panic(err)
	}
}

// A Migration performs some reversible change to the database.
type Migration interface {
	// The ID identifies a particular migration and distinguishes it
	// from other migrations in order to avoid performing the
	// migration repeatedly or attempting to undo a migration that
	// was never performed.
	ID() string
	// Up performs the migration, creating a table, changing a
	// column type, or performing some other database schema
	// operation.
	Up(MigrationTx)
	// Down undoes the migration, for example by deleting a table
	// or reversing a change to a column.
	Down(MigrationTx)
}

// ApplyMigration performs a migration and marks it as having been
// performed in the Migrations table. All changes are performed in a
// transaction: any error will cause all changes to be rolled back and
// all changes should occur atomically. If the Migration has already
// been performed, then ApplyMigration does nothing and returns a nil
// error.
func ApplyMigration(db *sql.DB, mi Migration) (err error) {
	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			err = recovered.(error) // TODO check this
		}
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	if alreadyApplied, err := wasMigrationApplied(tx, mi); err != nil {
		return err
	} else if alreadyApplied {
		return nil
	}
	mi.Up(MigrationTx{tx})
	_, err = tx.Exec("insert into Migrations (id) values (?)", mi.ID())
	if err != nil {
		return nil
	}
	return
}

// UndoMigration undoes a migration and marks it as not having been
// applied. Similarly to ApplyMigration, all changes are performed in a
// transaction. The changes will only be performed if the migration has
// previously been applied.
func UndoMigration(db *sql.DB, mi Migration) (err error) {
	var tx *sql.Tx
	tx, err = db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			err = recovered.(error) // TODO check
		}
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	if alreadyApplied, err := wasMigrationApplied(tx, mi); err != nil {
		return err
	} else if !alreadyApplied {
		return nil
	}
	mi.Down(MigrationTx{tx})
	_, err = tx.Exec("delete from Migrations where id = ?", mi.ID())
	if err != nil {
		return nil
	}
	return
}

func wasMigrationApplied(tx *sql.Tx, mi Migration) (bool, error) {
	row := tx.QueryRow("select count(*) from Migrations where id = ?", mi.ID())
	var alreadyApplied int64
	if err := row.Scan(&alreadyApplied); err != nil {
		return false, err
	}
	return alreadyApplied != 0, nil
}
