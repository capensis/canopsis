<template lang="pug">
  v-layout.chart-metrics-widget(column, align-center)
    h4.chart-metrics-widget__title {{ title }}
    line-chart.chart-metrics-widget__chart(
      :chart-id="chartId",
      :options="chartOptions",
      :datasets="datasets",
      :width="width",
      :height="height",
      :dark="$system.dark"
    )
      template(#actions="{ chart }")
        kpi-chart-export-actions.mt-4(:downloading="downloading", :chart="chart", v-on="$listeners")
</template>

<script>
import { X_AXES_IDS, SAMPLINGS } from '@/constants';

import { colorToRgba, getMetricColor } from '@/helpers/color';

import { chartMetricsOptionsMixin } from '@/mixins/chart/metrics-options';

import KpiChartExportActions from '@/components/other/kpi/charts/partials/kpi-chart-export-actions.vue';

const LineChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/line-chart.vue');

export default {
  inject: ['$system'],
  components: { KpiChartExportActions, LineChart },
  mixins: [chartMetricsOptionsMixin],
  props: {
    chartId: {
      type: String,
      required: false,
    },
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
      default: 560,
      type: Number,
    },
    sampling: {
      type: String,
      default: SAMPLINGS.day,
    },
    downloading: {
      type: Boolean,
      default: false,
    },
    responsive: {
      type: Boolean,
      default: true,
    },
    animation: {
      type: Boolean,
      default: true,
    },
  },
  computed: {
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
        responsive: this.responsive,
        animation: this.animation,
        interaction: {
          mode: 'x',
          intersect: false,
        },
        scales: {
          ...this.xAxes,
          ...this.yAxes,
        },
        plugins: {
          legend: this.legend,
          tooltip: {
            bodyFont: {
              size: this.tooltipBodyFontSize,
            },
            callbacks: {
              title: this.getChartTooltipTitle,
            },
          },
        },
      };
    },

    datasets() {
      return this.preparedMetrics.reduce((acc, { title: metric, label, data, color }) => {
        const metricColor = color ?? getMetricColor(metric);

        const datasetLabel = label ?? this.getMetricLabel(metric);

        const defaultDataset = {
          metric,
          backgroundColor: metricColor,
          borderColor: metricColor,
          xAxisID: X_AXES_IDS.default,
          yAxisID: this.getMetricYAxisId(metric),
          label: datasetLabel,
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
            label: `${datasetLabel} (${this.$t('common.previous')})`,
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
  methods: {
    getMetricLabel(metric) {
      const metricMessageKey = `alarm.metrics.${metric}`;

      return this.$te(metricMessageKey) ? this.$t(metricMessageKey) : metric;
    },
  },
};
</script>
