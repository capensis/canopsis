<template>
  <c-tooltip location="right">
    <template #activator="{ props: activatorProps }">
      <v-icon v-bind="activatorProps" :color="statusColor" :icon="statusIcon" />
    </template>
    <span>{{ label }}</span>
  </c-tooltip>
</template>

<script>
import { computed } from 'vue';

import { JOB_EXECUTION_STATUS } from '@/constants';

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
    const label = computed(() => t(`jobs.executionStatus.${props.status}`));
    const statusColor = computed(
      () => ({
          [JOB_EXECUTION_STATUS.running]: 'warning',
          [JOB_EXECUTION_STATUS.succeeded]: 'success',
          [JOB_EXECUTION_STATUS.failed]: 'error',
          [JOB_EXECUTION_STATUS.aborted]: 'grey',
        })[props.status] ?? 'grey',
    );
    const statusIcon = computed(
      () => ({
          [JOB_EXECUTION_STATUS.running]: 'play_circle',
          [JOB_EXECUTION_STATUS.succeeded]: 'check_circle',
          [JOB_EXECUTION_STATUS.failed]: 'cancel',
          [JOB_EXECUTION_STATUS.aborted]: 'cancel',
        })[props.status] ?? 'help_outline',
    );

    return {
      label,
      statusColor,
      statusIcon,
    };
  },
};
</script>
