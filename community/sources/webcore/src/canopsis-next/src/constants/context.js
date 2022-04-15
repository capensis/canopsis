import { COLORS } from '@/config';

import { PBEHAVIOR_TYPE_TYPES } from './pbehavior';

export const EVENT_ENTITY_TYPES = {
  ack: 'ack',
  fastAck: 'fastAck',
  ackRemove: 'ackremove',
  pbehaviorAdd: 'pbehaviorAdd',
  pbehaviorList: 'pbehaviorList',
  assocTicket: 'assocticket',
  cancel: 'cancel',
  delete: 'delete',
  changeState: 'changestate',
  declareTicket: 'declareticket',
  snooze: 'snooze',
  done: 'done',
  validate: 'validate',
  invalidate: 'invalidate',
  pause: 'pause',
  play: 'play',
  groupRequest: 'groupRequest',
  group: 'group',
  pbhenter: 'pbhenter',
  pbhleave: 'pbhleave',
  comment: 'comment',
  manualMetaAlarmGroup: 'manual_metaalarm_group',
  manualMetaAlarmUngroup: 'manual_metaalarm_ungroup',
  manualMetaAlarmUpdate: 'manual_metaalarm_update',
  stateinc: 'stateinc',
  statedec: 'statedec',
  statusinc: 'statusinc',
  statusdec: 'statusdec',
  unsooze: 'unsooze',
  metaalarmattach: 'metaalarmattach',
  executeInstruction: 'executeInstruction',
  instructionStart: 'instructionstart',
  instructionPause: 'instructionpause',
  instructionResume: 'instructionresume',
  instructionComplete: 'instructioncomplete',
  instructionAbort: 'instructionabort',
  instructionFail: 'instructionfail',
  instructionJobStart: 'instructionjobstart',
  instructionJobComplete: 'instructionjobcomplete',
  instructionJobAbort: 'instructionjobabort',
  instructionJobFail: 'instructionjobfail',
  autoInstructionStart: 'autoinstructionstart',
  autoInstructionComplete: 'autoinstructioncomplete',
  autoInstructionFail: 'autoinstructionfail',
  autoInstructionAlreadyRunning: 'autoinstructionalreadyrunning',
  junitTestSuiteUpdate: 'junittestsuiteupdate',
  junitTestCaseUpdate: 'junittestcaseupdate',
};

export const ENTITY_INFOS_TYPE = {
  state: 'state',
  status: 'status',
  action: 'action',
};

export const ENTITIES_STATES_KEYS = {
  ok: 'ok',
  minor: 'minor',
  major: 'major',
  critical: 'critical',
};

export const ENTITIES_STATES = {
  ok: 0,
  minor: 1,
  major: 2,
  critical: 3,
};

export const ENTITIES_STATUSES = {
  closed: 0,
  ongoing: 1,
  stealthy: 2,
  flapping: 3,
  cancelled: 4,
  noEvents: 5,
};

export const ENTITIES_STATES_STYLES = {
  [ENTITIES_STATES.ok]: {
    color: COLORS.state.ok,
    text: 'ok',
    icon: 'assistant_photo',
  },
  [ENTITIES_STATES.minor]: {
    color: COLORS.state.minor,
    text: 'minor',
    icon: 'assistant_photo',
  },
  [ENTITIES_STATES.major]: {
    color: COLORS.state.major,
    text: 'major',
    icon: 'assistant_photo',
  },
  [ENTITIES_STATES.critical]: {
    color: COLORS.state.critical,
    text: 'critical',
    icon: 'assistant_photo',
  },
};

export const SERVICE_STATES = {
  ok: 'ok',
  minor: 'minor',
  major: 'major',
  critical: 'critical',
  pause: 'pause',
};

export const SERVICE_STATES_COLORS = {
  [SERVICE_STATES.ok]: ENTITIES_STATES_STYLES[ENTITIES_STATES.ok].color,
  [SERVICE_STATES.minor]: ENTITIES_STATES_STYLES[ENTITIES_STATES.minor].color,
  [SERVICE_STATES.major]: ENTITIES_STATES_STYLES[ENTITIES_STATES.major].color,
  [SERVICE_STATES.critical]: ENTITIES_STATES_STYLES[ENTITIES_STATES.critical].color,
};

export const COUNTER_STATES_ICONS = {
  [ENTITIES_STATES_KEYS.ok]: 'wb_sunny',
  [ENTITIES_STATES_KEYS.minor]: 'person',
  [ENTITIES_STATES_KEYS.major]: 'person',
  [ENTITIES_STATES_KEYS.critical]: 'wb_cloudy',
};

export const WEATHER_ICONS = {
  [SERVICE_STATES.ok]: 'wb_sunny',
  [SERVICE_STATES.minor]: 'person',
  [SERVICE_STATES.major]: 'person',
  [SERVICE_STATES.critical]: 'wb_cloudy',
  [PBEHAVIOR_TYPE_TYPES.maintenance]: 'build',
  [PBEHAVIOR_TYPE_TYPES.inactive]: 'brightness_3',
  [PBEHAVIOR_TYPE_TYPES.pause]: 'pause',
};

export const ENTITY_STATUS_STYLES = {
  [ENTITIES_STATUSES.closed]: {
    color: COLORS.status.closed,
    text: 'closed',
    icon: 'check_circle_outline',
  },
  [ENTITIES_STATUSES.ongoing]: {
    color: COLORS.status.ongoing,
    text: 'ongoing',
    icon: 'warning',
  },
  [ENTITIES_STATUSES.stealthy]: {
    color: COLORS.status.stealthy,
    text: 'stealthy',
    icon: 'swap_vert',
  },
  [ENTITIES_STATUSES.flapping]: {
    color: COLORS.status.flapping,
    text: 'flapping',
    icon: 'swap_vert',
  },
  [ENTITIES_STATUSES.cancelled]: {
    color: COLORS.status.cancelled,
    text: 'cancelled',
    icon: 'highlight_off',
  },
  [ENTITIES_STATUSES.noEvents]: {
    color: COLORS.status.noEvents,
    text: 'no events',
    icon: 'sync_problem',
  },
};

export const EVENT_ENTITY_STYLE = {
  [EVENT_ENTITY_TYPES.ack]: {
    color: COLORS.entitiesEvents.ack,
    icon: 'playlist_add_check',
  },
  [EVENT_ENTITY_TYPES.fastAck]: {
    icon: 'check',
  },
  [EVENT_ENTITY_TYPES.pbehaviorAdd]: {
    icon: 'pause',
  },
  [EVENT_ENTITY_TYPES.ackRemove]: {
    color: COLORS.entitiesEvents.ackRemove,
    icon: 'not_interested',
  },
  [EVENT_ENTITY_TYPES.declareTicket]: {
    color: COLORS.entitiesEvents.declareTicket,
    icon: 'report_problem',
  },
  [EVENT_ENTITY_TYPES.assocTicket]: {
    icon: 'local_play',
  },
  [EVENT_ENTITY_TYPES.delete]: {
    icon: 'delete',
  },
  [EVENT_ENTITY_TYPES.changeState]: {
    icon: 'thumbs_up_down',
  },
  [EVENT_ENTITY_TYPES.snooze]: {
    color: COLORS.entitiesEvents.snooze,
    icon: 'alarm',
  },
  [EVENT_ENTITY_TYPES.done]: {
    color: COLORS.entitiesEvents.done,
    icon: 'assignment_turned_in',
  },
  [EVENT_ENTITY_TYPES.validate]: {
    icon: 'thumb_up',
  },
  [EVENT_ENTITY_TYPES.invalidate]: {
    icon: 'thumb_down',
  },
  [EVENT_ENTITY_TYPES.pause]: {
    icon: 'pause',
  },
  [EVENT_ENTITY_TYPES.play]: {
    icon: 'play_arrow',
  },
  [EVENT_ENTITY_TYPES.groupRequest]: {
    icon: 'note_add',
  },
  [EVENT_ENTITY_TYPES.pbhenter]: {
    color: COLORS.entitiesEvents.pbhenter,
    icon: 'pause',
  },
  [EVENT_ENTITY_TYPES.pbhleave]: {
    color: COLORS.entitiesEvents.pbhleave,
    icon: 'play_arrow',
  },
  groupConsequences: {
    icon: 'center_focus_strong',
  },
  groupCauses: {
    icon: 'center_focus_weak',
  },
  [EVENT_ENTITY_TYPES.comment]: {
    color: COLORS.entitiesEvents.comment,
    icon: 'comment',
  },
  [EVENT_ENTITY_TYPES.manualMetaAlarmGroup]: {
    icon: 'center_focus_strong',
  },
  [EVENT_ENTITY_TYPES.manualMetaAlarmUngroup]: {
    icon: 'link_off',
  },
  [EVENT_ENTITY_TYPES.metaalarmattach]: {
    color: COLORS.entitiesEvents.metaalarmattach,
    icon: 'center_focus_weak',
  },
  [EVENT_ENTITY_TYPES.executeInstruction]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionStart]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionPause]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionResume]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionComplete]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionAbort]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionFail]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionJobStart]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionJobComplete]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionJobAbort]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionJobFail]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.autoInstructionStart]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.autoInstructionComplete]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.autoInstructionFail]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.autoInstructionAlreadyRunning]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.junitTestSuiteUpdate]: {
    icon: 'keyboard_arrow_up',
  },
  [EVENT_ENTITY_TYPES.junitTestCaseUpdate]: {
    icon: 'keyboard_arrow_up',
  },
  [EVENT_ENTITY_TYPES.executeInstruction]: {
    icon: 'assignment',
  },
};

export const UNKNOWN_VALUE_STYLE = {
  color: COLORS.status.unknown,
  text: 'Invalid val',
  icon: 'clear',
};

export const SERVICE_WEATHER_WIDGET_MODAL_TYPES = {
  moreInfo: 'more-info',
  alarmList: 'alarm-list',
  both: 'both',
};

export const WEATHER_EVENT_DEFAULT_ENTITY = 'engine';

export const WEATHER_ACK_EVENT_OUTPUT = {
  ack: 'MDS_ACKNOWLEDGE',
  validateOk: 'MDS_VALIDATEOK',
  validateCancel: 'MDS_VALIDATECANCEL',
};

export const DEFAULT_SERVICE_DEPENDENCIES_COLUMNS = [
  {
    labelKey: 'common.name',
    value: 'entity.name',
  },
  {
    labelKey: 'common.type',
    value: 'entity.type',
  },
];

export const AVAILABLE_COUNTERS = {
  total: 'total',
  total_active: 'total_active',
  snooze: 'snooze',
  ack: 'ack',
  ticket: 'ticket',
  pbehavior_active: 'pbehavior_active',
};

export const DEFAULT_COUNTER_BLOCK_TEMPLATE = `<h2 style="text-align: justify;">
  <strong>{{ counter.filter.title }}</strong></h2><br>
  <center><strong><span style="font-size: 18px;">{{ counter.total_active }} alarmes actives</span></strong></center>
  <br>Seuil mineur à {{ levels.values.minor }}, seuil critique à {{ levels.values.critical }}
  <p style="text-align: justify;">{{ counter.ack }} acquittées, {{ counter.ticket}} avec ticket</p>`;

export const COUNTERS_LIMIT = 3;

export const BASIC_ENTITY_TYPES = {
  connector: 'connector',
  component: 'component',
  resource: 'resource',
};

export const ENTITY_TYPES = {
  service: 'service',

  ...BASIC_ENTITY_TYPES,
};

export const COLOR_INDICATOR_TYPES = {
  state: 'state',
  impactState: 'impact_state',
};

export const STATE_SETTING_METHODS = {
  worst: 'worst',
  worstOfShare: 'worst_of_share',
};

export const STATE_SETTING_THRESHOLD_TYPES = {
  number: 0,
  percent: 1,
};

export const CONTEXT_ACTIONS_TYPES = {
  createEntity: 'createEntity',
  editEntity: 'editEntity',
  duplicateEntity: 'duplicateEntity',
  deleteEntity: 'deleteEntity',
  pbehaviorAdd: 'pbehaviorAdd',
  pbehaviorList: 'pbehaviorList',
  pbehaviorDelete: 'pbehaviorDelete',
  variablesHelp: 'variablesHelp',

  listFilters: 'listFilters',
  editFilter: 'editFilter',
  addFilter: 'addFilter',
};

export const WEATHER_ACTIONS_TYPES = {
  entityAck: 'entityAck',
  entityAssocTicket: 'entityAssocTicket',
  entityValidate: 'entityValidate',
  entityInvalidate: 'entityInvalidate',
  entityPause: 'entityPause',
  entityPlay: 'entityPlay',
  entityCancel: 'entityCancel',
  entityComment: 'entityComment',

  entityLinks: 'entityLinks',

  moreInfos: 'moreInfos',
  alarmsList: 'alarmsList',
  pbehaviorList: 'pbehaviorList',
  variablesHelp: 'variablesHelp',
  executeInstruction: 'executeInstruction',
};

export const COUNTER_ACTIONS_TYPES = {
  alarmsList: 'alarmsList',
  variablesHelp: 'variablesHelp',
};

export const CONTEXT_COLUMN_INFOS_PREFIX = 'infos.';

export const CONTEXT_COLUMNS_WITH_SORTABLE = [ // TODO: We should receive it from backend side in the future
  '_id',
  'name',
  'type',
  'category',
  'impact_level',
  'category.name',
  'idle_since',
];
