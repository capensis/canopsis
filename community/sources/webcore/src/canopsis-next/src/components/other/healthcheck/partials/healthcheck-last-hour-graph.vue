<template>
  <div class="position-relative">
    <c-progress-overlay :pending="pending" />

    <limited-time-line-chart
      :datasets="datasets"
      :options="lastHourChartOptions"
      class="text-primary"
    >
      <template #actions="{ chart }">
        <chart-export-actions
          :chart="chart"
          class="mt-4"
          hide-csv
          @export:png="$emit('export:png', $event)"
        />
      </template>
    </limited-time-line-chart>
  </div>
</template>

<script>
import { COLORS } from '@/config';
import { DATETIME_FORMATS, HEALTHCHECK_HISTORY_GRAPH_RECEIVED_FACTOR } from '@/constants';

import { colorToRgba } from '@/helpers/color';

import ChartExportActions from '@/components/common/chart/chart-export-actions.vue';
import LimitedTimeLineChart from '@/components/common/chart/limited-time-line-chart.vue';

export default {
  components: {
    ChartExportActions,
    LimitedTimeLineChart,
  },
  props: {
    messagesRates: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    maxMessagesPerMinute: {
      type: Number,
      required: true,
    },
  },

  emits: ['export:png'],

  computed: {
    suggestedMax() {
      return (
        Math.max(0, ...this.messagesRates.map(({ rate }) => rate)) *
        HEALTHCHECK_HISTORY_GRAPH_RECEIVED_FACTOR
      );
    },

    datasets() {
      return [
        {
          backgroundColor: colorToRgba(COLORS.primary, 0.7),
          fill: true,
          data: this.messagesRates.map(({ time, rate }) => ({
            x: time * 1000,
            y: rate,
          })),
        },
      ];
    },

    lastHourChartOptions() {
      return {
        animation: false,
        parsing: false,
        scales: {
          x: {
            time: {
              unit: 'minute',
              tooltipFormat: DATETIME_FORMATS.timePicker,
              displayFormats: {
                minute: DATETIME_FORMATS.timePicker,
              },
            },
          },
          y: {
            suggestedMax: this.suggestedMax,
          },
        },
        plugins: {
          tooltip: {
            callbacks: {
              label: tooltip => `${tooltip.formattedValue} ${this.$t('healthcheck.messagesPerMinute')}`,
            },
          },
          limit: {
            enabled: Number.isFinite(this.maxMessagesPerMinute),
            scaleID: 'x',
            value: this.maxMessagesPerMinute,
            backgroundColor: COLORS.healthcheck.error,
            borderWidth: 1,
          },
        },
      };
    },
  },
};
</script>
