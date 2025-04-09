package prometheus

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const canopsisMetricsNamespace = "canopsis"

const (
	EventfilterErrorsGauge = iota
	OpenedAlarmsGauge
	ResolvedAlarmsGauge
	ActiveEntitiesGauge
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
)

type GaugeVec interface {
	prometheus.Collector
	WithLabelValues(lvs ...string) prometheus.Gauge
}

type MetricsValue struct {
	gauges   map[int]float64
	gaugesMx sync.Mutex

	gaugeVectors   map[int]map[string]float64
	gaugeVectorsMx sync.Mutex
}

func newMetricsValues() *MetricsValue {
	return &MetricsValue{
		gauges: map[int]float64{
			EventfilterErrorsGauge: 0,
			OpenedAlarmsGauge:      0,
			ResolvedAlarmsGauge:    0,
			ActiveEntitiesGauge:    0,
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
		},
		gaugesMx:       sync.Mutex{},
		gaugeVectorsMx: sync.Mutex{},
	}
}

func (v *MetricsValue) SetGauge(gauge int, val float64) {
	v.gaugesMx.Lock()
	defer v.gaugesMx.Unlock()

	_, ok := v.gauges[gauge]
	if ok {
		v.gauges[gauge] = val
	}
}

func (v *MetricsValue) SetGaugeVector(gaugeVector int, label string, val float64) {
	v.gaugeVectorsMx.Lock()
	defer v.gaugeVectorsMx.Unlock()

	vector, ok := v.gaugeVectors[gaugeVector]
	if ok {
		vector[label] = val
	}
}

type Metrics struct {
	gauges       map[int]prometheus.Gauge
	gaugeVectors map[int]GaugeVec

	// collUpdMx is needed to sync Collect and Set methods to avoid sending outdated metrics to prometheus.
	collUpdMx sync.Mutex
}

func NewMetrics() *Metrics {
	return &Metrics{
		gauges: map[int]prometheus.Gauge{
			EventfilterErrorsGauge: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: canopsisMetricsNamespace,
				Name:      "eventfilter_errors",
				Help:      "Number of event filter errors",
			}),
			OpenedAlarmsGauge: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: canopsisMetricsNamespace,
				Name:      "opened_alarms",
				Help:      "Number of opened alarms",
			}),
			ResolvedAlarmsGauge: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: canopsisMetricsNamespace,
				Name:      "resolved_alarms",
				Help:      "Number of resolved alarms",
			}),
			ActiveEntitiesGauge: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: canopsisMetricsNamespace,
				Name:      "active_entities",
				Help:      "Number of active entities",
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
		gaugeVectors: map[int]GaugeVec{
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
		},
		collUpdMx: sync.Mutex{},
	}
}

func (m *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, gaugeVec := range m.gaugeVectors {
		gaugeVec.Describe(ch)
	}

	for _, gauge := range m.gauges {
		gauge.Describe(ch)
	}
}

func (m *Metrics) Collect(ch chan<- prometheus.Metric) {
	m.collUpdMx.Lock()
	defer m.collUpdMx.Unlock()

	for _, gaugeVec := range m.gaugeVectors {
		gaugeVec.Collect(ch)
	}

	for _, gauge := range m.gauges {
		gauge.Collect(ch)
	}
}

func (m *Metrics) Set(value *MetricsValue) {
	m.collUpdMx.Lock()
	value.gaugesMx.Lock()
	value.gaugeVectorsMx.Lock()

	defer func() {
		m.collUpdMx.Unlock()
		value.gaugesMx.Unlock()
		value.gaugeVectorsMx.Unlock()
	}()

	for k, v := range value.gauges {
		gauge, ok := m.gauges[k]
		if ok {
			gauge.Set(v)
		}
	}

	for k, v := range value.gaugeVectors {
		gaugeVector, ok := m.gaugeVectors[k]
		if ok {
			for label, val := range v {
				gaugeVector.WithLabelValues(label).Set(val)
			}
		}
	}
}
