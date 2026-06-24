package main


import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	scr "github.com/moera-sudo/backend/scripts/database/internal/scripts"
)

func main() {
	cmd := flag.String("cmd", "", "Command to run: truncate, drop_schema, info, superadmin")
	table := flag.String("table", "", "Specific table to clear(optional)")
	isDocker := flag.Bool("docker", false, "Switch to inDocker version of DB")
	// superAdminUsername := flag.String("username", "superadmin", "Superadmin username")
	// superAdminEmail := flag.String("email", "superadmin@local.dev", "Superadmin email")
	// superAdminPassword := flag.String("password", "change-me-superadmin", "Superadmin password")
	// superAdminDisplayName := flag.String("display-name", "Super Admin", "Superadmin display name")
	flag.Parse()

	if *cmd == "" {
		fmt.Println("Empty command. Please use one of these commands: -cmd=[truncate|drop_schema|info|superadmin]")
		os.Exit(1)
	}

	conn := connectDB(*isDocker)
	defer conn.Close(context.Background())

	switch *cmd {
	case "truncate":
		scr.Truncate(conn, *table)
	case "drop_schema":
		scr.DropAll(conn)
	case "info":
		scr.ShowInfo(conn)
	// case "superadmin":
	// 	scr.RegisterSuperAdmin(conn, scr.SuperAdminInput{
	// 		Username:    *superAdminUsername,
	// 		Email:       *superAdminEmail,
	// 		Password:    *superAdminPassword,
	// 		DisplayName: *superAdminDisplayName,
	// 	})
	default:
		fmt.Println("Incorrect command. Nothing changed")
	}
}

func connectDB(isDocker bool) *pgx.Conn {
	godotenv.Load("../../../../.env") 
	dbURL := os.Getenv("DATABASE_URL")
	if isDocker{
		dbURL = "postgres://postgres:123@localhost:55432/suda?sslmode=disable" 
	} 

	if dbURL == "" {
		log.Fatal(".env config is not available. Using default db string")
		dbURL = "postgres://postgres:123@localhost:5432/messenger_db?sslmode=disable"
	}
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}
	return conn
}
