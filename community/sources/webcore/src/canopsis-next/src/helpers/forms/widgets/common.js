import { cloneDeep } from 'lodash';

import {
  WIDGET_TYPES,
  WIDGET_GRID_SIZES_KEYS,
  DEFAULT_WIDGET_GRID_PARAMETERS,
} from '@/constants';

import { generateWidgetId } from '@/helpers/entities';

import { formToAlarmListWidget } from './alarm';

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

/**
 * Convert widget to form object
 *
 * @param {Widget} [widget = {}]
 * @returns {WidgetForm}
 */
export const widgetToForm = (widget = { type: WIDGET_TYPES.alarmList }) => ({
  _id: widget._id || generateWidgetId(widget.type),
  type: widget.type,
  title: widget.title || '',
  parameters: widget.parameters ? cloneDeep(widget.parameters) : {},
  grid_parameters: Object.values(WIDGET_GRID_SIZES_KEYS).reduce((acc, size) => {
    acc[size] = { ...DEFAULT_WIDGET_GRID_PARAMETERS };

    return acc;
  }, {}),
});

/**
 * Convert form object to widget
 *
 * @param {WidgetForm} form
 * @returns {Widget}
 */
export const formToWidget = (form) => {
  const method = {
    [WIDGET_TYPES.alarmList]: formToAlarmListWidget,
  }[form.type];

  return method ? method(form) : form;
};
