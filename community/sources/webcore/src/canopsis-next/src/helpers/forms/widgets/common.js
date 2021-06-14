import { cloneDeep } from 'lodash';

import {
  DEFAULT_WIDGET_GRID_PARAMETERS,
  WIDGET_GRID_SIZES_KEYS,
  WIDGET_TYPES,
} from '@/constants';

import uuid from '@/helpers/uuid';

/**
 * @typedef { 'AlarmsList' | 'Context' | 'ServiceWeather' | 'StatsCalendar' | 'Text' | 'Counter' | 'Junit' } WidgetType
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
  _id: widget._id || uuid(`widget_${widget.type}`),
  type: widget.type,
  title: widget.title || '',
  parameters: widget.parameters ? cloneDeep(widget.parameters) : {},
  grid_parameters: Object.values(WIDGET_GRID_SIZES_KEYS).reduce((acc, size) => {
    acc[size] = { ...DEFAULT_WIDGET_GRID_PARAMETERS };

    return acc;
  }, {}),
});
