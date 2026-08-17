<template>
  <v-tooltip location="start">
    <template #activator="{ props: tooltipProps }">
      <v-icon v-bind="tooltipProps" :color="color">
        {{ icon }}
      </v-icon>
    </template>
    <span>{{ t(`remediation.instructionExecute.status.tooltips.${status}`) }}</span>
  </v-tooltip>
</template>

<script>
import { computed } from 'vue';

import {
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES,
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES_COLORS,
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES_ICONS,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    status: {
      type: Number,
      default: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running,
    },
  },

  setup(props) {
    const { t } = useI18n();

    const icon = computed(() => REMEDIATION_INSTRUCTION_EXECUTION_STATUSES_ICONS[props.status]);
    const color = computed(() => REMEDIATION_INSTRUCTION_EXECUTION_STATUSES_COLORS[props.status]);

    return {
      color,
      icon,
      t,
    };
  },
};
</script>
