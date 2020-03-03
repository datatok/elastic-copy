package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

// EnvSettings describes all of the environment settings.
type EnvSettings struct {

	// Logging stuff
	LogLevel string

	Debug bool
}

func New() *EnvSettings {

	env := EnvSettings{
		LogLevel:                envOr("LOG_LEVEL", "debug"),
	}

	env.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))

	return &env
}

// logrus.SetLevel(logrus.DebugLevel)

// AddFlags binds flags to the given flagset.
func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&s.LogLevel, "log_level", "", s.LogLevel, "Log level (debug, info, warn, error)")
	fs.BoolVar(&s.Debug, "debug", s.Debug, "enable verbose output")
}

func envOr(name, def string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return def
}

func (s *EnvSettings) EnvVars() map[string]string {
	envvars := map[string]string{
		"APP_BIN":          os.Args[0],
		"DEBUG":            fmt.Sprint(s.Debug),
		"LOG_LEVEL":        s.LogLevel,
	}

	return envvars
}
