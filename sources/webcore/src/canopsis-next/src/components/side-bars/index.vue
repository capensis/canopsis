<template lang="pug">
  div
    side-bar-wrapper(:navigationDrawerProps="navigationDrawerProps[sideBarName] || defaultNavigationDrawerProps")
      component(:is="sideBarName", :config="sideBarConfig")
</template>

<script>
import sideBarInnerMixin from '@/mixins/side-bar/side-bar-inner';

import SideBarWrapper from './side-bar-wrapper.vue';
import AlarmSettings from './settings/widgets/alarm.vue';
import ContextSettings from './settings/widgets/context.vue';
import WeatherSettings from './settings/widgets/weather.vue';
import StatsHistogramSettings from './settings/widgets/stats-histogram.vue';
import StatsTableSettings from './settings/widgets/stats-table.vue';

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
    StatsTableSettings,
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
