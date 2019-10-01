<template lang="pug">
  grid-layout(
    :layout="tab.layout",
    @layout-updated="$emit('input', $event)",
    :autoSize="true",
    :row-height="12"
  )
    grid-item(
      v-for="item in tab.layout",
      :x="item.x",
      :y="item.y",
      :h="item.h",
      :w="item.w",
      :i="item.i",
      :key="item.i"
    )
      component(
        :is="component(item.i).is",
        :widget="component(item.i).widget"
      )
</template>

<script>
import { GridLayout, GridItem } from 'vue-grid-layout';

import { WIDGET_TYPES } from '@/constants';

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

export default {
  components: {
    GridLayout,
    GridItem,
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
  },
  props: {
    tab: {
      type: Object,
      required: true,
    },
    hasUpdateAccess: {
      type: Boolean,
      default: false,
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
  data() {
    return {
      widgetsComponentsMap: {
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
      },
    };
  },
  computed: {
    // TODO: Improve this widget computation
    component() {
      return (id) => {
        const foundWidget = this.tab.widgets.find(widget => widget._id === id);

        return {
          is: this.widgetsComponentsMap[foundWidget.type],
          widget: foundWidget,
        };
      };
    },
    minHeight() {
      return 20;
    },
  },
};
</script>
