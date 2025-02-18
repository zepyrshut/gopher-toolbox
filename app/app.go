package app

import (
	"bufio"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"gopher-toolbox/mail"
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
	UsernameAlReadyExists string = "username_already_exists"
	UserSessionKey        string = "user_session_key"
	EmailAlreadyExists    string = "email_already_exists"
	PhoneNumberKey        string = "phone_number_key"
	PhoneAlreadyExists    string = "phone_already_exists"
	IncorrectPassword     string = "incorrect_password"
	ErrorGeneratingToken  string = "error_generating_token"
	LoggedIn              string = "logged_in"
)

type App struct {
	Database AppDatabase
	Security AppSecurity
	AppInfo  AppInfo
	Mailer   mail.Mailer
}

type AppDatabase struct {
	DriverName string
	DataSource string
	Migrate    bool
}

type AppInfo struct {
	Name    string
	Version string
}

type AppSecurity struct {
	AsymmetricKey paseto.V4AsymmetricSecretKey
	PublicKey     paseto.V4AsymmetricPublicKey
	Duration      time.Duration
}

func New(name, version, envDirectory string) *App {
	var err error

	err = loadEnvFile(envDirectory)
	if err != nil {
		slog.Error("error loading env file, using default values", "error", err)
	}

	var durationTime time.Duration
	var ak paseto.V4AsymmetricSecretKey

	if os.Getenv("ASYMMETRIC_KEY") != "" {
		ak, err = paseto.NewV4AsymmetricSecretKeyFromHex(os.Getenv("ASYMMETRIC_KEY"))
		if err != nil {
			slog.Error("error creating asymmetric key", "error", err)
		}
	} else {
		ak = paseto.NewV4AsymmetricSecretKey()
	}

	slog.Info("asymmetric key", "key", ak.ExportHex())

	pk := ak.Public()

	duration := os.Getenv("DURATION")
	durationTime = time.Hour * 24 * 7
	if duration != "" {
		if parsed, err := time.ParseDuration(duration); err == nil {
			durationTime = parsed
		}
	}

	return &App{
		Mailer: mail.New(
			os.Getenv("SMTP_HOST"),
			os.Getenv("SMTP_PORT"),
			os.Getenv("SMTP_USER"),
			os.Getenv("SMTP_PASS"),
		),
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
			Name:    name,
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
