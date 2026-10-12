<template>
  <div class="position-relative">
    <c-progress-overlay :pending="pending" />

    <limited-time-line-chart
      :datasets="datasets"
      :options="lastHourChartOptions"
      :dark="$system.dark"
      class="text--primary"
    >
      <template #actions="{ chart }">
        <chart-export-actions
          :chart="chart"
          class="mt-4"
          hide-csv
          v-on="$listeners"
        />
      </template>
    </limited-time-line-chart>
  </div>
</template>

<script>
import { COLORS } from '@/config';
import { DATETIME_FORMATS, HEALTHCHECK_HISTORY_GRAPH_RECEIVED_FACTOR } from '@/constants';

import { colorToRgba } from '@/helpers/color';

import { entitiesMessageRateStatsMixin } from '@/mixins/entities/message-rate-stats';
import { localQueryMixin } from '@/mixins/query/query';

import ChartExportActions from '@/components/common/chart/chart-export-actions.vue';

const LimitedTimeLineChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/limited-time-line-chart.vue');

export default {
  inject: ['$system'],
  components: { ChartExportActions, LimitedTimeLineChart },
  mixins: [entitiesMessageRateStatsMixin, localQueryMixin],
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
  computed: {
    suggestedMax() {
      return Math.max.apply(null, this.messagesRates.map(({ rate }) => rate))
        * HEALTHCHECK_HISTORY_GRAPH_RECEIVED_FACTOR;
    },

    datasets() {
      return [{
        backgroundColor: colorToRgba(COLORS.primary, 0.7),
        fill: true,
        data: this.messagesRates.map(({ time, rate }) => ({
          x: time * 1000,
          y: rate,
        })),
      }];
    },

    lastHourChartOptions() {
      return {
        animation: false,
        scales: {
          x: {
            time: {
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
            enabled: true,
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
