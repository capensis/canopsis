import { ALARM_METRIC_PARAMETERS, USER_METRIC_PARAMETERS, AGGREGATE_FUNCTIONS } from '@/constants';

export default {
  alarmMetrics: 'Alarm metrics',
  sli: 'SLI',
  metricsNotAvailable: 'TimescaleDB not running. Metrics are not available.',
  noData: 'No data available',
  selectMetric: 'Select metric to display',
  customColor: 'Custom color',
  calculationMethod: 'Calculation method',
  periodTrend: '{count} for the period\n{from} - {to}',
  largeCountOfMetrics: 'The list of metrics to display is too large and truncated.',
  onlyDisplayed: 'Only {count} metrics are displayed.',
  tabs: {
    collectionSettings: 'Collection settings',
    ratingSettings: 'Rating settings',
  },

  aggregateFunctions: {
    [AGGREGATE_FUNCTIONS.sum]: 'Sum',
    [AGGREGATE_FUNCTIONS.avg]: 'Average',
    [AGGREGATE_FUNCTIONS.min]: 'Min',
    [AGGREGATE_FUNCTIONS.max]: 'Max',
  },

  errors: {
    metricsMinLength: 'At least {count} metrics must be added',
    emptyMetrics: 'No data to display',
  },

  popups: {
    exportFailed: 'Failed to export chart in CSV format',
  },

  metrics: {
    parameter: 'Parameter to compare',
    tooltip: {
      [USER_METRIC_PARAMETERS.totalUserActivity]: '{value} total activity time',

      [ALARM_METRIC_PARAMETERS.createdAlarms]: '{value} created alarms',
      [ALARM_METRIC_PARAMETERS.activeAlarms]: '{value} active alarms',
      [ALARM_METRIC_PARAMETERS.nonDisplayedAlarms]: '{value} non-displayed alarms',
      [ALARM_METRIC_PARAMETERS.instructionAlarms]: '{value} alarms under auto remediation',
      [ALARM_METRIC_PARAMETERS.pbehaviorAlarms]: '{value} alarms under PBehavior',
      [ALARM_METRIC_PARAMETERS.correlationAlarms]: '{value} alarms with correlation',
      [ALARM_METRIC_PARAMETERS.ackAlarms]: '{value} alarms with acks',
      [ALARM_METRIC_PARAMETERS.ackActiveAlarms]: '{value} active alarms with acks',
      [ALARM_METRIC_PARAMETERS.cancelAckAlarms]: '{value} alarms with cancelled acks',
      [ALARM_METRIC_PARAMETERS.ticketActiveAlarms]: '{value} active alarms with tickets',
      [ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms]: '{value} active alarms without tickets',
      [ALARM_METRIC_PARAMETERS.ratioCorrelation]: '{value}% of alarms with auto remediation',
      [ALARM_METRIC_PARAMETERS.ratioInstructions]: '{value}% alarms with instructions',
      [ALARM_METRIC_PARAMETERS.ratioTickets]: '{value}% of alarms with tickets created',
      [ALARM_METRIC_PARAMETERS.ratioNonDisplayed]: '{value}% of non-displayed alarms',
      [ALARM_METRIC_PARAMETERS.ratioRemediatedAlarms]: '{value}% of manually remediated alarms',
      [ALARM_METRIC_PARAMETERS.averageAck]: '{value} to ack alarms',
      [ALARM_METRIC_PARAMETERS.averageResolve]: '{value} to resolve alarms',
      [ALARM_METRIC_PARAMETERS.timeToAck]: '{value} to ack alarms',
      [ALARM_METRIC_PARAMETERS.timeToResolve]: '{value} to resolve alarms',
      [ALARM_METRIC_PARAMETERS.minAck]: '{value} - minimum time to ack alarms',
      [ALARM_METRIC_PARAMETERS.maxAck]: '{value} - maximum time to ack alarms',
      [ALARM_METRIC_PARAMETERS.manualInstructionExecutedAlarms]: '{value} manually remediated alarms',
      [ALARM_METRIC_PARAMETERS.manualInstructionAssignedAlarms]: '{value} alarms with manual instructions',
      [ALARM_METRIC_PARAMETERS.notAckedAlarms]: '{value} not acked alarms',
      [ALARM_METRIC_PARAMETERS.notAckedInHourAlarms]: '{value} not acked alarms with duration 1-4h',
      [ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms]: '{value} not acked alarms with duration 4-24h',
      [ALARM_METRIC_PARAMETERS.notAckedInDayAlarms]: '{value} not acked alarms older than 24h',
    },
  },

  filters: {
    helpInformation: 'Here the filter patterns for additional slices of data for counters and ratings can be added.',
  },

  ratingSettings: {
    helpInformation: 'The list of parameters to use for rating.',
  },

  collectionSetting: {
    basicMetrics: 'Basic metrics',
    optionalMetrics: 'Optional metrics',
    manualInstructions: 'Number of alarms with manual instructions',
    notAckedMetrics: 'Number of active not acked alarms of different durations',
  },
};
