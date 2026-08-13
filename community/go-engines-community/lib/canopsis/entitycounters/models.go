package entitycounters

type StateSettingsInfo struct {
	ComponentStateSettings         bool `bson:"component_state_settings"`
	ComponentStateSettingsToAdd    bool `bson:"component_state_settings_to_add"`
	ComponentStateSettingsToRemove bool `bson:"component_state_settings_to_remove"`
}

type ServicesInfo struct {
	Services         []string `bson:"services"`
	ServicesToAdd    []string `bson:"services_to_add"`
	ServicesToRemove []string `bson:"services_to_remove"`

	InheritedServices []string `bson:"inherited_services"`
}

type ComponentCountersCalcData struct {
	Info        StateSettingsInfo
	Counters    EntityCounters
	PrevState   int
	CurState    int
	PrevActive  bool
	CurActive   bool
	AlarmExists bool
}

type EntityServiceCountersCalcData struct {
	ServicesToAdd     map[string]bool
	ServicesToRemove  map[string]bool
	InheritedServices map[string]bool
	Counters          EntityCounters
	PrevPbhTypeID     string
	CurPbhTypeID      string
	PrevState         int
	CurState          int
	PrevActive        bool
	CurActive         bool
	IsAcked           bool
	PrevInherited     bool
	CurInherited      bool
	AlarmExists       bool
	EntityEnabled     bool
}

type UpdatedServicesInfo struct {
	State  int
	Output string
}
