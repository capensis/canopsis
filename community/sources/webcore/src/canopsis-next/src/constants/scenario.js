export const ACTION_TYPES = {
  ack: 'ack',
  ackremove: 'ackremove',
  assocticket: 'assocticket',
  cancel: 'cancel',
  snooze: 'snooze',
  pbehavior: 'pbehavior',
  changeState: 'changestate',
  webhook: 'webhook',
};

export const CAT_ACTION_TYPES = [ACTION_TYPES.webhook];

export const SCENARIO_TRIGGERS = {
  create: 'create',
  stateinc: 'stateinc',
  statedec: 'statedec',
  changestate: 'changestate',
  changestatus: 'changestatus',
  ack: 'ack',
  ackremove: 'ackremove',
  cancel: 'cancel',
  uncancel: 'uncancel',
  comment: 'comment',
  done: 'done',
  declareticket: 'declareticket',
  declareticketwebhook: 'declareticketwebhook',
  assocticket: 'assocticket',
  snooze: 'snooze',
  unsnooze: 'unsnooze',
  resolve: 'resolve',
  activate: 'activate',
  pbhenter: 'pbhenter',
  pbhleave: 'pbhleave',
  instructionfail: 'instructionfail',
  autoinstructionfail: 'autoinstructionfail',
  instructionjobcomplete: 'instructionjobcomplete',
  instructionjobfail: 'instructionjobfail',
  instructioncomplete: 'instructioncomplete',
  autoinstructioncomplete: 'autoinstructioncomplete',
};

export const CAT_SCENARIO_TRIGGERS = [
  SCENARIO_TRIGGERS.declareticket,
  SCENARIO_TRIGGERS.declareticketwebhook,
  SCENARIO_TRIGGERS.instructionfail,
  SCENARIO_TRIGGERS.autoinstructionfail,
  SCENARIO_TRIGGERS.instructionjobcomplete,
  SCENARIO_TRIGGERS.instructionjobfail,
  SCENARIO_TRIGGERS.instructioncomplete,
  SCENARIO_TRIGGERS.autoinstructioncomplete,
];
