export const ROUTER_MODE = process.env.NODE_ENV === 'development' ? 'history' : 'hash';

export const COOKIE_SESSION_KEY = process.env.VUE_APP_COOKIE_SESSION_KEY;

export const PAGINATION_LIMIT = parseInt(process.env.VUE_APP_PAGINATION_LIMIT, 10);

export const PAGINATION_PER_PAGE_VALUES = [5, 10, 20, 50, 100];

export const DEFAULT_LOCALE = 'fr';

export const POPUP_AUTO_CLOSE_DELAY = 3000;

export const VUETIFY_ANIMATION_DELAY = 300;

export const SIDE_BAR_WIDTH = 250;

export const EXPAND_DEFAULT_MAX_LETTERS = 50;

export const MEDIA_QUERIES_BREAKPOINTS = {
  mobile: 450,
  tablet: 900,
  laptop: 'Infinity',
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
  view: '/rest/object/view',
};
