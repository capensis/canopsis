<template>
  <v-tooltip right>
    <template #activator="{ on }">
      <v-icon
        :color="statusColor"
        v-on="on"
      >
        {{ statusIcon }}
      </v-icon>
    </template>
    <span>{{ label }}</span>
  </v-tooltip>
</template>

<script>
import { computed } from 'vue';

import { JOB_STATUS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    status: {
      type: Number,
      default: null,
    },
  },
  setup(props) {
    const { t } = useI18n();

    const label = computed(() => t(`jobs.status.${props.status}`));

    const statusColor = computed(() => (props.status === JOB_STATUS.running ? 'success' : 'grey darken-1'));

    const statusIcon = computed(() => ({
      [JOB_STATUS.running]: 'play_circle',
      [JOB_STATUS.paused]: 'pause_circle',
      [JOB_STATUS.stopped]: 'stop_circle',
    }[props.status] ?? 'stop_circle'));

    return {
      label,
      statusColor,
      statusIcon,
    };
  },
};
</script>
