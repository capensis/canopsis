<template lang="pug">
  div
    side-bar-wrapper(:navigationDrawerProps="navigationDrawerProps[sideBarName] || defaultNavigationDrawerProps")
      component(:is="sideBarName", :config="sideBarConfig")
</template>

<script>
import sideBarInnerMixin from '@/mixins/side-bar/side-bar-inner';

import AlarmSettings from '@/components/widgets/alarm.vue';
import ContextSettings from '@/components/widgets/context.vue';
import WeatherSettings from '@/components/widgets/weather.vue';
import StatsHistogramSettings from '@/components/widgets/stats-histogram.vue';
import StatsCurvesSettings from '@/components/widgets/stats-curves.vue';
import StatsTableSettings from '@/components/widgets/stats-table.vue';
import StatsCalendarSettings from '@/components/widgets/stats-calendar.vue';
import StatsNumberSettings from '@/components/widgets/stats-number.vue';
import StatsParetoSettings from '@/components/widgets/stats-pareto.vue';
import TextSettings from '@/components/widgets/text.vue';

import SideBarWrapper from './side-bar-wrapper.vue';
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
