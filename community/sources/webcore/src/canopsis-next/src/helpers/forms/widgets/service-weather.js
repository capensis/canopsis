import { cloneDeep } from 'lodash';

import { DEFAULT_WEATHER_LIMIT, PAGINATION_LIMIT } from '@/config';
import {
  COLOR_INDICATOR_TYPES,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
} from '@/constants';

import { defaultColumnsToColumns } from '@/helpers/entities';
import { alarmListWidgetDefaultParametersToFormParameters, widgetSortToForm } from './alarm';
import { widgetFiltersParametersToFormParameters } from './context';

export const serviceWeatherWidgetParametersToFormParameters = (parameters = {}) => ({
  ...widgetFiltersParametersToFormParameters(parameters),

  sort: widgetSortToForm(parameters.sort),
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
  isCountersEnabled: parameters.isCountersEnabled ?? false,
  heightFactor: parameters.heightFactor ?? 6,
  modalType: parameters.modalType ?? SERVICE_WEATHER_WIDGET_MODAL_TYPES.both,
  alarmsList: alarmListWidgetDefaultParametersToFormParameters(parameters.alarmList), // TODO: put this to func below
  modalItemsPerPage: parameters.modalItemsPerPage ?? PAGINATION_LIMIT,
});

export const formParametersToServiceWeatherWidgetParameters = (parameters = {}) => ({
  ...parameters,

});
