package logredis

import (
	"os"
	"testing"
)

func TestNewHookFunc(t *testing.T) {
	hook, err := NewHook("redishost", "key", "format", "appname", "hostname", "password", 1)
	if hook != nil {
		t.Fatalf("TestNewHookFunc, expected no hook, got hook: %s", hook)
	}

	if err != nil {
		t.Fatalf("TestNewHookFunc, expected %q, got %s.", "unknown message format", err)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
