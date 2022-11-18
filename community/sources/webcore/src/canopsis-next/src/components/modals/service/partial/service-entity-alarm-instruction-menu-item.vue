<template lang="pug">
  v-list-tile(
    :disabled="disabled",
    @click.stop.prevent="$emit('execute', assignedInstruction)"
  )
    v-list-tile-title {{ label }}
</template>

<script>
import { REMEDIATION_INSTRUCTION_EXECUTION_STATUSES } from '@/constants';

import { isInstructionExecutionIconInProgress } from '@/helpers/forms/remediation-instruction-execution';

export default {
  props: {
    assignedInstruction: {
      type: Object,
      required: true,
    },
  },
  computed: {
    disabled() {
      return this.assignedInstruction?.execution?.status === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running
        || isInstructionExecutionIconInProgress(this.assignedInstruction?.instruction_execution_icon);
    },

    label() {
      const { execution } = this.assignedInstruction;
      let titlePrefix = 'execute';

      if (execution) {
        titlePrefix = execution.status === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running
          ? 'inProgress'
          : 'resume';
      }

      return this.$t(`remediationInstructions.${titlePrefix}Instruction`, {
        instructionName: this.assignedInstruction.name,
      });
    },
  },
};
</script>
