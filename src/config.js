export const API_HOST = process.env.VUE_APP_API_HOST;

export const PAGINATION_LIMIT = parseInt(process.env.VUE_APP_PAGINATION_LIMIT, 10);
export const MOBILE_BREAKPOINT = parseInt(process.env.VUE_APP_MOBILE_BREAKPOINT, 10);
export const TABLET_BREAKPOINT = parseInt(process.env.VUE_APP_TABLET_BREAKPOINT, 10);
export const LAPTOP_BREAKPOINT = process.env.VUE_APP_LAPTOP_BREAKPOINT;
export const ALARM_LIST_LOADER_WIDTH = parseInt(process.env.VUE_APP_ALARM_LIST_LOADER_WIDTH, 10);
export const ALARM_LIST_LOADER_HEIGHT = parseInt(process.env.VUE_APP_ALARM_LIST_LOADER_HEIGHT, 10);

export const AUTH_TOKEN_KEY = 'access-token';

export const DEFAULT_LOCALE = 'fr';

export const EVENT_TYPES = {
  ack: 'ack',
  ackRemove: 'ackremove',
  associateTicket: 'assocticket',
  cancel: 'cancel',
  changeState: 'changestate',
  declareTicket: 'declareticket',
  snooze: 'snooze',
};

export const STATES = {
  info: 0,
  minor: 1,
  major: 2,
  critical: 3,
};

export const API_ROUTES = {
  login: '/login',
  alarmList: '/alerts/get-alarms',
  pbehavior: '/api/v2/pbehavior',
  event: '/event',
  userPreferences: '/rest/userpreferences/userpreferences',
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
