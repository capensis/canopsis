<template lang="pug">
  div
    template(v-if="widget.title")
      v-card-title.lighten-1.pa-1
        v-layout(justify-space-between, align-center)
          v-flex
            h4.ml-2.font-weight-regular(:data-test="`widgetTitle-${widget._id}`") {{ widget.title }}
          v-spacer
      v-divider
    v-card-text.pa-0
      component(
        v-bind="widgetsComponentsMap(widget.type).bind",
        :widget="widget",
        :tabId="tab._id",
        :isEditingMode="isEditingMode"
      )
    copy-widget-id(v-show="isEditingMode", :widgetId="widget._id")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { WIDGET_TYPES, WIDGET_TYPES_RULES } from '@/constants';

import AlarmsListWidget from '@/components/other/alarm/alarms-list.vue';
import EntitiesListWidget from '@/components/other/context/entities-list.vue';
import WeatherWidget from '@/components/other/service-weather/weather.vue';
import StatsHistogramWidget from '@/components/other/stats/histogram/stats-histogram.vue';
import StatsCurvesWidget from '@/components/other/stats/curves/stats-curves.vue';
import StatsTableWidget from '@/components/other/stats/stats-table.vue';
import StatsCalendarWidget from '@/components/other/stats/calendar/stats-calendar.vue';
import StatsNumberWidget from '@/components/other/stats/stats-number.vue';
import StatsParetoWidget from '@/components/other/stats/pareto/stats-pareto.vue';
import TextWidget from '@/components/other/text/text.vue';
import CounterWidget from '@/components/other/counter/counter.vue';
import AlertOverlay from '@/components/layout/alert/alert-overlay.vue';
import CopyWidgetId from '@/components/side-bars/settings/widgets/fields/common/copy-widget-id.vue';

const { mapGetters } = createNamespacedHelpers('info');

export default {
  components: {
    CopyWidgetId,
    AlarmsListWidget,
    EntitiesListWidget,
    WeatherWidget,
    StatsHistogramWidget,
    StatsCurvesWidget,
    StatsTableWidget,
    StatsCalendarWidget,
    StatsNumberWidget,
    StatsParetoWidget,
    TextWidget,
    CounterWidget,
    AlertOverlay,
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
    updateTabMethod: {
      type: Function,
      required: true,
    },
  },
  computed: {
    ...mapGetters(['edition']),

    widgetsComponentsMap() {
      return (widgetType) => {
        const baseMap = {
          [WIDGET_TYPES.alarmList]: 'alarms-list-widget',
          [WIDGET_TYPES.context]: 'entities-list-widget',
          [WIDGET_TYPES.weather]: 'weather-widget',
          [WIDGET_TYPES.statsHistogram]: 'stats-histogram-widget',
          [WIDGET_TYPES.statsCurves]: 'stats-curves-widget',
          [WIDGET_TYPES.statsTable]: 'stats-table-widget',
          [WIDGET_TYPES.statsCalendar]: 'stats-calendar-widget',
          [WIDGET_TYPES.statsNumber]: 'stats-number-widget',
          [WIDGET_TYPES.statsPareto]: 'stats-pareto-widget',
          [WIDGET_TYPES.text]: 'text-widget',
          [WIDGET_TYPES.counter]: 'counter-widget',
        };

        let widgetSpecificsProp = {};

        Object.entries(WIDGET_TYPES_RULES).forEach(([key, rule]) => {
          if (rule.edition !== this.edition) {
            baseMap[key] = 'alert-overlay';
            widgetSpecificsProp = {
              message: this.$t('errors.statsWrongEditionError'),
              value: true,
            };
          }
        });

        return {
          bind: {
            ...widgetSpecificsProp,

            is: baseMap[widgetType],
          },
        };
      };
    },
  },
};
</script>
