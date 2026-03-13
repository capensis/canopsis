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

import { JOB_STATE } from '@/constants';

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

    const label = computed(() => t(`jobs.activeState.${props.status}`) || t('jobs.status.unknown'));

    const statusColor = computed(() => ({
      [JOB_STATE.running]: 'primary',
      [JOB_STATE.paused]: 'warning',
      [JOB_STATE.stopped]: 'error',
    }[props.status] ?? 'grey'));

    const statusIcon = computed(() => ({
      [JOB_STATE.running]: 'play_arrow',
      [JOB_STATE.paused]: 'pause',
      [JOB_STATE.stopped]: 'close',
    }[props.status] ?? 'help_outline'));

    return {
      label,
      statusColor,
      statusIcon,
    };
  },
};
</script>
