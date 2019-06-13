<template lang="pug">
  div.stats-wrapper
    progress-overlay(:pending="pending")
    stats-curves-chart(:labels="labels", :datasets="datasets", :options="options")
</template>

<script>
import { get } from 'lodash';

import { STATS_DEFAULT_COLOR, STATS_TYPES, STATS_CURVES_POINTS_STYLES } from '@/constants';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import widgetStatsChartWrapperMixin from '@/mixins/widget/stats/stats-chart-wrapper';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';

import StatsCurvesChart from './stats-curves-chart.vue';

export default {
  components: {
    ProgressOverlay,
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
            cubicInterpolationMode: 'monotone',
            pointStyle: get(this.widget.parameters, `statsPointsStyles.${stat}`, STATS_CURVES_POINTS_STYLES.rect),
          });

          return acc;
        }, []);
      }

      return [];
    },

    options() {
      return {
        annotation: this.annotationLine,
        tooltips: {
          callbacks: {
            label: this.tooltipLabel,
          },
        },
      };
    },
  },
  methods: {
    tooltipLabel(tooltipItem, data) {
      const PROPERTIES_FILTERS_MAP = {
        [STATS_TYPES.stateRate.value]: value => this.$options.filters.percentage(value),
        [STATS_TYPES.ackTimeSla.value]: value => this.$options.filters.percentage(value),
        [STATS_TYPES.resolveTimeSla.value]: value => this.$options.filters.percentage(value),
      };

      const { stats } = this.query;

      const statObject = stats ? stats[data.datasets[tooltipItem.datasetIndex].label] : null;
      let label = data.datasets[tooltipItem.datasetIndex].label || '';

      if (label) {
        label += ': ';
      }

      if (statObject && PROPERTIES_FILTERS_MAP[statObject.stat]) {
        label += PROPERTIES_FILTERS_MAP[statObject.stat](tooltipItem.yLabel);
      } else {
        label += tooltipItem.yLabel;
      }

      return label;
    },
  },
};
</script>

<style lang="scss" scoped>
  .stats-wrapper {
    position: relative;
  }
</style>
