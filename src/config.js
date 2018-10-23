export const ROUTER_MODE = process.env.NODE_ENV === 'development' ? 'history' : 'hash';

export const COOKIE_SESSION_KEY = process.env.VUE_APP_COOKIE_SESSION_KEY;

export const PAGINATION_LIMIT = parseInt(process.env.VUE_APP_PAGINATION_LIMIT, 10);

export const PAGINATION_PER_PAGE_VALUES = [5, 10, 20, 50, 100];

export const PAGINATION_TOTAL_VISIBLE = 7;

export const DEFAULT_LOCALE = 'en';

export const POPUP_AUTO_CLOSE_DELAY = 3000;

export const VUETIFY_ANIMATION_DELAY = 300;

export const SIDE_BAR_WIDTH = 250;

export const EXPAND_DEFAULT_MAX_LETTERS = 50;

export const MEDIA_QUERIES_BREAKPOINTS = {
  m: 450,
  t: 900,
  l: 'Infinity',
};

export const API_ROUTES = {
  auth: '/auth',
  currentUser: '/account/me',
  alarmList: '/alerts/get-alarms',
  context: '/context',
  createEntity: '/put-entities',
  watcher: '/api/v2/watchers',
  weatherWatcher: '/api/v2/weather/watchers',
  pbehavior: '/api/v2/pbehavior',
  pbehaviorById: '/api/v2/pbehavior_byeid',
  event: '/event',
  userPreferences: '/rest/userpreferences/userpreferences',
  view: '/api/v2/views',
  viewGroup: '/api/v2/views/groups',
  stats: '/api/v2/stats',
  role: '/rest/default_rights/role',
};
