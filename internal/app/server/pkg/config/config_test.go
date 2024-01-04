package config

import (
	"testing"
)

func TestConfig_Load(t *testing.T) {
	t.Setenv(portEnvVarName, "1")
}
