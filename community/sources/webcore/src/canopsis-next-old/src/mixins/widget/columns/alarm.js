import { RESIZING_CELLS_CONTENTS_BEHAVIORS, COLOR_INDICATOR_TYPES } from '@/constants';

import {
  getAlarmsListWidgetColumnValueFilter,
  getAlarmsListWidgetColumnComponentGetter,
} from '@/helpers/entities/widget/forms/alarm';

import { entitiesAlarmColumnsFiltersMixin } from '../../entities/associative-table/alarm-columns-filters';

import { widgetColumnsEntityInfoPropertyMixin } from './entity-info-property';

export const widgetColumnsAlarmMixin = {
  mixins: [entitiesAlarmColumnsFiltersMixin, widgetColumnsEntityInfoPropertyMixin],
  data() {
    return {
      columnsFilters: [],
      columnsFiltersPending: false,
    };
  },
  computed: {
    infoPopupsMap() {
      return (this.widget.parameters?.infoPopups ?? []).reduce((acc, { column, template }) => {
        acc[column] = template;

        return acc;
      }, {});
    },

    columnsFiltersMap() {
      return (this.columnsFilters ?? []).reduce((acc, { column, filter, attributes = [] }) => {
        acc[column] = this.getFilter(filter, attributes);

        return acc;
      }, {});
    },

    preparedColumns() {
      const {
        showRootCauseByStateClick = true,
        columns: columnsSettings = {},
      } = this.widget.parameters ?? {};

      return (this.columns ?? []).map(column => ({
        ...column,
        text: column.label || this.findAliasByColumnValue(column.value, 'entity') || column.text,
        popupTemplate: this.infoPopupsMap[column.value] ?? this.infoPopupsMap[`alarm.${column.value}`],
        filter: this.$i18n.locale && this.getColumnFilter(column.value),
        getComponent: getAlarmsListWidgetColumnComponentGetter(
          column,
          {
            showRootCauseByStateClick,
            isCellContentTruncated:
              columnsSettings.cells_content_behavior === RESIZING_CELLS_CONTENTS_BEHAVIORS.truncate,
          },
        ),
        colorIndicatorEnabled: Object.values(COLOR_INDICATOR_TYPES).includes(column.colorIndicator),
      }));
    },
  },
  mounted() {
    this.fetchColumnFilters();
    this.fetchEntityInfoPropertiesList({ params: { paginate: false } });
  },
  methods: {
    getColumnFilter(value) {
      return this.columnsFiltersMap[value] ?? getAlarmsListWidgetColumnValueFilter(value);
    },

    getFilter(filter, attributes = []) {
      const filterFunc = this.$options.filters[filter];

      return value => (filterFunc ? filterFunc(value, ...attributes) : value);
    },

    async fetchColumnFilters() {
      try {
        this.columnsFiltersPending = true;
        this.columnsFilters = await this.fetchAlarmColumnsFiltersList();
      } catch (err) {
        console.warn(err);
      } finally {
        this.columnsFiltersPending = false;
      }
    },
  },
};
