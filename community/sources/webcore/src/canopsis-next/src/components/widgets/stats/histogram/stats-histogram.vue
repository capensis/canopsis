<template lang="pug">
  div.position-relative
    c-progress-overlay(:pending="pending")
    c-alert-overlay(
      :value="hasError",
      :message="serverErrorMessage"
    )
    stats-histogram-chart(:labels="labels", :datasets="datasets", :options="options")
</template>

<script>
import { get } from 'lodash';

import { STATS_DEFAULT_COLOR } from '@/constants';

import entitiesStatsMixin from '@/mixins/entities/stats';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import widgetStatsWrapperMixin from '@/mixins/widget/stats/stats-wrapper';
import widgetStatsChartWrapperMixin from '@/mixins/widget/stats/stats-chart-wrapper';

const StatsHistogramChart = () => import(/* webpackChunkName: "Charts" */ './stats-histogram-chart.vue');

export default {
  components: {
    StatsHistogramChart,
  },
  mixins: [
    entitiesStatsMixin,
    widgetFetchQueryMixin,
    widgetStatsWrapperMixin,
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
