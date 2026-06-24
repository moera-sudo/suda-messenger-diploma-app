package scripts

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func Truncate(conn *pgx.Conn, tableName string) {
	ctx := context.Background()

	if tableName != "" {
		// Очистка одной таблицы
		fmt.Printf("Clearing table: %s...\n", tableName)
		_, err := conn.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", tableName))
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	} else {
		// Очистка ВСЕГО
		fmt.Println("Clearing ALL tables...")
		query := `
			DO $$
			DECLARE
				r RECORD;
			BEGIN
				FOR r IN (
					SELECT tablename
					FROM pg_tables
					WHERE schemaname = 'public'
				)
				LOOP
					EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' RESTART IDENTITY CASCADE';
				END LOOP;
			END $$;
		`
		_, err := conn.Exec(ctx, query)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
	fmt.Println("✅ Done!")
}

func DropAll(conn *pgx.Conn) {
	fmt.Println("Dropping public schema ")
	_, err := conn.Exec(context.Background(), "DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println("Schema reset. Better run migrations")
}

func ShowInfo(conn *pgx.Conn) {
	var count int
	conn.QueryRow(context.Background(), "SELECT count(*) FROM messenger_users").Scan(&count)
	fmt.Printf("Users count: %d\n", count)
}
