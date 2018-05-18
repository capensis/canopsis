export const API_HOST = process.env.VUE_APP_API_HOST;

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
