package eventfilter

import "errors"

var ErrDropOutcome = errors.New("drop event")

// ErrRulesetChanged is returned by ProcessEvent when a suspended event is resumed
// but the event filter ruleset has changed since it was parked.
var ErrRulesetChanged = errors.New("event filter ruleset changed during external data fetch")
