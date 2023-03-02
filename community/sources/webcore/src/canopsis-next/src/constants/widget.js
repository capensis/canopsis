import { MEDIA_QUERIES_BREAKPOINTS } from '@/config';

import featuresService from '@/services/features';

import { CANOPSIS_EDITION } from './permission';

export const WIDGET_TYPES = {
  alarmList: 'AlarmsList',
  context: 'Context',
  serviceWeather: 'ServiceWeather',
  statsCalendar: 'StatsCalendar',
  text: 'Text',
  counter: 'Counter',
  testingWeather: 'Junit',
  map: 'Map',

  ...featuresService.get('constants.WIDGET_TYPES'),
};

export const COMPONENTS_BY_WIDGET_TYPES = {
  [WIDGET_TYPES.alarmList]: 'alarms-list-widget',
  [WIDGET_TYPES.context]: 'entities-list-widget',
  [WIDGET_TYPES.serviceWeather]: 'service-weather-widget',
  [WIDGET_TYPES.statsCalendar]: 'stats-calendar-widget',
  [WIDGET_TYPES.text]: 'text-widget',
  [WIDGET_TYPES.counter]: 'counter-widget',
  [WIDGET_TYPES.testingWeather]: 'testing-weather-widget',
  [WIDGET_TYPES.map]: 'map-widget',

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

  ...featuresService.get('constants.SIDE_BARS_BY_WIDGET_TYPES'),
};

export const WIDGET_TYPES_RULES = {
  [WIDGET_TYPES.statsCalendar]: { edition: CANOPSIS_EDITION.pro },
  [WIDGET_TYPES.testingWeather]: { edition: CANOPSIS_EDITION.pro },
  [WIDGET_TYPES.map]: { edition: CANOPSIS_EDITION.pro },

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
  weatherItem: 'weather_item',
  weatherModal: 'weather_modal',
  weatherEntity: 'weather_entity',
};

export const COLUMNS_WIDGET_TEMPLATES_TYPES = [
  WIDGET_TEMPLATES_TYPES.alarmColumns,
  WIDGET_TEMPLATES_TYPES.entityColumns,
];

export const CUSTOM_WIDGET_TEMPLATE = Symbol('custom').toString();
