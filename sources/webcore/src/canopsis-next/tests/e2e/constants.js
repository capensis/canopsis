const RANDOM = Math.random().toString(36).substring(6);

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
};

module.exports.CONTEXT_WIDGET_SORT_FIELD = {
  name: 1,
  type: 2,
};


module.exports.SORT_ORDERS = {
  asc: 1,
  desc: 2,
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
