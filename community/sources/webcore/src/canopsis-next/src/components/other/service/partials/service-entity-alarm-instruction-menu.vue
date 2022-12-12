<template lang="pug">
  v-menu(offset-y, @click.prevent.stop="")
    v-btn(slot="activator", v-on="$listeners", depressed, small, light)
      v-icon {{ icon }}
    v-list
      v-list-tile(
        v-for="assignedInstruction in assignedInstructions",
        :key="assignedInstruction._id",
        :disabled="isDisabledAssignedInstruction(assignedInstruction)",
        @click.stop.prevent="$emit('execute', assignedInstruction)"
      )
        v-list-tile-title {{ getAssignedInstructionLabel(assignedInstruction) }}
</template>

<script>
import { get } from 'lodash';

import { REMEDIATION_INSTRUCTION_EXECUTION_STATUSES } from '@/constants';

export default {
  props: {
    assignedInstructions: {
      type: Array,
      default: () => [],
    },
    icon: {
      type: String,
      required: true,
    },
  },
  methods: {
    isDisabledAssignedInstruction(assignedInstruction) {
      return get(assignedInstruction, 'execution.status') === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running;
    },

    getAssignedInstructionLabel(assignedInstruction) {
      const { execution } = assignedInstruction;
      const titlePrefix = execution ? 'resume' : 'execute';

      return this.$t(`alarm.actions.titles.${titlePrefix}Instruction`, {
        instructionName: assignedInstruction.name,
      });
    },
  },
};
</script>
