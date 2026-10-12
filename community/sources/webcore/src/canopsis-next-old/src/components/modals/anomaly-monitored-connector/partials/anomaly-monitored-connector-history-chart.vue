<template>
  <div class="anomaly-monitored-connector-history-chart">
    <line-chart
      :chart-id="chartId"
      :options="chartOptions"
      :datasets="datasets"
      :dark="$system.dark"
    />
  </div>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

const LineChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/line-chart.vue');

const ANOMALY_COLOR = '#E53935';
const MAIN_COLOR = '#1E88E5';
const AVERAGE_COLOR = '#9E9E9E';

export default {
  inject: ['$system'],
  components: { LineChart },
  props: {
    history: {
      type: Object,
      default: null,
    },
  },
  setup(props) {
    const { t } = useI18n();

    const chartId = 'anomaly-connector-history-chart';

    const datasets = computed(() => {
      const { data = [], average_events: average = 0 } = props.history ?? {};

      if (!data?.length) {
        return [];
      }

      const points = data.map(point => ({
        x: point.time * 1000,
        y: point.count,
      }));

      const pointBackgroundColors = data.map(p => (p.anomaly ? ANOMALY_COLOR : MAIN_COLOR));
      const pointBorderColors = pointBackgroundColors;
      const pointStyles = data.map(p => (p.anomaly ? 'rect' : 'circle'));

      const averagePoints = data.map(p => ({
        x: p.time * 1000,
        y: average,
      }));

      return [
        {
          label: t('healthcheck.connectorsBlocks.eventsSeries'),
          data: points,
          borderColor: MAIN_COLOR,
          backgroundColor: 'transparent',
          tension: 0.2,
          pointRadius: 5,
          pointHoverRadius: 7,
          pointBackgroundColor: pointBackgroundColors,
          pointBorderColor: pointBorderColors,
          pointStyle: pointStyles,
          borderWidth: 2,
        },
        {
          label: t('healthcheck.connectorsBlocks.averageEvents'),
          data: averagePoints,
          borderColor: AVERAGE_COLOR,
          backgroundColor: 'transparent',
          pointRadius: 0,
          borderWidth: 2,
          tooltip: {
            enabled: false,
          },
        },
      ];
    });

    const chartOptions = computed(() => ({
      responsive: true,
      maintainAspectRatio: false,

      scales: {
        x: {
          type: 'time',
          time: {
            unit: 'minute',
            displayFormats: {
              minute: 'HH:mm',
            },
          },
          ticks: {
            maxRotation: 45,
            minRotation: 0,
          },
        },
        y: {
          beginAtZero: true,
          title: {
            display: true,
            text: t('healthcheck.connectorsBlocks.chartYAxis'),
          },
        },
      },
      plugins: {
        legend: {
          display: false,
        },
        tooltip: {
          displayColors: false,
          caretSize: 0,
          caretPadding: 10,
          yAlign: 'bottom',
          backgroundColor: 'rgba(0, 0, 0, 0.6)',
          callbacks: {
            title: () => '',
            label: (ctx) => {
              const value = ctx.parsed?.y;

              return t('healthcheck.connectorsBlocks.tooltipEvents', { count: value });
            },
          },
        },
      },
    }));

    return {
      chartId,
      chartOptions,
      datasets,
    };
  },
};
</script>

<style lang="scss" scoped>
.anomaly-monitored-connector-history-chart {
  position: relative;
  width: 100%;
  height: 340px;

  ::v-deep {
    canvas, div {
      width: 100% !important;
      height: 100% !important;
      position: relative !important;
    }
  }
}
</style>
