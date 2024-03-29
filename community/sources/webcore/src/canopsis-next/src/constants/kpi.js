import { ALARM_METRIC_PARAMETERS } from './alarm';
import { USER_METRIC_PARAMETERS } from './user';

export const KPI_SLI_GRAPH_BAR_PERCENTAGE = 0.5;

export const KPI_ALARMS_GRAPH_BAR_PERCENTAGE = 0.75;

export const KPI_SLI_GRAPH_DATA_TYPE = {
  percent: 'percent',
  time: 'time',
};

export const KPI_PIE_CHART_SHOW_MODS = {
  percent: 'percent',
  numbers: 'numbers',
};

export const KPI_RATING_CRITERIA = {
  user: 'username',
  role: 'role',
  category: 'category',
  impactLevel: 'impact_level',
};

export const KPI_RATING_USER_CRITERIA = [KPI_RATING_CRITERIA.user, KPI_RATING_CRITERIA.role];

export const KPI_RATING_USER_METRICS = [
  ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
  ALARM_METRIC_PARAMETERS.ackAlarms,
  ALARM_METRIC_PARAMETERS.cancelAckAlarms,
  ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms,
  ALARM_METRIC_PARAMETERS.averageAck,
  USER_METRIC_PARAMETERS.totalUserActivity,
];

export const KPI_RATING_ENTITY_METRICS = [
  ALARM_METRIC_PARAMETERS.createdAlarms,
  ALARM_METRIC_PARAMETERS.instructionAlarms,
  ALARM_METRIC_PARAMETERS.pbehaviorAlarms,
  ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
  ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms,
  ALARM_METRIC_PARAMETERS.nonDisplayedAlarms,
  ALARM_METRIC_PARAMETERS.ratioCorrelation,
  ALARM_METRIC_PARAMETERS.ratioInstructions,
  ALARM_METRIC_PARAMETERS.ratioTickets,
  ALARM_METRIC_PARAMETERS.ratioNonDisplayed,
  ALARM_METRIC_PARAMETERS.averageAck,
];

export const KPI_TABS = {
  graphs: 'graphs',
  filters: 'filters',
  ratingSettings: 'ratingSettings',
  collectionSettings: 'collectionSettings',
};

export const KPI_METRICS_MAX_ALARM_YEAR_INTERVAL_DIFF_IN_YEARS = 1;

export const AGGREGATE_FUNCTIONS = {
  sum: 'sum',
  avg: 'avg',
  min: 'min',
  max: 'max',
  last: 'last',
};

export const Y_AXES_IDS = {
  default: 'y',
  percent: 'yPercent',
  time: 'yTime',
  bytes: 'yBytes',
};

export const EXTERNAL_METRIC_UNITS = {
  millisecond: 'ms',
  microsecond: 'us',
  second: 's',
  continuousCounter: 'c',
  byte: 'B',
  kilobyte: 'KB',
  megabyte: 'MB',
  gigabyte: 'GB',
  terabyte: 'TB',
  percent: '%',
};

export const X_AXES_IDS = {
  default: 'x',
  history: 'xHistory',
};

export const MAX_METRICS_DISPLAY_COUNT = 40;

export const KPI_CHART_DEFAULT_HEIGHT = 560;

export const KPI_RATING_SETTINGS_TYPES = {
  entity: 0,
  user: 1,
};

export const KPI_ENTITY_RATING_SETTINGS_CUSTOM_CRITERIA = Symbol('custom').toString();

export const STATISTICS_WIDGETS_USER_METRICS = [
  ALARM_METRIC_PARAMETERS.ackAlarms,
  ALARM_METRIC_PARAMETERS.cancelAckAlarms,
  ALARM_METRIC_PARAMETERS.minAck,
  ALARM_METRIC_PARAMETERS.maxAck,
  ALARM_METRIC_PARAMETERS.averageAck,
  USER_METRIC_PARAMETERS.ackAlarmWithoutCancel,
  USER_METRIC_PARAMETERS.tickets,
  USER_METRIC_PARAMETERS.totalUserActivity,
  USER_METRIC_PARAMETERS.averageUserSession,
  USER_METRIC_PARAMETERS.minUserSession,
  USER_METRIC_PARAMETERS.maxUserSession,
];

export const STATISTICS_WIDGETS_ENTITY_METRICS = [
  ALARM_METRIC_PARAMETERS.createdAlarms,
  ALARM_METRIC_PARAMETERS.activeAlarms,
  ALARM_METRIC_PARAMETERS.instructionAlarms,
  ALARM_METRIC_PARAMETERS.manualInstructionAssignedAlarms,
  ALARM_METRIC_PARAMETERS.manualInstructionExecutedAlarms,
  ALARM_METRIC_PARAMETERS.nonDisplayedAlarms,
  ALARM_METRIC_PARAMETERS.pbehaviorAlarms,
  ALARM_METRIC_PARAMETERS.correlationAlarms,
  ALARM_METRIC_PARAMETERS.ackActiveAlarms,
  ALARM_METRIC_PARAMETERS.notAckedAlarms,
  ALARM_METRIC_PARAMETERS.notAckedInHourAlarms,
  ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms,
  ALARM_METRIC_PARAMETERS.notAckedInDayAlarms,
  ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
  ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms,
  ALARM_METRIC_PARAMETERS.ratioCorrelation,
  ALARM_METRIC_PARAMETERS.ratioInstructions,
  ALARM_METRIC_PARAMETERS.ratioRemediatedAlarms,
  ALARM_METRIC_PARAMETERS.ratioTickets,
  ALARM_METRIC_PARAMETERS.ratioNonDisplayed,
  ALARM_METRIC_PARAMETERS.timeToAck,
  ALARM_METRIC_PARAMETERS.minAck,
  ALARM_METRIC_PARAMETERS.maxAck,
  ALARM_METRIC_PARAMETERS.averageResolve,
  ALARM_METRIC_PARAMETERS.minResolve,
  ALARM_METRIC_PARAMETERS.maxResolve,
];

export const STATISTICS_WIDGETS_USER_METRICS_WITH_ENTITY_TYPE = [
  ALARM_METRIC_PARAMETERS.ackAlarms,
  ALARM_METRIC_PARAMETERS.cancelAckAlarms,
  ALARM_METRIC_PARAMETERS.minAck,
  ALARM_METRIC_PARAMETERS.maxAck,
  ALARM_METRIC_PARAMETERS.averageAck,
  USER_METRIC_PARAMETERS.ackAlarmWithoutCancel,
  USER_METRIC_PARAMETERS.tickets,
];
