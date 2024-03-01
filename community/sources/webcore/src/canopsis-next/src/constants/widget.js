import { MEDIA_QUERIES_BREAKPOINTS } from '@/config';

import featuresService from '@/services/features';

import { CANOPSIS_EDITION } from './permission';
import { ALARM_METRIC_PARAMETERS } from './alarm';
import { QUICK_RANGES } from './common';
import { SAMPLINGS } from './date';
import { AGGREGATE_FUNCTIONS, KPI_PIE_CHART_SHOW_MODS } from './kpi';

export const WIDGET_TYPES = {
  alarmList: 'AlarmsList',
  context: 'Context',
  serviceWeather: 'ServiceWeather',
  statsCalendar: 'StatsCalendar',
  text: 'Text',
  counter: 'Counter',
  testingWeather: 'Junit',
  map: 'Map',
  barChart: 'BarChart',
  lineChart: 'LineChart',
  pieChart: 'PieChart',
  numbers: 'Numbers',
  userStatistics: 'UserStatistics',
  alarmStatistics: 'AlarmStatistics',
  availability: 'Availability',

  ...featuresService.get('constants.WIDGET_TYPES'),
};

export const TOP_LEVEL_WIDGET_TYPES = [
  WIDGET_TYPES.alarmList,
  WIDGET_TYPES.context,
  WIDGET_TYPES.serviceWeather,
  WIDGET_TYPES.statsCalendar,
  WIDGET_TYPES.text,
  WIDGET_TYPES.counter,
  WIDGET_TYPES.testingWeather,
  WIDGET_TYPES.map,
  WIDGET_TYPES.availability,
];

export const COMPONENTS_BY_WIDGET_TYPES = {
  [WIDGET_TYPES.alarmList]: 'alarms-list-widget',
  [WIDGET_TYPES.context]: 'entities-list-widget',
  [WIDGET_TYPES.serviceWeather]: 'service-weather-widget',
  [WIDGET_TYPES.statsCalendar]: 'stats-calendar-widget',
  [WIDGET_TYPES.text]: 'text-widget',
  [WIDGET_TYPES.counter]: 'counter-widget',
  [WIDGET_TYPES.testingWeather]: 'testing-weather-widget',
  [WIDGET_TYPES.map]: 'map-widget',
  [WIDGET_TYPES.barChart]: 'bar-chart-widget',
  [WIDGET_TYPES.lineChart]: 'line-chart-widget',
  [WIDGET_TYPES.pieChart]: 'pie-chart-widget',
  [WIDGET_TYPES.numbers]: 'numbers-widget',
  [WIDGET_TYPES.userStatistics]: 'user-statistics-widget',
  [WIDGET_TYPES.alarmStatistics]: 'alarm-statistics-widget',
  [WIDGET_TYPES.availability]: 'availability-widget',

  ...featuresService.get('constants.COMPONENTS_BY_WIDGET_TYPES'),
};

export const WIDGET_ICONS = {
  [WIDGET_TYPES.alarmList]: 'view_list',
  [WIDGET_TYPES.context]: 'view_list',
  [WIDGET_TYPES.serviceWeather]: 'view_module',
  [WIDGET_TYPES.statsCalendar]: 'calendar_today',
  [WIDGET_TYPES.text]: 'view_headline',
  [WIDGET_TYPES.counter]: 'view_module',
  [WIDGET_TYPES.testingWeather]: 'view_module',
  [WIDGET_TYPES.map]: 'location_on',
  [WIDGET_TYPES.barChart]: 'stacked_bar_chart',
  [WIDGET_TYPES.lineChart]: 'show_chart',
  [WIDGET_TYPES.pieChart]: 'donut_small',
  [WIDGET_TYPES.numbers]: 'functions',
  [WIDGET_TYPES.numbers]: 'functions',
  [WIDGET_TYPES.userStatistics]: 'perm_identity',
  [WIDGET_TYPES.alarmStatistics]: 'notification_important',
  /**
   * TODO: Icon should be changed
   */
  [WIDGET_TYPES.availability]: 'notification_important',

  ...featuresService.get('constants.WIDGET_ICONS'),
};

export const SIDE_BARS = {
  alarmSettings: 'alarm-settings',
  contextSettings: 'context-settings',
  serviceWeatherSettings: 'service-weather-settings',
  statsCalendarSettings: 'stats-calendar-settings',
  textSettings: 'text-settings',
  counterSettings: 'counter-settings',
  testingWeatherSettings: 'testing-weather-settings',
  mapSettings: 'map-settings',
  barChartSettings: 'bar-chart-settings',
  lineChartSettings: 'line-chart-settings',
  pieChartSettings: 'pie-chart-settings',
  numbersSettings: 'numbers-settings',
  userStatisticsSettings: 'user-statistics-settings',
  alarmStatisticsSettings: 'alarm-statistics-settings',
  availabilitySettings: 'availability-settings',

  ...featuresService.get('constants.SIDE_BARS'),
};

export const SIDE_BARS_BY_WIDGET_TYPES = {
  [WIDGET_TYPES.alarmList]: SIDE_BARS.alarmSettings,
  [WIDGET_TYPES.context]: SIDE_BARS.contextSettings,
  [WIDGET_TYPES.serviceWeather]: SIDE_BARS.serviceWeatherSettings,
  [WIDGET_TYPES.statsCalendar]: SIDE_BARS.statsCalendarSettings,
  [WIDGET_TYPES.text]: SIDE_BARS.textSettings,
  [WIDGET_TYPES.counter]: SIDE_BARS.counterSettings,
  [WIDGET_TYPES.testingWeather]: SIDE_BARS.testingWeatherSettings,
  [WIDGET_TYPES.map]: SIDE_BARS.mapSettings,
  [WIDGET_TYPES.barChart]: SIDE_BARS.barChartSettings,
  [WIDGET_TYPES.lineChart]: SIDE_BARS.lineChartSettings,
  [WIDGET_TYPES.pieChart]: SIDE_BARS.pieChartSettings,
  [WIDGET_TYPES.numbers]: SIDE_BARS.numbersSettings,
  [WIDGET_TYPES.userStatistics]: SIDE_BARS.userStatisticsSettings,
  [WIDGET_TYPES.alarmStatistics]: SIDE_BARS.alarmStatisticsSettings,
  [WIDGET_TYPES.availability]: SIDE_BARS.availabilitySettings,

  ...featuresService.get('constants.SIDE_BARS_BY_WIDGET_TYPES'),
};

export const WIDGET_TYPES_RULES = {
  [WIDGET_TYPES.statsCalendar]: { edition: CANOPSIS_EDITION.pro },
  [WIDGET_TYPES.testingWeather]: { edition: CANOPSIS_EDITION.pro },
  [WIDGET_TYPES.map]: { edition: CANOPSIS_EDITION.pro },
  [WIDGET_TYPES.userStatistics]: { edition: CANOPSIS_EDITION.pro },
  [WIDGET_TYPES.alarmStatistics]: { edition: CANOPSIS_EDITION.pro },
  [WIDGET_TYPES.availability]: { edition: CANOPSIS_EDITION.pro },

  ...featuresService.get('constants.WIDGET_TYPES_RULES'),
};

export const WIDGET_GRID_SIZES_KEYS = {
  mobile: 'mobile',
  tablet: 'tablet',
  desktop: 'desktop',
};

export const MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP = {
  m: WIDGET_GRID_SIZES_KEYS.mobile,
  t: WIDGET_GRID_SIZES_KEYS.tablet,
  l: WIDGET_GRID_SIZES_KEYS.desktop,
  xl: WIDGET_GRID_SIZES_KEYS.desktop,
};

export const WIDGET_LAYOUT_MAX_WIDTHS = {
  [WIDGET_GRID_SIZES_KEYS.desktop]: '100%',
  [WIDGET_GRID_SIZES_KEYS.tablet]: `${MEDIA_QUERIES_BREAKPOINTS.t}px`,
  [WIDGET_GRID_SIZES_KEYS.mobile]: `${MEDIA_QUERIES_BREAKPOINTS.m}px`,
};

export const WIDGET_GRID_SIZES_STYLES = {
  [WIDGET_GRID_SIZES_KEYS.mobile]: {
    value: WIDGET_GRID_SIZES_KEYS.mobile,
    icon: 'stay_primary_portrait',
  },
  [WIDGET_GRID_SIZES_KEYS.tablet]: {
    value: WIDGET_GRID_SIZES_KEYS.tablet,
    icon: 'tablet_mac',
  },
  [WIDGET_GRID_SIZES_KEYS.desktop]: {
    value: WIDGET_GRID_SIZES_KEYS.desktop,
    icon: 'desktop_windows',
  },
};

export const WIDGET_GRID_ROW_HEIGHT = 20;

export const WIDGET_GRID_COLUMNS_COUNT = 12;

export const DEFAULT_WIDGET_GRID_PARAMETERS = {
  x: 0,
  y: 0,
  h: 1,
  w: WIDGET_GRID_COLUMNS_COUNT,
  autoHeight: true,
};

export const EXPORT_CSV_SEPARATORS = {
  comma: 'comma',
  semicolon: 'semicolon',
  tab: 'tab',
  space: 'space',
};

export const EXPORT_CSV_DATETIME_FORMATS = {
  datetimeSecondsLocal: {
    value: 'YYYY-MM-DDThh:mm:ss',
    text: '2021-01-02T15:04:05',
  },
  datetimeSeconds: {
    value: 'YYYY-MM-DDThh:mm:ssZ',
    text: '2021-01-02T15:04:05+0200',
  },
  dayOfMonthMonthNameYearTimeLocal: {
    value: 'DD MMM YYYY hh:mm:ss',
    text: '02 Jan 2021 15:04',
  },
  dayOfMonthMonthNameYearTime: {
    value: 'DD MMM YYYY hh:mm:ss ZZ',
    text: '02 Jan 2021 15:04 CEST',
  },
  weekdayDayOfMonthMonthNameYearTime: {
    value: 'W, DD MMM YYYY hh:mm:ss ZZ',
    text: 'Mon, 02 Jan 2021 15:04:05 CEST',
  },
};

export const EXPORT_STATUSES = {
  running: 0,
  completed: 1,
  failed: 2,
};

export const WIDGET_TEMPLATES_TYPES = {
  alarmColumns: 'alarm_columns',
  entityColumns: 'entity_columns',
  alarmMoreInfos: 'alarm_more_infos',
  alarmExportToPdf: 'alarm_export_to_pdf',
  weatherItem: 'weather_item',
  weatherModal: 'weather_modal',
  weatherEntity: 'weather_entity',
};

export const COLUMNS_WIDGET_TEMPLATES_TYPES = [
  WIDGET_TEMPLATES_TYPES.alarmColumns,
  WIDGET_TEMPLATES_TYPES.entityColumns,
];

export const CUSTOM_WIDGET_TEMPLATE = Symbol('custom').toString();

export const CHART_WIDGET_PRESET_TYPES = {
  numberOfActiveAlarms: 'number_of_active_alarms',
  ackStatistics: 'ack_statistics',
  ticketsStatistics: 'tickets_statistics',
  ackCancellation: 'ack_cancellation',
  activeAck: 'active_ack',
  notAckedAlarms: 'not_acked_alarms',
  nonDisplayedAlarms: 'non_displayed_alarms',
  manualInstruction: 'manual_instruction',
  numberOfCreatedAlarms: 'number_of_created_alarms',
};

export const BAR_CHART_WIDGET_PRESET_TYPES = [
  'numberOfActiveAlarms',
  'ackStatistics',
  'ticketsStatistics',
  'ackCancellation',
];

export const PIE_CHART_WIDGET_PRESET_TYPES = [
  'ackStatistics',
  'ticketsStatistics',
  'activeAck',
  'notAckedAlarms',
  'nonDisplayedAlarms',
  'manualInstruction',
];

export const NUMBERS_CHART_WIDGET_PRESET_TYPES = [
  'ackStatistics',
  'ticketsStatistics',
  'activeAck',
  'numberOfCreatedAlarms',
  'notAckedAlarms',
  'nonDisplayedAlarms',
  'manualInstruction',
];

export const BAR_CHART_WIDGET_PRESET_PARAMETERS_BY_TYPE = {
  [CHART_WIDGET_PRESET_TYPES.numberOfActiveAlarms]: {
    stacked: false,
    comparison: true,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.activeAlarms },
    ],
    default_time_range: QUICK_RANGES.last7Days.value,
    default_sampling: SAMPLINGS.day,
  },
  [CHART_WIDGET_PRESET_TYPES.ackStatistics]: {
    stacked: true,
    comparison: false,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.ackActiveAlarms },
      { metric: ALARM_METRIC_PARAMETERS.notAckedAlarms },
    ],
    default_time_range: QUICK_RANGES.last30Days.value,
    default_sampling: SAMPLINGS.day,
  },
  [CHART_WIDGET_PRESET_TYPES.ticketsStatistics]: {
    stacked: true,
    comparison: false,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.ticketActiveAlarms },
      { metric: ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms },
    ],
    default_time_range: QUICK_RANGES.last30Days.value,
    default_sampling: SAMPLINGS.day,
  },
  [CHART_WIDGET_PRESET_TYPES.ackCancellation]: {
    stacked: true,
    comparison: false,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.ackAlarms },
      { metric: ALARM_METRIC_PARAMETERS.cancelAckAlarms },
    ],
    default_time_range: QUICK_RANGES.last30Days.value,
    default_sampling: SAMPLINGS.day,
  },
};

export const PIE_CHART_WIDGET_PRESET_PARAMETERS_BY_TYPE = {
  [CHART_WIDGET_PRESET_TYPES.ackStatistics]: {
    aggregate_func: AGGREGATE_FUNCTIONS.avg,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.ackActiveAlarms },
      { metric: ALARM_METRIC_PARAMETERS.notAckedAlarms },
    ],
  },
  [CHART_WIDGET_PRESET_TYPES.ticketsStatistics]: {
    aggregate_func: AGGREGATE_FUNCTIONS.avg,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.ticketActiveAlarms },
      { metric: ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms },
    ],
  },
  [CHART_WIDGET_PRESET_TYPES.activeAck]: {
    aggregate_func: AGGREGATE_FUNCTIONS.sum,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.ackAlarms },
      { metric: ALARM_METRIC_PARAMETERS.cancelAckAlarms },
    ],
  },
  [CHART_WIDGET_PRESET_TYPES.notAckedAlarms]: {
    show_mode: KPI_PIE_CHART_SHOW_MODS.numbers,
    aggregate_func: AGGREGATE_FUNCTIONS.avg,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.notAckedInHourAlarms },
      { metric: ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms },
      { metric: ALARM_METRIC_PARAMETERS.notAckedInDayAlarms },
    ],
  },
  [CHART_WIDGET_PRESET_TYPES.nonDisplayedAlarms]: {
    show_mode: KPI_PIE_CHART_SHOW_MODS.numbers,
    aggregate_func: AGGREGATE_FUNCTIONS.avg,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.instructionAlarms },
      { metric: ALARM_METRIC_PARAMETERS.pbehaviorAlarms },
      { metric: ALARM_METRIC_PARAMETERS.correlationAlarms },
    ],
  },
  [CHART_WIDGET_PRESET_TYPES.manualInstruction]: {
    show_mode: KPI_PIE_CHART_SHOW_MODS.numbers,
    aggregate_func: AGGREGATE_FUNCTIONS.avg,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.manualInstructionAssignedAlarms },
      { metric: ALARM_METRIC_PARAMETERS.manualInstructionExecutedAlarms },
    ],
  },
};

export const NUMBERS_CHART_WIDGET_PRESET_PARAMETERS_BY_TYPE = {
  [CHART_WIDGET_PRESET_TYPES.ackStatistics]: {
    aggregate_func: AGGREGATE_FUNCTIONS.avg,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.activeAlarms },
      { metric: ALARM_METRIC_PARAMETERS.ackActiveAlarms },
      { metric: ALARM_METRIC_PARAMETERS.notAckedAlarms },
    ],
  },
  [CHART_WIDGET_PRESET_TYPES.ticketsStatistics]: {
    aggregate_func: AGGREGATE_FUNCTIONS.avg,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.activeAlarms },
      { metric: ALARM_METRIC_PARAMETERS.ticketActiveAlarms },
      { metric: ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms },
    ],
  },
  [CHART_WIDGET_PRESET_TYPES.activeAck]: {
    aggregate_func: AGGREGATE_FUNCTIONS.sum,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.createdAlarms },
      { metric: ALARM_METRIC_PARAMETERS.ackAlarms },
      { metric: ALARM_METRIC_PARAMETERS.cancelAckAlarms },
    ],
  },
  [CHART_WIDGET_PRESET_TYPES.numberOfCreatedAlarms]: {
    aggregate_func: AGGREGATE_FUNCTIONS.sum,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.createdAlarms },
    ],
  },
  [CHART_WIDGET_PRESET_TYPES.notAckedAlarms]: {
    aggregate_func: AGGREGATE_FUNCTIONS.avg,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.activeAlarms },
      { metric: ALARM_METRIC_PARAMETERS.notAckedAlarms },
      { metric: ALARM_METRIC_PARAMETERS.notAckedInHourAlarms },
      { metric: ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms },
      { metric: ALARM_METRIC_PARAMETERS.notAckedInDayAlarms },
    ],
  },
  [CHART_WIDGET_PRESET_TYPES.nonDisplayedAlarms]: {
    aggregate_func: AGGREGATE_FUNCTIONS.avg,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.activeAlarms },
      { metric: ALARM_METRIC_PARAMETERS.nonDisplayedAlarms },
      { metric: ALARM_METRIC_PARAMETERS.instructionAlarms },
      { metric: ALARM_METRIC_PARAMETERS.pbehaviorAlarms },
      { metric: ALARM_METRIC_PARAMETERS.correlationAlarms },
    ],
  },
  [CHART_WIDGET_PRESET_TYPES.manualInstruction]: {
    aggregate_func: AGGREGATE_FUNCTIONS.avg,
    metrics: [
      { metric: ALARM_METRIC_PARAMETERS.activeAlarms },
      { metric: ALARM_METRIC_PARAMETERS.manualInstructionAssignedAlarms },
      { metric: ALARM_METRIC_PARAMETERS.manualInstructionExecutedAlarms },
    ],
  },
};

export const CHART_PRESET_CUSTOM_ITEM_VALUE = Symbol('custom');

export const NUMBERS_CHART_DEFAULT_FONT_SIZE = 110;

export const NUMBERS_CHART_MAX_AUTO_FONT_SIZE = 130;

export const NUMBERS_CHART_MIN_AUTO_FONT_SIZE = 50;

export const NUMBERS_CHART_FONT_SIZE_WIDTH_COEFFICIENT = 20;
