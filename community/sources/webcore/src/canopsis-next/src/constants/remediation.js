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

export const REMEDIATION_TABS = {
  instructions: 'instructions',
  configurations: 'configurations',
  jobs: 'jobs',
};

export const REMEDIATION_INSTRUCTION_TYPES = {
  manual: 0,
  auto: 1,
};

export const REMEDIATION_INSTRUCTION_APPROVAL_TYPES = {
  role: 0,
  user: 1,
};

export const REMEDIATION_INSTRUCTION_EXECUTION_STATUSES = {
  running: 0,
  paused: 1,
  completed: 2,
  aborted: 3,
  failed: 4,
};

export const REMEDIATION_JOB_EXECUTION_STATUSES = {
  running: 0,
  succeeded: 1,
  failed: 2,
  canceled: 3,
};

/**
 * 19/01/2038 @ 3:14am (UTC) in unix timestamp
 *
 * @type {number}
 */
export const MAX_PBEHAVIOR_DEFAULT_TSTOP = 2147483647;

export const REMEDIATION_CONFIGURATION_JOBS_AUTH_TYPES_WITH_USERNAME = ['basic-auth'];

export const WEATHER_ENTITY_PBEHAVIOR_DEFAULT_TITLE = 'downtime';

export const ASSOCIATIVE_TABLES_NAMES = {
  alarmColumnsFilters: 'alarm-column-filters',
  dynamicInfoTemplates: 'dynamic-info-templates',
  filterHints: 'filterhints',
};

export const PBEHAVIOR_TYPE_TYPES = {
  active: 'active',
  inactive: 'inactive',
  maintenance: 'maintenance',
  pause: 'pause',
};
