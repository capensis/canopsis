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

export const SCHEMA_EMBEDDED_KEY = '_embedded';

export const MEDIA_QUERIES_BREAKPOINTS = {
  m: 450,
  t: 900,
  l: 'Infinity',
};

export const API_ROUTES = {
  auth: '/auth',
  currentUser: '/account/me',
  version: '/api/v2/version',
  alarmList: '/alerts/get-alarms',
  context: '/context',
  createEntity: '/put-entities',
  watcher: '/api/v2/watchers',
  weatherWatcher: '/api/v2/weather/watchers',
  pbehavior: {
    pbehavior: '/api/v2/pbehavior',
    list: '/pbehavior/read',
    comment: {
      create: '/pbehavior/comment/create',
      update: '/pbehavior/comment/update',
      delete: '/pbehavior/comment/delete',
    },
  },
  pbehaviorById: '/api/v2/pbehavior_byeid',
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
  webhook: '/api/v2/webhook',
  infos: {
    login: '/api/internal/login/login_info',
  },
};
