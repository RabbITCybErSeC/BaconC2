package config

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"log"
	"os"

	"github.com/RabbITCybErSeC/BaconC2/server/db"
	"github.com/RabbITCybErSeC/BaconC2/server/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ServerConfig struct {
	DB              *gorm.DB
	Port            string
	Env             string
	DBPath          string
	MaxAgents       int
	JWTSecret       string
	AgentHTTPConfig AgentHTTPConfig
	FrontHTTPConfig FrontEndHTTPConfig
	UDPConfig       UDPConfig
}

type FrontEndHTTPConfig struct {
	Port int
}

type AgentHTTPConfig struct {
	Port    int
	Enabled bool
}

type UDPConfig struct {
	Port    int
	Enabled bool
}

func NewServerConfig() *ServerConfig {
	httpPort := flag.Int("http-port", 8080, "HTTP server port")
	udpPort := flag.Int("udp-port", 8081, "UDP server port")
	enableUDP := flag.Bool("enable-udp", false, "Enable UDP transport")
	flag.Parse()

	// Load JWT_SECRET from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	config := &ServerConfig{
		Port:      ":8080",
		Env:       "development",
		DBPath:    "agents.db",
		MaxAgents: 100,
		JWTSecret: jwtSecret,
		FrontHTTPConfig: FrontEndHTTPConfig{
			Port: *httpPort,
		},
		AgentHTTPConfig: AgentHTTPConfig{Port: 8081, Enabled: true},
		UDPConfig: UDPConfig{
			Port:    *udpPort,
			Enabled: *enableUDP,
		},
	}

	db, err := initializeDatabase(config.DBPath)
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	config.DB = db

	return config
}

func initializeDatabase(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = runMigrations(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Command{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Agent{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	if err := seedStaticUser(db); err != nil {
		return err
	}

	return nil
}

func seedStaticUser(gormDB *gorm.DB) error {
	userRepo := db.NewUserRepository(gormDB)

	staticUserName := "admin"
	if _, err := userRepo.FindByUsername(staticUserName); err == nil {
		log.Println("Static user already exists, skipping creation")
		return nil
	}

	staticPassword, err := generateRandomPassword(16)
	if err != nil {
		return err
	}

	// Hash the random password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(staticPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	staticUser := &models.User{
		Username: staticUserName,
		Password: string(hashedPassword),
	}

	if err := userRepo.Save(staticUser); err != nil {
		return err
	}

	log.Printf("Created static user: username=%s, password=%s (change this in production!)", staticUserName, staticPassword)
	return nil
}

func generateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (c *ServerConfig) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
