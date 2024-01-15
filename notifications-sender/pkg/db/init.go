package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	m "github.com/rubenv/sql-migrate"
)

type StartUpOptions struct {
	DBHost         string
	DBPort         int
	DBName         string
	DBUsername     string
	DBPassword     string
	SkipMigrations bool
}

func StartDBStore(opts StartUpOptions) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		opts.DBUsername, opts.DBPassword, opts.DBHost, opts.DBPort, opts.DBName))
	if err != nil {
		return nil, err
	}

	if opts.SkipMigrations {
		return db, nil
	}

	migrations := &m.MemoryMigrationSource{
		Migrations: []*m.Migration{
			{
				Id: "100",
				Up: []string{
					`CREATE TABLE IF NOT EXISTS notifications (
						id VARCHAR(200),
						type VARCHAR(200),
						message VARCHAR(200),
						receiver VARCHAR(200),
						status VARCHAR(200),
						retries INT,
						created_at datetime DEFAULT CURRENT_TIMESTAMP,
						updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
						PRIMARY KEY (id))`,
				},
			},
		},
	}

	_, err = m.Exec(db.DB, "mysql", migrations, m.Up)
	if err != nil {
		return nil, err
	}

	return db, nil
}
