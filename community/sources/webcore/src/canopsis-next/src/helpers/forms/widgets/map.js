import {
  COLOR_INDICATOR_TYPES,
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_CONTEXT_WIDGET_COLUMNS,
  DEFAULT_MAP_ENTITY_TEMPLATE,
  DEFAULT_PERIODIC_REFRESH,
} from '@/constants';

import { durationWithEnabledToForm } from '@/helpers/date/duration';
import { defaultColumnsToColumns } from '@/helpers/entities';

import { widgetColumnsToForm, formToWidgetColumns } from '../shared/widget-column';

/**
 * @typedef {Object} MapWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {string} map
 * @property {string} color_indicator
 * @property {string} entity_info_template
 * @property {boolean} entities_under_pbehavior_enabled
 * @property {string | null} mainFilter
 * @property {WidgetColumn[]} alarmsColumns
 * @property {WidgetColumn[]} entitiesColumns
 */

/**
 * @typedef {MapWidgetParameters} MapWidgetParametersForm
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
  alarmsColumns: parameters.alarmsColumns
    ? widgetColumnsToForm(parameters.alarmsColumns)
    : defaultColumnsToColumns(DEFAULT_ALARMS_WIDGET_COLUMNS),
  entitiesColumns: parameters.entitiesColumns
    ? widgetColumnsToForm(parameters.entitiesColumns)
    : defaultColumnsToColumns(DEFAULT_CONTEXT_WIDGET_COLUMNS),
});

/**
 * Convert form to map widget parameters
 *
 * @param {MapWidgetParametersForm} form
 * @return {MapWidgetParameters}
 */
export const formToMapWidgetParameters = form => ({
  ...form,

  alarmsColumns: formToWidgetColumns(form.alarmsColumns),
  entitiesColumns: formToWidgetColumns(form.entitiesColumns),
});
