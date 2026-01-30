export const PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES = {
  selected: 0,
  all: 1,
};

export const PLANNING_TABS = {
  types: 'types',
  reasons: 'reasons',
  exceptions: 'exceptions',
};

export const PBEHAVIOR_RRULE_PERIODS_RANGES = {
  thisWeek: 'thisWeek',
  nextWeek: 'nextWeek',
  next2Weeks: 'next2Weeks',
  thisMonth: 'thisMonth',
  nextMonth: 'nextMonth',
};

export const PBEHAVIOR_TYPE_TYPES = {
  active: 'active',
  inactive: 'inactive',
  maintenance: 'maintenance',
  pause: 'pause',
};

export const PBEHAVIOR_TYPE_TYPES_ICONS = {
  [PBEHAVIOR_TYPE_TYPES.active]: 'wb_sunny',
  [PBEHAVIOR_TYPE_TYPES.inactive]: 'nightlight_round',
  [PBEHAVIOR_TYPE_TYPES.maintenance]: 'build',
  [PBEHAVIOR_TYPE_TYPES.pause]: 'pause',
};

export const WEATHER_ENTITY_PBEHAVIOR_DEFAULT_TITLE = 'downtime';

export const PBEHAVIOR_PATTERN_PREFIX = 'pbehavior_info.';

export const PBEHAVIOR_PATTERN_FIELDS = {
  name: `${PBEHAVIOR_PATTERN_PREFIX}id`,
  reason: `${PBEHAVIOR_PATTERN_PREFIX}reason`,
  type: `${PBEHAVIOR_PATTERN_PREFIX}type`,
  canonicalType: `${PBEHAVIOR_PATTERN_PREFIX}canonical_type`,
};

export const PBEHAVIOR_INFO_FIELDS = {
  typeName: 'type_name',
  reason: 'reason',
  name: 'name',
  canonicalType: 'canonical_type',
};

export const PBEHAVIOR_ORIGINS = {
  alarmList: 'AlarmList',
  serviceWeather: 'ServiceWeather',
};

export const PBEHAVIOR_CANONICAL_TYPES = {
  pause: 'pause',
  active: 'active',
};

export const PBEHAVIOR_LIST_FIELDS = {
  name: 'name',
  author: 'author.display_name',
  enabled: 'enabled',
  begins: 'tstart',
  ends: 'tstop',
  rruleEnd: 'rrule_end',
  rrule: 'rrule',
  type: 'type.name',
  reason: 'reason.name',
  created: 'created',
  updated: 'updated',
  lastAlarmDate: 'last_alarm_date',
  alarmCount: 'alarm_count',
  patternMs: 'pattern_ms',
  patternExecAt: 'pattern_exec_at',
  typeIcon: 'type.icon_name',
  status: 'is_active_status',
  actions: 'actions',
};
