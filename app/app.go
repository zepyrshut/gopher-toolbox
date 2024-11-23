package app

import (
	"bufio"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"gopher-toolbox/utils"
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

const (
	// Handlers keys
	InvalidRequest  string = "invalid_request"
	InternalError   string = "internal_error"
	NotFound        string = "not_found"
	Created         string = "created"
	Updated         string = "updated"
	Deleted         string = "deleted"
	Enabled         string = "enabled"
	Disabled        string = "disabled"
	Retrieved       string = "retrieved"
	ErrorCreating   string = "error_creating"
	ErrorUpdating   string = "error_updating"
	ErrorEnabling   string = "error_enabling"
	ErrorDisabling  string = "error_disabling"
	ErrorGetting    string = "error_getting"
	ErrorGettingAll string = "error_getting_all"
	InvalidEntityID string = "invalid_entity_id"
	NotImplemented  string = "not_implemented"

	// User keys
	UserUsernameKey       string = "user_username_key"
	UserEmailKey          string = "user_email_key"
	UsernameAlReadyExists string = "username_already_exists"
	EmailAlreadyExists    string = "email_already_exists"
	IncorrectPassword     string = "incorrect_password"
	ErrorGeneratingToken  string = "error_generating_token"
	LoggedIn              string = "logged_in"
)

type App struct {
	Database AppDatabase
	Security AppSecurity
	AppInfo  AppInfo
}

type AppDatabase struct {
	DriverName string
	DataSource string
	Migrate    bool
}

type AppInfo struct {
	Version string
}

type AppSecurity struct {
	AsymmetricKey paseto.V4AsymmetricSecretKey
	PublicKey     paseto.V4AsymmetricPublicKey
	Duration      time.Duration
}

func New(version string) *App {
	var err error

	err = loadEnvFile()
	if err != nil {
		slog.Error("error loading env file", "error", err)
		panic(err)
	}

	var durationTime time.Duration
	var ak paseto.V4AsymmetricSecretKey

	ak, err = paseto.NewV4AsymmetricSecretKeyFromHex(os.Getenv("ASYMMETRICKEY"))
	if err != nil {
		ak = paseto.NewV4AsymmetricSecretKey()
	}
	pk := ak.Public()

	duration := os.Getenv("DURATION")
	if duration != "" {
		durationTime, err = time.ParseDuration(duration)
		if err != nil {
			durationTime = time.Hour * 24 * 7
		}
	}

	return &App{
		Database: AppDatabase{
			Migrate:    utils.GetBool(os.Getenv("MIGRATE")),
			DriverName: os.Getenv("DRIVERNAME"),
			DataSource: os.Getenv("DATASOURCE"),
		},
		Security: AppSecurity{
			AsymmetricKey: ak,
			PublicKey:     pk,
			Duration:      durationTime,
		},
		AppInfo: AppInfo{
			Version: version,
		},
	}
}

// MigrateDB migrates the database. The migrations must stored in the
// "database/migrations" directory inside cmd directory along with the main.go.
//
// cmd/main.go
//
// cmd/database/migrations/*.sql
func (a *App) Migrate(database embed.FS) {
	if a.Database.Migrate == false {
		slog.Info("migration disabled")
		return
	}
	dbConn, err := sql.Open(a.Database.DriverName, a.Database.DataSource)
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

	m, err := migrate.NewWithSourceInstance("iofs", d, a.Database.DataSource)
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

func loadEnvFile() error {
	file, err := os.Open(".env")
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
