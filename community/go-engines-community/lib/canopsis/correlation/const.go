package correlation

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"

const (
	DefaultMetaAlarmEntityPrefix  = "meta-alarm-entity-"
	DefaultMetaAlarmComponent     = "metaalarm"
	DefaultMetaAlarmConnector     = canopsis.CorrelationConnector
	DefaultMetaAlarmConnectorName = canopsis.CorrelationConnector

	// InfosKeyMetaAlarmMeta
	// todo remove : it was required for event filter rules to set entity infos for entity of meta alarm but now it can be done by rule params
	InfosKeyMetaAlarmMeta = "Meta"

	DefaultConfigTimeInterval = 86400
)
