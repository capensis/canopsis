export const API_HOST = process.env.VUE_APP_API_HOST;

export const PAGINATION_LIMIT = parseInt(process.env.VUE_APP_PAGINATION_LIMIT, 10);

export const DEFAULT_LOCALE = 'fr';

export const POPUP_AUTO_CLOSE_DELAY = 3000;

export const VUETIFY_ANIMATION_DELAY = 300;

export const MEDIA_QUERIES_BREAKPOINTS = {
  mobile: 450,
  tablet: 900,
  laptop: 'Infinity',
};

export const API_ROUTES = {
  auth: '/auth',
  currentUser: '/account/me',
  login: '/login',
  alarmList: '/alerts/get-alarms',
  context: '/context',
  pbehavior: '/api/v2/pbehavior',
  event: '/event',
  userPreferences: '/rest/userpreferences/userpreferences',
};
