export const ROUTES_NAMES = {
  home: 'home',
  error: 'error',
  login: 'login',
  view: 'view',
  viewKiosk: 'view-kiosk',
  alarms: 'alarms',
  adminRights: 'admin-rights',
  adminUsers: 'admin-users',
  adminRoles: 'admin-roles',
  adminParameters: 'admin-parameters',
  adminBroadcastMessages: 'admin-broadcast-messages',
  adminPlaylists: 'admin-playlists',
  adminPlanning: 'admin-planning-administration',
  adminRemediation: 'admin-remediation-administration',
  adminHealthcheck: 'admin-healthcheck',
  adminKPI: 'admin-kpi',
  adminMaps: 'admin-maps',
  adminTags: 'admin-tags',
  adminStorageSettings: 'admin-storage-settings',
  adminStateSettings: 'admin-state-settings',
  adminEventsRecordings: 'admin-events-recordings',
  exploitationPbehaviors: 'exploitation-pbehaviors',
  exploitationEventFilters: 'exploitation-event-filters',
  exploitationSnmpRules: 'exploitation-snmp-rules',
  exploitationDynamicInfos: 'exploitation-dynamic-infos',
  playlist: 'playlist',
  exploitationMetaAlarmRules: 'exploitation-meta-alarm-rules',
  exploitationScenarios: 'exploitation-scenarios',
  exploitationIdleRules: 'exploitation-idle-rules',
  exploitationFlappingRules: 'exploitation-flapping-rules',
  exploitationResolveRules: 'exploitation-resolve-rules',
  exploitationDeclareTicketRules: 'exploitation-declare-ticket-rules',
  exploitationLinkRules: 'exploitation-link-rules',
  notificationInstructionStats: 'notification-instruction-stats',
  profilePatterns: 'profile-patterns',
  profileThemes: 'profile-themes',
};

export const ROUTES = {
  home: '/',
  login: '/login',
  error: '/error',
  view: '/view/:id',
  viewKiosk: '/kiosk-views/:id/:tabId',
  alarms: '/alarms/:id',
  adminRights: '/admin/rights',
  adminUsers: '/admin/users',
  adminRoles: '/admin/roles',
  adminParameters: '/admin/parameters',
  adminBroadcastMessages: '/admin/broadcast-messages',
  adminPlaylists: '/admin/playlists',
  adminPlanning: '/admin/planning',
  adminRemediation: '/admin/remediation',
  adminHealthcheck: '/admin/healthcheck',
  adminKPI: '/admin/kpi',
  adminMaps: '/admin/maps',
  adminTags: '/admin/tags',
  adminStorageSettings: '/admin/storage-settings',
  adminStateSettings: '/admin/state-settings',
  adminEventsRecordings: '/admin/events-recordings',
  exploitationPbehaviors: '/exploitation/pbehaviors',
  exploitationEventFilters: '/exploitation/event-filters',
  exploitationSnmpRules: '/exploitation/snmp-rules',
  exploitationDynamicInfos: '/exploitation/dynamic-infos',
  playlist: '/playlist/:id',
  exploitationMetaAlarmRules: '/exploitation/meta-alarm-rule',
  exploitationScenarios: '/exploitation/scenarios',
  exploitationIdleRules: '/exploitation/idle-rules',
  exploitationFlappingRules: '/exploitation/flapping-rules',
  exploitationResolveRules: '/exploitation/resolve-rules',
  exploitationDeclareTicketRules: '/exploitation/declare-ticket-rules',
  exploitationLinkRules: '/exploitation/link-rules',
  notificationInstructionStats: '/notification/instruction-stats',
  profilePatterns: '/filters',
  profileThemes: '/themes',
};

export const GROUPS_NAVIGATION_TYPES = {
  sideBar: 'side-bar',
  topBar: 'top-bar',
};

export const SORT_ORDERS = {
  asc: 'ASC',
  desc: 'DESC',
};

export const POPUP_TYPES = {
  success: 'success',
  info: 'info',
  warning: 'warning',
  error: 'error',
};

export const PORTALS_NAMES = {
  additionalTopBarItems: 'additional-top-bar-items',
};

export const MAX_LIMIT = 10000;

export const REQUEST_METHODS = {
  post: 'POST',
  get: 'GET',
  put: 'PUT',
  patch: 'PATCH',
  delete: 'DELETE',
  head: 'HEAD',
  connect: 'CONNECT',
  options: 'OPTIONS',
  trace: 'TRACE',
};

export const DAYS_IN_YEAR = 365;

export const DAYS_IN_MONTH = 30;

export const MONTHS_IN_YEAR = 12;

export const DAYS_IN_WEEK = 7;

export const CANOPSIS_DOCUMENTATION = 'https://doc.canopsis.net';

export const CANOPSIS_WEBSITE = 'https://www.capensis.fr/canopsis/';

export const CANOPSIS_FORUM = 'https://community.capensis.org/';

export const QUICK_RANGES = {
  last15Minutes: {
    value: 'last15Minutes',
    start: 'now-15m',
    stop: 'now',
  },
  last30Minutes: {
    value: 'last30Minutes',
    start: 'now-30m',
    stop: 'now',
  },
  last1Hour: {
    value: 'last1Hour',
    start: 'now-1h',
    stop: 'now',
  },
  last3Hour: {
    value: 'last3Hour',
    start: 'now-3h',
    stop: 'now',
  },
  last6Hour: {
    value: 'last6Hour',
    start: 'now-6h',
    stop: 'now',
  },
  last12Hour: {
    value: 'last12Hour',
    start: 'now-12h',
    stop: 'now',
  },
  last24Hour: {
    value: 'last24Hour',
    start: 'now-24h',
    stop: 'now',
  },
  /**
   * We've changed start from 'today-2d' to current to avoid problem with more bars on KPI
   */
  last2Days: {
    value: 'last2Days',
    start: 'today-1d',
    stop: 'today',
  },
  /**
   * We've changed start from 'today-7d' to current to avoid problem with more bars on KPI
   */
  last7Days: {
    value: 'last7Days',
    start: 'today-6d',
    stop: 'today',
  },
  /**
   * We've changed start from 'today-30d' to current to avoid problem with more bars on KPI
   */
  last30Days: {
    value: 'last30Days',
    start: 'today-29d',
    stop: 'today',
  },
  last1Year: {
    value: 'last1Year',
    start: 'today-1y',
    stop: 'today',
  },
  yesterday: {
    value: 'yesterday',
    start: 'today-1d/d',
    stop: 'today-1d/d',
  },
  previousWeek: {
    value: 'previousWeek',
    start: 'today-1w/w',
    stop: 'today-1w/w',
  },
  previousMonth: {
    value: 'previousMonth',
    start: 'today-1M/M',
    stop: 'today-1M/M',
  },
  today: {
    value: 'today',
    start: 'today/d',
    stop: 'today/d',
  },
  todaySoFar: {
    value: 'todaySoFar',
    start: 'today/d',
    stop: 'now',
  },
  thisWeek: {
    value: 'thisWeek',
    start: 'today/w',
    stop: 'today/w',
  },
  thisWeekSoFar: {
    value: 'thisWeekSoFar',
    start: 'today/w',
    stop: 'now',
  },
  thisMonth: {
    value: 'thisMonth',
    start: 'today/M',
    stop: 'today/M',
  },
  thisMonthSoFar: {
    value: 'thisMonthSoFar',
    start: 'today/M',
    stop: 'now',
  },
  last3Months: {
    value: 'last3Months',
    start: 'today-3M/M',
    stop: 'today-1M/M',
  },
  last6Months: {
    value: 'last6Months',
    start: 'today-6M/M',
    stop: 'today-1M/M',
  },
  custom: {
    value: 'custom',
  },
};

export const LIVE_REPORTING_QUICK_RANGES = {
  ...QUICK_RANGES,

  [QUICK_RANGES.last2Days.value]: {
    ...QUICK_RANGES.last2Days,
    stop: 'now',
  },
  [QUICK_RANGES.last7Days.value]: {
    ...QUICK_RANGES.last7Days,
    stop: 'now',
  },
  [QUICK_RANGES.last30Days.value]: {
    ...QUICK_RANGES.last30Days,
    stop: 'now',
  },
  [QUICK_RANGES.last1Year.value]: {
    ...QUICK_RANGES.last1Year,
    stop: 'now',
  },
};

export const METRICS_QUICK_RANGES = {
  [QUICK_RANGES.last2Days.value]: QUICK_RANGES.last2Days,
  [QUICK_RANGES.last7Days.value]: QUICK_RANGES.last7Days,
  [QUICK_RANGES.last30Days.value]: QUICK_RANGES.last30Days,
  [QUICK_RANGES.last1Year.value]: QUICK_RANGES.last1Year,
  [QUICK_RANGES.yesterday.value]: QUICK_RANGES.yesterday,
  [QUICK_RANGES.previousWeek.value]: QUICK_RANGES.previousWeek,
  [QUICK_RANGES.previousMonth.value]: QUICK_RANGES.previousMonth,
  [QUICK_RANGES.today.value]: QUICK_RANGES.today,
  [QUICK_RANGES.thisWeek.value]: QUICK_RANGES.thisWeek,
  [QUICK_RANGES.thisWeekSoFar.value]: QUICK_RANGES.thisWeekSoFar,
  [QUICK_RANGES.thisMonth.value]: QUICK_RANGES.thisMonth,
  [QUICK_RANGES.thisMonthSoFar.value]: QUICK_RANGES.thisMonthSoFar,
  [QUICK_RANGES.custom.value]: QUICK_RANGES.custom,
};

export const HEALTHCHECK_QUICK_RANGES = {
  [QUICK_RANGES.last3Hour.value]: QUICK_RANGES.last3Hour,
  [QUICK_RANGES.last6Hour.value]: QUICK_RANGES.last6Hour,
  [QUICK_RANGES.last12Hour.value]: QUICK_RANGES.last12Hour,
  [QUICK_RANGES.last24Hour.value]: QUICK_RANGES.last24Hour,

  ...METRICS_QUICK_RANGES,
};

export const GRID_SIZES = {
  min: 0,
  max: 12,
  step: 1,
};

export const DEFAULT_TIMEZONE = 'Europe/Paris';

export const WORKFLOW_TYPES = {
  stop: true,
  continue: false,
};

export const PAYLOAD_VARIABLE_REGEXP = /({{(.){2,}}})/g;

export const PAYLOAD_OPERATORS = ['range', 'end'];

export const PATTERNS_TYPES = {
  alarm: 'alarm',
  event: 'event',
  entity: 'entity',
  totalEntity: 'totalEntity',
};

export const AUTH_METHODS = {
  local: 'local',
};

export const AUTH_SOURCES_WITH_PASSWORD_CHANGING = {
  basic: '',
  ldap: 'ldap',
  cas: 'cas',
  saml: 'saml',
};

export const ASSOCIATIVE_TABLES_NAMES = {
  alarmColumnsFilters: 'alarm-column-filters',
  dynamicInfoTemplates: 'dynamic-info-templates',
  filterHints: 'filterhints',
};

export const DEFAULT_COLUMN_TEMPLATE_VALUE = '{{ value }}';

export const VALIDATION_DELAY = 300;

export const TRIGGERS_TYPES = {
  create: 'create',
  stateinc: 'stateinc',
  statedec: 'statedec',
  changestate: 'changestate',
  changestatus: 'changestatus',
  ack: 'ack',
  ackremove: 'ackremove',
  cancel: 'cancel',
  uncancel: 'uncancel',
  comment: 'comment',
  declareticket: 'declareticket',
  declareticketwebhook: 'declareticketwebhook',
  assocticket: 'assocticket',
  snooze: 'snooze',
  unsnooze: 'unsnooze',
  resolve: 'resolve',
  activate: 'activate',
  pbhenter: 'pbhenter',
  pbhleave: 'pbhleave',
  instructionfail: 'instructionfail',
  autoinstructionfail: 'autoinstructionfail',
  instructionjobcomplete: 'instructionjobcomplete',
  instructionjobfail: 'instructionjobfail',
  instructioncomplete: 'instructioncomplete',
  autoinstructioncomplete: 'autoinstructioncomplete',
  autoinstructionresultok: 'autoinstructionresultok',
  autoinstructionresultfail: 'autoinstructionresultfail',
  eventscount: 'eventscount',
};

export const PRO_TRIGGERS = [
  TRIGGERS_TYPES.declareticket,
  TRIGGERS_TYPES.declareticketwebhook,
  TRIGGERS_TYPES.instructionfail,
  TRIGGERS_TYPES.autoinstructionfail,
  TRIGGERS_TYPES.instructionjobcomplete,
  TRIGGERS_TYPES.instructionjobfail,
  TRIGGERS_TYPES.instructioncomplete,
  TRIGGERS_TYPES.autoinstructioncomplete,
  TRIGGERS_TYPES.autoinstructionresultok,
  TRIGGERS_TYPES.autoinstructionresultfail,
];

export const DEPRECATED_TRIGGERS = [
  TRIGGERS_TYPES.declareticket,
];

export const HEADERS = {
  authorization: 'Authorization',
  contentType: 'Content-Type',
};

export const CONTENT_TYPES = {
  json: 'application/json',
  javascript: 'application/javascript',
  xWwwFormUrlencoded: 'application/x-www-form-urlencoded',
};

export const LOGIN_APP_INFO_POLLING_DELAY = 5000;

export const RESPONSE_STATUSES = {
  badRequest: 400,
  unauthorized: 401,
  notFound: 404,
  forbidden: 403,
  timeout: 408,

  internalServerError: 500,
  badGateway: 502,
  serviceUnavailable: 503,
};

export const LOGOUT_RESPONSE_STATUSES = [
  RESPONSE_STATUSES.unauthorized,
  RESPONSE_STATUSES.badGateway,
  RESPONSE_STATUSES.serviceUnavailable,
];

export const EXCLUDED_SERVER_ERROR_STATUSES = [
  RESPONSE_STATUSES.unauthorized,
  RESPONSE_STATUSES.forbidden,
];

export const PAYLOADS_INFO_VARIABLES = {
  value: '.Value',
  description: '.Description',
};

export const WEBHOOK_EXECUTION_STATUSES = {
  waiting: 0,
  running: 1,
  succeeded: 2,
  failed: 3,
};

export const KEY_CODES = {
  backspace: 8,
  enter: 13,
  shift: 16,
  esc: 27,
  left: 37,
  up: 38,
  right: 39,
  down: 40,
  end: 35,
  home: 36,
  delete: 46,
  c: 67,
  v: 86,
};

export const VARIABLES_STEP_FIELDS = {
  timestamp: 't',
  value: 'val',
  message: 'm',
  author: 'a',
};
