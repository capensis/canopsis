export const API_HOST = process.env.VUE_APP_API_HOST;

export const PAGINATION_LIMIT = parseInt(process.env.VUE_APP_PAGINATION_LIMIT, 10);

export const AUTH_TOKEN_KEY = 'access-token';

export const DEFAULT_LOCALE = 'fr';

export const API_ROUTES = {
  login: '/login',
  alarmList: '/alerts/get-alarms',
};

export const NON_AUTH_API_ROUTES = [
  API_ROUTES.login,
];
