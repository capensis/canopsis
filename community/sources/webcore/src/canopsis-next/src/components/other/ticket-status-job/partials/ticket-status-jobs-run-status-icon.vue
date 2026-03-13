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

import { JOB_RUN_STATUS } from '@/constants';

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

    const label = computed(() => (
      t(`jobs.runStatus.${props.status}`) ? t(`jobs.runStatus.${props.status}`) : t('jobs.runStatus.inProgress')
    ));

    const statusColor = computed(() => (
      {
        [JOB_RUN_STATUS.succeed]: 'success',
        [JOB_RUN_STATUS.failed]: 'error',
        [JOB_RUN_STATUS.aborted]: 'grey',
      }[props.status] ?? 'grey'
    ));

    const statusIcon = computed(() => (
      {
        [JOB_RUN_STATUS.succeed]: 'check',
        [JOB_RUN_STATUS.failed]: 'close',
        [JOB_RUN_STATUS.aborted]: 'close',
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
