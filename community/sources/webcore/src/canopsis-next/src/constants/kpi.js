import { ALARM_METRIC_PARAMETERS } from '@/constants/alarm';
import { USER_METRIC_PARAMETERS } from '@/constants/user';

export const KPI_SLI_GRAPH_BAR_PERCENTAGE = 0.5;

export const KPI_ALARMS_GRAPH_BAR_PERCENTAGE = 0.75;

export const KPI_SLI_GRAPH_DATA_TYPE = {
  percent: 'percent',
  time: 'time',
};

export const KPI_RATING_CRITERIA = {
  user: 1,
  role: 2,
  category: 3,
  impactLevel: 4,
};

export const KPI_RATING_USER_CRITERIA = [KPI_RATING_CRITERIA.user, KPI_RATING_CRITERIA.role];

export const KPI_RATING_ENTITY_CRITERIA = [KPI_RATING_CRITERIA.category, KPI_RATING_CRITERIA.impactLevel];

export const KPI_RATING_USER_METRICS = [
  ALARM_METRIC_PARAMETERS.ticketAlarms,
  ALARM_METRIC_PARAMETERS.ackAlarms,
  ALARM_METRIC_PARAMETERS.cancelAckAlarms,
  ALARM_METRIC_PARAMETERS.ackWithoutCancelAlarms,
  ALARM_METRIC_PARAMETERS.averageAck,
  USER_METRIC_PARAMETERS.totalUserActivity,
];

export const KPI_RATING_ENTITY_METRICS = [
  ALARM_METRIC_PARAMETERS.totalAlarms,
  ALARM_METRIC_PARAMETERS.instructionAlarms,
  ALARM_METRIC_PARAMETERS.pbehaviorAlarms,
  ALARM_METRIC_PARAMETERS.ticketAlarms,
  ALARM_METRIC_PARAMETERS.withoutTicketAlarms,
  ALARM_METRIC_PARAMETERS.nonDisplayedAlarms,
  ALARM_METRIC_PARAMETERS.ratioCorrelation,
  ALARM_METRIC_PARAMETERS.ratioInstructions,
  ALARM_METRIC_PARAMETERS.ratioTickets,
  ALARM_METRIC_PARAMETERS.ratioNonDisplayed,
  ALARM_METRIC_PARAMETERS.averageAck,
  ALARM_METRIC_PARAMETERS.averageResolve,
];

export const KPI_TABS = {
  graphs: 'graphs',
  filters: 'filters',
};
