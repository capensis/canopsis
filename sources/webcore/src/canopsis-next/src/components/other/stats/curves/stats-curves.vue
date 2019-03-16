<template lang="pug">
  div.stats-wrapper
    stats-curves-chart(v-if="!pending", :labels="labels", :datasets="datasets")
    v-layout(v-else, justify-center)
      v-progress-circular(
      indeterminate,
      color="primary",
      )
</template>

<script>
import { get } from 'lodash';

import { STATS_DEFAULT_COLOR } from '@/constants';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import widgetStatsChartWrapperMixin from '@/mixins/widget/stats/stats-chart-wrapper';

import StatsCurvesChart from './stats-curves-chart.vue';

export default {
  components: {
    StatsCurvesChart,
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
            borderColor: get(this.widget.parameters, `statsColors.${stat}`, STATS_DEFAULT_COLOR),
            backgroundColor: 'transparent',
          });

          return acc;
        }, []);
      }

      return [];
    },
  },
};
</script>
