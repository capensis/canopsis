import { get } from 'lodash';

import {
  ALARM_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_UNSORTABLE_FIELDS,
  COLOR_INDICATOR_TYPES,
} from '@/constants';

import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';
import { convertDurationToString } from '@/helpers/date/duration';

import { entitiesAlarmColumnsFiltersMixin } from '@/mixins/entities/associative-table/alarm-columns-filters';

import { widgetColumnsMixin } from './common';

export const widgetColumnsAlarmMixin = {
  mixins: [
    widgetColumnsMixin,
    entitiesAlarmColumnsFiltersMixin,
  ],
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

    componentGettersMap() {
      return {
        [ALARM_FIELDS.state]: context => ({
          bind: {
            is: 'alarm-column-value-state',
            alarm: context.alarm,
          },
        }),
        [ALARM_FIELDS.status]: context => ({
          bind: {
            is: 'alarm-column-value-status',
            alarm: context.alarm,
          },
        }),
        [ALARM_FIELDS.impactState]: context => ({
          bind: {
            is: 'color-indicator-wrapper',
            type: COLOR_INDICATOR_TYPES.impactState,
            entity: context.alarm.entity,
            alarm: context.alarm,
          },
        }),
        [ALARM_FIELDS.links]: context => ({
          bind: {
            is: 'alarm-column-value-categories',
            asList: get(this.widget.parameters, 'linksCategoriesAsList.enabled', false),
            limit: get(this.widget.parameters, 'linksCategoriesAsList.limit'),
            links: context.alarm.links ?? {},
          },
          on: {
            activate: context.$listeners.activate,
          },
        }),
        [ALARM_FIELDS.extraDetails]: context => ({
          bind: {
            is: 'alarm-column-value-extra-details',
            alarm: context.alarm,
          },
        }),
        [ALARM_FIELDS.tags]: context => ({
          bind: {
            is: 'c-alarm-tags-chips',
            alarm: context.alarm,
            selectedTag: context.selectedTag,
          },
          on: {
            select: context.$listeners['select:tag'],
          },
        }),
      };
    },

    columnPropertiesFiltersMap() {
      return {
        [ALARM_FIELDS.lastUpdateDate]: convertDateToStringWithFormatForToday,
        [ALARM_FIELDS.creationDate]: convertDateToStringWithFormatForToday,
        [ALARM_FIELDS.lastEventDate]: convertDateToStringWithFormatForToday,
        [ALARM_FIELDS.activationDate]: convertDateToStringWithFormatForToday,
        [ALARM_FIELDS.ackAt]: convertDateToStringWithFormatForToday,
        [ALARM_FIELDS.stateAt]: convertDateToStringWithFormatForToday,
        [ALARM_FIELDS.statusAt]: convertDateToStringWithFormatForToday,
        [ALARM_FIELDS.resolved]: convertDateToStringWithFormatForToday,
        [ALARM_FIELDS.timestamp]: convertDateToStringWithFormatForToday,
        [ALARM_FIELDS.duration]: convertDurationToString,
        [ALARM_FIELDS.currentStateDuration]: convertDurationToString,
        [ALARM_FIELDS.activeDuration]: convertDurationToString,
        [ALARM_FIELDS.snoozeDuration]: convertDurationToString,
        [ALARM_FIELDS.pbhInactiveDuration]: convertDurationToString,

        ...this.columnsFiltersMap,
      };
    },

    columns() {
      return (this.widget.parameters?.widgetColumns ?? []).map(column => ({
        ...column,

        popupTemplate: this.infoPopupsMap[column.value],
        text: this.getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
        sortable: this.getSortable(column, ALARM_UNSORTABLE_FIELDS),
        filter: this.$i18n.locale && this.columnPropertiesFiltersMap[column.value],
        getComponent: this.getComponentGetter(column),
        colorIndicatorEnabled: Object.values(COLOR_INDICATOR_TYPES).includes(column.colorIndicator),
      }));
    },
  },
  mounted() {
    this.fetchColumnFilters();
  },
  methods: {
    getComponentGetter(column) {
      const getCell = this.componentGettersMap[column.value];

      if (getCell) {
        return getCell;
      }

      if (column.value.startsWith('links.')) {
        return context => ({
          bind: {
            links: get(context.alarm, column.value, []),

            is: 'alarm-column-value-links',
          },
        });
      }

      const prepareFunc = column.filter ?? String;

      return context => ({
        bind: {
          is: 'c-ellipsis',
          text: prepareFunc(get(context.alarm, column.value, '')),
        },
      });
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
