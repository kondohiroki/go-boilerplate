package migrations

import (
	"context"

	"github.com/kondohiroki/go-boilerplate/internal/db/pgx"
)

func init() {
	Migrations = append(Migrations, createUserTable)
}

var createUserTable = &Migration{
	Name: "20230407151155_create_user_table",
	Up: func() error {
		_, err := pgx.GetPgxPool().Exec(context.Background(), `
			CREATE TABLE users (
				id SERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				email VARCHAR(255) NOT NULL UNIQUE
			);

			INSERT INTO users (name, email) VALUES ('Default User', 'user@example.com');
		`)

		if err != nil {
			return err
		}
		return nil

	},
	Down: func() error {
		_, err := pgx.GetPgxPool().Exec(context.Background(), `
			// code here
		`)
		if err != nil {
			return err
		}

		return nil
	},
}
