import { removeTrailingSlashes } from '@/helpers/url';

export const { BASE_URL, VUE_APP_COOKIE_SESSION_KEY } = process.env;

export const API_BASE_URL = '/backend';

export const APP_HOST = removeTrailingSlashes(`${window.location.origin}${BASE_URL}`);

export const API_HOST = removeTrailingSlashes(`${APP_HOST}${API_BASE_URL}`);

export const ROUTER_MODE = 'history';

export const COOKIE_SESSION_KEY = VUE_APP_COOKIE_SESSION_KEY;

export const PAGINATION_LIMIT = parseInt(process.env.VUE_APP_PAGINATION_LIMIT, 10);

export const PAGINATION_PER_PAGE_VALUES = [5, 10, 20, 50, 100];

export const PAGINATION_TOTAL_VISIBLE = 7;

export const DEFAULT_WEATHER_LIMIT = 120;

export const DEFAULT_LOCALE = 'en';

export const DEFAULT_KEEP_ALIVE_INTERVAL = 120000;

export const LOCALE_PRIORITIES = {
  default: 1,
  global: 2,
  personal: 3,
};

export const POPUP_AUTO_CLOSE_DELAY = 3000;

export const ACTIVE_SESSIONS_COUNT_FETCHING_INTERVAL = 60000;

export const ACTIVE_BROADCAST_MESSAGE_FETCHING_INTERVAL = 60000;

export const VUETIFY_ANIMATION_DELAY = 300;

export const POPUP_TICK_DELAY = 100;

export const SIDE_BAR_WIDTH = 250;

export const EXPAND_DEFAULT_MAX_LETTERS = 50;

export const SCHEMA_EMBEDDED_KEY = '_embedded';

export const MEDIA_QUERIES_BREAKPOINTS = {
  m: 450,
  t: 900,
  l: 1400,
  xl: 'Infinity',
};

export const MAX_LOGO_SIZE_IN_KB = 16384;

export const DEFAULT_APP_TITLE = 'Canopsis';

export const EXPORT_VIEWS_AND_GROUPS_SUFFIX = 'canopsis_groups_views-';

export const API_ROUTES = {
  auth: '/auth',
  currentUser: '/api/v4/account/me',
  version: '/api/v2/version',
  alarmList: '/api/v4/alarms',
  alarmListExport: '/api/v4/alarm-export',
  context: '/context',
  contextExport: '/api/v4/entity-export',
  createEntity: '/put-entities',
  watcher: '/api/v2/watchers',
  watcherng: '/api/v2/watcherng',
  weatherWatcher: '/api/v4/weather-watchers',
  actions: '/api/v2/actions',
  event: '/event',
  userPreferences: '/rest/userpreferences/userpreferences',
  view: '/api/v2/views',
  viewGroup: '/api/v2/views/groups',
  stats: '/api/v2/stats',
  action: '/rest/default_rights/action',
  user: {
    list: '/rest/default_rights/user',
    create: '/account/user',
    remove: '/account/delete/user',
  },
  role: {
    list: '/rest/default_rights/role',
    create: '/account/role',
    remove: '/account/delete/role',
  },
  eventFilterRules: '/api/v2/eventfilter/rules',
  metaAlarmRule: '/api/v2/metaalarmrule',
  snmpRule: {
    list: '/snmprule',
    create: '/snmprule/put',
  },
  snmpMib: {
    list: '/snmpmib',
    distinct: '/snmpmibdistinct',
    upload: '/uploadmib',
  },
  infos: {
    login: '/api/internal/login/login_info',
    app: '/api/internal/app_info',
    userInterface: '/api/internal/user_interface',
  },
  heartbeat: '/api/v4/heartbeats',
  dynamicInfo: '/api/v2/dynamic-infos',
  filterHints: '/api/v2/associativetable/filterhints',
  sessionStart: '/api/v2/sessionstart',
  sessionTracePath: '/api/v2/session_tracepath',
  keepalive: '/api/v2/keepalive',
  alarmColumnFilters: '/api/v2/associativetable/alarm-column-filters',
  dynamicInfoTemplates: '/api/v2/associativetable/dynamic-info-templates',
  sessionsCount: '/api/v4/sessions-count',
  broadcastMessage: {
    list: '/api/v2/broadcast-message',
    activeList: '/api/v2/broadcast-message/active',
  },
  counter: '/api/v4/alarm-counters',
  playlist: '/api/v2/playlist',
  /**
   * TODO: rename to pbehavior
   */
  planning: {
    timespan: '/api/v4/pbehavior-timespans',
    exceptions: '/api/v4/pbehavior-exceptions',
    types: '/api/v4/pbehavior-types',
    pbehaviors: '/api/v4/pbehaviors',
    pbehaviorComments: '/api/v4/pbehavior-comments',
    pbehaviorById: '/api/v4/entities/pbehaviors',
    reasons: '/api/v4/pbehavior-reasons',
  },
  remediation: {
    instructions: '/api/v4/cat/instructions',
    jobs: '/api/v4/cat/jobs',
    configurations: '/api/v4/cat/job-configs',
    executions: '/api/v4/cat/executions',
    jobExecutions: '/api/v4/cat/job-executions',
  },
  file: '/api/v4/cat/file',
  engineRunInfo: '/api/v4/engine-runinfo',
  cas: {
    login: '/cas/login',
    loggedin: '/cas/loggedin',
  },
  scenarios: '/api/v4/scenarios',
  logout: '/logout',
};

export const COLORS = {
  primary: '#2fab63',
  secondary: '#2b3e4f',
  state: {
    ok: '#00a65a',
    minor: '#fcdc00',
    major: '#ff9900',
    critical: '#f56954',
    pause: '#808080',
  },
  status: {
    off: 'black',
    ongoing: 'grey',
    stealthy: 'gold',
    flapping: 'orange',
    cancelled: 'red',
    unknown: 'black',
  },
  entitiesEvents: {
    ack: '#9c27b0',
    ackRemove: '#9c27b0',
    declareTicket: '#2196f3',
    snooze: '#e91e63',
    done: '#00a65a',
    pbhenter: '#18ffff',
    pbhleave: '#18ffff',
    comment: '#ba68c8',
    metaalarmattach: '#6d4c41',
  },
  statsDefault: '#dddddd',
};

export const FILE_BASE_URL = `${API_BASE_URL}${API_ROUTES.file}`;

export const INSTRUCTION_EXECUTE_FETCHING_INTERVAL = 10000;

export const EXPORT_FETCHING_INTERVAL = 2000;

export const DEFAULT_JOB_EXECUTOR_FETCH_TIMEOUT_SECONDS = 30;

export const DEFAULT_CATEGORIES_LIMIT = 3;

