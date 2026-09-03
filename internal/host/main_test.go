package host

import (
	"os"
	"testing"
)

// Host's upstream fixtures assert Chinese compatibility strings. Tests for
// Vietnamese display output opt into vi explicitly.
func TestMain(m *testing.M) {
	_ = os.Setenv("AINOVEL_LOCALE", "zh")
	os.Exit(m.Run())
}
