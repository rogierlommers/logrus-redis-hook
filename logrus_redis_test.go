package logredis

import (
	"os"
	"testing"
)

func TestNewHookFunc(t *testing.T) {
	config := HookConfig{
		Host:     "localhost",
		Key:      "key",
		Format:   "format",
		App:      "appname",
		Password: "password",
		Port:     1,
		DB:       1,
		TTL:      3600,
	}

	hook, err := NewHook(config)
	if hook != nil {
		t.Fatalf("TestNewHookFunc, expected no hook, got hook: %v", hook)
	}

	if err == nil {
		t.Fatalf("TestNewHookFunc, expected %q, got %s.", "unknown message format", err)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
