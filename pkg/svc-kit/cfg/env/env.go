package env

import (
	"flag"
	"fmt"

	"github.com/joho/godotenv"
)

// Env - App Environment
type Env string

const (
	// Local env
	Local Env = "local"
	// Compose  env for docker
	Compose Env = "compose"
	// Testing env
	Testing Env = "testing"
	// Ropsten test env
	Ropsten Env = "ropsten"
	// Production env
	Production Env = "production"
)

var instance Env

func Get() Env {
	return instance
}

// Detect env
func Detect() (Env, error) {
	envFlag := flag.String("env", string(Local), "'local', 'compose', 'ropsten', or 'production'. 'testing' is automatically detected")
	flag.Parse()

	instance = Env(*envFlag)
	filename := instance
	if flag.Lookup("test.v") != nil {
		instance = Testing
	}
	err := godotenv.Load(fmt.Sprintf(`%s.env`, filename))
	return instance, err
}
