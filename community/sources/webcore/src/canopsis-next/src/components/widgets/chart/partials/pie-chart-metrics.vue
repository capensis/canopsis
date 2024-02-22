<template>
  <v-layout
    class="kpi-widget"
    column
    align-center
  >
    <h4 class="kpi-widget__title">
      {{ title }}
    </h4>
    <pie-chart
      :chart-id="chartId"
      :datasets="datasets"
      :labels="labels"
      :options="chartOptions"
      :width="700"
      :height="400"
      class="pie-chart-metrics__chart kpi-widget__chart"
    >
      <template #actions="{ chart }">
        <kpi-chart-export-actions
          :downloading="downloading"
          :chart="chart"
          class="mt-4"
          v-on="$listeners"
        />
      </template>
    </pie-chart>
  </v-layout>
</template>

<script>
import { COLORS } from '@/config';
import { KPI_PIE_CHART_SHOW_MODS } from '@/constants';

import { getMostReadableTextColor } from '@/helpers/color';
import { getMetricColor } from '@/helpers/entities/metric/color';
import { convertNumberToRoundedPercentString } from '@/helpers/string';

import { chartMetricsOptionsMixin } from '@/mixins/chart/metrics-options';

import KpiChartExportActions from '@/components/other/kpi/charts/partials/kpi-chart-export-actions.vue';

const PieChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/pie-chart.vue');

export default {
  inject: ['$system'],
  components: { KpiChartExportActions, PieChart },
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
    showMode: {
      type: String,
      default: KPI_PIE_CHART_SHOW_MODS.numbers,
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
    totalDataValue() {
      return this.metrics.reduce((acc, metric) => acc + metric.value, 0);
    },

    isPercentMode() {
      return this.showMode === KPI_PIE_CHART_SHOW_MODS.percent;
    },

    datasets() {
      const { data, backgroundColor } = this.preparedMetrics.reduce((acc, metric) => {
        acc.data.push(metric.value);
        acc.backgroundColor.push(metric.color ?? getMetricColor(metric.title));

        return acc;
      }, {
        backgroundColor: [],
        data: [],
      });

      return [{
        label: this.title,
        borderWidth: 0,
        backgroundColor,
        data,
      }];
    },

    labels() {
      return this.preparedMetrics.map(metric => metric.label ?? this.getMetricLabel(metric.title));
    },

    chartOptions() {
      return {
        responsive: this.responsive,
        animation: this.animation,
        plugins: {
          datalabels: {
            color: ({ dataIndex, dataset }) => {
              const dataColor = dataset.backgroundColor[dataIndex];

              return getMostReadableTextColor(dataColor, { level: 'AA', size: 'large' });
            },
            formatter: this.formatDataLabel,
            font: {
              size: 11,
              family: 'Arial, sans-serif',
            },
          },
          legend: this.legend,
          emptyPie: {
            width: 2,
            color: COLORS.error,
            text: this.$t('kpi.errors.emptyMetrics'),
            font: {
              size: 15,
              family: 'Arial, sans-serif',
              color: COLORS.error,
            },
          },
          tooltip: {
            callbacks: {
              title: () => this.title,
              label: this.formatTooltipLabel,
            },
          },
        },
      };
    },
  },
  methods: {
    getMetricLabel(metric) {
      const metricMessageKey = `alarm.metrics.${metric}`;

      return this.$te(metricMessageKey) ? this.$t(metricMessageKey) : metric;
    },

    formatDataLabel(value) {
      if (!value) {
        return null;
      }

      return this.isPercentMode
        ? convertNumberToRoundedPercentString(value / this.totalDataValue)
        : value;
    },

    formatTooltipLabel({ label, raw }) {
      const percent = convertNumberToRoundedPercentString(raw / this.totalDataValue);

      const value = this.isPercentMode
        ? `${percent} (${raw})`
        : `${raw} (${percent})`;

      return `${label}: ${value}`;
    },
  },
};
</script>

<style lang="scss">
.pie-chart-metrics {
  &__chart #pie-chart {
    height: 400px !important;
  }
}
</style>
