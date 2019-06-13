<template lang="pug">
  div.stats-wrapper
    progress-overlay(:pending="pending")
    stats-histogram-chart(:labels="labels", :datasets="datasets", :options="options")
</template>

<script>
import { get } from 'lodash';

import { STATS_DEFAULT_COLOR } from '@/constants';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import widgetStatsChartWrapperMixin from '@/mixins/widget/stats/stats-chart-wrapper';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';

import StatsHistogramChart from './stats-histogram-chart.vue';

export default {
  components: {
    ProgressOverlay,
    StatsHistogramChart,
  },
  mixins: [
    entitiesStatsMixin,
    widgetQueryMixin,
    entitiesUserPreferenceMixin,
    widgetStatsChartWrapperMixin,
  ],
  computed: {
    datasets() {
      if (this.stats) {
        return Object.keys(this.stats).reduce((acc, stat) => {
          acc.push({
            label: stat,
            data: this.stats[stat].sum.map(value => value.value),
            backgroundColor: get(this.widget.parameters, `statsColors.${stat}`, STATS_DEFAULT_COLOR),
          });

          return acc;
        }, []);
      }

      return [];
    },
  },
};
</script>

<style lang="scss" scoped>
  .stats-wrapper {
    position: relative;
  }
</style>
