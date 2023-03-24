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
import { X_AXES_IDS, SAMPLINGS } from '@/constants';

import { colorToRgba, getMetricColor } from '@/helpers/color';
import { getDateLabelBySampling, hasHistoryData } from '@/helpers/metrics';

import { chartMetricsOptionsMixin } from '@/mixins/chart/metrics-options';

const BarChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/bar-chart.vue');

export default {
  inject: ['$system'],
  components: { BarChart },
  mixins: [chartMetricsOptionsMixin],
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
    hasHistoryData() {
      return hasHistoryData(this.metrics);
    },

    xAxes() {
      const xAxes = {
        [X_AXES_IDS.default]: {
          stacked: this.stacked,
          type: 'time',
          ticks: {
            source: 'data',
            callback: this.getChartTimeTickLabel,
            font: this.labelsFont,
          },
        },
      };

      if (this.hasHistoryData) {
        const [metric] = this.metrics;

        const labels = metric?.data.map(({ history_timestamp: timestamp }) => (
          getDateLabelBySampling(timestamp, this.sampling).split('\n')
        )) ?? [];

        xAxes[X_AXES_IDS.history] = {
          type: 'category',
          labels,
          ticks: {
            font: this.labelsFont,
          },
        };
      }

      return xAxes;
    },

    chartOptions() {
      return {
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

        if (this.hasHistoryData) {
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
