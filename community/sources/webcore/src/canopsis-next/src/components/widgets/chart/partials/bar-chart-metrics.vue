<template>
  <v-layout
    class="kpi-widget"
    column
    align-center
  >
    <h4 class="kpi-widget__title">
      {{ title }}
    </h4>
    <bar-chart
      :chart-id="chartId"
      :options="chartOptions"
      :datasets="datasets"
      :width="width"
      :height="height"
      :dark="$system.dark"
      class="kpi-widget__chart"
    >
      <template #actions="{ chart }">
        <kpi-chart-export-actions
          :downloading="downloading"
          :chart="chart"
          class="mt-4"
          v-on="$listeners"
        />
      </template>
    </bar-chart>
  </v-layout>
</template>

<script>
import { X_AXES_IDS, SAMPLINGS, KPI_CHART_DEFAULT_HEIGHT } from '@/constants';

import { colorToRgba } from '@/helpers/color';
import { getMetricColor } from '@/helpers/entities/metric/color';
import { convertMetricValueByUnit, getDateLabelBySampling } from '@/helpers/entities/metric/list';

import { chartMetricsOptionsMixin } from '@/mixins/chart/metrics-options';

import KpiChartExportActions from '@/components/other/kpi/charts/partials/kpi-chart-export-actions.vue';

const BarChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/bar-chart.vue');

export default {
  inject: ['$system'],
  components: { KpiChartExportActions, BarChart },
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
      default: KPI_CHART_DEFAULT_HEIGHT,
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
            callbacks: {
              title: this.getChartTooltipTitle,
              label: this.getChartTooltipLabel,
            },
          },
        },
      };
    },

    datasets() {
      return this.metrics.reduce((acc, { title: metric, label, unit, data, color }) => {
        const metricColor = color ?? getMetricColor(metric);
        const datasetLabel = label ?? this.getMetricLabel(metric);
        const yAxisID = this.getMetricYAxisId(metric, unit);

        const defaultDataset = {
          metric,
          yAxisID,
          backgroundColor: metricColor,
          label: datasetLabel,
          unit,
          data: data.map(({ timestamp, value }) => ({
            x: timestamp * 1000,
            y: convertMetricValueByUnit(value, unit),
          })),
        };

        if (this.stacked) {
          defaultDataset.stack = 'default';
        }

        acc.push(defaultDataset);

        if (this.hasHistoryData) {
          const historyDataset = {
            metric,
            yAxisID,
            backgroundColor: colorToRgba(metricColor, 0.5),
            label: `${datasetLabel} (${this.$t('common.previous')})`,
            unit,
            data: data.map(({ timestamp, history_timestamp: historyTimestamp, history_value: historyValue }) => ({
              x: timestamp * 1000,
              y: convertMetricValueByUnit(historyValue, unit),
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
  methods: {
    getMetricLabel(metric) {
      const metricMessageKey = `alarm.metrics.${metric}`;

      return this.$te(metricMessageKey) ? this.$t(metricMessageKey) : metric;
    },
  },
};
</script>
