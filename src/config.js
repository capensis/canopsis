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

export const STATES_CHIPS_AND_FLAGS_STYLE = {
  0: {
    color: 'green',
    text: 'ok',
    icon: 'assistant_photo',
  },
  1: {
    color: 'yellow darken-1',
    text: 'minor',
    icon: 'assistant_photo',
  },
  2: {
    color: 'orange',
    text: 'major',
    icon: 'assistant_photo',
  },
  3: {
    color: 'red',
    text: 'critical',
    icon: 'assistant_photo',
  },
};

export const STATUS_CHIPS_AND_FLAGS_STYLE = {
  0: {
    color: 'black',
    text: 'off',
    icon: 'keyboard_arrow_up',
  },
  1: {
    color: 'grey',
    text: 'ongoing',
    icon: 'keyboard_arrow_up',
  },
  2: {
    color: 'yellow darken-1',
    text: 'stealthy',
    icon: 'keyboard_arrow_up',
  },
  3: {
    color: 'orange',
    text: 'flapping',
    icon: 'keyboard_arrow_up',
  },
  4: {
    color: 'red',
    text: 'cancelled',
    icon: 'keyboard_arrow_up',
  },
};
