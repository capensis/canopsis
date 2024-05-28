import { cloneDeep } from 'lodash';

import { DEFAULT_WEATHER_LIMIT, PAGINATION_LIMIT } from '@/config';
import {
  SORT_ORDERS,
  COLOR_INDICATOR_TYPES,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
  DEFAULT_PERIODIC_REFRESH,
  DEFAULT_SERVICE_WEATHER_BLOCK_TEMPLATE,
  DEFAULT_SERVICE_WEATHER_MODAL_TEMPLATE,
  DEFAULT_SERVICE_WEATHER_ENTITY_TEMPLATE,
  DEFAULT_WIDGET_MARGIN,
} from '@/constants';

import { durationWithEnabledToForm } from '@/helpers/date/duration';

import { alarmListBaseParametersToForm, formToAlarmListBaseParameters } from './alarm';
import { widgetColumnsToForm, formToWidgetColumns } from '../shared/widget-column';
import { widgetTemplateValueToForm, formToWidgetTemplateValue } from '../widget-template';

/**
 * @typedef {'more-info' | 'alarm-list' | 'both'} ServiceWeatherWidgetModalType
 */

/**
 * @typedef {Object} ServiceWeatherWidgetCounters
 * @property {boolean} enabled
 * @property {string[]} types
 */

/**
 * @typedef {Object} ServiceWeatherWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {string | null} mainFilter
 * @property {WidgetSort} sort
 * @property {string} blockTemplateTemplate
 * @property {string} modalTemplateTemplate
 * @property {string} entityTemplateTemplate
 * @property {string} blockTemplate
 * @property {string} modalTemplate
 * @property {string} entityTemplate
 * @property {number} columnMobile
 * @property {number} columnTablet
 * @property {number} columnDesktop
 * @property {number} limit
 * @property {ColorIndicator} colorIndicator
 * @property {ServiceWeatherWidgetModalType} modalType
 * @property {string} serviceDependenciesColumnsTemplate
 * @property {WidgetColumn[]} serviceDependenciesColumns
 * @property {WidgetMargin} margin
 * @property {ServiceWeatherWidgetCounters} counters
 * @property {number} heightFactor
 * @property {number} modalItemsPerPage
 * @property {boolean} isPriorityEnabled
 * @property {boolean} isHideGrayEnabled
 * @property {AlarmListBaseParameters} alarmsList
 * @property {boolean} entitiesActionsInQueue
 */

/**
 * @typedef {ServiceWeatherWidgetParameters} ServiceWeatherWidgetParametersForm
 * @property {string | Symbol} serviceDependenciesColumnsTemplate
 * @property {WidgetColumnForm[]} serviceDependenciesColumns
 */

/**
 * Convert service weather widget parameters to form
 *
 * @param {ServiceWeatherWidgetParameters} [parameters = {}]
 * @return {ServiceWeatherWidgetParametersForm}
 */
export const serviceWeatherWidgetParametersToForm = (parameters = {}) => ({
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
  mainFilter: parameters.mainFilter ?? null,
  sort: parameters.sort ? { ...parameters.sort } : { order: SORT_ORDERS.asc },
  blockTemplateTemplate: widgetTemplateValueToForm(parameters.blockTemplateTemplate),
  modalTemplateTemplate: widgetTemplateValueToForm(parameters.modalTemplateTemplate),
  entityTemplateTemplate: widgetTemplateValueToForm(parameters.entityTemplateTemplate),
  blockTemplate: parameters.blockTemplate ?? DEFAULT_SERVICE_WEATHER_BLOCK_TEMPLATE,
  modalTemplate: parameters.modalTemplate ?? DEFAULT_SERVICE_WEATHER_MODAL_TEMPLATE,
  entityTemplate: parameters.entityTemplate ?? DEFAULT_SERVICE_WEATHER_ENTITY_TEMPLATE,
  columnMobile: parameters.columnMobile ?? 2,
  columnTablet: parameters.columnTablet ?? 3,
  columnDesktop: parameters.columnDesktop ?? 4,
  limit: parameters.limit ?? DEFAULT_WEATHER_LIMIT,
  colorIndicator: parameters.colorIndicator ?? COLOR_INDICATOR_TYPES.state,
  serviceDependenciesColumnsTemplate: widgetTemplateValueToForm(parameters.serviceDependenciesColumnsTemplate),
  serviceDependenciesColumns:
    widgetColumnsToForm(parameters.serviceDependenciesColumns ?? DEFAULT_SERVICE_DEPENDENCIES_COLUMNS),
  margin: parameters.margin
    ? { ...parameters.margin }
    : { ...DEFAULT_WIDGET_MARGIN },
  heightFactor: parameters.heightFactor ?? 6,
  modalType: parameters.modalType ?? SERVICE_WEATHER_WIDGET_MODAL_TYPES.both,
  modalItemsPerPage: parameters.modalItemsPerPage ?? PAGINATION_LIMIT,
  alarmsList: alarmListBaseParametersToForm(parameters.alarmsList),
  counters: parameters.counters
    ? cloneDeep(parameters.counters)
    : {
      pbehavior_enabled: false,
      pbehavior_types: [],
      state_enabled: false,
      state_types: [],
    },
  isPriorityEnabled: parameters.isPriorityEnabled ?? true,
  isHideGrayEnabled: parameters.isHideGrayEnabled ?? true,
  entitiesActionsInQueue: parameters.entitiesActionsInQueue ?? false,
});

/**
 * Convert form to service weather widget parameters
 *
 * @param {ServiceWeatherWidgetParametersForm} form
 * @return {ServiceWeatherWidgetParameters}
 */
export const formToServiceWeatherWidgetParameters = form => ({
  ...form,

  blockTemplateTemplate: formToWidgetTemplateValue(form.blockTemplateTemplate),
  modalTemplateTemplate: formToWidgetTemplateValue(form.modalTemplateTemplate),
  entityTemplateTemplate: formToWidgetTemplateValue(form.entityTemplateTemplate),
  serviceDependenciesColumnsTemplate: formToWidgetTemplateValue(form.serviceDependenciesColumnsTemplate),
  serviceDependenciesColumns: formToWidgetColumns(form.serviceDependenciesColumns),
  alarmsList: formToAlarmListBaseParameters(form.alarmsList),
});
