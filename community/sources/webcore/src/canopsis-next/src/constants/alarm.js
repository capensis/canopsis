import { COLORS } from '@/config';

export const ALARM_FIELDS = {
  id: '_id',
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
  canceled: 'v.canceled',
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

export const ALARM_LIST_ACTIONS_TYPES = {
  ack: 'ack',
  fastAck: 'fastAck',
  ackRemove: 'ackRemove',
  pbehaviorAdd: 'pbehaviorAdd',
  moreInfos: 'moreInfos',
  snooze: 'snooze',
  declareTicket: 'declareTicket',
  associateTicket: 'associateTicket',
  cancel: 'cancel',
  changeState: 'changeState',
  variablesHelp: 'variablesHelp',
  history: 'history',
  groupRequest: 'groupRequest',
  createManualMetaAlarm: 'createManualMetaAlarm',
  removeAlarmsFromManualMetaAlarm: 'removeAlarmsFromManualMetaAlarm',
  updateManualMetaAlarm: 'updateManualMetaAlarm',
  comment: 'comment',

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
  { value: ALARM_FIELDS.connector },
  { value: ALARM_FIELDS.connectorName },
  { value: ALARM_FIELDS.component },
  { value: ALARM_FIELDS.resource },
  { value: ALARM_FIELDS.output },
  { value: ALARM_FIELDS.extraDetails },
  { value: ALARM_FIELDS.state },
  { value: ALARM_FIELDS.status },
];

export const DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS = [
  { value: ALARM_FIELDS.connector },
  { value: ALARM_FIELDS.connectorName },
  { value: ALARM_FIELDS.resource },
  { value: ALARM_FIELDS.output },
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
};

export const ALARM_OPTIONAL_METRIC_PARAMETERS = {
  manualInstructionAssignedAlarms: 'manual_instruction_assigned_alarms',
  manualInstructionExecutedAlarms: 'manual_instruction_executed_alarms',
  notAckedAlarms: 'not_acked_alarms',
  notAckedInHourAlarms: 'not_acked_in_hour_alarms',
  notAckedInFourHoursAlarms: 'not_acked_in_four_hours_alarms',
  notAckedInDayAlarms: 'not_acked_in_day_alarms',
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
  snooze: ALARM_FIELDS.snooze,
  lastComment: ALARM_FIELDS.lastComment,
  longOutput: ALARM_FIELDS.longOutput,
  initialOutput: ALARM_FIELDS.initialOutput,
  initialLongOutput: ALARM_FIELDS.initialLongOutput,
  totalStateChanges: ALARM_FIELDS.totalStateChanges,
  tags: ALARM_FIELDS.tags,
  activated: ALARM_FIELDS.activated,
};

export const ALARM_ACK_INITIATORS = {
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

export const ALARM_FIELDS_TO_LABELS_KEYS = {
  [ALARM_FIELDS.id]: 'common.id',
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

  /**
   * OBJECTS
   */
  [ALARM_FIELDS.ack]: 'common.ack',
  [ALARM_FIELDS.ticket]: 'common.ticket',
  [ALARM_FIELDS.canceled]: 'common.canceled',
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

  entityName: '.Entity.Name',
  entityInfosValue: '(index .Entity.Infos.%infos_name%).Value',
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
};
