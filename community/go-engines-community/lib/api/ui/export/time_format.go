package export

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

// ResolveTimeFormat returns a whitelisted Go layout for an export request.
// It uses the user interface format when the request does not specify one, then
// falls back to the default when the stored format is empty or unrecognized.
func ResolveTimeFormat(ctx context.Context, requested string, adapter config.UserInterfaceAdapter) (string, error) {
	timeFormat := validation.GetRealFormatTime(requested)
	if requested == "" {
		conf, err := adapter.GetConfig(ctx)
		if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
			return "", err
		}
		timeFormat = validation.GetRealFormatTime(conf.DateFormat)
	}
	if timeFormat == "" {
		timeFormat = validation.DefaultTimeFormat
	}
	return timeFormat, nil
}
