import {
  COLOR_INDICATOR_TYPES,
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_CONTEXT_WIDGET_COLUMNS,
  DEFAULT_MAP_ENTITY_TEMPLATE,
  DEFAULT_PERIODIC_REFRESH,
} from '@/constants';

import { durationWithEnabledToForm } from '@/helpers/date/duration';

import { widgetColumnsToForm, formToWidgetColumns } from '../shared/widget-column';
import { widgetTemplateValueToForm, formToWidgetTemplateValue } from '../widget-template';

/**
 * @typedef {Object} MapWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {string} map
 * @property {string} color_indicator
 * @property {string} entity_info_template
 * @property {boolean} entities_under_pbehavior_enabled
 * @property {string | null} mainFilter
 * @property {string} alarmsColumnsTemplate
 * @property {string} entitiesColumnsTemplate
 * @property {WidgetColumn[]} alarmsColumns
 * @property {WidgetColumn[]} entitiesColumns
 */

/**
 * @typedef {MapWidgetParameters} MapWidgetParametersForm
 * @property {string | Symbol} alarmsColumnsTemplate
 * @property {string | Symbol} entitiesColumnsTemplate
 * @property {WidgetColumnForm[]} alarmsColumns
 * @property {WidgetColumnForm[]} entitiesColumns
 */

/**
 * Convert map widget parameters to form
 *
 * @param {MapWidgetParameters} parameters
 * @return {MapWidgetParametersForm}
 */
export const mapWidgetParametersToForm = (parameters = {}) => ({
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
  map: parameters.map ?? '',
  color_indicator: parameters.color_indicator ?? COLOR_INDICATOR_TYPES.state,
  entities_under_pbehavior_enabled: parameters.entities_under_pbehavior_enabled ?? true,
  mainFilter: parameters.mainFilter ?? null,
  entity_info_template: parameters.entity_info_template ?? DEFAULT_MAP_ENTITY_TEMPLATE,
  alarmsColumnsTemplate: widgetTemplateValueToForm(parameters.alarmsColumnsTemplate),
  entitiesColumnsTemplate: widgetTemplateValueToForm(parameters.entitiesColumnsTemplate),
  alarmsColumns:
    widgetColumnsToForm(parameters.alarmsColumns ?? DEFAULT_ALARMS_WIDGET_COLUMNS),
  entitiesColumns:
    widgetColumnsToForm(parameters.entitiesColumns ?? DEFAULT_CONTEXT_WIDGET_COLUMNS),
});

/**
 * Convert form to map widget parameters
 *
 * @param {MapWidgetParametersForm} form
 * @return {MapWidgetParameters}
 */
export const formToMapWidgetParameters = form => ({
  ...form,

  alarmsColumnsTemplate: formToWidgetTemplateValue(form.alarmsColumnsTemplate),
  entitiesColumnsTemplate: formToWidgetTemplateValue(form.entitiesColumnsTemplate),
  alarmsColumns: formToWidgetColumns(form.alarmsColumns),
  entitiesColumns: formToWidgetColumns(form.entitiesColumns),
});
