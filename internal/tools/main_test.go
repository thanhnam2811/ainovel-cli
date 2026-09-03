package tools

import (
	"os"
	"testing"
)

// Most upstream tool fixtures intentionally use Chinese prose. Keep them in
// compatibility mode; locale-specific tests opt back into vi explicitly.
func TestMain(m *testing.M) {
	_ = os.Setenv("AINOVEL_LOCALE", "zh")
	os.Exit(m.Run())
}
