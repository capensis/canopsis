<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ title }}</span>
    </template>

    <template #text>
      <div class="position-relative">
        <c-progress-overlay :pending="pending" />

        <anomaly-monitored-connector-history-body
          :enabled="connector.enabled"
          :connector-name="connector.name"
          :history="history"
          :interval="query.interval"
          @update:enabled="updateEnabled"
          @update:interval="updateQueryInterval"
        />
      </div>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import { TIME_UNITS } from '@/constants';

import {
  convertStartDateIntervalToTimestamp,
  convertStopDateIntervalToTimestamp,
  getNowIntervalValueForHours,
} from '@/helpers/date/date-intervals';

import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useAnomalyMonitoredConnector } from '@/hooks/store/modules/anomaly-monitored-connector';

import AnomalyMonitoredConnectorHistoryBody from './partials/anomaly-monitored-connector-history-body.vue';

export default {
  components: { AnomalyMonitoredConnectorHistoryBody },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const config = computed(() => props.modal.config ?? {});
    const connector = ref({ ...config.value.connector });
    const history = ref({});

    const {
      updateAnomalyMonitoredConnectorEnabled,
      fetchAnomalyMonitoredConnectorHistoryWithoutStore,
    } = useAnomalyMonitoredConnector();

    const {
      pending,
      query,
      updateQueryInterval,
      fetchHandlerWithQuery: fetchHistory,
    } = usePendingWithLocalQuery({
      initialQuery: {
        interval: getNowIntervalValueForHours(0, TIME_UNITS.hour),
      },
      fetchHandler: async (fetchQuery) => {
        history.value = await fetchAnomalyMonitoredConnectorHistoryWithoutStore({
          id: connector.value.id,
          params: {
            from: convertStartDateIntervalToTimestamp(fetchQuery.interval),
            to: convertStopDateIntervalToTimestamp(fetchQuery.interval),
          },
        });
      },
    });

    const title = computed(() => connector.value.name ?? '');

    const updateEnabled = async (value) => {
      connector.value = await updateAnomalyMonitoredConnectorEnabled({
        id: connector.value.id,
        data: { enabled: value },
      });

      config.value.fetchList?.();
    };

    onMounted(fetchHistory);

    return {
      connector,
      history,
      pending,
      query,
      title,
      updateEnabled,
      updateQueryInterval,
    };
  },
};
</script>
