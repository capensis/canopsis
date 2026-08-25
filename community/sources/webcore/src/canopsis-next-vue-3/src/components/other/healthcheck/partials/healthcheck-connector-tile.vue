<template>
  <v-card
    :color="color"
    :class="{ 'healthcheck-connector-tile--clickable': !isDisabled }"
    class="healthcheck-connector-tile pa-4 text-white"
    variant="flat"
    v-on="cardListeners"
  >
    <h4 class="text-subtitle-1 font-weight-bold text-truncate">
      {{ connector.name }}
    </h4>
    <v-divider class="my-2" />
    <div v-if="isDisabled" class="text-body-2">
      {{ $t('healthcheck.connectorsBlocks.disabled') }}
    </div>
    <div v-else-if="isInactive" class="text-body-2">
      <span>{{ $t('healthcheck.connectorsBlocks.notEnoughData') }}</span>
      <c-help-icon
        :text="$t('healthcheck.connectorsBlocks.description')"
        icon-class="ml-1"
        max-width="280"
        top
        small
        @click.stop
      />
    </div>
    <template v-else>
      <div class="text-body-2">
        {{ $t('healthcheck.connectorsBlocks.numberOfEvents') }}:
        <strong>{{ connector.event_count }}</strong>
      </div>
      <div v-if="connector.time" class="text-body-2">
        {{ dateToString(connector.time) }}
      </div>
    </template>
  </v-card>
</template>

<script>
import { ANOMALY_MONITORED_CONNECTOR_STATUSES, MODALS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

const STATUS_COLORS = {
  [ANOMALY_MONITORED_CONNECTOR_STATUSES.ok]: 'success',
  [ANOMALY_MONITORED_CONNECTOR_STATUSES.anomaly]: 'error',
  [ANOMALY_MONITORED_CONNECTOR_STATUSES.inactive]: 'grey-darken-1',
  [ANOMALY_MONITORED_CONNECTOR_STATUSES.disabled]: 'grey',
};

export default {
  props: {
    connector: {
      type: Object,
      required: true,
    },
  },

  emits: ['refresh'],

  computed: {
    isDisabled() {
      return this.connector.status === ANOMALY_MONITORED_CONNECTOR_STATUSES.disabled;
    },

    isInactive() {
      return this.connector.status === ANOMALY_MONITORED_CONNECTOR_STATUSES.inactive;
    },

    color() {
      return STATUS_COLORS[this.connector.status] ?? 'grey-darken-1';
    },

    cardListeners() {
      if (this.isDisabled) {
        return {};
      }

      return {
        click: this.showAnomalyMonitoredConnectorHistoryModal,
      };
    },
  },

  methods: {
    dateToString(value) {
      return convertDateToString(value);
    },

    showAnomalyMonitoredConnectorHistoryModal() {
      this.$modals.show({
        name: MODALS.anomalyMonitoredConnectorHistory,
        config: {
          connector: this.connector,
          fetchList: () => this.$emit('refresh'),
        },
      });
    },
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
