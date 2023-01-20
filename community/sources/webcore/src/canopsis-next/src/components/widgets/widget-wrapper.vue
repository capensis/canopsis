<template lang="pug">
  div(:style="widgetWrapperStyles")
    template(v-if="widget.title || editing")
      v-card-title.widget-title.white.pa-2
        v-layout(justify-space-between, align-center)
          v-flex
            h4.ml-2.font-weight-regular {{ widget.title }}
          v-spacer
      v-divider
    v-card-text.pa-0.position-relative
      component(
        v-bind="widgetProps",
        :widget="preparedWidget",
        :tab-id="tab._id",
        :editing="editing"
      )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import {
  WIDGET_TYPES,
  WIDGET_TYPES_RULES,
  WIDGET_GRID_ROW_HEIGHT,
  ALARM_UNSORTABLE_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ENTITY_FIELDS_TO_LABELS_KEYS, DEFAULT_SERVICE_DEPENDENCIES_COLUMNS, DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS,
} from '@/constants';

import { setSeveralFields } from '@/helpers/immutable';

import { widgetColumnsMixin } from '@/mixins/widget/columns/common';

import AlarmsListWidget from './alarm/alarms-list.vue';
import EntitiesListWidget from './context/entities-list.vue';
import ServiceWeatherWidget from './service-weather/service-weather.vue';
import TestingWeatherWidget from './testing-weather/testing-weather.vue';
import StatsCalendarWidget from './stats/calendar/stats-calendar.vue';
import TextWidget from './text/text.vue';
import CounterWidget from './counter/counter.vue';
import MapWidget from './map/map.vue';

const { mapGetters } = createNamespacedHelpers('info');

export default {
  components: {
    AlarmsListWidget,
    EntitiesListWidget,
    ServiceWeatherWidget,
    TestingWeatherWidget,
    StatsCalendarWidget,
    TextWidget,
    CounterWidget,
    MapWidget,
  },
  mixins: [widgetColumnsMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    tab: {
      type: Object,
      required: true,
    },
    editing: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    ...mapGetters(['edition']),

    /**
     * Height is multiplying on 2 for loader correct displaying
     *
     * @returns {{minHeight: string}}
     */
    widgetWrapperStyles() {
      return { minHeight: `${WIDGET_GRID_ROW_HEIGHT * 2}px` };
    },

    preparedWidget() {
      switch (this.widget.type) {
        case WIDGET_TYPES.alarmList:
          return this.prepareAlarmWidget(this.widget);

        case WIDGET_TYPES.context: // TODO: finish it
          return {
            ...this.widget,

            parameters: {
              ...this.widget.parameters,

              widgetColumns: [],
              widgetGroupColumns: [],
              serviceDependenciesColumns: [],
              widgetExportColumns: [],
            },
          };

        case WIDGET_TYPES.serviceWeather: // TODO: finish it
          return {
            ...this.widget,

            parameters: {
              ...this.widget.parameters,

              widgetColumns: [],
              widgetGroupColumns: [],
              serviceDependenciesColumns: [],
              widgetExportColumns: [],
            },
          };

        case WIDGET_TYPES.statsCalendar: // TODO: finish it
          return {
            ...this.widget,

            parameters: {
              ...this.widget.parameters,

              widgetColumns: [],
              widgetGroupColumns: [],
              serviceDependenciesColumns: [],
              widgetExportColumns: [],
            },
          };

        default:
          return this.widget;
      }
    },

    widgetProps() {
      const { type } = this.widget;
      const widgetComponentsMap = {
        [WIDGET_TYPES.alarmList]: 'alarms-list-widget',
        [WIDGET_TYPES.context]: 'entities-list-widget',
        [WIDGET_TYPES.serviceWeather]: 'service-weather-widget',
        [WIDGET_TYPES.statsCalendar]: 'stats-calendar-widget',
        [WIDGET_TYPES.text]: 'text-widget',
        [WIDGET_TYPES.counter]: 'counter-widget',
        [WIDGET_TYPES.testingWeather]: 'testing-weather-widget',
        [WIDGET_TYPES.map]: 'map-widget',
      };
      let widgetSpecificsProp = {};

      Object.entries(WIDGET_TYPES_RULES).forEach(([key, rule]) => {
        if (rule.edition !== this.edition) {
          widgetComponentsMap[key] = 'c-alert-overlay';
          widgetSpecificsProp = {
            message: this.$t('errors.statsWrongEditionError'),
            value: true,
          };
        }
      });

      const component = widgetComponentsMap[type];

      if (!component) {
        return {
          is: 'c-alert-overlay',
          message: this.$t('errors.unknownWidgetType', { type }),
          value: true,
        };
      }

      return {
        ...widgetSpecificsProp,

        is: component,
      };
    },
  },
  methods: {
    prepareAlarmWidget(widget) {
      return setSeveralFields(widget, {
        'parameters.widgetColumns': (columns = []) => (
          columns.map(column => ({
            ...column,

            sortable: this.getSortable(column, ALARM_UNSORTABLE_FIELDS),
            text: this.getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
          }))
        ),

        'parameters.widgetGroupColumns': (columns = DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS) => (
          columns.map(column => ({
            ...column,

            sortable: this.getSortable(column, ALARM_UNSORTABLE_FIELDS),
            text: this.getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
          }))
        ),

        'parameters.serviceDependenciesColumns': (columns = DEFAULT_SERVICE_DEPENDENCIES_COLUMNS) => (
          columns.map(column => ({
            ...column,

            sortable: false,
            text: this.getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
            value: column.value.startsWith('entity.') ? column.value : `entity.${column.value}`,
          }))
        ),

        'parameters.widgetExportColumns': (columns = []) => (
          columns.map(column => ({
            ...column,

            text: this.getColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
          }))
        ),
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.widget-title {
  height: 37px;
}
</style>
