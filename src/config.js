export const API_HOST = process.env.VUE_APP_API_HOST;

export const POPUP_AUTO_CLOSE_DELAY = parseInt(process.env.VUE_APP_POPUP_AUTO_CLOSE_DELAY, 10);

export const VUETIFY_ANIMATION_DELAY = 300;

export const AUTH_TOKEN_KEY = 'access-token';

export const DEFAULT_LOCALE = 'fr';

export const API_ROUTES = {
  login: '/login',
  alarmList: '/alerts/get-alarms',
  pbehavior: '/api/v2/pbehavior',
};

export const NON_AUTH_API_ROUTES = [
  API_ROUTES.login,
];
