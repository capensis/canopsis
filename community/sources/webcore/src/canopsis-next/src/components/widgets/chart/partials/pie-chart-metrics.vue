<template lang="pug">
  v-layout.chart-metrics-widget(column, align-center)
    h4.chart-metrics-widget__title {{ title }}
    pie-chart.chart-metrics-widget__chart(
      :datasets="datasets",
      :labels="labels",
      :options="chartOptions",
      :width="700",
      :height="400"
    )
</template>

<script>
import { COLORS } from '@/config';

import { KPI_PIE_CHART_SHOW_MODS } from '@/constants';

import { getMetricColor, getMostReadableTextColor } from '@/helpers/color';
import { convertNumberToRoundedPercentString } from '@/helpers/string';

const PieChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/pie-chart.vue');

export default {
  inject: ['$system'],
  components: { PieChart },
  props: {
    metrics: {
      type: Array,
      default: () => [],
    },
    colorsByMetrics: {
      type: Object,
      default: () => ({}),
    },
    title: {
      type: String,
      required: false,
    },
    showMode: {
      type: String,
      default: KPI_PIE_CHART_SHOW_MODS.numbers,
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
      const { data, backgroundColor } = this.metrics.reduce((acc, metric) => {
        acc.data.push(metric.value);
        acc.backgroundColor.push(this.colorsByMetrics[metric.title] ?? getMetricColor(metric.title));

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
      return this.metrics.map(metric => this.$t(`alarm.metrics.${metric.title}`));
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
          legend: {
            position: 'right',
            maxWidth: 300,
            labels: {
              font: {
                size: 11,
                family: 'Arial, sans-serif',
              },
              boxWidth: 15,
              boxHeight: 15,
              padding: 8,
            },
          },
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
