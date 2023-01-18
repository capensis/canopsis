import { get } from 'lodash';

import { ALARM_ENTITY_FIELDS, COLOR_INDICATOR_TYPES } from '@/constants';

import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';
import { convertDurationToString } from '@/helpers/date/duration';

export const widgetColumnsFiltersMixin = { // TODO: rename for alarm
  computed: {
    columnsFiltersMap() {
      return this.columnsFilters.reduce((acc, { column, filter, attributes = [] }) => {
        acc[column] = this.getFilter(filter, attributes);

        return acc;
      }, {});
    },

    componentGettersMap() {
      return {
        [ALARM_ENTITY_FIELDS.state]: context => ({
          bind: {
            is: 'alarm-column-value-state',
            alarm: context.alarm,
          },
        }),
        [ALARM_ENTITY_FIELDS.status]: context => ({
          bind: {
            is: 'alarm-column-value-status',
            alarm: context.alarm,
          },
        }),
        [ALARM_ENTITY_FIELDS.priority]: context => ({
          bind: {
            is: 'color-indicator-wrapper',
            type: COLOR_INDICATOR_TYPES.impactState,
            entity: context.alarm.entity,
            alarm: context.alarm,
          },
        }),
        [ALARM_ENTITY_FIELDS.impactState]: context => ({
          bind: {
            is: 'color-indicator-wrapper',
            type: COLOR_INDICATOR_TYPES.impactState,
            entity: context.alarm.entity,
            alarm: context.alarm,
          },
        }),
        links: context => ({
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
        [ALARM_ENTITY_FIELDS.extraDetails]: context => ({
          bind: {
            is: 'alarm-column-value-extra-details',
            alarm: context.alarm,
          },
        }),
        [ALARM_ENTITY_FIELDS.tags]: context => ({
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
        'v.last_update_date': convertDateToStringWithFormatForToday,
        'v.creation_date': convertDateToStringWithFormatForToday,
        'v.last_event_date': convertDateToStringWithFormatForToday,
        'v.activation_date': convertDateToStringWithFormatForToday,
        'v.ack.t': convertDateToStringWithFormatForToday,
        'v.state.t': convertDateToStringWithFormatForToday,
        'v.status.t': convertDateToStringWithFormatForToday,
        'v.resolved': convertDateToStringWithFormatForToday,
        'v.duration': convertDurationToString,
        'v.current_state_duration': convertDurationToString,
        t: convertDateToStringWithFormatForToday,
        'v.active_duration': convertDurationToString,
        'v.snooze_duration': convertDurationToString,
        'v.pbh_inactive_duration': convertDurationToString,

        ...this.columnsFiltersMap,
      };
    },

    preparedColumns() {
      return this.columns.map(column => ({
        ...column,

        filter: this.$i18n.locale && this.columnPropertiesFiltersMap[column.value],
        getComponent: this.getComponentGetter(column),
        colorIndicatorEnabled: Object.values(COLOR_INDICATOR_TYPES).includes(column.colorIndicator),
      }));
    },
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
  },
};
