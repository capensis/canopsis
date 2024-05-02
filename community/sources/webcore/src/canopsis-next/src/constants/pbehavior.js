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

export const PBEHAVIOR_PATTERN_FIELDS = {
  name: 'pbehavior_info.id',
  reason: 'pbehavior_info.reason',
  type: 'pbehavior_info.type',
  canonicalType: 'pbehavior_info.canonical_type',
};

export const PBEHAVIOR_INFO_FIELDS = {
  typeName: 'type_name',
  reason: 'reason',
  name: 'name',
  canonicalType: 'canonical_type',
};

export const PBEHAVIOR_ORIGINS = {
  serviceWeather: 'ServiceWeather',
};
