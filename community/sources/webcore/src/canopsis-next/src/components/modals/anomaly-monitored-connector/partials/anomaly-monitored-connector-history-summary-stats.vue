<template>
  <v-layout
    class="anomaly-monitored-connector-history-summary-stats gap-3 mb-4"
    justify-center
    wrap
  >
    <v-sheet
      v-for="tile in tiles"
      :key="tile.color"
      :color="tile.color"
      class="anomaly-monitored-connector-history-summary-stats__tile pa-4 text-center"
    >
      <div class="text--secondary text-body-1">
        {{ tile.label }}
      </div>
      <div class="text-subtitle-1 font-weight-medium">
        {{ tile.value }}
      </div>
    </v-sheet>
  </v-layout>
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
    const { t, tc } = useI18n();

    const tiles = computed(() => [
      {
        value: props.history.total_anomalies ?? '-',
        label: t('healthcheck.connectorsBlocks.anomaliesDetected'),
        color: 'error lighten-5',
      },
      {
        value: props.history.average_events ?? '-',
        label: t('healthcheck.connectorsBlocks.averageEvents'),
        color: 'info lighten-4',
      },
      {
        value: props.history.interval_minutes
          ? `${props.history.interval_minutes} ${tc('common.times.minute', props.history.interval_minutes)}`
          : '-',
        label: t('common.interval'),
        color: 'grey lighten-2',
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
