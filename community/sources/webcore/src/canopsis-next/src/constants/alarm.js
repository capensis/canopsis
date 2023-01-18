import { COLORS } from '@/config';

import featuresService from '@/services/features';

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
  manualMetaAlarmGroup: 'manualMetaAlarmGroup',
  manualMetaAlarmUngroup: 'manualMetaAlarmUngroup',
  manualMetaAlarmUpdate: 'manualMetaAlarmUpdate',
  comment: 'comment',

  ...featuresService.get('constants.ALARM_LIST_ACTIONS_TYPES'),

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

export const ALARM_ENTITY_FIELDS = { // TODO: update fields
  connector: 'v.connector',
  connectorName: 'v.connector_name',
  component: 'v.component',
  resource: 'v.resource',
  output: 'v.output',
  extraDetails: 'extra_details',
  priority: 'priority',
  impactState: 'impact_state',
  state: 'v.state.val',
  status: 'v.status.val',
  tags: 'tags',
};

export const DEFAULT_ALARMS_WIDGET_COLUMNS = [
  { value: ALARM_ENTITY_FIELDS.connector },
  { value: ALARM_ENTITY_FIELDS.connectorName },
  { value: ALARM_ENTITY_FIELDS.component },
  { value: ALARM_ENTITY_FIELDS.resource },
  { value: ALARM_ENTITY_FIELDS.output },
  { value: ALARM_ENTITY_FIELDS.extraDetails },
  { value: ALARM_ENTITY_FIELDS.state },
  { value: ALARM_ENTITY_FIELDS.status },
];

export const DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS = [
  { value: ALARM_ENTITY_FIELDS.connector },
  { value: ALARM_ENTITY_FIELDS.connectorName },
  { value: ALARM_ENTITY_FIELDS.resource },
  { value: ALARM_ENTITY_FIELDS.output },
  { value: ALARM_ENTITY_FIELDS.extraDetails },
  { value: ALARM_ENTITY_FIELDS.state },
  { value: ALARM_ENTITY_FIELDS.status },
];

export const MANUAL_META_ALARM_EVENT_DEFAULT_FIELDS = {
  component: 'metaalarm',
  connector: 'engine',
  connector_name: 'correlation',
  source_type: 'metaalarm',
};

export const ALARMS_OPENED_VALUES = {
  opened: true,
  all: null,
  resolved: false,
};

export const ALARMS_LIST_WIDGET_ACTIVE_COLUMNS_MAP = {
  priority: 'impact_state',
};

export const ALARM_METRIC_PARAMETERS = {
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
  manualInstructionAssignedAlarms: 'manual_instruction_assigned_alarms',
  manualInstructionExecutedAlarms: 'manual_instruction_executed_alarms',
};

export const ALARMS_LIST_HEADER_OPACITY_DELAY = 500;

export const ALARM_INTERVAL_FIELDS = {
  timestamp: 't',
  resolved: 'v.resolved',
  lastUpdateDate: 'v.last_update_date',
  lastEventDate: 'v.last_event_date',
};

export const ALARM_PATTERN_FIELDS = {
  displayName: 'v.display_name',
  state: 'v.state.val',
  status: 'v.status.val',
  component: 'v.component',
  resource: 'v.resource',
  connector: 'v.connector',
  connectorName: 'v.connector_name',
  creationDate: 'v.creation_date',
  duration: 'v.duration',
  infos: 'v.infos',
  output: 'v.output',
  lastEventDate: 'v.last_event_date',
  lastUpdateDate: 'v.last_update_date',
  activationDate: 'v.activation_date',
  ack: 'v.ack',
  ackAt: 'v.ack.t',
  ackBy: 'v.ack.a',
  ackMessage: 'v.ack.m',
  ackInitiator: 'v.ack.initiator',
  resolvedAt: 'v.resolved',
  ticket: 'v.ticket',
  canceled: 'v.canceled',
  snooze: 'v.snooze',
  lastComment: 'v.last_comment.m',
  longOutput: 'v.long_output',
  initialOutput: 'v.initial_output',
  initialLongOutput: 'v.initial_long_output',
  totalStateChanges: 'v.total_state_changes',
  tags: 'tags',
  activated: 'activated',
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
  id: 'alarm._id',
  ack: 'alarm.v.ack',
  state: 'alarm.v.state',
  status: 'alarm.v.status',
  ticket: 'alarm.v.ticket',
  component: 'alarm.v.component',
  connector: 'alarm.v.connector',
  connectorName: 'alarm.v.connector_name',
  resource: 'alarm.v.resource',
  creationDate: 'alarm.v.creation_date',
  displayName: 'alarm.v.display_name',
  output: 'alarm.v.output',
  lastUpdateDate: 'alarm.v.last_update_date',
  lastEventDate: 'alarm.v.last_event_date',
  pbehaviorInfo: 'alarm.v.pbehavior_info',
  duration: 'alarm.v.duration',
  eventsCount: 'alarm.v.events_count',
};

export const ALARM_LIST_WIDGET_COLUMNS = {
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
  lastEventDate: 'v.last_event_date',
  lastUpdateDate: 'v.last_update_date',
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
};

export const ALARM_LIST_WIDGET_COLUMNS_TO_LABELS_KEYS = {
  [ALARM_LIST_WIDGET_COLUMNS.id]: 'common.id',
  [ALARM_LIST_WIDGET_COLUMNS.displayName]: 'common.displayName',
  [ALARM_LIST_WIDGET_COLUMNS.output]: 'common.output',
  [ALARM_LIST_WIDGET_COLUMNS.longOutput]: 'common.longOutput',
  [ALARM_LIST_WIDGET_COLUMNS.initialOutput]: 'common.initialOutput',
  [ALARM_LIST_WIDGET_COLUMNS.initialLongOutput]: 'common.initialLongOutput',
  [ALARM_LIST_WIDGET_COLUMNS.connector]: 'common.connector',
  [ALARM_LIST_WIDGET_COLUMNS.connectorName]: 'common.connectorName',
  [ALARM_LIST_WIDGET_COLUMNS.component]: 'common.component',
  [ALARM_LIST_WIDGET_COLUMNS.resource]: 'common.resource',
  [ALARM_LIST_WIDGET_COLUMNS.lastComment]: 'common.lastComment',
  [ALARM_LIST_WIDGET_COLUMNS.ackBy]: 'common.ackBy',
  [ALARM_LIST_WIDGET_COLUMNS.ackMessage]: 'common.ackMessage',
  [ALARM_LIST_WIDGET_COLUMNS.ackInitiator]: 'common.ackInitiator',
  [ALARM_LIST_WIDGET_COLUMNS.stateMessage]: 'common.stateMessage',
  [ALARM_LIST_WIDGET_COLUMNS.statusMessage]: 'common.statusMessage',
  [ALARM_LIST_WIDGET_COLUMNS.state]: 'common.state',
  [ALARM_LIST_WIDGET_COLUMNS.status]: 'common.status',
  [ALARM_LIST_WIDGET_COLUMNS.totalStateChanges]: 'common.totalStateChanges',
  [ALARM_LIST_WIDGET_COLUMNS.timestamp]: 'common.timestamp',
  [ALARM_LIST_WIDGET_COLUMNS.creationDate]: 'common.creationDate',
  [ALARM_LIST_WIDGET_COLUMNS.lastEventDate]: 'common.lastEventDate',
  [ALARM_LIST_WIDGET_COLUMNS.lastUpdateDate]: 'common.lastUpdateDate',
  [ALARM_LIST_WIDGET_COLUMNS.ackAt]: 'common.ackAt',
  [ALARM_LIST_WIDGET_COLUMNS.stateAt]: 'common.stateAt',
  [ALARM_LIST_WIDGET_COLUMNS.statusAt]: 'common.statusAt',
  [ALARM_LIST_WIDGET_COLUMNS.resolved]: 'common.resolved',
  [ALARM_LIST_WIDGET_COLUMNS.activationDate]: 'common.activationDate',
  [ALARM_LIST_WIDGET_COLUMNS.duration]: 'common.duration',
  [ALARM_LIST_WIDGET_COLUMNS.currentStateDuration]: 'common.currentStateDuration',
  [ALARM_LIST_WIDGET_COLUMNS.snoozeDuration]: 'common.snoozeDuration',
  [ALARM_LIST_WIDGET_COLUMNS.pbhInactiveDuration]: 'common.pbhInactiveDuration',
  [ALARM_LIST_WIDGET_COLUMNS.activeDuration]: 'common.activeDuration',
  [ALARM_LIST_WIDGET_COLUMNS.tags]: 'common.tags',
  [ALARM_LIST_WIDGET_COLUMNS.extraDetails]: 'common.extraDetails',
  [ALARM_LIST_WIDGET_COLUMNS.impactState]: 'common.impactState',
  [ALARM_LIST_WIDGET_COLUMNS.infos]: 'common.infos',
  [ALARM_LIST_WIDGET_COLUMNS.links]: 'common.links',
  [ALARM_LIST_WIDGET_COLUMNS.entityId]: 'common.entityId',
  [ALARM_LIST_WIDGET_COLUMNS.entityName]: 'common.entityName',
  [ALARM_LIST_WIDGET_COLUMNS.entityCategoryName]: 'common.entityCategoryName',
  [ALARM_LIST_WIDGET_COLUMNS.entityType]: 'common.entityType',
  [ALARM_LIST_WIDGET_COLUMNS.entityComponent]: 'common.entityComponent',
  [ALARM_LIST_WIDGET_COLUMNS.entityConnector]: 'common.entityConnector',
  [ALARM_LIST_WIDGET_COLUMNS.entityImpactLevel]: 'common.entityImpactLevel',
  [ALARM_LIST_WIDGET_COLUMNS.entityKoEvents]: 'common.entityKoEvents',
  [ALARM_LIST_WIDGET_COLUMNS.entityOkEvents]: 'common.entityOkEvents',
  [ALARM_LIST_WIDGET_COLUMNS.entityInfos]: 'common.entityInfos',
  [ALARM_LIST_WIDGET_COLUMNS.entityComponentInfos]: 'common.entityComponentInfos',
};
