package database

import (
	"fmt"
	"log"
	"tether-server/config"
	"tether-server/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.AppConfig.DBHost,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
		config.AppConfig.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	// Force drop all tables and recreate
	log.Println("Dropping all existing tables...")
	db.Exec("DROP TABLE IF EXISTS messages CASCADE")
	db.Exec("DROP TABLE IF EXISTS chats CASCADE")
	db.Exec("DROP TABLE IF EXISTS verification_codes CASCADE")
	db.Exec("DROP TABLE IF EXISTS users CASCADE")

	// Also drop any sequences that might exist
	db.Exec("DROP SEQUENCE IF EXISTS users_id_seq CASCADE")
	db.Exec("DROP SEQUENCE IF EXISTS chats_id_seq CASCADE")
	db.Exec("DROP SEQUENCE IF EXISTS messages_id_seq CASCADE")

	log.Println("Creating new UUID-based structure...")

	// Auto Migrate to create new structure
	err = db.AutoMigrate(
		&models.User{},
		&models.Chat{},
		&models.Message{},
		&models.EmailVerification{},
		&models.RefreshToken{},
		&models.Workspace{},
		&models.WorkspaceMember{},
		&models.Board{},
		&models.Column{},
		&models.Card{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database. \n", err)
	}

	DB = db
	log.Println("Database connected successfully")
}
