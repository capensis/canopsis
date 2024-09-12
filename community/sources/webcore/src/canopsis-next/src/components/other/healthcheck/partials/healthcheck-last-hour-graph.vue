<template lang="pug">
  limited-time-line-chart.text--primary(
    :datasets="datasets",
    :options="lastHourChartOptions",
    :dark="$system.dark"
  )
    template(#actions="{ chart }")
      v-layout.mt-4(justify-end)
        v-btn.ma-0(color="primary", @click="exportChart(chart)")
          v-icon(left) file_download
          span {{ $t('common.downloadAsPng') }}
</template>

<script>
import { COLORS, HEALTHCHECK_LAST_HOUR_FILENAME_PREFIX, SOCKET_ROOMS } from '@/config';
import {
  DATETIME_FORMATS,
  HEALTHCHECK_HISTORY_GRAPH_RECEIVED_FACTOR,
  MESSAGE_STATS_INTERVALS,
  TIME_UNITS,
} from '@/constants';

import { colorToRgba } from '@/helpers/color';
import { convertDateToString, getNowTimestamp, subtractUnitFromDate } from '@/helpers/date/date';
import { canvasToBlob } from '@/helpers/charts/canvas';
import { saveFile } from '@/helpers/file/files';

import { entitiesMessageRateStatsMixin } from '@/mixins/entities/message-rate-stats';
import { localQueryMixin } from '@/mixins/query-local/query';

const LimitedTimeLineChart = () => import(/* webpackChunkName: "Charts" */ '@/components/common/chart/limited-time-line-chart.vue');

export default {
  inject: ['$system'],
  components: { LimitedTimeLineChart },
  mixins: [entitiesMessageRateStatsMixin, localQueryMixin],
  props: {
    maxMessagesPerMinute: {
      type: Number,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      messagesRates: [],
    };
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
  mounted() {
    this.fetchList();

    this.$socket
      .join(SOCKET_ROOMS.messageRates)
      .addListener(this.setMessageRates);
  },
  beforeDestroy() {
    this.$socket
      .leave(SOCKET_ROOMS.messageRates)
      .removeListener(this.setMessageRates);
  },
  methods: {
    setMessageRates(messagesRates) {
      this.messagesRates = messagesRates;
    },

    async exportChart(chart) {
      try {
        const firstStatsMessage = this.messagesRates[0];
        const lastStatsMessage = this.messagesRates[this.messagesRates.length - 1];

        const day = convertDateToString(firstStatsMessage.time, DATETIME_FORMATS.short);
        const fromTime = convertDateToString(firstStatsMessage.time, DATETIME_FORMATS.timePicker);
        const toTime = convertDateToString(lastStatsMessage.time, DATETIME_FORMATS.timePicker);

        const chartBlob = await canvasToBlob(chart.canvas);

        await saveFile(chartBlob, `${HEALTHCHECK_LAST_HOUR_FILENAME_PREFIX}${day}(${fromTime}-${toTime})`);
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: err.message || this.$t('errors.default') });
      }
    },

    async fetchList() {
      this.pending = true;

      const to = getNowTimestamp();
      const from = subtractUnitFromDate(to, 1, TIME_UNITS.hour);

      const { data: messagesRates } = await this.fetchMessageRateStatsWithoutStore({
        params: {
          from,
          to,
          interval: MESSAGE_STATS_INTERVALS.minute,
        },
      });

      this.setMessageRates(messagesRates);

      this.pending = false;
    },
  },
};
</script>
