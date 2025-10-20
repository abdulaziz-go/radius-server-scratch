package tests

import (
	"os"
	"path/filepath"
	"testing"

	"radius-server/src/common/logger"
	"radius-server/src/config"
	"radius-server/src/database"
)

// TestMain sets up global test initialization. It runs before any tests in this package.
func TestMain(m *testing.M) {
	// Ensure working directory is the module root so ./.env is discoverable
	chdirToModuleRoot()
	config.LoadConfig()
	logger.InitializeLogger()
	InitlizeValues()
	logger.Logger.Info().Msg("Starting tests...")
	if err := database.Connect(); err != nil {
		logger.Logger.Fatal().Msgf("Connection to database error. %s", err.Error())
	}
	if err := CreateNas(); err != nil {
		logger.Logger.Fatal().Msgf("Creating NAS error. %s", err.Error())
	}
	defer func() {
		if err := DeleteNas(); err != nil {
			logger.Logger.Error().Msgf("Deleting NAS error. %s", err.Error())
		}
	}()

	code := m.Run()
	os.Exit(code)
}

// chdirToModuleRoot walks up from the current directory until it finds go.mod,
// then changes the working directory to that path. If not found, it stays put.
func chdirToModuleRoot() {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	prev := ""
	for cwd != prev {
		if fileExists(filepath.Join(cwd, "go.mod")) {
			_ = os.Chdir(cwd)
			return
		}
		prev = cwd
		cwd = filepath.Dir(cwd)
	}
}

func fileExists(p string) bool {
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
