import { COLORS } from '@/config';

export const ALARM_FIELDS = {
  id: '_id',
  assignedInstructions: 'assigned_instructions',
  displayName: 'v.display_name',
  output: 'v.output',
  longOutput: 'v.long_output',
  initialOutput: 'v.initial_output',
  initialLongOutput: 'v.initial_long_output',
  connector: 'v.connector',
  connectorName: 'v.connector_name',
  component: 'v.component',
  resource: 'v.resource',
  lastComment: 'v.last_comment.m',
  lastCommentInitiator: 'v.last_comment.initiator',
  ackBy: 'v.ack.a',
  ackMessage: 'v.ack.m',
  ackInitiator: 'v.ack.initiator',
  stateMessage: 'v.state.m',
  statusMessage: 'v.status.m',
  state: 'v.state.val',
  status: 'v.status.val',
  totalStateChanges: 'v.total_state_changes',
  timestamp: 't',
  creationDate: 'v.creation_date',
  lastUpdateDate: 'v.last_update_date',
  lastEventDate: 'v.last_event_date',
  ackAt: 'v.ack.t',
  stateAt: 'v.state.t',
  statusAt: 'v.status.t',
  resolved: 'v.resolved',
  activationDate: 'v.activation_date',
  duration: 'v.duration',
  currentStateDuration: 'v.current_state_duration',
  snoozeDuration: 'v.snooze_duration',
  pbhInactiveDuration: 'v.pbh_inactive_duration',
  activeDuration: 'v.active_duration',
  eventsCount: 'v.events_count',
  tags: 'tags',
  extraDetails: 'extra_details',
  impactState: 'impact_state',
  infos: 'v.infos',
  links: 'links',
  ticketAuthor: 'v.ticket.a',
  ticketMessage: 'v.ticket.m',
  ticketInitiator: 'v.ticket.initiator',
  ticketValue: 'v.ticket.ticket',
  ticketCreatedAt: 'v.ticket.t',
  changeState: 'v.change_state',
  meta: 'v.meta',
  entityId: 'entity._id',
  entityName: 'entity.name',
  entityCategoryName: 'entity.category.name',
  entityType: 'entity.type',
  entityComponent: 'entity.component',
  entityConnector: 'entity.connector',
  entityImpactLevel: 'entity.impact_level',
  entityKoEvents: 'entity.ko_events',
  entityOkEvents: 'entity.ok_events',
  entityInfos: 'entity.infos',
  entityComponentInfos: 'entity.component_infos',
  entityLastPbehaviorDate: 'entity.last_pbehavior_date',

  /**
   * OBJECTS
   */
  ack: 'v.ack',
  ticket: 'v.ticket',
  ticketData: 'v.ticket.ticket_data',
  canceled: 'v.canceled',
  canceledInitiator: 'v.canceled.initiator',
  snooze: 'v.snooze',
  pbehaviorInfo: 'v.pbehavior_info',

  /**
   * VIRTUAL
   */
  activated: 'activated',
};

export const ALARM_INFOS_FIELDS = [
  ALARM_FIELDS.infos,
  ALARM_FIELDS.entityInfos,
  ALARM_FIELDS.entityComponentInfos,
];

export const ALARM_OUTPUT_FIELDS = [
  ALARM_FIELDS.output,
  ALARM_FIELDS.longOutput,
  ALARM_FIELDS.initialOutput,
  ALARM_FIELDS.initialLongOutput,
];

export const ALARM_LEVELS = {
  minor: 20,
  major: 30,
  critical: 40,
};

export const ALARM_LEVELS_COLORS = {
  ok: COLORS.state.ok,
  minor: COLORS.state.minor,
  major: COLORS.state.major,
  critical: COLORS.state.critical,
};

export const ALARM_STATES = {
  ok: 0,
  minor: 1,
  major: 2,
  critical: 3,
};

export const ALARM_STATES_KEYS = {
  ok: 'ok',
  minor: 'minor',
  major: 'major',
  critical: 'critical',
};

export const ALARM_STATES_TEXTS = {
  [ALARM_STATES.ok]: 'ok',
  [ALARM_STATES.minor]: 'minor',
  [ALARM_STATES.major]: 'major',
  [ALARM_STATES.critical]: 'critical',
};

export const ALARM_STATES_ICONS = {
  [ALARM_STATES.ok]: 'assistant_photo',
  [ALARM_STATES.minor]: 'assistant_photo',
  [ALARM_STATES.major]: 'assistant_photo',
  [ALARM_STATES.critical]: 'assistant_photo',
};

export const ALARM_STATUSES = {
  closed: 0,
  ongoing: 1,
  stealthy: 2,
  flapping: 3,
  cancelled: 4,
  noEvents: 5,
};

export const ALARM_STATUSES_ICONS = {
  [ALARM_STATUSES.closed]: 'check_circle_outline',
  [ALARM_STATUSES.ongoing]: 'warning',
  [ALARM_STATUSES.stealthy]: 'swap_vert',
  [ALARM_STATUSES.flapping]: 'swap_vert',
  [ALARM_STATUSES.cancelled]: 'highlight_off',
  [ALARM_STATUSES.noEvents]: 'sync_problem',
};

export const ALARM_STATUSES_TEXTS = {
  [ALARM_STATUSES.closed]: 'closed',
  [ALARM_STATUSES.ongoing]: 'ongoing',
  [ALARM_STATUSES.stealthy]: 'stealthy',
  [ALARM_STATUSES.flapping]: 'flapping',
  [ALARM_STATUSES.cancelled]: 'cancelled',
  [ALARM_STATUSES.noEvents]: 'no events',
};

export const ALARM_LIST_ACTIONS_TYPES = {
  pbehaviorAdd: 'pbehaviorAdd',
  moreInfos: 'moreInfos',
  snooze: 'snooze',
  declareTicket: 'declareTicket',
  associateTicket: 'associateTicket',
  changeState: 'changeState',
  variablesHelp: 'variablesHelp',
  history: 'history',
  exportPdf: 'exportPdf',
  createManualMetaAlarm: 'createManualMetaAlarm',
  removeAlarmsFromManualMetaAlarm: 'removeAlarmsFromManualMetaAlarm',
  removeAlarmsFromAutoMetaAlarm: 'removeAlarmsFromAutoMetaAlarm',
  comment: 'comment',

  ack: 'ack',
  fastAck: 'fastAck',
  ackRemove: 'ackRemove',

  cancel: 'cancel',
  fastCancel: 'fastCancel',
  unCancel: 'unCancel',

  links: 'links',

  correlation: 'correlation',

  listFilters: 'listFilters',
  editFilter: 'editFilter',
  addFilter: 'addFilter',
  userFilter: 'userFilter',

  listRemediationInstructionsFilters: 'listRemediationInstructionsFilters',
  editRemediationInstructionsFilter: 'editRemediationInstructionsFilter',
  addRemediationInstructionsFilter: 'addRemediationInstructionsFilter',
  userRemediationInstructionsFilter: 'userRemediationInstructionsFilter',

  executeInstruction: 'executeInstruction',

  addBookmark: 'addBookmark',
  removeBookmark: 'removeBookmark',
};

export const ALARM_LIST_ACTIONS_TYPES_ICONS = {
  [ALARM_LIST_ACTIONS_TYPES.pbehaviorAdd]: 'pause',
  [ALARM_LIST_ACTIONS_TYPES.snooze]: 'alarm',
  [ALARM_LIST_ACTIONS_TYPES.declareTicket]: 'note_add',
  [ALARM_LIST_ACTIONS_TYPES.associateTicket]: '$vuetify.icons.sticky_note_2',
  [ALARM_LIST_ACTIONS_TYPES.changeState]: 'thumbs_up_down',
  [ALARM_LIST_ACTIONS_TYPES.variablesHelp]: 'help',
  [ALARM_LIST_ACTIONS_TYPES.history]: 'history',
  [ALARM_LIST_ACTIONS_TYPES.exportPdf]: 'assignment_returned',
  [ALARM_LIST_ACTIONS_TYPES.createManualMetaAlarm]: 'center_focus_strong',
  [ALARM_LIST_ACTIONS_TYPES.removeAlarmsFromManualMetaAlarm]: 'link_off',
  [ALARM_LIST_ACTIONS_TYPES.removeAlarmsFromAutoMetaAlarm]: 'link_off',
  [ALARM_LIST_ACTIONS_TYPES.comment]: 'comment',

  [ALARM_LIST_ACTIONS_TYPES.ack]: 'check',
  [ALARM_LIST_ACTIONS_TYPES.fastAck]: 'done_all',
  [ALARM_LIST_ACTIONS_TYPES.ackRemove]: 'remove_done',

  [ALARM_LIST_ACTIONS_TYPES.cancel]: 'block',
  [ALARM_LIST_ACTIONS_TYPES.fastCancel]: 'delete',
  [ALARM_LIST_ACTIONS_TYPES.unCancel]: 'delete_forever',

  [ALARM_LIST_ACTIONS_TYPES.executeInstruction]: 'assignment',

  [ALARM_LIST_ACTIONS_TYPES.addBookmark]: '$vuetify.icons.bookmark_add',
  [ALARM_LIST_ACTIONS_TYPES.removeBookmark]: '$vuetify.icons.bookmark_remove',
};

export const ALARM_LIST_STEPS = {
  stateinc: 'stateinc',
  statedec: 'statedec',
  changeState: 'changestate',
  stateCounter: 'statecounter',

  statusinc: 'statusinc',
  statusdec: 'statusdec',
  resolve: 'resolve',
  activate: 'activate',

  ack: 'ack',
  ackRemove: 'ackremove',

  pbhenter: 'pbhenter',
  pbhleave: 'pbhleave',

  assocTicket: 'assocticket',

  webhookStart: 'webhookstart',
  webhookProgress: 'webhookprogress', // TODO: NEW
  webhookComplete: 'webhookcomplete',
  webhookFail: 'webhookfail',

  declareTicket: 'declareticket',
  declareTicketFail: 'declareticketfail',
  ticketDeclarationRuleProgress: 'ticketdeclatetionruleprogress', // TODO: NEW
  ticketDeclarationRuleCompleted: 'ticketdeclatetionrulecomplete', // TODO: NEW
  ticketDeclarationRuleFailed: 'ticketdeclatetionrulefailed', // TODO: NEW

  snooze: 'snooze',
  unsnooze: 'unsnooze',

  comment: 'comment',

  metaalarmattach: 'metaalarmattach',
  metaalarmdetach: 'metaalarmdetach',

  instructionStart: 'instructionstart',
  instructionPause: 'instructionpause',
  instructionResume: 'instructionresume',
  instructionComplete: 'instructioncomplete',
  instructionAbort: 'instructionabort',
  instructionFail: 'instructionfail',

  autoInstructionStart: 'autoinstructionstart',
  autoInstructionComplete: 'autoinstructioncomplete',
  autoInstructionFail: 'autoinstructionfail',

  junitTestSuiteUpdate: 'junittestsuiteupdate',
  junitTestCaseUpdate: 'junittestcaseupdate',
};

export const ALARM_STEPS_ICONS = {
  [ALARM_LIST_STEPS.ack]: 'check',
  [ALARM_LIST_STEPS.ackRemove]: 'remove_done',
  [ALARM_LIST_STEPS.assocTicket]: '$vuetify.icons.sticky_note_2',
  [ALARM_LIST_STEPS.changeState]: 'warining',
  [ALARM_LIST_STEPS.declareTicket]: 'note_add',
  [ALARM_LIST_STEPS.declareTicketFail]: 'note_add',
  [ALARM_LIST_STEPS.webhookStart]: '$vuetify.icons.webhook',
  [ALARM_LIST_STEPS.webhookProgress]: '$vuetify.icons.webhook',
  [ALARM_LIST_STEPS.webhookComplete]: '$vuetify.icons.webhook',
  [ALARM_LIST_STEPS.webhookFail]: '$vuetify.icons.webhook',
  [ALARM_LIST_STEPS.ticketDeclarationRuleProgress]: 'note_add',
  [ALARM_LIST_STEPS.ticketDeclarationRuleCompleted]: 'note_add',
  [ALARM_LIST_STEPS.ticketDeclarationRuleFailed]: 'note_add',
  [ALARM_LIST_STEPS.snooze]: 'alarm',
  [ALARM_LIST_STEPS.unsnooze]: 'alarm_off',
  [ALARM_LIST_STEPS.pbhenter]: 'pause',
  [ALARM_LIST_STEPS.pbhleave]: 'pause',
  [ALARM_LIST_STEPS.comment]: 'comment',
  [ALARM_LIST_STEPS.metaalarmattach]: 'center_focus_weak',
  [ALARM_LIST_STEPS.instructionStart]: '$vuetify.icons.manual_instruction',
  [ALARM_LIST_STEPS.instructionPause]: '$vuetify.icons.manual_instruction',
  [ALARM_LIST_STEPS.instructionResume]: '$vuetify.icons.manual_instruction',
  [ALARM_LIST_STEPS.instructionComplete]: '$vuetify.icons.manual_instruction',
  [ALARM_LIST_STEPS.instructionAbort]: '$vuetify.icons.manual_instruction',
  [ALARM_LIST_STEPS.instructionFail]: '$vuetify.icons.manual_instruction',
  [ALARM_LIST_STEPS.autoInstructionStart]: 'assignment',
  [ALARM_LIST_STEPS.autoInstructionComplete]: 'assignment',
  [ALARM_LIST_STEPS.autoInstructionFail]: 'assignment',
  [ALARM_LIST_STEPS.activate]: 'notifications',
  [ALARM_LIST_STEPS.resolve]: 'check_circle_outline',
};

export const ALARM_STEPS_COLORS = {
  [ALARM_LIST_STEPS.ack]: COLORS.entitiesEvents.ack,
  [ALARM_LIST_STEPS.ackRemove]: COLORS.entitiesEvents.ackRemove,
  [ALARM_LIST_STEPS.snooze]: COLORS.entitiesEvents.snooze,
  [ALARM_LIST_STEPS.pbhenter]: COLORS.entitiesEvents.pbhenter,
  [ALARM_LIST_STEPS.pbhleave]: COLORS.entitiesEvents.pbhleave,
  [ALARM_LIST_STEPS.comment]: COLORS.entitiesEvents.comment,
  [ALARM_LIST_STEPS.metaalarmattach]: COLORS.entitiesEvents.metaalarmattach,
};

export const ALARM_STEPS_WITH_AUTHOR_IN_TITLE = [
  ALARM_LIST_STEPS.statedec,
  ALARM_LIST_STEPS.stateinc,
  ALARM_LIST_STEPS.changeState,
  ALARM_LIST_STEPS.statusdec,
  ALARM_LIST_STEPS.statusinc,
  ALARM_LIST_STEPS.resolve,
  ALARM_LIST_STEPS.snooze,
  ALARM_LIST_STEPS.unsnooze,
  ALARM_LIST_STEPS.ack,
  ALARM_LIST_STEPS.ackRemove,
  ALARM_LIST_STEPS.pbhenter,
  ALARM_LIST_STEPS.pbhleave,
  ALARM_LIST_STEPS.assocTicket,
  ALARM_LIST_STEPS.webhookStart,
  ALARM_LIST_STEPS.webhookComplete,
  ALARM_LIST_STEPS.webhookFail,
  ALARM_LIST_STEPS.ticketDeclarationRuleCompleted,
  ALARM_LIST_STEPS.ticketDeclarationRuleFailed,
  ALARM_LIST_STEPS.instructionStart,
  ALARM_LIST_STEPS.instructionPause,
  ALARM_LIST_STEPS.instructionResume,
  ALARM_LIST_STEPS.instructionComplete,
  ALARM_LIST_STEPS.instructionAbort,
  ALARM_LIST_STEPS.instructionFail,
  ALARM_LIST_STEPS.metaalarmattach,
  ALARM_LIST_STEPS.metaalarmdetach,
  ALARM_LIST_STEPS.comment,
];

export const ALARM_STEPS_WITH_LAUNCHED_IN_TITLE = [
  ALARM_LIST_STEPS.webhookComplete,
  ALARM_LIST_STEPS.webhookFail,
  ALARM_LIST_STEPS.ticketDeclarationRuleCompleted,
  ALARM_LIST_STEPS.ticketDeclarationRuleFailed,
];

export const ALARM_STEPS_WITH_CONDITION_FOR_AUTHOR_IN_TITLE = [
  ALARM_LIST_STEPS.webhookStart,
  ALARM_LIST_STEPS.webhookComplete,
  ALARM_LIST_STEPS.webhookFail,
];

export const ALARM_UNKNOWN_VALUE = {
  color: COLORS.status.unknown,
  text: 'Invalid val',
  icon: 'clear',
};

export const META_ALARMS_RULE_TYPES = {
  relation: 'relation',
  timebased: 'timebased',
  attribute: 'attribute',
  complex: 'complex',
  valuegroup: 'valuegroup',
  corel: 'corel',

  /**
   * Manual group type doesn't use in the form
   * We are using it only inside alarms list widget
   */
  manualgroup: 'manualgroup',
};

export const META_ALARMS_THRESHOLD_TYPES = {
  thresholdRate: 'thresholdRate',
  thresholdCount: 'thresholdCount',
};

export const ALARMS_LIST_TIME_LINE_SYSTEM_AUTHOR = 'canopsis.engine';

export const DEFAULT_ALARMS_WIDGET_COLUMNS = [
  { value: ALARM_FIELDS.connector, isHtml: true },
  { value: ALARM_FIELDS.connectorName, isHtml: true },
  { value: ALARM_FIELDS.component, isHtml: true },
  { value: ALARM_FIELDS.resource, isHtml: true },
  { value: ALARM_FIELDS.output, isHtml: true },
  { value: ALARM_FIELDS.extraDetails },
  { value: ALARM_FIELDS.state },
  { value: ALARM_FIELDS.status },
];

export const DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS = [
  { value: ALARM_FIELDS.connector, isHtml: true },
  { value: ALARM_FIELDS.connectorName, isHtml: true },
  { value: ALARM_FIELDS.resource, isHtml: true },
  { value: ALARM_FIELDS.output, isHtml: true },
  { value: ALARM_FIELDS.extraDetails },
  { value: ALARM_FIELDS.state },
  { value: ALARM_FIELDS.status },
];

export const DEFAULT_CONTEXT_WIDGET_RESOLVED_ALARMS_COLUMNS = [
  { value: ALARM_FIELDS.id },
  { value: ALARM_FIELDS.creationDate },
  { value: ALARM_FIELDS.lastUpdateDate },
];

export const DEFAULT_CONTEXT_WIDGET_ACTIVE_ALARM_COLUMNS = [
  { value: ALARM_FIELDS.displayName },
  { value: ALARM_FIELDS.output },
  { value: ALARM_FIELDS.state },
  { value: ALARM_FIELDS.status },
];

export const ALARMS_OPENED_VALUES = {
  opened: true,
  all: null,
  resolved: false,
};

export const ALARMS_RESIZING_CELLS_CONTENTS_BEHAVIORS = {
  wrap: 'wrap',
  truncate: 'truncate',
};

export const ALARM_BASIC_METRIC_PARAMETERS = {
  activeAlarms: 'active_alarms',
  createdAlarms: 'created_alarms',
  nonDisplayedAlarms: 'non_displayed_alarms',
  instructionAlarms: 'instruction_alarms',
  pbehaviorAlarms: 'pbehavior_alarms',
  correlationAlarms: 'correlation_alarms',
  ackAlarms: 'ack_alarms',
  cancelAckAlarms: 'cancel_ack_alarms',
  ackActiveAlarms: 'ack_active_alarms',
  ticketActiveAlarms: 'ticket_active_alarms',
  withoutTicketActiveAlarms: 'without_ticket_active_alarms',
  ratioCorrelation: 'ratio_correlation',
  ratioInstructions: 'ratio_instructions',
  ratioTickets: 'ratio_tickets',
  ratioNonDisplayed: 'ratio_non_displayed',
  ratioRemediatedAlarms: 'ratio_remediated_alarms',
  averageAck: 'average_ack',
  averageResolve: 'average_resolve',
  timeToAck: 'time_to_ack',
  timeToResolve: 'time_to_resolve',
  minAck: 'min_ack',
  maxAck: 'max_ack',
};

export const ALARM_OPTIONAL_METRIC_PARAMETERS = {
  manualInstructionAssignedAlarms: 'manual_instruction_assigned_alarms',
  manualInstructionExecutedAlarms: 'manual_instruction_executed_alarms',
  notAckedAlarms: 'not_acked_alarms',
  notAckedInHourAlarms: 'not_acked_in_hour_alarms',
  notAckedInFourHoursAlarms: 'not_acked_in_four_hours_alarms',
  notAckedInDayAlarms: 'not_acked_in_day_alarms',
  minResolve: 'min_resolve',
  maxResolve: 'max_resolve',
};

export const ALARM_METRIC_PARAMETERS = {
  ...ALARM_BASIC_METRIC_PARAMETERS,
  ...ALARM_OPTIONAL_METRIC_PARAMETERS,
};

export const ALARMS_LIST_HEADER_OPACITY_DELAY = 500;

export const ALARM_INTERVAL_FIELDS = {
  timestamp: 't',
  resolved: 'v.resolved',
  lastUpdateDate: 'v.last_update_date',
  lastEventDate: 'v.last_event_date',
};

const {
  ack,
  ticket,
  canceled,
  snooze,
  pbehaviorInfo,
  activated,

  ...alarmListWidgetColumns
} = ALARM_FIELDS;

export const ALARM_LIST_WIDGET_COLUMNS = alarmListWidgetColumns;

export const ALARM_PATTERN_FIELDS = {
  displayName: ALARM_FIELDS.displayName,
  state: ALARM_FIELDS.state,
  status: ALARM_FIELDS.status,
  component: ALARM_FIELDS.component,
  resource: ALARM_FIELDS.resource,
  connector: ALARM_FIELDS.connector,
  connectorName: ALARM_FIELDS.connectorName,
  creationDate: ALARM_FIELDS.creationDate,
  duration: ALARM_FIELDS.duration,
  infos: ALARM_FIELDS.infos,
  output: ALARM_FIELDS.output,
  lastEventDate: ALARM_FIELDS.lastEventDate,
  lastUpdateDate: ALARM_FIELDS.lastUpdateDate,
  activationDate: ALARM_FIELDS.activationDate,
  ack: ALARM_FIELDS.ack,
  ackAt: ALARM_FIELDS.ackAt,
  ackBy: ALARM_FIELDS.ackBy,
  ackMessage: ALARM_FIELDS.ackMessage,
  ackInitiator: ALARM_FIELDS.ackInitiator,
  resolved: ALARM_FIELDS.resolved,
  ticket: ALARM_FIELDS.ticket,
  canceled: ALARM_FIELDS.canceled,
  canceledInitiator: ALARM_FIELDS.canceledInitiator,
  snooze: ALARM_FIELDS.snooze,
  lastComment: ALARM_FIELDS.lastComment,
  lastCommentInitiator: ALARM_FIELDS.lastCommentInitiator,
  longOutput: ALARM_FIELDS.longOutput,
  initialOutput: ALARM_FIELDS.initialOutput,
  initialLongOutput: ALARM_FIELDS.initialLongOutput,
  totalStateChanges: ALARM_FIELDS.totalStateChanges,
  tags: ALARM_FIELDS.tags,
  activated: ALARM_FIELDS.activated,
  ticketValue: ALARM_FIELDS.ticketValue,
  ticketMessage: ALARM_FIELDS.ticketMessage,
  ticketInitiator: ALARM_FIELDS.ticketInitiator,
  ticketData: ALARM_FIELDS.ticketData,
  changeState: ALARM_FIELDS.changeState,
  meta: ALARM_FIELDS.meta,
};

export const ALARM_EVENT_INITIATORS = {
  user: 'user',
  system: 'system',
  external: 'external',
};

export const ALARM_STEP_FIELDS = {
  timestamp: 't',
  value: 'val',
  message: 'm',
  author: 'a',
};

export const ALARM_TEMPLATE_FIELDS = {
  id: `alarm.${ALARM_FIELDS.id}`,
  ack: `alarm.${ALARM_FIELDS.ack}`,
  state: `alarm.${ALARM_FIELDS.state}`,
  status: `alarm.${ALARM_FIELDS.status}`,
  ticket: `alarm.${ALARM_FIELDS.ticket}`,
  component: `alarm.${ALARM_FIELDS.component}`,
  connector: `alarm.${ALARM_FIELDS.connector}`,
  connectorName: `alarm.${ALARM_FIELDS.connectorName}`,
  resource: `alarm.${ALARM_FIELDS.resource}`,
  creationDate: `alarm.${ALARM_FIELDS.creationDate}`,
  displayName: `alarm.${ALARM_FIELDS.displayName}`,
  output: `alarm.${ALARM_FIELDS.output}`,
  lastUpdateDate: `alarm.${ALARM_FIELDS.lastUpdateDate}`,
  lastEventDate: `alarm.${ALARM_FIELDS.lastEventDate}`,
  pbehaviorInfo: `alarm.${ALARM_FIELDS.pbehaviorInfo}`,
  duration: `alarm.${ALARM_FIELDS.duration}`,
  eventsCount: `alarm.${ALARM_FIELDS.eventsCount}`,
};

export const ALARM_EXPORT_PDF_FIELDS = {
  currentDate: 'current_date',
  displayName: 'display_name',
  state: 'state',
  status: 'status',
  connector: 'connector',
  connectorName: 'connector_name',
  component: 'component',
  resource: 'resource',
  initialOutput: 'initial_output',
  output: 'output',
  eventsCount: 'events_count',
  duration: 'duration',
  creationDate: 'creation_date',
  lastEventDate: 'last_event_date',
  lastUpdateDate: 'last_update_date',
  acknowledgeDate: 'acknowledge_date',
  resolved: 'resolved',
  activationDate: 'activation_date',
  infos: 'infos',
  pbehaviorInfo: 'pbehavior_info',
  ticket: 'ticket',
  comments: 'comments',
  tags: 'tags',
  links: 'links',
};

export const ALARM_EXPORT_PDF_FIELDS_TO_ORIGINAL_FIELDS = {
  [ALARM_EXPORT_PDF_FIELDS.currentDate]: 'current_date',
  [ALARM_EXPORT_PDF_FIELDS.displayName]: 'alarm.v.display_name',
  [ALARM_EXPORT_PDF_FIELDS.state]: 'alarm.v.state',
  [ALARM_EXPORT_PDF_FIELDS.status]: 'alarm.v.status',
  [ALARM_EXPORT_PDF_FIELDS.connector]: 'alarm.v.connector',
  [ALARM_EXPORT_PDF_FIELDS.connectorName]: 'alarm.v.connector_name',
  [ALARM_EXPORT_PDF_FIELDS.component]: 'alarm.v.component',
  [ALARM_EXPORT_PDF_FIELDS.resource]: 'alarm.v.resource',
  [ALARM_EXPORT_PDF_FIELDS.initialOutput]: 'alarm.v.initial_output',
  [ALARM_EXPORT_PDF_FIELDS.output]: 'alarm.v.output',
  [ALARM_EXPORT_PDF_FIELDS.eventsCount]: 'alarm.v.events_count',
  [ALARM_EXPORT_PDF_FIELDS.duration]: 'alarm.v.duration',
  [ALARM_EXPORT_PDF_FIELDS.creationDate]: 'alarm.v.creation_date',
  [ALARM_EXPORT_PDF_FIELDS.lastEventDate]: 'alarm.v.last_event_date',
  [ALARM_EXPORT_PDF_FIELDS.lastUpdateDate]: 'alarm.v.last_update_date',
  [ALARM_EXPORT_PDF_FIELDS.acknowledgeDate]: 'alarm.v.acknowledge_date',
  [ALARM_EXPORT_PDF_FIELDS.resolved]: 'alarm.v.resolved',
  [ALARM_EXPORT_PDF_FIELDS.activationDate]: 'alarm.v.activation_date',
  [ALARM_EXPORT_PDF_FIELDS.infos]: 'alarm.v.infos',
  [ALARM_EXPORT_PDF_FIELDS.pbehaviorInfo]: 'alarm.v.pbehavior_info',
  [ALARM_EXPORT_PDF_FIELDS.ticket]: 'alarm.v.ticket',
  [ALARM_EXPORT_PDF_FIELDS.comments]: 'alarm.comments',
  [ALARM_EXPORT_PDF_FIELDS.tags]: 'alarm.tags',
  [ALARM_EXPORT_PDF_FIELDS.links]: 'alarm.links',
};

export const ALARM_FIELDS_TO_LABELS_KEYS = {
  [ALARM_FIELDS.id]: 'common.id',
  [ALARM_FIELDS.assignedInstructions]: 'alarm.fields.assignedInstructions',
  [ALARM_FIELDS.displayName]: 'alarm.fields.displayName',
  [ALARM_FIELDS.output]: 'common.output',
  [ALARM_FIELDS.longOutput]: 'common.longOutput',
  [ALARM_FIELDS.initialOutput]: 'alarm.fields.initialOutput',
  [ALARM_FIELDS.initialLongOutput]: 'alarm.fields.initialLongOutput',
  [ALARM_FIELDS.connector]: 'common.connector',
  [ALARM_FIELDS.connectorName]: 'common.connectorName',
  [ALARM_FIELDS.component]: 'common.component',
  [ALARM_FIELDS.resource]: 'common.resource',
  [ALARM_FIELDS.lastComment]: 'alarm.fields.lastComment',
  [ALARM_FIELDS.lastCommentInitiator]: 'alarm.fields.lastCommentInitiator',
  [ALARM_FIELDS.ackBy]: 'alarm.fields.ackBy',
  [ALARM_FIELDS.ackMessage]: 'alarm.fields.ackMessage',
  [ALARM_FIELDS.ackInitiator]: 'alarm.fields.ackInitiator',
  [ALARM_FIELDS.stateMessage]: 'alarm.fields.stateMessage',
  [ALARM_FIELDS.statusMessage]: 'alarm.fields.statusMessage',
  [ALARM_FIELDS.state]: 'common.state',
  [ALARM_FIELDS.status]: 'common.status',
  [ALARM_FIELDS.totalStateChanges]: 'alarm.fields.totalStateChanges',
  [ALARM_FIELDS.timestamp]: 'common.timestamp',
  [ALARM_FIELDS.creationDate]: 'common.created',
  [ALARM_FIELDS.lastUpdateDate]: 'common.updated',
  [ALARM_FIELDS.lastEventDate]: 'common.lastEventDate',
  [ALARM_FIELDS.ackAt]: 'alarm.fields.ackAt',
  [ALARM_FIELDS.stateAt]: 'alarm.fields.stateAt',
  [ALARM_FIELDS.statusAt]: 'alarm.fields.statusAt',
  [ALARM_FIELDS.resolved]: 'alarm.fields.resolved',
  [ALARM_FIELDS.activationDate]: 'alarm.fields.activationDate',
  [ALARM_FIELDS.duration]: 'common.duration',
  [ALARM_FIELDS.currentStateDuration]: 'alarm.fields.currentStateDuration',
  [ALARM_FIELDS.snoozeDuration]: 'alarm.fields.snoozeDuration',
  [ALARM_FIELDS.pbhInactiveDuration]: 'alarm.fields.pbhInactiveDuration',
  [ALARM_FIELDS.activeDuration]: 'alarm.fields.activeDuration',
  [ALARM_FIELDS.eventsCount]: 'alarm.fields.eventsCount',
  [ALARM_FIELDS.tags]: 'common.tag',
  [ALARM_FIELDS.extraDetails]: 'alarm.fields.extraDetails',
  [ALARM_FIELDS.impactState]: 'common.impactState',
  [ALARM_FIELDS.infos]: 'common.infos',
  [ALARM_FIELDS.links]: 'common.link',
  [ALARM_FIELDS.ticketAuthor]: 'alarm.fields.ticketAuthor',
  [ALARM_FIELDS.ticketMessage]: 'alarm.fields.ticketMessage',
  [ALARM_FIELDS.ticketValue]: 'alarm.fields.ticketId',
  [ALARM_FIELDS.ticketInitiator]: 'alarm.fields.ticketInitiator',
  [ALARM_FIELDS.ticketCreatedAt]: 'alarm.fields.ticketCreatedAt',
  [ALARM_FIELDS.entityId]: 'alarm.fields.entityId',
  [ALARM_FIELDS.entityName]: 'alarm.fields.entityName',
  [ALARM_FIELDS.entityCategoryName]: 'alarm.fields.entityCategoryName',
  [ALARM_FIELDS.entityType]: 'alarm.fields.entityType',
  [ALARM_FIELDS.entityComponent]: 'alarm.fields.entityComponent',
  [ALARM_FIELDS.entityConnector]: 'alarm.fields.entityConnector',
  [ALARM_FIELDS.entityImpactLevel]: 'alarm.fields.entityImpactLevel',
  [ALARM_FIELDS.entityKoEvents]: 'alarm.fields.entityKoEvents',
  [ALARM_FIELDS.entityOkEvents]: 'alarm.fields.entityOkEvents',
  [ALARM_FIELDS.entityInfos]: 'alarm.fields.entityInfos',
  [ALARM_FIELDS.entityComponentInfos]: 'alarm.fields.entityComponentInfos',
  [ALARM_FIELDS.entityLastPbehaviorDate]: 'alarm.fields.entityLastPbehaviorDate',
  [ALARM_FIELDS.changeState]: 'alarm.fields.changeState',
  [ALARM_FIELDS.meta]: 'alarm.metaAlarm',

  /**
   * OBJECTS
   */
  [ALARM_FIELDS.ack]: 'common.ack',
  [ALARM_FIELDS.ticket]: 'common.ticket',
  [ALARM_FIELDS.ticketData]: 'alarm.fields.ticketData',
  [ALARM_FIELDS.canceled]: 'common.canceled',
  [ALARM_FIELDS.canceledInitiator]: 'alarm.fields.canceledInitiator',
  [ALARM_FIELDS.snooze]: 'common.snooze',
  [ALARM_FIELDS.pbehaviorInfo]: 'pbehavior.pbehaviorInfo',

  /**
   * VIRTUAL
   */
  [ALARM_FIELDS.activated]: 'common.activated',
};

export const ALARM_UNSORTABLE_FIELDS = [
  ALARM_FIELDS.extraDetails,
  ALARM_FIELDS.links,
  ALARM_FIELDS.tags,
];

export const ALARM_DENSE_TYPES = {
  large: 0,
  medium: 1,
  small: 2,
};

export const ALARM_PAYLOADS_VARIABLES = {
  alarm: '.Alarm',
  alarms: '.Alarms',
  component: '.Value.Component',
  resource: '.Value.Resource',
  stateMessage: '.Value.State.Message',
  stateValue: '.Value.State.Value',
  statusValue: '.Value.Status.Value',
  ticketAuthor: '.Value.Ticket.Author',
  ticketValue: '.Value.Ticket.Ticket',
  ticketMessage: '.Value.Ticket.Message',
  ackAuthor: '.Value.ACK.Author',
  ackMessage: '.Value.ACK.Message',
  lastCommentAuthor: '.Value.LastComment.Author',
  lastCommentMessage: '.Value.LastComment.Message',
  infosValue: '(index (index .Alarm.Value.Infos "%rule_id%") "%infos_name%")',

  entityName: '.Entity.Name',
  entityInfosValue: '(index .Entity.Infos "%infos_name%").Value',
};

export const ACK_MODAL_ACTIONS_TYPES = {
  ack: 0,
  ackAndAssociateTicket: 1,
  ackAndDeclareTicket: 2,
};

export const ALARMS_EXPAND_PANEL_TABS = {
  moreInfos: 'moreInfos',
  timeLine: 'timeLine',
  ticketsDeclared: 'ticketsDeclared',
  pbehavior: 'pbehavior',
  alarmsChildren: 'alarmsChildren',
  trackSource: 'trackSource',
  impactChain: 'impactChain',
  entityGantt: 'entityGantt',
  charts: 'charts',
};

export const ALARM_LINK_ICON_CHIP_WIDTH = 24;

export const ALARM_LINK_ICON_CHIP_COLUMN_GAP = {
  [ALARM_DENSE_TYPES.large]: 8,
  [ALARM_DENSE_TYPES.medium]: 8,
  [ALARM_DENSE_TYPES.small]: 4,
};

export const ALARM_LINK_TD_PADDINGS = {
  [ALARM_DENSE_TYPES.large]: 8,
  [ALARM_DENSE_TYPES.medium]: 6,
  [ALARM_DENSE_TYPES.small]: 6,
};

export const ALARM_EXPORT_FILE_NAME_PREFIX = 'alarm';

export const ALARM_ACTION_BUTTON_WIDTHS = {
  [ALARM_DENSE_TYPES.large]: 36,
  [ALARM_DENSE_TYPES.medium]: 36,
  [ALARM_DENSE_TYPES.small]: 24,
};

export const ALARM_ACTION_BUTTON_MARGINS = {
  [ALARM_DENSE_TYPES.large]: 8,
  [ALARM_DENSE_TYPES.medium]: 0,
  [ALARM_DENSE_TYPES.small]: 0,
};

export const ALARM_ACTIONS_PADDINGS = {
  [ALARM_DENSE_TYPES.large]: 8,
  [ALARM_DENSE_TYPES.medium]: 6,
  [ALARM_DENSE_TYPES.small]: 6,
};

export const DEFAULT_ALARM_ACTIONS_INLINE_COUNT = 4;
