<template>
  <v-layout column>
    <healthcheck-last-hour-filters
      :event-types="query.eventTypes"
      :connector-names="query.connectorNames"
      @update:event-types="updateQueryField('eventTypes', $event)"
      @update:connector-names="updateQueryField('connectorNames', $event)"
    />
    <healthcheck-last-hour-graph
      :max-messages-per-minute="maxMessagesPerMinute"
      :messages-rates="messagesRates"
      :pending="pending"
      @export:png="exportChartAsPng"
    />
  </v-layout>
</template>

<script>
import { HEALTHCHECK_LAST_HOUR_FILENAME_PREFIX, SOCKET_ROOMS } from '@/config';
import { DATETIME_FORMATS, MESSAGE_STATS_INTERVALS, TIME_UNITS } from '@/constants';

import { convertDateToString, getNowTimestamp, subtractUnitFromDate } from '@/helpers/date/date';
import { saveFile } from '@/helpers/file/files';

import { entitiesMessageRateStatsMixin } from '@/mixins/entities/message-rate-stats';
import { localQueryMixin } from '@/mixins/query/query';

import HealthcheckLastHourGraph from '@/components/other/healthcheck/partials/healthcheck-last-hour-graph.vue';

import HealthcheckLastHourFilters from './healthcheck-last-hour-filters.vue';

export default {
  inject: ['$system'],
  components: { HealthcheckLastHourFilters, HealthcheckLastHourGraph },
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
      query: {
        eventTypes: [],
        connectorNames: [],
      },
    };
  },
  mounted() {
    this.fetchList();
  },
  beforeDestroy() {
    this.leaveFromSocketRoom();
  },
  methods: {
    setMessageRates(messagesRates) {
      this.messagesRates = messagesRates;
    },

    joinToSocketRoom() {
      this.$socket
        .join(SOCKET_ROOMS.messageRates, this.getQuery())
        .addListener(this.setMessageRates);
    },

    leaveFromSocketRoom() {
      this.$socket
        .leave(SOCKET_ROOMS.messageRates)
        .removeListener(this.setMessageRates);
    },

    async exportChartAsPng(chartBlob) {
      try {
        const firstStatsMessage = this.messagesRates.at(0);
        const lastStatsMessage = this.messagesRates.at(-1);

        const day = convertDateToString(firstStatsMessage.time, DATETIME_FORMATS.short);
        const fromTime = convertDateToString(firstStatsMessage.time, DATETIME_FORMATS.timePicker);
        const toTime = convertDateToString(lastStatsMessage.time, DATETIME_FORMATS.timePicker);

        await saveFile(chartBlob, `${HEALTHCHECK_LAST_HOUR_FILENAME_PREFIX}${day}(${fromTime}-${toTime})`);
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: err.message || this.$t('errors.default') });
      }
    },

    getQuery() {
      const to = getNowTimestamp();
      const from = subtractUnitFromDate(to, 1, TIME_UNITS.hour);

      return {
        from,
        to,
        event_types: this.query.eventTypes,
        interval: MESSAGE_STATS_INTERVALS.minute,
      };
    },

    async fetchList() {
      this.pending = true;

      this.leaveFromSocketRoom();
      this.joinToSocketRoom();

      const { data: messagesRates } = await this.fetchMessageRateStatsWithoutStore({
        params: this.getQuery(),
      });

      this.setMessageRates(messagesRates);

      this.pending = false;
    },
  },
};
</script>
