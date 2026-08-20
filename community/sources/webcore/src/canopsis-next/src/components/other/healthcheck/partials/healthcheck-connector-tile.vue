<template>
  <v-card
    :color="color"
    :class="{ 'healthcheck-connector-tile--clickable': !isDisabled }"
    class="healthcheck-connector-tile pa-4"
    dark
    v-on="cardListeners"
  >
    <h4 class="text-subtitle-1 font-weight-bold text-truncate">
      {{ connector.name }}
    </h4>
    <v-divider class="my-2" dark />
    <div
      v-if="isDisabled"
      class="text-body-2"
    >
      {{ $t('healthcheck.connectorsBlocks.disabled') }}
    </div>
    <div
      v-else-if="isInactive"
      class="text-body-2"
    >
      <span>{{ $t('healthcheck.connectorsBlocks.notEnoughData') }}</span>
      <c-help-icon
        :text="$t('healthcheck.connectorsBlocks.description')"
        icon-class="ml-1"
        max-width="280"
        top
        small
        @click.native.stop
      />
    </div>
    <template v-else>
      <div class="text-body-2">
        {{ $t('healthcheck.connectorsBlocks.numberOfEvents') }}: <strong>{{ connector.event_count }}</strong>
      </div>
      <div
        v-if="connector.time"
        class="text-body-2"
      >
        {{ connector.time | date }}
      </div>
    </template>
  </v-card>
</template>

<script>
import { computed } from 'vue';

import { ANOMALY_MONITORED_CONNECTOR_STATUSES, MODALS } from '@/constants';

import { useModals } from '@/hooks/modals';

const STATUS_COLORS = {
  [ANOMALY_MONITORED_CONNECTOR_STATUSES.ok]: 'success',
  [ANOMALY_MONITORED_CONNECTOR_STATUSES.anomaly]: 'error',
  [ANOMALY_MONITORED_CONNECTOR_STATUSES.inactive]: 'grey darken-1',
  [ANOMALY_MONITORED_CONNECTOR_STATUSES.disabled]: 'grey',
};

export default {
  props: {
    connector: {
      type: Object,
      required: true,
    },
  },
  setup(props, { emit }) {
    const modals = useModals();

    const isDisabled = computed(() => props.connector.status === ANOMALY_MONITORED_CONNECTOR_STATUSES.disabled);
    const isInactive = computed(() => props.connector.status === ANOMALY_MONITORED_CONNECTOR_STATUSES.inactive);

    const color = computed(() => STATUS_COLORS[props.connector.status] ?? 'grey darken-1');

    /**
     * Opens the anomaly monitored connector history modal for this tile’s connector.
     * No-op when the connector status is disabled (tile is not interactive).
     */
    const showAnomalyMonitoredConnectorHistoryModal = () => modals.show({
      name: MODALS.anomalyMonitoredConnectorHistory,
      config: {
        connector: props.connector,
        fetchList: () => emit('refresh'),
      },
    });

    const cardListeners = computed(() => (
      isDisabled.value ? {} : { click: showAnomalyMonitoredConnectorHistoryModal }
    ));

    return {
      isDisabled,
      isInactive,
      color,
      cardListeners,
    };
  },
};
</script>

<style lang="scss" scoped>
.healthcheck-connector-tile {
  width: 340px;
  min-height: 120px;

  &--clickable {
    cursor: pointer;
  }
}
</style>
