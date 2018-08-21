export const ENTITIES_TYPES = {
  alarm: 'alarm',
  entity: 'entity',
  watcher: 'watcher',
  watcherEntity: 'watcherEntity',
  weather: 'weather',
  pbehavior: 'pbehavior',
  event: 'event',
  userPreference: 'userPreference',
  view: 'view',
  widgetWrapper: 'widgetWrapper',
  widget: 'widget',
};

export const MODALS = {
  createAckEvent: 'create-ack-event',
  createAssociateTicketEvent: 'create-associate-ticket-event',
  createCancelEvent: 'create-cancel-event',
  createChangeStateEvent: 'create-change-state-event',
  createDeclareTicketEvent: 'create-declare-ticket-event',
  createSnoozeEvent: 'create-snooze-event',
  createPbehavior: 'create-pbehavior',
  createEntity: 'create-entity',
  createWatcher: 'create-watcher',
  watcher: 'watcher',
  pbehaviorList: 'pbehavior-list',
  editLiveReporting: 'edit-live-reporting',
  moreInfos: 'more-infos',
  confirmation: 'confirmation',
  createWidget: 'create-widget',
};

export const EVENT_ENTITY_TYPES = {
  ack: 'ack',
  ackRemove: 'ackremove',
  associateTicket: 'assocticket',
  cancel: 'cancel',
  changeState: 'changestate',
  declareTicket: 'declareticket',
  snooze: 'snooze',
  done: 'done',
};

export const ENTITY_INFOS_TYPE = {
  state: 'state',
  status: 'status',
  action: 'action',
};

export const ENTITIES_STATES = {
  ok: 0,
  minor: 1,
  major: 2,
  critical: 3,
};

export const ENTITIES_STATUSES = {
  off: 0,
  ongoing: 1,
  stealthy: 2,
  flapping: 3,
  cancelled: 4,
};

export const ENTITIES_STATES_STYLES = {
  [ENTITIES_STATES.ok]: {
    color: 'green',
    text: 'ok',
    icon: 'assistant_photo',
  },
  [ENTITIES_STATES.minor]: {
    color: 'gold',
    text: 'minor',
    icon: 'assistant_photo',
  },
  [ENTITIES_STATES.major]: {
    color: 'orange',
    text: 'major',
    icon: 'assistant_photo',
  },
  [ENTITIES_STATES.critical]: {
    color: 'red',
    text: 'critical',
    icon: 'assistant_photo',
  },
};

export const WATCHER_STATES_COLORS = {
  [ENTITIES_STATES.ok]: 'green darken-1',
  [ENTITIES_STATES.minor]: 'yellow darken-1',
  [ENTITIES_STATES.major]: 'orange darken-1',
  [ENTITIES_STATES.critical]: 'red darken-1',
};

export const WATCHER_PBEHAVIOR_COLOR = 'grey lighten-1';

export const ENTITY_STATUS_STYLES = {
  [ENTITIES_STATUSES.off]: {
    color: 'black',
    text: 'off',
    icon: 'keyboard_arrow_up',
  },
  [ENTITIES_STATUSES.ongoing]: {
    color: 'grey',
    text: 'ongoing',
    icon: 'keyboard_arrow_up',
  },
  [ENTITIES_STATUSES.stealthy]: {
    color: 'gold',
    text: 'stealthy',
    icon: 'keyboard_arrow_up',
  },
  [ENTITIES_STATUSES.flapping]: {
    color: 'orange',
    text: 'flapping',
    icon: 'keyboard_arrow_up',
  },
  [ENTITIES_STATUSES.cancelled]: {
    color: 'red',
    text: 'cancelled',
    icon: 'keyboard_arrow_up',
  },
};

export const WIDGET_TYPES = {
  alarmList: 'listalarm',
  context: 'crudcontext',
  weather: 'serviceweather',
  widgetWrapper: 'widgetwrapper',
};

export const EVENT_ENTITY_STYLE = {
  [EVENT_ENTITY_TYPES.ack]: {
    color: 'purple',
    text: 'Acknowledged',
    icon: 'done',
  },
  [EVENT_ENTITY_TYPES.ackRemove]: {
    color: 'purple',
    text: 'Ack removed',
    icon: 'not_interested',
  },
  [EVENT_ENTITY_TYPES.declareTicket]: {
    color: 'blue',
    text: 'Ticket declared',
    icon: 'local_play',
  },
  [EVENT_ENTITY_TYPES.snooze]: {
    color: 'pink',
    text: 'Snoozed',
    icon: 'alarm',
  },
  [EVENT_ENTITY_TYPES.done]: {
    color: 'green',
    text: 'Done',
    icon: 'assignment_turned_in',
  },
};

export const UNKNOWN_VALUE_STYLE = {
  color: 'black',
  text: 'Invalid val',
  icon: 'clear',
};

export const FILTER_OPERATORS = {
  equal: 'equal',
  notEqual: 'not equal',
  in: 'in',
  notIn: 'not in',
  beginsWith: 'begins with',
  doesntBeginWith: 'doesn\'t begin with',
  contains: 'contains',
  doesntContains: 'doesn\'t contain',
  endsWith: 'ends with',
  doesntEndWith: 'doesn\'t end with',
  isEmpty: 'is empty',
  isNotEmpty: 'is not empty',
  isNull: 'is null',
  isNotNull: 'is not null',
};

export const FILTER_DEFAULT_VALUES = {
  condition: '$and',
  rule: {
    field: '',
    operator: '',
    input: '',
  },
  group: {
    condition: '$and',
    groups: {},
    rules: {},
  },
};

export const DATETIME_FORMATS = {
  datetime: 'DD/MM/YYYY H:mm:ss',
  date: 'DD/MM/YYYY',
  time: 'H:mm:ss',
};
