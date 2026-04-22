<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ title }}</span>
    </template>
    <template #text="">
      <div class="position-relative">
        <c-progress-overlay :pending="pending" />
        <anomaly-monitored-connector-history-body
          :enabled="connector.enabled"
          :connector-name="connector.name"
          :history="history"
          :interval="query.interval"
          @update:interval="updateQueryInterval"
          @update:enabled="updateEnabled"
        />
      </div>
    </template>
  </modal-wrapper>
</template>

<script>
import { computed, ref, onMounted } from 'vue';

import { MODALS, TIME_UNITS } from '@/constants';

import {
  getNowIntervalValueForHours,
  convertStartDateIntervalToTimestamp,
  convertStopDateIntervalToTimestamp,
} from '@/helpers/date/date-intervals';

import { useInnerModal } from '@/hooks/modals';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useAnomalyMonitoredConnectors } from '@/hooks/store/modules/anomaly-monitored-connector';

import ModalWrapper from '../modal-wrapper.vue';

import AnomalyMonitoredConnectorHistoryBody from './partials/anomaly-monitored-connector-history-body.vue';

export default {
  name: MODALS.anomalyMonitoredConnectorHistory,
  components: {
    ModalWrapper,
    AnomalyMonitoredConnectorHistoryBody,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config } = useInnerModal(props);
    const {
      updateAnomalyMonitoredConnectorEnabled,
      fetchAnomalyMonitoredConnectorHistoryWithoutStore,
    } = useAnomalyMonitoredConnectors();

    const connector = ref({ ...config.value.connector });
    const history = ref({});

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
          id: config.value.connector?.id,
          params: {
            from: convertStartDateIntervalToTimestamp(fetchQuery.interval),
            to: convertStopDateIntervalToTimestamp(fetchQuery.interval),
          },
        });
      },
    });

    const title = computed(() => config.value.connector?.name ?? '');

    /**
     * Persists the connector enabled flag (PATCH), assigns the update result to local `connector`,
     * and calls `config.fetchList()` so the parent list view stays in sync.
     *
     * @param {boolean} value - Enabled state emitted from `anomaly-monitored-connector-history-body`.
     * @returns {Promise<*>} Promise from `updateAnomalyMonitoredConnectorEnabled` (axios/store layer).
     */
    const updateEnabled = async (value) => {
      connector.value = await updateAnomalyMonitoredConnectorEnabled({
        id: connector.value.id,
        data: { enabled: value },
      });

      config.value.fetchList();
    };

    onMounted(fetchHistory);

    return {
      connector,
      history,
      title,
      pending,
      query,
      updateQueryInterval,
      updateEnabled,
    };
  },
};
</script>
