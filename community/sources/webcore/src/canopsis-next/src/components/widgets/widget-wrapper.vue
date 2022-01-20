<template lang="pug">
  div(:style="widgetWrapperStyles")
    template(v-if="widget.title || isEditingMode")
      v-card-title.widget-title.white.pa-2
        v-layout(justify-space-between, align-center)
          v-flex
            h4.ml-2.font-weight-regular {{ widget.title }}
          v-spacer
      v-divider
    v-card-text.pa-0.position-relative
      component(
        v-bind="getWidgetPropsByType(widget.type)",
        :widget="widget",
        :tab-id="tab._id",
        :is-editing-mode="isEditingMode"
      )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { WIDGET_TYPES, WIDGET_TYPES_RULES, WIDGET_GRID_ROW_HEIGHT } from '@/constants';

import AlarmsListWidget from './alarm/alarms-list.vue';
import EntitiesListWidget from './context/entities-list.vue';
import ServiceWeatherWidget from './service-weather/service-weather.vue';
import TestingWeatherWidget from './testing-weather/testing-weather.vue';
import StatsHistogramWidget from './stats/histogram/stats-histogram.vue';
import StatsCurvesWidget from './stats/curves/stats-curves.vue';
import StatsTableWidget from './stats/stats-table.vue';
import StatsCalendarWidget from './stats/calendar/stats-calendar.vue';
import StatsNumberWidget from './stats/stats-number.vue';
import StatsParetoWidget from './stats/pareto/stats-pareto.vue';
import TextWidget from './text/text.vue';
import CounterWidget from './counter/counter.vue';

const { mapGetters } = createNamespacedHelpers('info');

export default {
  components: {
    AlarmsListWidget,
    EntitiesListWidget,
    ServiceWeatherWidget,
    TestingWeatherWidget,
    StatsHistogramWidget,
    StatsCurvesWidget,
    StatsTableWidget,
    StatsCalendarWidget,
    StatsNumberWidget,
    StatsParetoWidget,
    TextWidget,
    CounterWidget,
  },
  props: {
    widget: {
      type: Object,
      required: true,
    },
    tab: {
      type: Object,
      required: true,
    },
    isEditingMode: {
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
  },
  methods: {
    getWidgetPropsByType(type) {
      const widgetComponentsMap = {
        [WIDGET_TYPES.alarmList]: 'alarms-list-widget',
        [WIDGET_TYPES.context]: 'entities-list-widget',
        [WIDGET_TYPES.serviceWeather]: 'service-weather-widget',
        [WIDGET_TYPES.statsHistogram]: 'stats-histogram-widget',
        [WIDGET_TYPES.statsCurves]: 'stats-curves-widget',
        [WIDGET_TYPES.statsTable]: 'stats-table-widget',
        [WIDGET_TYPES.statsCalendar]: 'stats-calendar-widget',
        [WIDGET_TYPES.statsNumber]: 'stats-number-widget',
        [WIDGET_TYPES.statsPareto]: 'stats-pareto-widget',
        [WIDGET_TYPES.text]: 'text-widget',
        [WIDGET_TYPES.counter]: 'counter-widget',
        [WIDGET_TYPES.testingWeather]: 'testing-weather-widget',
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

      return {
        ...widgetSpecificsProp,

        is: widgetComponentsMap[type],
      };
    },
  },
};
</script>

<style lang="scss" scoped>
.widget-title {
  height: 37px;
}

.copy-widget-id {
  z-index: 2;
  position: relative;
}
</style>
