<template lang="pug">
  v-layout.chart-metrics-widget(column, align-center)
    h4.chart-metrics-widget__title {{ title }}
    bar-chart.chart-metrics-widget__chart(
      :options="chartOptions",
      :datasets="datasets",
      :width="width",
      :height="height",
      :dark="$system.dark"
    )
</template>

<script>
import { isUndefined } from 'lodash';

import { BAR_X_AXES_IDS, SAMPLINGS } from '@/constants';

import { colorToRgba, getMetricColor } from '@/helpers/color';
import { getDateLabelBySampling } from '@/helpers/metrics';

import { chartBarMetricsOptionsMixin } from '@/mixins/chart/bar-metrics-options';

const BarChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/bar-chart.vue');

export default {
  inject: ['$system'],
  components: { BarChart },
  mixins: [chartBarMetricsOptionsMixin],
  props: {
    metrics: {
      type: Array,
      default: () => [],
    },
    title: {
      type: String,
      required: false,
    },
    width: {
      default: 600,
      type: Number,
    },
    height: {
      default: 440,
      type: Number,
    },
    sampling: {
      type: String,
      default: SAMPLINGS.day,
    },
    stacked: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    xAxes() {
      const [metric] = this.metrics;

      const labels = metric?.data.map(({ history_timestamp: timestamp }) => (
        getDateLabelBySampling(timestamp, this.sampling).split('\n')
      )) ?? [];

      return {
        [BAR_X_AXES_IDS.default]: {
          stacked: this.stacked,
          type: 'time',
          ticks: {
            source: 'data',
            callback: this.getChartTimeTickLabel,
            font: this.labelsFont,
          },
        },
        [BAR_X_AXES_IDS.history]: {
          type: 'category',
          labels,
          ticks: {
            font: this.labelsFont,
          },
        },
      };
    },

    chartOptions() {
      return {
        animation: false,
        responsive: true,
        interaction: {
          mode: 'x',
          intersect: false,
        },
        scales: {
          ...this.xAxes,
          ...this.yAxes,
        },
        plugins: {
          legend: {
            position: 'right',
            align: 'start',
            maxWidth: 600,
            labels: {
              font: this.labelsFont,
              boxWidth: 15,
              boxHeight: 15,
            },
          },
          tooltip: {
            callbacks: {
              title: this.getChartTooltipTitle,
            },
          },
        },
      };
    },

    datasets() {
      return this.metrics.reduce((acc, { title: metric, data, color }) => {
        const metricColor = color ?? getMetricColor(metric);
        const defaultDataset = {
          metric,
          backgroundColor: metricColor,
          yAxisID: this.getMetricYAxisId(metric),
          label: this.$t(`alarm.metrics.${metric}`),
          data: data.map(({ timestamp, value }) => ({
            x: timestamp * 1000,
            y: value,
          })),
        };

        if (this.stacked) {
          defaultDataset.stack = 'default';
        }

        acc.push(defaultDataset);

        const hasHistoryData = data.some(({ history_timestamp: historyTimestamp }) => !isUndefined(historyTimestamp));

        if (hasHistoryData) {
          const historyDataset = {
            metric,
            backgroundColor: colorToRgba(metricColor, 0.5),
            yAxisID: this.getMetricYAxisId(metric),
            label: `${this.$t(`alarm.metrics.${metric}`)} (${this.$t('common.previous')})`,
            data: data.map(({ timestamp, history_timestamp: historyTimestamp, history_value: historyValue }) => ({
              x: timestamp * 1000,
              y: historyValue,
              originalX: historyTimestamp * 1000,
            })),
          };

          if (this.stacked) {
            historyDataset.stack = 'history';
          }

          acc.push(historyDataset);
        }

        return acc;
      }, []);
    },
  },
};
</script>

<style lang="scss" scoped>
.chart-metrics-widget__chart {
  width: 100%;
}
</style>
