package test

import (
	"DoushengABCD/service"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	service.InitRedis()
	service.InitDatabase()
	code := m.Run()
	os.Exit(code)
}
