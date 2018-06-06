export const ENTITIES_TYPES = {
  alarm: 'alarm',
  pbehavior: 'pbehavior',
  event: 'event',
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

export const EVENT_ENTITY_TYPES = {
  ack: 'ack',
  ackRemove: 'ackremove',
  associateTicket: 'assocticket',
  cancel: 'cancel',
  changeState: 'changestate',
  declareTicket: 'declareticket',
  snooze: 'snooze',
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
    color: 'yellow darken-1',
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
    color: 'yellow darken-1',
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
