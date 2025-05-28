package prometheus

import (
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/prometheus/client_golang/prometheus"
)

const canopsisMetricsNamespace = "canopsis"

const (
	EventfilterErrorsGauge = iota
	ResolvedAlarmsGauge
	DisabledEntitiesGauge
	UserConnectionsGauge
	EnabledUsersGauge
	EventFiltersGauge
	ActivePBehaviorsGauge
	MetaAlarmsRulesGauge
	DynamicInfosRulesGauge
)

const (
	EngineStatusGaugeVector = iota
	LastExploitationModTimeGaugeVector
	OpenedAlarmsGaugeVector
	ActiveEntitiesGaugeVector
	InstructionsGaugeVector
)

const (
	EventsRateCounter = iota
)

type Metrics struct {
	gauges         map[int]prometheus.Gauge
	gaugeVectors   map[int]*prometheus.GaugeVec
	counterVectors map[int]*prometheus.CounterVec
}

func (m *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, gaugeVec := range m.gaugeVectors {
		gaugeVec.Describe(ch)
	}

	for _, gauge := range m.gauges {
		gauge.Describe(ch)
	}

	for _, counterVec := range m.counterVectors {
		counterVec.Describe(ch)
	}
}

func (m *Metrics) Collect(ch chan<- prometheus.Metric) {
	for _, gaugeVec := range m.gaugeVectors {
		gaugeVec.Collect(ch)
	}

	for _, gauge := range m.gauges {
		gauge.Collect(ch)
	}

	for _, counterVec := range m.counterVectors {
		counterVec.Collect(ch)
	}
}

func (m *Metrics) CounterVectorAdd(vector int, label string, val float64) {
	counterVec, ok := m.counterVectors[vector]
	if ok {
		counterVec.WithLabelValues(label).Add(val)
	}
}

func (m *Metrics) CounterVectorInc(vector int, label string) {
	counterVec, ok := m.counterVectors[vector]
	if ok {
		counterVec.WithLabelValues(label).Inc()
	}
}

func NewFifoMetrics() *Metrics {
	eventTypes := []string{
		types.EventTypeAck,
		types.EventTypeAckremove,
		types.EventTypeAssocTicket,
		types.EventTypeCancel,
		types.EventTypeCheck,
		types.EventTypeComment,
		types.EventTypeChangestate,
		types.EventTypeSnooze,
		types.EventTypeUnsnooze,
		types.EventTypeUncancel,
		types.EventTypeContextUpdate,
		types.EventTypeDeclareTicketWebhook,
		types.EventTypePbhEnter,
		types.EventTypePbhLeaveAndEnter,
		types.EventTypePbhLeave,
		types.EventTypeResolveCancel,
		types.EventTypeResolveClose,
		types.EventTypeResolveDeleted,
		types.EventTypeUpdateStatus,
		types.EventTypeRunDelayedScenario,
		types.EventTypeMetaAlarm,
		types.EventTypeMetaAlarmAttachChildren,
		types.EventTypeMetaAlarmDetachChildren,
		types.EventTypeInstructionStarted,
		types.EventTypeInstructionPaused,
		types.EventTypeInstructionResumed,
		types.EventTypeInstructionCompleted,
		types.EventTypeInstructionFailed,
		types.EventTypeInstructionAborted,
		types.EventTypeRecomputeEntityService,
		types.EventTypeEntityUpdated,
		types.EventTypeEntityToggled,
		types.EventTypeJunitTestSuiteUpdated,
		types.EventTypeJunitTestCaseUpdated,
		types.EventTypeNoEvents,
		types.EventTypeTrigger,
		types.EventTypeAutoInstructionActivate,
		types.EventTypeMetaAlarmChildActivate,
		types.EventTypeMetaAlarmChildDeactivate,
	}

	m := &Metrics{
		counterVectors: map[int]*prometheus.CounterVec{
			EventsRateCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: canopsisMetricsNamespace,
				Name:      "events_counter",
				Help:      "Number of events",
			}, []string{"type"}),
		},
	}

	for _, eventType := range eventTypes {
		m.CounterVectorAdd(EventsRateCounter, eventType, 0)
	}

	return m
}

type DbCollectionsMetrics struct {
	Metrics
	// collUpdMx is needed to sync Collect and Set methods to avoid sending outdated metrics to prometheus.
	collUpdMx sync.Mutex
}

func NewDbCollectionsMetrics() *DbCollectionsMetrics {
	return &DbCollectionsMetrics{
		Metrics: Metrics{
			gauges: map[int]prometheus.Gauge{
				EventfilterErrorsGauge: prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "eventfilter_errors",
					Help:      "Number of event filter errors",
				}),
				ResolvedAlarmsGauge: prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "resolved_alarms",
					Help:      "Number of resolved alarms",
				}),
				DisabledEntitiesGauge: prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "disabled_entities",
					Help:      "Number of disabled entities",
				}),
				UserConnectionsGauge: prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "user_connections",
					Help:      "Number of user connections",
				}),
				EnabledUsersGauge: prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "enabled_users",
					Help:      "Number of enabled users",
				}),
				EventFiltersGauge: prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "event_filters",
					Help:      "Number of event filters",
				}),
				ActivePBehaviorsGauge: prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "active_pbehavior",
					Help:      "Number of active pbehaviors",
				}),
				MetaAlarmsRulesGauge: prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "meta_alarms_rules",
					Help:      "Number of meta alarm rules",
				}),
				DynamicInfosRulesGauge: prometheus.NewGauge(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "dynamic_infos_rules",
					Help:      "Number of dynamic infos rules",
				}),
			},
			gaugeVectors: map[int]*prometheus.GaugeVec{
				EngineStatusGaugeVector: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "engine_status",
					Help:      "Status of the engine (1 for running, 0 for stopped)",
				}, []string{"engine_name"}),
				LastExploitationModTimeGaugeVector: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "last_exploitation_mod_time",
					Help:      "Last modification time of exploitation menu elements (Unix timestamp)",
				}, []string{"type"}),
				OpenedAlarmsGaugeVector: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "opened_alarms",
					Help:      "Number of opened alarms",
				}, []string{"active"}),
				ActiveEntitiesGaugeVector: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "active_entities",
					Help:      "Number of active entities",
				}, []string{"type"}),
				InstructionsGaugeVector: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: canopsisMetricsNamespace,
					Name:      "instructions",
					Help:      "Number of instructions",
				}, []string{"type"}),
			},
		},
		collUpdMx: sync.Mutex{},
	}
}

func (m *DbCollectionsMetrics) Collect(ch chan<- prometheus.Metric) {
	m.collUpdMx.Lock()
	defer m.collUpdMx.Unlock()

	m.Metrics.Collect(ch)
}

func (m *DbCollectionsMetrics) Set(value *DbCollectionsMetricsValue) {
	m.collUpdMx.Lock()
	defer m.collUpdMx.Unlock()

	value.gaugesMx.Lock()
	for k, v := range value.gauges {
		gauge, ok := m.gauges[k]
		if ok {
			gauge.Set(v)
		}
	}
	value.gaugesMx.Unlock()

	value.gaugeVectorsMx.Lock()
	for k, v := range value.gaugeVectors {
		gaugeVector, ok := m.gaugeVectors[k]
		if ok {
			for label, val := range v {
				gaugeVector.WithLabelValues(label).Set(val)
			}
		}
	}
	value.gaugeVectorsMx.Unlock()
}

type DbCollectionsMetricsValue struct {
	gauges   map[int]float64
	gaugesMx sync.Mutex

	gaugeVectors   map[int]map[string]float64
	gaugeVectorsMx sync.Mutex
}

func newDbCollectionsMetricsValues() *DbCollectionsMetricsValue {
	return &DbCollectionsMetricsValue{
		gauges: map[int]float64{
			EventfilterErrorsGauge: 0,
			ResolvedAlarmsGauge:    0,
			DisabledEntitiesGauge:  0,
			UserConnectionsGauge:   0,
			EnabledUsersGauge:      0,
			EventFiltersGauge:      0,
			ActivePBehaviorsGauge:  0,
			MetaAlarmsRulesGauge:   0,
			DynamicInfosRulesGauge: 0,
		},
		gaugeVectors: map[int]map[string]float64{
			EngineStatusGaugeVector:            make(map[string]float64),
			LastExploitationModTimeGaugeVector: make(map[string]float64),
			OpenedAlarmsGaugeVector:            make(map[string]float64),
			ActiveEntitiesGaugeVector:          make(map[string]float64),
			InstructionsGaugeVector:            make(map[string]float64),
		},
		gaugesMx:       sync.Mutex{},
		gaugeVectorsMx: sync.Mutex{},
	}
}

func (v *DbCollectionsMetricsValue) SetGauge(gauge int, val float64) {
	v.gaugesMx.Lock()
	defer v.gaugesMx.Unlock()

	_, ok := v.gauges[gauge]
	if ok {
		v.gauges[gauge] = val
	}
}

func (v *DbCollectionsMetricsValue) SetGaugeVector(gaugeVector int, label string, val float64) {
	v.gaugeVectorsMx.Lock()
	defer v.gaugeVectorsMx.Unlock()

	vector, ok := v.gaugeVectors[gaugeVector]
	if ok {
		vector[label] = val
	}
}
