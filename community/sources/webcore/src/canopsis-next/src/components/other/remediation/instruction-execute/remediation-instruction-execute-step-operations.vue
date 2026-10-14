<template>
  <c-layout column>
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
  </c-layout>
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

  emits: ['finish', 'next-operation', 'previous-operation', 'execute-job', 'cancel-job-execution'],

  setup(props, { emit }) {
    const getOperationNumber = index => `${props.stepNumber}${getLetterByIndex(index)}`;

    const nextOperation = (index) => {
      const event = index === props.operations.length - 1
        ? 'finish'
        : 'next-operation';

      emit(event);
    };

    const previousOperation = () => emit('previous-operation');
    const executeJob = data => emit('execute-job', data);
    const cancelJobExecution = data => emit('cancel-job-execution', data);

    return {
      getOperationNumber,
      nextOperation,
      previousOperation,
      executeJob,
      cancelJobExecution,
    };
  },
};
</script>
