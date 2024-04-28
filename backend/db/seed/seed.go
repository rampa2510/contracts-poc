package seed

import "database/sql"

func SeedToDb(db *sql.DB) error {
	const createTable = `
    CREATE TABLE IF NOT EXISTS users (
      id INTEGER NOT NULL PRIMARY KEY,
      name TEXT NOT NULL,
      email TEXT NOT NULL
    );
    `

	if _, err := db.Exec(createTable); err != nil {
		return err
	}

	return nil
}
