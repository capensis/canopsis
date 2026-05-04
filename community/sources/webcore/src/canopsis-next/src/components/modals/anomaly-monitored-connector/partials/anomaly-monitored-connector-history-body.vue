<template>
  <v-layout class="gap-4" column>
    <c-enabled-field
      :value="enabled"
      :confirm-disable="confirmDisable"
      hide-details
      with-background
      @input="updateEnabled"
    />

    <anomaly-monitored-connector-history-graph-placeholder v-if="!history.data?.length" />

    <template v-else>
      <v-layout justify-end>
        <div>
          <c-hours-interval-field
            :value="interval"
            @input="updateInterval"
          />
        </div>
      </v-layout>

      <anomaly-monitored-connector-history-summary-stats :history="history" />

      <anomaly-monitored-connector-history-chart :history="history" />
    </template>
  </v-layout>
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
  setup(props, { emit }) {
    const { t } = useI18n();
    const modals = useModals();

    /**
     * Forwards the hours-interval field value to the parent (`update:interval`).
     *
     * @param {string} value - Interval string from `c-hours-interval-field`.
     */
    const updateInterval = value => emit('update:interval', value);

    /**
     * Forwards the enabled toggle value to the parent (`update:enabled`).
     *
     * @param {boolean} value - Enabled flag from `c-enabled-field`.
     */
    const updateEnabled = value => emit('update:enabled', value);

    /**
     * Shows a confirmation modal before disabling the connector; used by `c-enabled-field` via
     * `confirm-disable`.
     *
     * @returns {Promise<boolean>} `true` if the user confirms, `false` if they cancel or dismiss.
     */
    const confirmDisable = () => new Promise((resolve) => {
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
      updateInterval,
      updateEnabled,
      confirmDisable,
    };
  },
};
</script>
