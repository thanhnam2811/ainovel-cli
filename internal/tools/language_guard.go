package tools

import (
	"encoding/json"
	"fmt"

	"github.com/voocel/ainovel-cli/internal/errs"
	"github.com/voocel/ainovel-cli/internal/localization"
)

func requireVietnameseArgs(label string, args json.RawMessage) error {
	if err := localization.RequireVietnamese(label, string(args)); err != nil {
		return fmt.Errorf("%w: %w", err, errs.ErrToolArgs)
	}
	return nil
}
