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

import { defaultColumnsToColumns } from '@/helpers/entities';

import {
  alarmListBaseParametersToForm,
  formToAlarmListBaseParameters,
} from './alarm';

/**
 * @typedef {'more-info' | 'alarm-list' | 'both'} ServiceWeatherWidgetModalType
 */

/**
 * @typedef {'impact-state' | 'state'} ServiceWeatherWidgetColorIndicator
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
 * @property {string} blockTemplate
 * @property {string} modalTemplate
 * @property {string} entityTemplate
 * @property {number} columnMobile
 * @property {number} columnTablet
 * @property {number} columnDesktop
 * @property {number} limit
 * @property {ServiceWeatherWidgetColorIndicator} colorIndicator
 * @property {ServiceWeatherWidgetModalType} modalType
 * @property {WidgetColumn[]} serviceDependenciesColumns
 * @property {WidgetMargin} margin
 * @property {ServiceWeatherWidgetCounters} counters
 * @property {number} heightFactor
 * @property {number} modalItemsPerPage
 * @property {boolean} isPriorityEnabled
 * @property {AlarmListBaseParameters} alarmsList
 */

/**
 * Convert service weather widget parameters to form
 *
 * @param {ServiceWeatherWidgetParameters} [parameters = {}]
 * @return {ServiceWeatherWidgetParameters}
 */
export const serviceWeatherWidgetParametersToForm = (parameters = {}) => ({
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
  mainFilter: parameters.mainFilter ?? null,
  sort: parameters.sort ? { ...parameters.sort } : { order: SORT_ORDERS.asc },
  blockTemplate: parameters.blockTemplate ?? DEFAULT_SERVICE_WEATHER_BLOCK_TEMPLATE,
  modalTemplate: parameters.modalTemplate ?? DEFAULT_SERVICE_WEATHER_MODAL_TEMPLATE,
  entityTemplate: parameters.entityTemplate ?? DEFAULT_SERVICE_WEATHER_ENTITY_TEMPLATE,
  columnMobile: parameters.columnMobile ?? 2,
  columnTablet: parameters.columnTablet ?? 4,
  columnDesktop: parameters.columnDesktop ?? 4,
  limit: parameters.limit ?? DEFAULT_WEATHER_LIMIT,
  colorIndicator: parameters.colorIndicator ?? COLOR_INDICATOR_TYPES.state,
  serviceDependenciesColumns: parameters.serviceDependenciesColumns
    ? cloneDeep(parameters.serviceDependenciesColumns)
    : defaultColumnsToColumns(DEFAULT_SERVICE_DEPENDENCIES_COLUMNS),
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
});

/**
 * Convert form to service weather widget parameters
 *
 * @param {ServiceWeatherWidgetParameters} form
 * @return {ServiceWeatherWidgetParameters}
 */
export const formToServiceWeatherWidgetParameters = form => ({
  ...form,

  alarmsList: formToAlarmListBaseParameters(form.alarmsList),
});
