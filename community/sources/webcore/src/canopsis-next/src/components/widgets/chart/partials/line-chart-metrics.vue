<template lang="pug">
  v-layout.chart-metrics-widget(column, align-center)
    h4.chart-metrics-widget__title {{ title }}
    line-chart.chart-metrics-widget__chart(
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
import { hasHistoryData } from '@/helpers/metrics';

import { chartMetricsOptionsMixin } from '@/mixins/chart/metrics-options';

const LineChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/line-chart.vue');

export default {
  inject: ['$system'],
  components: { LineChart },
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
  },
  computed: {
    hasHistoryData() {
      return hasHistoryData(this.metrics);
    },

    xAxes() {
      const xAxes = {
        [X_AXES_IDS.default]: {
          type: 'time',
          ticks: {
            source: 'data',
            callback: this.getChartTimeTickLabel,
            font: this.labelsFont,
          },
        },
      };

      if (this.hasHistoryData) {
        xAxes[X_AXES_IDS.history] = {
          type: 'time',
          ticks: {
            source: 'data',
            callback: this.getChartTimeTickLabel,
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
          borderColor: metricColor,
          xAxisID: X_AXES_IDS.default,
          yAxisID: this.getMetricYAxisId(metric),
          label: this.$t(`alarm.metrics.${metric}`),
          data: data.map(({ timestamp, value }) => ({
            x: timestamp * 1000,
            y: value,
          })),
        };

        acc.push(defaultDataset);

        if (this.hasHistoryData) {
          const historyMetricColor = colorToRgba(metricColor, 0.5);
          const historyDataset = {
            metric,
            backgroundColor: historyMetricColor,
            borderColor: historyMetricColor,
            xAxisID: X_AXES_IDS.history,
            yAxisID: this.getMetricYAxisId(metric),
            label: `${this.$t(`alarm.metrics.${metric}`)} (${this.$t('common.previous')})`,
            data: data.map(({ history_timestamp: historyTimestamp, history_value: historyValue }) => ({
              x: historyTimestamp * 1000,
              y: historyValue,
            })),
          };

          acc.push(historyDataset);
        }

        return acc;
      }, []);
    },
  },
};
</script>
