<template>
  <v-layout column>
    <remediation-instruction-execute-step-operation
      v-for="(operation, index) in operations"
      :key="operation.operation_id"
      :operation="operation"
      :operation-number="getOperationNumber(index)"
      :is-first-operation="!index"
      :is-first-step="isFirstStep"
      :next-pending="nextPending"
      :previous-pending="previousPending"
      @next="nextOperation(index)"
      @previous="previousOperation"
      @execute-job="executeJob"
      @cancel-job-execution="cancelJobExecution"
    />
  </v-layout>
</template>

<script>
import { getLetterByIndex } from '@/helpers/string';

import RemediationInstructionExecuteStepOperation from './remediation-instruction-execute-step-operation.vue';

export default {
  components: { RemediationInstructionExecuteStepOperation },
  props: {
    operations: {
      type: Array,
      default: () => [],
    },
    stepNumber: {
      type: [Number, String],
      required: true,
    },
    isFirstStep: {
      type: Boolean,
      default: false,
    },
    previousPending: {
      type: Boolean,
      default: false,
    },
    nextPending: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    getOperationNumber(index) {
      return `${this.stepNumber}${getLetterByIndex(index)}`;
    },

    nextOperation(index) {
      const event = index === this.operations.length - 1
        ? 'finish'
        : 'next-operation';

      this.$emit(event);
    },

    previousOperation() {
      this.$emit('previous-operation');
    },

    executeJob(data) {
      this.$emit('execute-job', data);
    },

    cancelJobExecution(data) {
      this.$emit('cancel-job-execution', data);
    },
  },
};
</script>
