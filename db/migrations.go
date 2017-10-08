package db

import "database/sql"

type MigrationTx struct {
	tx *sql.Tx
}

func (tx MigrationTx) Exec(q string, vs ...interface{}) error {
	_, err := tx.tx.Exec(q, vs...)
	return err
}

type Migration interface {
	ID() string
	Up(MigrationTx)
	Down(MigrationTx)
}

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
	row := tx.QueryRow("select count(*) from Migrations where id = ?", mi.ID())
	var alreadyApplied int64
	if err = row.Scan(&alreadyApplied); err != nil {
		return err
	}
	if alreadyApplied != 0 {
		return nil
	}
	mi.Up(MigrationTx{tx})
	_, err = tx.Exec("insert into Migrations (id) values (?)", mi.ID())
	if err != nil {
		return nil
	}
	return
}
func UndoMigration(db *sql.DB, mi Migration) (err error) {
	// TODO check the migration is already there
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
	mi.Down(MigrationTx{tx})
	_, err = tx.Exec("delete from Migrations where id = ?", mi.ID())
	if err != nil {
		return nil
	}
	return
}
