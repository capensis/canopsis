import { cloneDeep } from 'lodash';

import { DEFAULT_WEATHER_LIMIT, PAGINATION_LIMIT } from '@/config';
import {
  SORT_ORDERS,
  COLOR_INDICATOR_TYPES,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
} from '@/constants';

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
 * @property {WidgetFilter[]} filters
 * @property {string | null} main_filter
 * @property {WidgetSort} sort
 * @property {string} blockTemplate
 * @property {string} modalTemplate
 * @property {string} entityTemplate
 * @property {number} columnSM
 * @property {number} columnMD
 * @property {number} columnLG
 * @property {number} limit
 * @property {ServiceWeatherWidgetColorIndicator} colorIndicator
 * @property {ServiceWeatherWidgetModalType} modalType
 * @property {WidgetColumn[]} serviceDependenciesColumns
 * @property {WidgetMargin} margin
 * @property {ServiceWeatherWidgetCounters} counters
 * @property {number} heightFactor
 * @property {number} modalItemsPerPage
 * @property {AlarmListBaseParameters} alarmsList
 */

/**
 * Convert service weather widget parameters to form
 *
 * @param {ServiceWeatherWidgetParameters} [parameters = {}]
 * @return {ServiceWeatherWidgetParameters}
 */
export const serviceWeatherWidgetParametersToFormParameters = (parameters = {}) => ({
  filters: parameters.filters
    ? cloneDeep(parameters.filters)
    : [],
  main_filter: parameters.main_filter ?? null,
  sort: parameters.sort ? { ...parameters.sort } : { order: SORT_ORDERS.asc },
  blockTemplate: parameters.blockTemplate ?? `<p><strong><span style="font-size: 18px;">{{entity.name}}</span></strong></p>
<hr id="null">
<p>{{ entity.output }}</p>
<p> Dernière mise à jour : {{ timestamp entity.last_update_date }}</p>`,

  modalTemplate: parameters.modalTemplate ?? '{{ entities name="entity._id" }}',
  entityTemplate: parameters.entityTemplate ?? `<ul>
    <li><strong>Libellé</strong> : {{entity.name}}</li>
</ul>`,
  columnSM: parameters.columnSM ?? 6,
  columnMD: parameters.columnMD ?? 4,
  columnLG: parameters.columnLG ?? 3,
  limit: parameters.limit ?? DEFAULT_WEATHER_LIMIT,
  colorIndicator: parameters.colorIndicator ?? COLOR_INDICATOR_TYPES.state,
  serviceDependenciesColumns: parameters.serviceDependenciesColumns
    ? cloneDeep(parameters.serviceDependenciesColumns)
    : defaultColumnsToColumns(DEFAULT_SERVICE_DEPENDENCIES_COLUMNS),
  margin: parameters.margin
    ? { ...parameters.margin }
    : {
      top: 1,
      right: 1,
      bottom: 1,
      left: 1,
    },
  heightFactor: parameters.heightFactor ?? 6,
  modalType: parameters.modalType ?? SERVICE_WEATHER_WIDGET_MODAL_TYPES.both,
  modalItemsPerPage: parameters.modalItemsPerPage ?? PAGINATION_LIMIT,
  alarmsList: alarmListBaseParametersToForm(parameters.alarmsList),
});

/**
 * Convert form to service weather widget parameters
 *
 * @param {ServiceWeatherWidgetParameters} form
 * @return {ServiceWeatherWidgetParameters}
 */
export const formParametersToServiceWeatherWidgetParameters = form => ({
  ...form,

  alarmsList: formToAlarmListBaseParameters(form.alarmsList),
});
