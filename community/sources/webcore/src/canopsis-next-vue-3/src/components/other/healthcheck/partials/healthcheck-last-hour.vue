<template>
  <div class="d-flex flex-column">
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
  </div>
</template>

<script>
import { HEALTHCHECK_LAST_HOUR_FILENAME_PREFIX, SOCKET_ROOMS } from '@/config';
import { DATETIME_FORMATS, MESSAGE_STATS_INTERVALS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';
import { saveFile } from '@/helpers/file/files';

import { useMessageRateStats } from '@/hooks/store/modules/message-rate-stats';

import HealthcheckLastHourFilters from './healthcheck-last-hour-filters.vue';
import HealthcheckLastHourGraph from './healthcheck-last-hour-graph.vue';

export default {
  components: {
    HealthcheckLastHourFilters,
    HealthcheckLastHourGraph,
  },
  props: {
    maxMessagesPerMinute: {
      type: Number,
      required: true,
    },
  },

  setup() {
    const { fetchMessageRateStatsWithoutStore } = useMessageRateStats();

    return {
      fetchMessageRateStatsWithoutStore,
    };
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

  beforeUnmount() {
    this.leaveFromSocketRoom();
  },

  methods: {
    setMessageRates(messagesRates = []) {
      this.messagesRates = messagesRates;
    },

    updateQueryField(field, value) {
      this.query = {
        ...this.query,
        [field]: value,
      };
      this.fetchList();
    },

    joinToSocketRoom() {
      this.$socket
        .join(SOCKET_ROOMS.messageRates, this.getQuery())
        .addListener(this.setMessageRates);
    },

    leaveFromSocketRoom() {
      this.$socket.leave(SOCKET_ROOMS.messageRates).removeListener(this.setMessageRates);
    },

    async exportChartAsPng(chartBlob) {
      try {
        const firstStatsMessage = this.messagesRates.at(0);
        const lastStatsMessage = this.messagesRates.at(-1);

        const day = convertDateToString(firstStatsMessage?.time, DATETIME_FORMATS.short);
        const fromTime = convertDateToString(firstStatsMessage?.time, DATETIME_FORMATS.timePicker);
        const toTime = convertDateToString(lastStatsMessage?.time, DATETIME_FORMATS.timePicker);

        await saveFile(
          chartBlob,
          `${HEALTHCHECK_LAST_HOUR_FILENAME_PREFIX}${day}(${fromTime}-${toTime})`,
        );
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: err.message || this.$t('errors.default') });
      }
    },

    getQuery() {
      const query = {
        interval: MESSAGE_STATS_INTERVALS.minute,
      };

      if (this.query.eventTypes.length) {
        query.event_types = this.query.eventTypes;
      }

      if (this.query.connectorNames.length) {
        query.connector_names = this.query.connectorNames;
      }

      return query;
    },

    async fetchList() {
      this.pending = true;

      try {
        this.leaveFromSocketRoom();
        this.joinToSocketRoom();

        const { data: messagesRates } = await this.fetchMessageRateStatsWithoutStore({
          params: this.getQuery(),
        });

        this.setMessageRates(messagesRates);
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
