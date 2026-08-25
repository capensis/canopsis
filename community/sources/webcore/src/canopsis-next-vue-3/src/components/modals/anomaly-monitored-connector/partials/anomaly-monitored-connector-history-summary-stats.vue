<template>
  <c-layout
    class="anomaly-monitored-connector-history-summary-stats mb-4"
    gap="3"
    justify="center"
    wrap
  >
    <v-sheet
      v-for="tile in tiles"
      :key="tile.color"
      :color="tile.color"
      class="anomaly-monitored-connector-history-summary-stats__tile pa-4 text-center"
    >
      <div class="text-medium-emphasis text-body-1">
        {{ tile.label }}
      </div>
      <div class="text-subtitle-1 font-weight-medium">
        {{ tile.value }}
      </div>
    </v-sheet>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    history: {
      type: Object,
      default: () => ({}),
    },
  },

  setup(props) {
    const { t } = useI18n();

    const averageEvents = computed(() => {
      const value = Math.round(props.history.average_events);

      return Number.isFinite(value) ? value : '-';
    });

    const interval = computed(() => {
      if (!props.history.interval_minutes) {
        return '-';
      }

      return `${props.history.interval_minutes} ${t(
        'common.times.minute',
        props.history.interval_minutes,
      )}`;
    });

    const tiles = computed(() => [
      {
        value: props.history.total_anomalies ?? '-',
        label: t('healthcheck.connectorsBlocks.anomaliesDetected'),
        color: 'error-lighten-5',
      },
      {
        value: averageEvents.value,
        label: t('healthcheck.connectorsBlocks.averageEvents'),
        color: 'info-lighten-4',
      },
      {
        value: interval.value,
        label: t('common.interval'),
        color: 'grey-lighten-2',
      },
    ]);

    return {
      tiles,
    };
  },
};
</script>

<style lang="scss" scoped>
.anomaly-monitored-connector-history-summary-stats__tile {
  min-width: 200px;
  border-radius: 5px;
}
</style>
