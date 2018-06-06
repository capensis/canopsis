export const API_HOST = process.env.VUE_APP_API_HOST;

export const POPUP_AUTO_CLOSE_DELAY = parseInt(process.env.VUE_APP_POPUP_AUTO_CLOSE_DELAY, 10);

export const VUETIFY_ANIMATION_DELAY = 300;
export const PAGINATION_LIMIT = parseInt(process.env.VUE_APP_PAGINATION_LIMIT, 10);
export const MOBILE_BREAKPOINT = parseInt(process.env.VUE_APP_MOBILE_BREAKPOINT, 10);
export const TABLET_BREAKPOINT = parseInt(process.env.VUE_APP_TABLET_BREAKPOINT, 10);
export const LAPTOP_BREAKPOINT = process.env.VUE_APP_LAPTOP_BREAKPOINT;
export const ALARM_LIST_LOADER_WIDTH = parseInt(process.env.VUE_APP_ALARM_LIST_LOADER_WIDTH, 10);
export const ALARM_LIST_LOADER_HEIGHT = parseInt(process.env.VUE_APP_ALARM_LIST_LOADER_HEIGHT, 10);

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

export const API_ROUTES = {
  auth: '/auth',
  currentUser: '/account/me',
  login: '/login',
  alarmList: '/alerts/get-alarms',
  pbehavior: '/api/v2/pbehavior',
  event: '/event',
  userPreferences: '/rest/userpreferences/userpreferences',
};

export const STATES = {
  info: 0,
  minor: 1,
  major: 2,
  critical: 3,
};

export const STATUSES = {
  off: 0,
  ongoing: 1,
  stealthy: 2,
  flapping: 3,
  cancelled: 4,
};

export const STATES_CHIPS_AND_FLAGS_STYLE = {
  [STATES.info]: {
    color: 'green',
    text: 'ok',
    icon: 'assistant_photo',
  },
  [STATES.minor]: {
    color: 'yellow darken-1',
    text: 'minor',
    icon: 'assistant_photo',
  },
  [STATES.major]: {
    color: 'orange',
    text: 'major',
    icon: 'assistant_photo',
  },
  [STATES.critical]: {
    color: 'red',
    text: 'critical',
    icon: 'assistant_photo',
  },
};

export const STATUS_CHIPS_AND_FLAGS_STYLE = {
  [STATUSES.off]: {
    color: 'black',
    text: 'off',
    icon: 'keyboard_arrow_up',
  },
  [STATUSES.ongoing]: {
    color: 'grey',
    text: 'ongoing',
    icon: 'keyboard_arrow_up',
  },
  [STATUSES.stealthy]: {
    color: 'yellow darken-1',
    text: 'stealthy',
    icon: 'keyboard_arrow_up',
  },
  [STATUSES.flapping]: {
    color: 'orange',
    text: 'flapping',
    icon: 'keyboard_arrow_up',
  },
  [STATUSES.cancelled]: {
    color: 'red',
    text: 'cancelled',
    icon: 'keyboard_arrow_up',
  },
};

export const ENTITIES_TYPES = {
  alarm: 'alarm',
  pbehavior: 'pbehavior',
  userPreference: 'userPreference',
};

export const MODALS = {
  createAckEvent: 'create-ack-event',
  createAssociateTicketEvent: 'create-associate-ticket-event',
  createCancelEvent: 'create-cancel-event',
  createChangeStateEvent: 'create-change-state-event',
  createDeclareTicketEvent: 'create-declare-ticket-event',
  createSnoozeEvent: 'create-snooze-event',
  createPbehavior: 'create-pbehavior',
  pbehaviorList: 'pbehavior-list',
};
