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

export const ENTITY_STATES = {
  info: 0,
  minor: 1,
  major: 2,
  critical: 3,
};

export const ENTITY_STATES_STYLES = {
  [ENTITY_STATES.info]: {
    color: 'green',
    text: 'ok',
    icon: 'assistant_photo',
  },
  [ENTITY_STATES.minor]: {
    color: 'yellow darken-1',
    text: 'minor',
    icon: 'assistant_photo',
  },
  [ENTITY_STATES.major]: {
    color: 'orange',
    text: 'major',
    icon: 'assistant_photo',
  },
  [ENTITY_STATES.critical]: {
    color: 'red',
    text: 'critical',
    icon: 'assistant_photo',
  },
};

export const ENTITY_STATUS_STYLES = {
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
