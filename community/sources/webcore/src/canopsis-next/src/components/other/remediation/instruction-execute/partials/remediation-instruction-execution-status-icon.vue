<template>
  <v-tooltip left>
    <template #activator="{ on }">
      <v-icon :color="color" v-on="on">
        {{ icon }}
      </v-icon>
    </template>
    <span>{{ $t(`remediation.instructionExecute.status.tooltips.${status}`) }}</span>
  </v-tooltip>
</template>

<script>
import { computed } from 'vue';

import {
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES,
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES_ICONS,
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES_COLORS,
} from '@/constants';

export default {
  props: {
    status: {
      type: Number,
      default: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running,
    },
  },
  setup(props) {
    const icon = computed(() => REMEDIATION_INSTRUCTION_EXECUTION_STATUSES_ICONS[props.status]);
    const color = computed(() => REMEDIATION_INSTRUCTION_EXECUTION_STATUSES_COLORS[props.status]);

    return {
      icon,
      color,
    };
  },
};
</script>
