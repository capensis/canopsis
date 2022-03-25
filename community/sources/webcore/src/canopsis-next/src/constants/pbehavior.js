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

/**
 * 19/01/2038 @ 3:14am (UTC) in unix timestamp
 *
 * @type {number}
 */
export const MAX_PBEHAVIOR_DEFAULT_TSTOP = 2147483647;

export const WEATHER_ENTITY_PBEHAVIOR_DEFAULT_TITLE = 'downtime';

export const PBEHAVIOR_PATTERN_FIELDS = {
  name: 'name',
  reason: 'reason',
  type: 'type',
};
