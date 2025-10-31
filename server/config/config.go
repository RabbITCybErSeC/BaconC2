package config

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/RabbITCybErSeC/BaconC2/server/db"
	"github.com/RabbITCybErSeC/BaconC2/server/models"
	"github.com/joho/godotenv"
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
	// Load .env file if it exists - try multiple locations
	_ = godotenv.Load()              // current directory
	_ = godotenv.Load(".env.local")  // local overrides
	_ = godotenv.Load("config/.env") // config directory

	httpPort := flag.Int("http-port", getIntEnv("HTTP_PORT", 8080), "HTTP server port")
	agentPort := flag.Int("agent-port", getIntEnv("AGENT_PORT", 8081), "Agent server port")
	udpPort := flag.Int("udp-port", getIntEnv("UDP_PORT", 8081), "UDP server port")
	enableUDP := flag.Bool("enable-udp", getBoolEnv("ENABLE_UDP", false), "Enable UDP transport")
	flag.Parse()

	// Load JWT_SECRET from environment variable
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters long")
	}

	config := &ServerConfig{
		Port:      ":" + strconv.Itoa(*httpPort),
		Env:       getEnv("ENVIRONMENT", "development"),
		DBPath:    getEnv("DB_PATH", "agents.db"),
		MaxAgents: getIntEnv("MAX_AGENTS", 100),
		JWTSecret: jwtSecret,
		FrontHTTPConfig: FrontEndHTTPConfig{
			Port: *httpPort,
		},
		AgentHTTPConfig: AgentHTTPConfig{
			Port:    *agentPort,
			Enabled: getBoolEnv("AGENT_HTTP_ENABLED", true),
		},
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

	// Log configuration summary
	log.Printf("Server starting - Port: %s, Environment: %s, Max Agents: %d",
		config.Port, config.Env, config.MaxAgents)
	if config.UDPConfig.Enabled {
		log.Printf("UDP transport enabled on port %d", config.UDPConfig.Port)
	}

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
	if err := db.AutoMigrate(&models.AgentCommand{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.ServerAgentModel{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.AgentCommandResult{}); err != nil {
		return err
	}

	if err := migrateArgsColumn(db); err != nil {
		return err
	}

	if err := seedStaticUser(db); err != nil {
		return err
	}

	return nil
}

func migrateArgsColumn(db *gorm.DB) error {
	var columnExists bool
	err := db.Raw("SELECT COUNT(*) FROM pragma_table_info('agent_commands') WHERE name='args'").Scan(&columnExists).Error
	if err != nil {
		return err
	}

	if columnExists {
		err = db.Exec("UPDATE agent_commands SET args = '[]' WHERE args IS NULL OR args = ''").Error
		if err != nil {
			return err
		}
	}

	return nil
}

func seedStaticUser(gormDB *gorm.DB) error {
	userRepo := db.NewUserRepository(gormDB)

	staticUserName := getEnv("ADMIN_USERNAME", "admin")
	if _, err := userRepo.FindByUsername(staticUserName); err == nil {
		log.Printf("Admin user '%s' already exists, skipping creation", staticUserName)
		return nil
	}

	var staticPassword string
	if envPassword := getEnv("ADMIN_PASSWORD", ""); envPassword != "" {
		staticPassword = envPassword
		log.Printf("Using admin password from environment")
	} else {
		var err error
		staticPassword, err = generateRandomPassword(16)
		if err != nil {
			return err
		}
		log.Printf("Generated admin credentials - username: %s, password: %s", staticUserName, staticPassword)
		log.Printf("WARNING: Change the admin password in production!")
	}

	// Hash the password
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

	log.Printf("Successfully created admin user: %s", staticUserName)
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
	if c.DB == nil {
		return nil
	}
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// IsProduction returns true if running in production environment
func (c *ServerConfig) IsProduction() bool {
	return c.Env == "production"
}

// IsDevelopment returns true if running in development environment
func (c *ServerConfig) IsDevelopment() bool {
	return c.Env == "development"
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if value := getEnv(key, ""); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
		log.Printf("Warning: Invalid integer value for %s: %s, using default %d", key, value, fallback)
	}
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	value := getEnv(key, "")
	if value != "" {
		return value == "true" || value == "1"
	}
	return fallback
}
