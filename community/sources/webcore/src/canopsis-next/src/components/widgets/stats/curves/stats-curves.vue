<template lang="pug">
  div.position-relative
    c-progress-overlay(:pending="pending")
    c-alert-overlay(
      :value="hasError",
      :message="serverErrorMessage"
    )
    stats-curves-chart(:labels="labels", :datasets="datasets", :options="options")
</template>

<script>
import { get } from 'lodash';

import { STATS_DEFAULT_COLOR, STATS_CURVES_POINTS_STYLES } from '@/constants';

import entitiesStatsMixin from '@/mixins/entities/stats';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import widgetStatsWrapperMixin from '@/mixins/widget/stats/stats-wrapper';
import widgetStatsChartWrapperMixin from '@/mixins/widget/stats/stats-chart-wrapper';

const StatsCurvesChart = () => import(/* webpackChunkName: "Charts" */ './stats-curves-chart.vue');

export default {
  components: {
    StatsCurvesChart,
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
            borderColor: get(this.widget.parameters, `statsColors.${stat}`, STATS_DEFAULT_COLOR),
            backgroundColor: 'transparent',
            cubicInterpolationMode: 'monotone',
            pointStyle: get(this.widget.parameters, `statsPointsStyles.${stat}`, STATS_CURVES_POINTS_STYLES.rect),
          });

          return acc;
        }, []);
      }

      return [];
    },
  },
};
</script>
