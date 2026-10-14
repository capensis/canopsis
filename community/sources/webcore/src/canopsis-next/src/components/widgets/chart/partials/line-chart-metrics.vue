<template>
  <c-layout class="kpi-widget" align="center" column>
    <h4 class="kpi-widget__title">
      {{ title }}
    </h4>
    <line-chart
      :options="chartOptions"
      :datasets="datasets"
      :dark="$system.dark"
      class="kpi-widget__chart"
    >
      <template v-if="!hideActions" #actions="{ chart }">
        <chart-export-actions
          :downloading="downloading"
          :chart="chart"
          class="mt-4"
          @export:csv="$emit('export:csv')"
          @export:png="$emit('export:png', $event)"
        />
      </template>
    </line-chart>
  </c-layout>
</template>

<script>
import { KPI_CHART_DEFAULT_HEIGHT, SAMPLINGS, X_AXES_IDS } from '@/constants';

import { colorToRgba } from '@/helpers/color';
import { getMetricColor } from '@/helpers/entities/metric/color';
import { convertMetricValueByUnit } from '@/helpers/entities/metric/list';

import { chartMetricsOptionsMixin } from '@/mixins/chart/metrics-options';

import ChartExportActions from '@/components/common/chart/chart-export-actions.vue';
import LineChart from '@/components/common/chart/line-chart.vue';

export default {
  inject: ['$system'],

  components: {
    ChartExportActions,
    LineChart,
  },

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
    hideActions: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['export:csv', 'export:png'],

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
              label: this.getChartTooltipLabel,
            },
          },
        },
      };
    },

    datasets() {
      return this.preparedMetrics.reduce((acc, { title: metric, label, unit, data, color }) => {
        const metricColor = color ?? getMetricColor(metric);
        const datasetLabel = label ?? this.getMetricLabel(metric);
        const yAxisID = this.getMetricYAxisId(metric, unit);
        const defaultDataset = {
          metric,
          backgroundColor: metricColor,
          borderColor: metricColor,
          xAxisID: X_AXES_IDS.default,
          yAxisID,
          label: datasetLabel,
          unit,
          data: data.map(({ timestamp, value }) => ({
            x: timestamp * 1000,
            y: convertMetricValueByUnit(value, unit),
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
            yAxisID,
            label: `${datasetLabel} (${this.$t('common.previous')})`,
            unit,
            data: data.map(({ history_timestamp: historyTimestamp, history_value: historyValue }) => ({
              x: historyTimestamp * 1000,
              y: convertMetricValueByUnit(historyValue, unit),
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
