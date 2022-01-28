import { cloneDeep } from 'lodash';

import {
  WIDGET_TYPES,
  WIDGET_GRID_SIZES_KEYS,
  DEFAULT_WIDGET_GRID_PARAMETERS,
} from '@/constants';

import {
  alarmListWidgetParametersToFormParameters,
  formParametersToAlarmListWidgetParameters,
} from './alarm';
import {
  contextWidgetParametersToFormParameters,
  formParametersToContextWidgetParameters,
} from '@/helpers/forms/widgets/context';
import {
  serviceWeatherWidgetParametersToFormParameters,
} from '@/helpers/forms/widgets/service-weather';

/**
 * @typedef { 'AlarmsList' | 'Context' | 'ServiceWeather' | 'StatsCalendar' | 'Text' | 'Counter' | 'Junit' } WidgetType
 */

/**
 * @typedef { 'comma' | 'semicolon' | 'tab' | 'space' } CsvSeparators
 */

/**
 * @typedef {Object} WidgetColumn
 * @property {string} value
 * @property {string} label
 * @property {boolean} isHtml
 * @property {boolean} colorIndicator
 */

/**
 * @typedef {Object} WidgetGridParameter
 * @property {number} y
 * @property {number} x
 * @property {number} w
 * @property {number} h
 * @property {number} h
 * @property {boolean} autoHeight
 */

/**
 * @typedef {Object} WidgetGridParameters
 * @property {WidgetGridParameter} mobile
 * @property {WidgetGridParameter} tablet
 * @property {WidgetGridParameter} desktop
 */

/**
 * @typedef {Object} Widget
 * @property {string} _id
 * @property {WidgetType} type
 * @property {string} title
 * @property {Object} parameters
 * @property {WidgetGridParameters} grid_parameters
 */

/**
 * @typedef {Widget} WidgetForm
 */

export const widgetParametersToFormParameters = ({ type, parameters } = {}) => {
  switch (type) {
    case WIDGET_TYPES.alarmList:
      return alarmListWidgetParametersToFormParameters(parameters);
    case WIDGET_TYPES.context:
      return contextWidgetParametersToFormParameters(parameters);
    case WIDGET_TYPES.serviceWeather:
      return serviceWeatherWidgetParametersToFormParameters(parameters);
    default:
      return parameters ? cloneDeep(parameters) : {};
  }
};

export const getDefaultGridParameters = () => Object.values(WIDGET_GRID_SIZES_KEYS).reduce((acc, size) => {
  acc[size] = { ...DEFAULT_WIDGET_GRID_PARAMETERS };

  return acc;
}, {});

export const getEmptyWidgetByType = type => ({
  type,
  title: '',
  parameters: {},
  grid_parameters: getDefaultGridParameters(),
});

/**
 * Convert widget to form object
 *
 * @param {Widget} [widget = {}]
 * @returns {WidgetForm}
 */
export const widgetToForm = (widget = { type: WIDGET_TYPES.alarmList }) => ({ // TODO: We've removed ID generation here
  type: widget.type,
  title: widget.title ?? '',
  parameters: widgetParametersToFormParameters(widget),
  grid_parameters: widget.grid_parameters
    ? cloneDeep(widget.grid_parameters)
    : getDefaultGridParameters(),
});

export const formParametersToWidgetParameters = ({ type, parameters = {} } = {}) => {
  switch (type) {
    case WIDGET_TYPES.alarmList:
      return formParametersToAlarmListWidgetParameters(parameters);
    case WIDGET_TYPES.context:
      return formParametersToContextWidgetParameters(parameters);
    default:
      return parameters;
  }
};

/**
 * Convert form object to widget
 *
 * @param {WidgetForm} form
 * @returns {Widget}
 */
export const formToWidget = form => ({
  ...form,

  parameters: formParametersToWidgetParameters(form),
});
