<template lang="pug">
  v-menu(offset-y, @click.prevent.stop="")
    template(#activator="{ on }")
      v-btn(v-on="on", depressed, small, light)
        v-icon {{ icon }}
    v-list
      v-list-tile(
        v-for="assignedInstruction in assignedInstructions",
        :key="assignedInstruction._id",
        :disabled="isDisabled(assignedInstruction)",
        @click.stop.prevent="$emit('execute', assignedInstruction)"
      )
        v-list-tile-title {{ getLabel(assignedInstruction) }}
</template>

<script>
import { find } from 'lodash';

import { REMEDIATION_INSTRUCTION_EXECUTION_STATUSES } from '@/constants';

import { isInstructionExecutionIconInProgress } from '@/helpers/forms/remediation-instruction-execution';

export default {
  props: {
    icon: {
      type: String,
      required: true,
    },
    entity: {
      type: Object,
      default: () => ({}),
    },
    assignedInstructions: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    hasRunningInstruction() {
      return isInstructionExecutionIconInProgress(this.entity.instruction_execution_icon);
    },

    pausedInstructions() {
      return this.assignedInstructions.filter(instruction => instruction.execution);
    },
  },
  methods: {
    isDisabled(instruction) {
      return this.hasRunningInstruction
        || (
          Boolean(this.pausedInstructions.length)
          && !find(this.pausedInstructions, { _id: instruction._id })
        );
    },

    getLabel(instruction) {
      const { execution, name } = instruction;
      let titlePrefix = 'execute';

      if (execution) {
        titlePrefix = execution.status === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running
          ? 'inProgress'
          : 'resume';
      }

      return this.$t(`remediation.instruction.${titlePrefix}Instruction`, {
        instructionName: name,
      });
    },
  },
};
</script>
