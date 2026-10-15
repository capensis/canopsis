const RANDOM = Math.random().toString(36).substring(6);

module.exports.WAIT_FOR_FIRST_XHR_TIME = 5000;

module.exports.CREDENTIALS = {
  admin: {
    username: 'root',
    password: 'root',
  },
};

module.exports.USERS = {
  text: 'user',
  counts: 5,
  create: {
    prefix: 'create',
  },
  edit: {
    prefix: 'edit',
  },
  mass: {
    prefix: 'mass',
  },
};

module.exports.ROLES = {
  text: 'role',
  counts: 3,
};

module.exports.NAVIGATION = {
  groups: {
    text: 'view',
    counts: 2,
    create: {
      prefix: 'create',
    },
    copy: {
      prefix: 'copy',
    },
    edit: {
      prefix: 'edit',
    },
  },
};

module.exports.USER = {
  username: 'test-username',
  firstname: 'test-firstname',
  lastname: 'test-lastname',
  email: 'test-email@example.com',
  password: 'test-password',
  role: 1,
  language: 2,
  navigationType: 1,
};

module.exports.VIEW = {
  name: `test-name-${RANDOM}`,
  title: `test-title-${RANDOM}`,
  description: `test-description-${RANDOM}`,
  group: `test-group-${RANDOM}`,
};

module.exports.SERVICE_WEATHER_WIDGET_MODAL_TYPES = {
  moreInfo: 1,
  alarmList: 2,
  both: 3,
};

module.exports.SERVICE_WEATHER_WIDGET_SORT_FIELD = {
  name: 1,
  state: 2,
  status: 3,
  criticity: 4,
  org: 5,
};

module.exports.ALARMS_WIDGET_SORT_FIELD = {
  connector: 1,
  connectorName: 2,
  component: 3,
  resource: 4,
  output: 5,
  extraDetails: 6,
  state: 7,
  status: 8,
};

module.exports.CONTEXT_WIDGET_SORT_FIELD = {
  name: 1,
  type: 2,
};

module.exports.SORT_ORDERS = {
  asc: 1,
  desc: 2,
};

module.exports.SORT_ORDERS_STRING = {
  asc: 'ascending',
  desc: 'descending',
};

module.exports.LANGUAGES_POSITIONS = {
  en: 1,
  fr: 2,
};

module.exports.PAGINATION_PER_PAGE_VALUES = {
  FIVE: 1,
  TEN: 2,
  TWENTY: 3,
  FIFTY: 4,
  HUNDRED: 5,
};

module.exports.INFO_POPUP_DEFAULT_COLUMNS = {
  connector: 1,
  connectorName: 2,
  component: 3,
  resource: 4,
  output: 5,
  extraDetails: 6,
  state: 7,
  status: 8,
};

module.exports.FILTERS_TYPE = {
  AND: 'and',
  OR: 'or',
};

module.exports.ROW_SIZE_KEYS = {
  SMARTPHONE: 'sm',
  TABLET: 'md',
  DESKTOP: 'lg',
};

module.exports.ROW_SIZE_CLASSES = {
  SMARTPHONE: 'xs',
  TABLET: 'md',
  DESKTOP: 'lg',
};

module.exports.VALUE_TYPES = {
  STRING: 1,
  NUMBER: 2,
  BOOLEAN: 3,
  NULL: 4,
};

module.exports.INTERVAL_RANGES = {
  CUSTOM: 1,
  LAST_TWO_DAY: 2,
  LAST_SEVEN_DAY: 3,
  LAST_THIRTY_DAY: 4,
  LAST_YEAR: 5,
  YESTERDAY: 6,
  PREVIOUS_WEEK: 7,
  PREVIOUS_MONTH: 8,
  TODAY: 9,
  TODAY_SO_FAR: 10,
  THIS_WEEK: 11,
  THIS_WEEK_SO_FAR: 12,
  THIS_MONTH: 13,
  THIS_MONTH_SO_FAR: 14,
  LAST_HOUR: 15,
  LAST_THREE_HOURS: 16,
  LAST_SIX_HOURS: 17,
  LAST_TWELVE_HOURS: 18,
  LAST_TWENTY_FOUR_HOURS: 19,
};

module.exports.FILTER_OPERATORS = {
  EQUAL: 1,
  NOT_EQUAL: 2,
  IN: 3,
  NOT_IN: 4,
  BEGINS_WITH: 5,
  DOES_NOT_BEGINS_WITH: 6,
  CONTAINS: 7,
  DOES_NOT_CONTAINS: 8,
  ENDS_WITH: 9,
  DOES_NOT_ENDS_WITH: 10,
  IS_EMPTY: 11,
  IS_NOT_EMPTY: 12,
  IS_NULL: 13,
  IS_NOT_NULL: 14,
};

module.exports.FILTER_COLUMNS = {
  CONNECTOR: 1,
  CONNECTOR_NAME: 2,
  COMPONENT: 3,
  RESOURCE: 4,
};

module.exports.NAVIGATION_TYPES = {
  sideBar: 1,
  topBar: 2,
};

module.exports.CONTEXT_FILTER_COLUMNS = {
  NAME: 1,
  TYPE: 2,
};

module.exports.CONTEXT_TYPE_OF_ENTITIES = {
  COMPONENT: 1,
  CONNECTOR: 2,
  RESOURCE: 3,
  WATCHER: 4,
};

module.exports.STAT_TYPES = {
  ALARMS_CREATED: 1,
  ALARMS_RESOLVED: 2,
  ALARMS_CANCELED: 3,
  ALARMS_ACKNOWLEDGED: 4,
  ACK_TIME_SLA: 5,
  RESOLVED_TIME_SLA: 6,
  TIME_IN_STATE: 7,
  STATE_RATE: 8,
  MTBF: 9,
  CURRENT_STATE: 10,
  ONGOING_ALARMS: 11,
  CURRENT_ONGOING_ALARMS: 12,
  CURRENT_ONGOING_ALARMS_WITH_ACK: 13,
  CURRENT_ONGOING_ALARMS_WITHOUT_ACK: 14,
};

module.exports.STAT_STATES = {
  OK: 1,
  MINOR: 2,
  MAJOR: 3,
  CRITICAL: 4,
};

module.exports.INTERVAL_PERIODS = {
  HOUR: 1,
  DAY: 2,
  WEEK: 3,
  MONTH: 4,
};

module.exports.STATS_CURVES_POINTS_STYLES = {
  CIRCLE: 1,
  CROSS: 2,
  CROSS_ROT: 3,
  DASH: 4,
  LINE: 5,
  RECT: 6,
  RECT_ROUNDED: 7,
  RECT_ROT: 8,
  STAR: 9,
  TRIANGLE: 10,
};

module.exports.STATS_DISPLAY_MODE = {
  VALUE: 1,
  CRITICITY: 2,
};

module.exports.PARETO_STATS_TITLES = {
  ACCUMULATION: 'Accumulation',
  ALARMS_CREATED: 'Alarmes créées',
};

module.exports.ALARMS_MASS_ACTIONS = {
  PERIODICAL_BEHAVIOR: 1,
  ACK: 2,
  FAST_ACK: 3,
  CANCEL_ACK: 4,
  CANCEL_ALARM: 5,
};

module.exports.ALARMS_SHARED_ACTIONS = {
  ACK: 1,
  FAST_ACK: 2,
  SNOOZE_ALARM: 3,
  CANCEL_ALARM: 4,
};

module.exports.ALARMS_RESOLVED_SHARED_ACTIONS = {
  MORE_INFOS: 1,
};

module.exports.ALARMS_SHARED_ACTIONS_WITH_ACK = {
  DECLARE_TICKET: 1,
  ASSOCIATE_TICKET: 2,
  CANCEL_ALARM: 3,
};

module.exports.ALARMS_SHARED_DROPDOWN_ACTIONS = {
  PERIODICAL_BEHAVIOR: 1,
  LIST_PERIODICAL_BEHAVIOR: 2,
  MORE_INFOS: 3,
};

module.exports.ALARMS_SHARED_DROPDOWN_ACTIONS_WITH_ACK = {
  CANCEL_ACK: 1,
  CHANGE_CRITICITY: 2,
  SNOOZE_ALARM: 3,
  PERIODICAL_BEHAVIOR: 4,
  LIST_PERIODICAL_BEHAVIOR: 5,
};

module.exports.DATE_INTERVAL_MINUTES = {
  ZERO: 1,
  FIVE: 2,
  TEN: 3,
  twenty: 4,
  TWENTY_FIVE: 5,
  THIRTY: 6,
  THIRTY_FIVE: 7,
  FORTY: 8,
  FORTY_FIVE: 9,
  FIFTY: 10,
  FIFTY_FIVE: 11,
};

module.exports.PERIODICAL_BEHAVIOR_REASONS = {
  REHABILITATION_PROBLEM: 1,
  ROBOT_PROBLEM: 2,
  SCENARIO_PROBLEM: 3,
  FUNCTIONAL_FLASH_ISSUE: 4,
  OTHER: 5,
};

module.exports.PERIODICAL_BEHAVIOR_FREQUENCY = {
  SECONDLY: 1,
  MINUTELY: 2,
  HOURLY: 3,
  DAILY: 4,
  WEEKLY: 5,
  MONTHLY: 6,
  YEARLY: 7,
};

module.exports.WEEK_DAYS = {
  MONDAY: 1,
  TUESDAY: 2,
  WEDNEDDAY: 3,
  THURSDAY: 4,
  FRIDAY: 5,
  SATURDAY: 6,
  SUNDAY: 7,
};

module.exports.MONTH = {
  JUNUARY: 1,
  FEBRARY: 2,
  MARCH: 3,
  APRIL: 4,
  MAY: 5,
  JUNE: 6,
  JULY: 7,
  AUGUST: 8,
  SEPTEMBER: 9,
  OCTOBER: 10,
  NOVEMBER: 11,
  DECEMBER: 12,
};

module.exports.PBEHAVIOR_TYPES = {
  MAINTENANCE: 1,
  UNMONITORED: 2,
  PAUSE: 3,
};

module.exports.SNOOZE_TYPES = {
  MINUTES: 1,
  HOURS: 2,
  DAYS: 3,
  WEEKS: 4,
  MONTHS: 5,
  YEARS: 6,
};

module.exports.PBEHAVIOR_STEPS = {
  GENERAL: 'general',
  FILTER: 'filter',
  RRULE: 'rrule',
  COMMENTS: 'comments',
};

module.exports.CRUD_ACTIONS = {
  DELETE: 'delete',
  UPDATE: 'update',
  CREATE: 'create',
};
