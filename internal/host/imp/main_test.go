package imp

import (
	"os"
	"testing"
)

// Import fixtures exercise the upstream Chinese-source compatibility path.
// Vietnamese runtime tests live at the localization/tool boundaries.
func TestMain(m *testing.M) {
	_ = os.Setenv("AINOVEL_LOCALE", "zh")
	os.Exit(m.Run())
}
