<template lang="pug">
  div
    side-bar-wrapper(:navigationDrawerProps="navigationDrawerProps[sideBarName] || defaultNavigationDrawerProps")
      component(:is="sideBarName", :config="sideBarConfig")
</template>

<script>
import sideBarInnerMixin from '@/mixins/side-bar/side-bar-inner';

import SideBarWrapper from './side-bar-wrapper.vue';
import AlarmSettings from '../widgets/alarm.vue';
import ContextSettings from '../widgets/context.vue';
import WeatherSettings from '../widgets/weather.vue';
import StatsHistogramSettings from '../widgets/stats-histogram.vue';
import StatsCurvesSettings from '../widgets/stats-curves.vue';
import StatsTableSettings from '../widgets/stats-table.vue';
import StatsCalendarSettings from '../widgets/stats-calendar.vue';
import StatsNumberSettings from '../widgets/stats-number.vue';
import StatsParetoSettings from '../widgets/stats-pareto.vue';
import TextSettings from '../widgets/text.vue';

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
