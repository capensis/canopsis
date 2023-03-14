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
  adminEngines: 'admin-engines',
  adminKPI: 'admin-kpi',
  adminMaps: 'admin-maps',
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
  adminEngines: '/admin/engines',
  adminKPI: '/admin/kpi',
  adminMaps: '/admin/maps',
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

export const EXCLUDED_SERVER_ERROR_STATUSES = [401, 403];

export const CANOPSIS_DOCUMENTATION = 'https://doc.canopsis.net';

export const CANOPSIS_WEBSITE = 'https://www.capensis.fr/canopsis/';

export const CANOPSIS_FORUM = 'https://community.capensis.org/';

export const TOURS = {
  alarmsExpandPanel: 'alarmsExpandPanel',
};

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
  last2Days: {
    value: 'last2Days',
    start: 'now-2d',
    stop: 'now',
  },
  last7Days: {
    value: 'last7Days',
    start: 'now-7d',
    stop: 'now',
  },
  last30Days: {
    value: 'last30Days',
    start: 'now-30d',
    stop: 'now',
  },
  last1Year: {
    value: 'last1Year',
    start: 'now-1y',
    stop: 'now',
  },
  yesterday: {
    value: 'yesterday',
    start: 'now-1d/d',
    stop: 'now-1d/d',
  },
  previousWeek: {
    value: 'previousWeek',
    start: 'now-1w/w',
    stop: 'now-1w/w',
  },
  previousMonth: {
    value: 'previousMonth',
    start: 'now-1M/M',
    stop: 'now-1M/M',
  },
  today: {
    value: 'today',
    start: 'now/d',
    stop: 'now/d',
  },
  todaySoFar: {
    value: 'todaySoFar',
    start: 'now/d',
    stop: 'now',
  },
  thisWeek: {
    value: 'thisWeek',
    start: 'now/w',
    stop: 'now/w',
  },
  thisWeekSoFar: {
    value: 'thisWeekSoFar',
    start: 'now/w',
    stop: 'now',
  },
  thisMonth: {
    value: 'thisMonth',
    start: 'now/M',
    stop: 'now/M',
  },
  thisMonthSoFar: {
    value: 'thisMonthSoFar',
    start: 'now/M',
    stop: 'now',
  },
  custom: {
    value: 'custom',
  },
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

export const ASSOCIATIVE_TABLES_NAMES = {
  alarmColumnsFilters: 'alarm-column-filters',
  dynamicInfoTemplates: 'dynamic-info-templates',
  filterHints: 'filterhints',
};

export const DEFAULT_COLUMN_TEMPLATE_VALUE = '{{ value }}';

export const VALIDATION_DELAY = 300;

export const TRIGGERS = {
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
};

export const PRO_TRIGGERS = [
  TRIGGERS.declareticket,
  TRIGGERS.declareticketwebhook,
  TRIGGERS.instructionfail,
  TRIGGERS.autoinstructionfail,
  TRIGGERS.instructionjobcomplete,
  TRIGGERS.instructionjobfail,
  TRIGGERS.instructioncomplete,
  TRIGGERS.autoinstructioncomplete,
];

export const DEPRECATED_TRIGGERS = [
  TRIGGERS.declareticket,
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
