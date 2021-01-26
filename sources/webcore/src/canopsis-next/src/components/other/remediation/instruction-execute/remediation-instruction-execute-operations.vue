<template lang="pug">
  v-layout(column)
    remediation-instruction-execute-operation(
      v-for="(operation, index) in operations",
      :key="operation.operation_id",
      :operation="operation",
      :operation-number="getOperationNumber(index)",
      :is-first-operation="!index",
      :is-first-step="isFirstStep",
      :execution-id="executionId",
      @next="nextOperation(index)",
      @previous="previousOperation"
    )
</template>

<script>
import { getLetterByIndex } from '@/helpers/string';

import RemediationInstructionExecuteOperation from './remediation-instruction-execute-operation.vue';

export default {
  components: { RemediationInstructionExecuteOperation },
  props: {
    executionId: {
      type: String,
      required: true,
    },
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
  },
  methods: {
    getOperationNumber(index) {
      return `${this.stepNumber}${getLetterByIndex(index)}`;
    },

    nextOperation(index) {
      const event = index === this.operations.length - 1
        ? 'finish'
        : 'next';

      this.$emit(event);
    },

    previousOperation() {
      this.$emit('previous');
    },
  },
};
</script>
