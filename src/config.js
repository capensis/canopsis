export const API_HOST = process.env.VUE_APP_API_HOST;

export const POPUP_AUTO_CLOSE_DELAY = parseInt(process.env.VUE_APP_POPUP_AUTO_CLOSE_DELAY, 10);

export const VUETIFY_ANIMATION_DELAY = 300;
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
};
