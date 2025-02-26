package app

import (
	"bufio"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// TODO: review consts
const (
	// Handlers keys
	InvalidRequest    string = "invalid_request"
	MalformedJSON     string = "malformed_json"
	TokenBlacklisted  string = "token_blacklisted"
	TokenInvalid      string = "token_invalid"
	ValidationFailed  string = "validation_failed"
	UntilBeforeTo     string = "until_before_to"
	InternalError     string = "internal_error"
	NotFound          string = "not_found"
	Created           string = "created"
	Updated           string = "updated"
	Deleted           string = "deleted"
	Enabled           string = "enabled"
	Disabled          string = "disabled"
	Retrieved         string = "retrieved"
	ErrorCreating     string = "error_creating"
	ErrorUpdating     string = "error_updating"
	ErrorEnabling     string = "error_enabling"
	ErrorDisabling    string = "error_disabling"
	ErrorGetting      string = "error_getting"
	ErrorGettingAll   string = "error_getting_all"
	ErrorMailing      string = "error_mailing"
	InvalidEntityID   string = "invalid_entity_id"
	NotImplemented    string = "not_implemented"
	NotPassValidation string = "not_pass_validation"

	// User keys
	UserUsernameKey       string = "username_key"
	UserEmailKey          string = "email_key"
	UsernameAlreadyExists string = "username_already_exists"
	UserSessionKey        string = "user_session_key"
	EmailAlreadyExists    string = "email_already_exists"
	PhoneNumberKey        string = "phone_number_key"
	PhoneAlreadyExists    string = "phone_already_exists"
	IncorrectPassword     string = "incorrect_password"
	ErrorGeneratingToken  string = "error_generating_token"
	LoggedIn              string = "logged_in"
)

var (
	logFile  *os.File
	logLevel string
)

type Config struct {
	// default ""
	Name string

	// default ""
	Version string

	// default ".env"
	EnvDirectory string

	// default "development"
	EnvMode string

	// default "debug"
	LogLevel string

	// default "UTC"
	Timezone string

	// default nil
	Paseto *Paseto

	// default ""
	SMTPHost string

	// default ""
	SMTPPort string

	// default ""
	SMTPUser string

	// default ""
	SMTPPass string

	// default ""
	DatabaseDriverName string

	// default ""
	DatabaseDataSource string

	// default false
	DatabaseMigrate bool
}

type App struct {
	config Config
}

type Paseto struct {
	AsymmetricKey paseto.V4AsymmetricSecretKey
	PublicKey     paseto.V4AsymmetricPublicKey
	Duration      time.Duration
}

func New(config ...Config) *App {
	cfg := Config{
		Name:               "",
		Version:            "",
		EnvDirectory:       ".env",
		EnvMode:            "development",
		LogLevel:           "debug",
		Timezone:           "UTC",
		Paseto:             nil,
		SMTPHost:           "",
		SMTPPort:           "",
		SMTPUser:           "",
		SMTPPass:           "",
		DatabaseDriverName: "pgx",
		DatabaseDataSource: "",
		DatabaseMigrate:    false,
	}

	if len(config) > 0 {
		cfg = config[0]

		if cfg.EnvDirectory == "" {
			cfg.EnvDirectory = ".env"
		}
		if cfg.EnvMode == "" {
			cfg.EnvMode = "development"
		}
		if cfg.LogLevel == "" {
			cfg.LogLevel = "debug"
		}
		if cfg.Timezone == "" {
			cfg.Timezone = "UTC"
		}
		if cfg.DatabaseDriverName == "" {
			cfg.DatabaseDriverName = "pgx"
		}
	}

	envDir := os.Getenv("ENV_DIRECTORY")
	if envDir == "" {
		envDir = cfg.EnvDirectory
	}

	err := loadEnvFile(envDir)
	if err != nil {
		slog.Error("error loading env file", "error", err, "directory", envDir)
	}

	if cfg.Name == "" && os.Getenv("APP_NAME") != "" {
		cfg.Name = os.Getenv("APP_NAME")
	}

	if cfg.Version == "" && os.Getenv("APP_VERSION") != "" {
		cfg.Version = os.Getenv("APP_VERSION")
	}

	if cfg.EnvMode == "" && os.Getenv("ENV_MODE") != "" {
		cfg.EnvMode = os.Getenv("ENV_MODE")
	}

	if cfg.LogLevel == "" && os.Getenv("LOG_LEVEL") != "" {
		cfg.LogLevel = os.Getenv("LOG_LEVEL")
		logLevel = cfg.LogLevel
	}

	if cfg.Timezone == "" && os.Getenv("TIMEZONE") != "" {
		cfg.Timezone = os.Getenv("TIMEZONE")
	}

	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		slog.Error("error loading timezone", "error", err, "timezone", cfg.Timezone)
		loc = time.UTC
	}
	time.Local = loc

	startRotativeLogger()

	if cfg.Paseto == nil {
		var ak paseto.V4AsymmetricSecretKey
		var err error

		if os.Getenv("PASETO_ASYMMETRIC_KEY") != "" {
			slog.Info("using paseto asymmetric key from env")
			ak, err = paseto.NewV4AsymmetricSecretKeyFromHex(os.Getenv("PASETO_ASYMMETRIC_KEY"))
			if err != nil {
				slog.Error("error creating asymmetric key", "error", err)
				ak = paseto.NewV4AsymmetricSecretKey()
			}
		} else {
			ak = paseto.NewV4AsymmetricSecretKey()
		}

		pk := ak.Public()

		duration := time.Hour * 24 * 7 // 7 días por defecto
		if os.Getenv("PASETO_DURATION") != "" {
			durationStr := os.Getenv("PASETO_DURATION")
			durationInt, err := time.ParseDuration(durationStr)
			if err != nil {
				slog.Error("error parsing PASETO_DURATION", "error", err, "duration", durationStr)
			} else {
				duration = durationInt
			}
		}

		cfg.Paseto = &Paseto{
			AsymmetricKey: ak,
			PublicKey:     pk,
			Duration:      duration,
		}
	}

	if cfg.SMTPHost == "" && os.Getenv("SMTP_HOST") != "" {
		cfg.SMTPHost = os.Getenv("SMTP_HOST")
	}

	if cfg.SMTPPort == "" && os.Getenv("SMTP_PORT") != "" {
		cfg.SMTPPort = os.Getenv("SMTP_PORT")
	}

	if cfg.SMTPUser == "" && os.Getenv("SMTP_USER") != "" {
		cfg.SMTPUser = os.Getenv("SMTP_USER")
	}

	if cfg.SMTPPass == "" && os.Getenv("SMTP_PASS") != "" {
		cfg.SMTPPass = os.Getenv("SMTP_PASS")
	}

	if cfg.DatabaseDriverName == "" && os.Getenv("DATABASE_DRIVER_NAME") != "" {
		cfg.DatabaseDriverName = os.Getenv("DATABASE_DRIVER_NAME")
	}

	if strings.HasPrefix(cfg.DatabaseDataSource, "OVERRIDE_") {
		envKey := strings.TrimPrefix(cfg.DatabaseDataSource, "OVERRIDE_")
		if envValue := os.Getenv(envKey); envValue != "" {
			slog.Info("using override database data source", "key", envKey)
			cfg.DatabaseDataSource = envValue
		} else {
			slog.Warn("override database data source key not found in environment", "key", envKey)
			if os.Getenv("DATABASE_DATA_SOURCE") != "" {
				cfg.DatabaseDataSource = os.Getenv("DATABASE_DATA_SOURCE")
			}
		}
	} else if cfg.DatabaseDataSource == "" && os.Getenv("DATABASE_DATA_SOURCE") != "" {
		cfg.DatabaseDataSource = os.Getenv("DATABASE_DATA_SOURCE")
	}

	if !cfg.DatabaseMigrate && os.Getenv("DATABASE_MIGRATE") == "true" {
		cfg.DatabaseMigrate = true
	}

	app := &App{
		config: cfg,
	}

	slog.Info(
		"app config",
		"name", cfg.Name,
		"version", cfg.Version,
		"env_directory", cfg.EnvDirectory,
		"env_mode", cfg.EnvMode,
		"log_level", cfg.LogLevel,
		"timezone", cfg.Timezone,
		"paseto_public_key", cfg.Paseto.PublicKey.ExportHex(),
		"paseto_duration", cfg.Paseto.Duration.String(),
		"smtp_host", cfg.SMTPHost,
		"smtp_port", cfg.SMTPPort,
		"smtp_user", cfg.SMTPUser,
		"smtp_pass", cfg.SMTPPass,
		"database_driver", cfg.DatabaseDriverName,
		"database_source", cfg.DatabaseDataSource,
		"database_migrate", cfg.DatabaseMigrate,
	)

	return app
}

func (a *App) Name() string {
	return a.config.Name
}

func (a *App) Version() string {
	return a.config.Version
}

func (a *App) EnvMode() string {
	return a.config.EnvMode
}

func (a *App) LogLevel() string {
	return a.config.LogLevel
}

func (a *App) Paseto() *Paseto {
	return a.config.Paseto
}

func (a *App) SMTPConfig() (host, port, user, pass string) {
	return a.config.SMTPHost, a.config.SMTPPort, a.config.SMTPUser, a.config.SMTPPass
}

func (a *App) DatabaseDataSource() string {
	return a.config.DatabaseDataSource
}

func (a *App) Timezone() string {
	return a.config.Timezone
}

// MigrateDB migrates the database. The migrations must stored in the
// "database/migrations" directory inside cmd directory along with the main.go.
//
// cmd/main.go
//
// cmd/database/migrations/*.sql
func (a *App) Migrate(database embed.FS) {
	if a.config.DatabaseMigrate == false {
		slog.Info("migration disabled")
		return
	}
	dbConn, err := sql.Open(a.config.DatabaseDriverName, a.config.DatabaseDataSource)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbConn.Close()

	d, err := iofs.New(database, "database/migrations")
	if err != nil {
		fmt.Println(err)
		return
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, a.config.DatabaseDataSource)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("cannot migrate", "error", err)
		panic(err)
	}
	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("migration has no changes")
		return
	}

	slog.Info("migration done")
}

func loadEnvFile(envDirectory string) error {
	file, err := os.Open(envDirectory)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}
	return scanner.Err()
}

func slogLevelType(level string) slog.Level {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	return logLevel
}

func newLogger(level slog.Level) {
	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Println("error creating logs directory:", err)
		return
	}

	now := time.Now().Format("2006-01-02")
	f, err := os.OpenFile(fmt.Sprintf("logs/log%s.log", now), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error opening log file:", err)
		return
	}

	mw := io.MultiWriter(os.Stdout, f)
	logger := slog.New(slog.NewTextHandler(mw, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}))

	if logFile != nil {
		logFile.Close() // Cierra el archivo anterior antes de rotar
	}

	logFile = f
	slog.SetDefault(logger)
}

func startRotativeLogger() {
	newLogger(slogLevelType(logLevel))

	ticker := time.NewTicker(time.Hour * 24)
	go func() {
		for range ticker.C {
			newLogger(slogLevelType(logLevel))
		}
	}()
}
