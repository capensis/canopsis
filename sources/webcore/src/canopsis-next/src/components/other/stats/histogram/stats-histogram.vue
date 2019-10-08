<template lang="pug">
  div.position-relative
    progress-overlay(:pending="pending")
    stats-alert-overlay(
      :value="hasError",
      :message="errorMessage",
      :errorMessage="serverErrorMessage",
      :stackError="stackError"
    )
    stats-histogram-chart(:labels="labels", :datasets="datasets", :options="options")
</template>

<script>
import { get } from 'lodash';

import { STATS_DEFAULT_COLOR } from '@/constants';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import widgetStatsWrapperMixin from '@/mixins/widget/stats/stats-wrapper';
import widgetStatsChartWrapperMixin from '@/mixins/widget/stats/stats-chart-wrapper';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';

import StatsAlertOverlay from '../partials/stats-alert-overlay.vue';
import StatsHistogramChart from './stats-histogram-chart.vue';

export default {
  components: {
    ProgressOverlay,
    StatsAlertOverlay,
    StatsHistogramChart,
  },
  mixins: [
    entitiesStatsMixin,
    widgetQueryMixin,
    entitiesUserPreferenceMixin,
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
