<template>
  <v-menu
    offset-y
    @click.prevent.stop=""
  >
    <template #activator="{ on }">
      <v-btn
        depressed
        small
        light
        v-on="on"
      >
        <v-icon>{{ icon }}</v-icon>
      </v-btn>
    </template>
    <v-list>
      <v-list-item
        v-for="assignedInstruction in assignedInstructions"
        :key="assignedInstruction._id"
        :disabled="isDisabled(assignedInstruction)"
        @click.stop.prevent="$emit('execute', assignedInstruction)"
      >
        <v-list-item-title>{{ getLabel(assignedInstruction) }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import { find } from 'lodash';

import { REMEDIATION_INSTRUCTION_EXECUTION_STATUSES } from '@/constants';

import { isInstructionExecutionIconInProgress } from '@/helpers/entities/remediation/instruction-execution/form';

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
