<template lang="pug">
  div
    side-bar-wrapper(:navigationDrawerProps="navigationDrawerProps[sideBarName] || defaultNavigationDrawerProps")
      component(:is="sideBarName", :config="sideBarConfig")
</template>

<script>
import sideBarInnerMixin from '@/mixins/side-bar/side-bar-inner';

import SideBarWrapper from './side-bar-wrapper.vue';
import AlarmSettings from './settings/alarm.vue';
import ContextSettings from './settings/context.vue';
import WeatherSettings from './settings/service-weather.vue';
import StatsHistogramSettings from './settings/stats-histogram.vue';
import StatsCurvesSettings from './settings/stats-curves.vue';
import StatsTableSettings from './settings/stats-table.vue';
import StatsCalendarSettings from './settings/stats-calendar.vue';
import StatsNumberSettings from './settings/stats-number.vue';
import StatsParetoSettings from './settings/stats-pareto.vue';
import TextSettings from './settings/text.vue';
import CounterSettings from './settings/counter.vue';

/**
 * Wrapper for all side-bars
 */
export default {
  components: {
    SideBarWrapper,
    AlarmSettings,
    ContextSettings,
    WeatherSettings,
    StatsHistogramSettings,
    StatsCurvesSettings,
    StatsTableSettings,
    StatsCalendarSettings,
    StatsNumberSettings,
    StatsParetoSettings,
    TextSettings,
    CounterSettings,
  },
  mixins: [sideBarInnerMixin],
  data() {
    return {
      navigationDrawerProps: {},
      defaultNavigationDrawerProps: {
        right: true,
        fixed: true,
        temporary: true,
        width: 400,
      },
    };
  },
};
</script>
