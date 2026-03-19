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

import { JOB_LAST_RUN_STATUS } from '@/constants';

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

    const label = computed(() => t(`jobs.lastRunStatus.${props.status}`));

    const statusColor = computed(() => (
      {
        [JOB_LAST_RUN_STATUS.succeed]: 'success',
        [JOB_LAST_RUN_STATUS.failed]: 'error',
        [JOB_LAST_RUN_STATUS.aborted]: 'grey',
      }[props.status] ?? 'grey'
    ));

    const statusIcon = computed(() => (
      {
        [JOB_LAST_RUN_STATUS.succeed]: 'check_circle',
        [JOB_LAST_RUN_STATUS.failed]: 'close_circle',
        [JOB_LAST_RUN_STATUS.aborted]: 'close_circle',
      }[props.status] ?? 'help_outline'
    ));

    return {
      label,
      statusColor,
      statusIcon,
    };
  },
};
</script>
