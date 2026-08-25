<template>
  <c-layout gap="4" column>
    <c-enabled-field
      :model-value="enabled"
      :confirm-disable="confirmDisable"
      hide-details
      with-background
      @update:model-value="updateEnabled"
    />

    <anomaly-monitored-connector-history-graph-placeholder v-if="!history.data?.length" />

    <template v-else>
      <c-layout justify="end">
        <div>
          <c-hours-interval-field :value="interval" @input="updateInterval" />
        </div>
      </c-layout>

      <anomaly-monitored-connector-history-summary-stats :history="history" />

      <anomaly-monitored-connector-history-chart :history="history" />
    </template>
  </c-layout>
</template>

<script>
import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import AnomalyMonitoredConnectorHistoryChart from './anomaly-monitored-connector-history-chart.vue';
import AnomalyMonitoredConnectorHistoryGraphPlaceholder from './anomaly-monitored-connector-history-graph-placeholder.vue';
import AnomalyMonitoredConnectorHistorySummaryStats from './anomaly-monitored-connector-history-summary-stats.vue';

export default {
  components: {
    AnomalyMonitoredConnectorHistoryChart,
    AnomalyMonitoredConnectorHistoryGraphPlaceholder,
    AnomalyMonitoredConnectorHistorySummaryStats,
  },

  props: {
    connectorName: {
      type: String,
      default: '',
    },
    enabled: {
      type: Boolean,
      default: false,
    },
    history: {
      type: Object,
      default: () => ({}),
    },
    interval: {
      type: String,
      default: '',
    },
  },

  emits: ['update:enabled', 'update:interval'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const modals = useModals();

    const updateInterval = value => emit('update:interval', value);
    const updateEnabled = value => emit('update:enabled', value);

    const confirmDisable = () =>
      new Promise((resolve) => {
        modals.show({
          name: MODALS.confirmation,
          config: {
            text: t('modals.anomalyMonitoredConnectorHistory.disableConfirm', {
              connectorName: props.connectorName,
            }),
            action: () => resolve(true),
            cancel: () => resolve(false),
          },
        });
      });

    return {
      confirmDisable,
      updateEnabled,
      updateInterval,
    };
  },
};
</script>
