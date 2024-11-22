package config

import (
	"bufio"
	"log/slog"
	"os"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
)

type App struct {
	DataSource string
	Security   Security
	AppInfo    AppInfo
}

type AppInfo struct {
	Version string
}

type Security struct {
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
		DataSource: os.Getenv("DATASOURCE"),
		Security: Security{
			AsymmetricKey: ak,
			PublicKey:     pk,
			Duration:      durationTime,
		},
		AppInfo: AppInfo{
			Version: version,
		},
	}
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
