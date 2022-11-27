package db

import (
	"fmt"
	"log"
	"tgbot/configs"

	"github.com/golang-migrate/migrate/v4"
	// needed migration packages
	_ "github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	"go.uber.org/zap"

	// database driver
	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Init initializes database connection then connect with postgres
func Init(config *configs.Configuration) (*gorm.DB, error) {

	fmt.Println("init db")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresHost, config.PostgresPort, config.PostgresDatabase)

	m, err := migrate.New("file://db/migrations", dbURL)
	if err != nil {
		log.Fatal("error in creating migrations: ", zap.Error(err))
	}
	if err := m.Up(); err != nil {
		log.Println("error updating migrations: ", zap.Error(err))
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}

	return db, nil
}
