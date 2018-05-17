export const API_HOST = process.env.VUE_APP_API_HOST;

export const AUTH_TOKEN_KEY = 'access-token';

export const DEFAULT_LOCALE = 'fr';

export const EVENT_TYPES = {
  ack: 'ack',
  ackRemove: 'ackremove',
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
