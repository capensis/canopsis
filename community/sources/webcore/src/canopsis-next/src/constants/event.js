export const EVENT_TYPES = {
  ack: 'ack',
  ackRemove: 'ackremove',
  assocTicket: 'assocticket',
  cancel: 'cancel',
  uncancel: 'uncancel',
  changeState: 'changestate',
  check: 'check',
  comment: 'comment',
  snooze: 'snooze',
};

export const HEALTHCHECK_EVENT_TYPES = {
  ...EVENT_TYPES,

  unsnooze: 'unsnooze',
  pbhenter: 'pbhenter',
  pbhleaveandenter: 'pbhleaveandenter',
  pbhleave: 'pbhleave',
  resolve_cancel: 'resolve_cancel',
  resolve_close: 'resolve_close',
  resolve_deleted: 'resolve_deleted',
  updatestatus: 'updatestatus',
  metaalarm: 'metaalarm',
  metaalarmattachchildren: 'metaalarmattachchildren',
  metaalarmdetachchildren: 'metaalarmdetachchildren',
  recomputeentityservice: 'recomputeentityservice',
  entityupdated: 'entityupdated',
  entitytoggled: 'entitytoggled',
  noevents: 'noevents',
  trigger: 'trigger',
};
