import { ALARM_METRIC_PARAMETERS } from '@/constants/alarm';
import { USER_METRIC_PARAMETERS } from '@/constants/user';

export const KPI_SLI_GRAPH_BAR_PERCENTAGE = 0.5;

export const KPI_ALARMS_GRAPH_BAR_PERCENTAGE = 0.75;

export const KPI_SLI_GRAPH_DATA_TYPE = {
  percent: 'percent',
  time: 'time',
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
