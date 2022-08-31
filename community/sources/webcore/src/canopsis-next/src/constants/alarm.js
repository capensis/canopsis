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

export const EVENT_FILTER_RULE_TYPES = {
  drop: 'drop',
  break: 'break',
  enrichment: 'enrichment',
};

export const META_ALARMS_RULE_TYPES = {
  relation: 'relation',
  timebased: 'timebased',
  attribute: 'attribute',
  complex: 'complex',
  valuegroup: 'valuegroup',
  corel: 'corel',

  /**
   * Manual group type doesn't using in the form
   * We are using it only inside alarms list widget
   */
  manualgroup: 'manualgroup',
};

export const META_ALARMS_THRESHOLD_TYPES = {
  thresholdRate: 'thresholdRate',
  thresholdCount: 'thresholdCount',
};

export const ALARMS_LIST_TIME_LINE_SYSTEM_AUTHOR = 'canopsis.engine';

export const ALARM_ENTITY_FIELDS = {
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
};

export const DEFAULT_ALARMS_WIDGET_COLUMNS = [
  {
    labelKey: 'tables.alarmGeneral.connector',
    value: ALARM_ENTITY_FIELDS.connector,
  },
  {
    labelKey: 'tables.alarmGeneral.connectorName',
    value: ALARM_ENTITY_FIELDS.connectorName,
  },
  {
    labelKey: 'tables.alarmGeneral.component',
    value: ALARM_ENTITY_FIELDS.component,
  },
  {
    labelKey: 'tables.alarmGeneral.resource',
    value: ALARM_ENTITY_FIELDS.resource,
  },
  {
    labelKey: 'tables.alarmGeneral.output',
    value: ALARM_ENTITY_FIELDS.output,
  },
  {
    labelKey: 'tables.alarmGeneral.extraDetails',
    value: ALARM_ENTITY_FIELDS.extraDetails,
  },
  {
    labelKey: 'tables.alarmGeneral.state',
    value: ALARM_ENTITY_FIELDS.state,
  },
  {
    labelKey: 'tables.alarmGeneral.status',
    value: ALARM_ENTITY_FIELDS.status,
  },
];

export const DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS = [
  {
    labelKey: 'tables.alarmGeneral.connector',
    value: ALARM_ENTITY_FIELDS.connector,
  },
  {
    labelKey: 'tables.alarmGeneral.connectorName',
    value: ALARM_ENTITY_FIELDS.connectorName,
  },
  {
    labelKey: 'tables.alarmGeneral.resource',
    value: ALARM_ENTITY_FIELDS.resource,
  },
  {
    labelKey: 'tables.alarmGeneral.output',
    value: ALARM_ENTITY_FIELDS.output,
  },
  {
    labelKey: 'tables.alarmGeneral.extraDetails',
    value: ALARM_ENTITY_FIELDS.extraDetails,
  },
  {
    labelKey: 'tables.alarmGeneral.state',
    value: ALARM_ENTITY_FIELDS.state,
  },
  {
    labelKey: 'tables.alarmGeneral.status',
    value: ALARM_ENTITY_FIELDS.status,
  },
];

export const ALARMS_GROUP_PREFIX = {
  CAUSES: 'causes_',
  CONSEQUENCES: 'consequences_',
};

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
  creationDate: 'creation_date',
  resolved: 'resolved',
  lastUpdateDate: 'last_update_date',
  lastEventDate: 'last_event_date',
};
