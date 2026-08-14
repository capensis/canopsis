<template>
  <c-layout column>
    <remediation-instruction-execute-step
      v-for="(step, index) in steps"
      :key="step.step_id"
      :step="step"
      :step-number="index + 1"
      :is-last-step="index === steps.length - 1"
      :is-first-step="index === 0"
      :execution-id="executionId"
      :next-pending="nextPending"
      :previous-pending="previousPending"
      @next-step="$emit('next-step', $event)"
      @next-operation="$emit('next-operation')"
      @previous-operation="$emit('previous-operation')"
      @execute-job="$emit('execute-job', $event)"
      @cancel-job-execution="$emit('cancel-job-execution', $event)"
    />
  </c-layout>
</template>

<script>
import RemediationInstructionExecuteStep from './remediation-instruction-execute-step.vue';

export default {
  components: { RemediationInstructionExecuteStep },

  props: {
    executionId: {
      type: String,
      required: true,
    },
    steps: {
      type: Array,
      default: () => [],
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

  emits: ['next-step', 'next-operation', 'previous-operation', 'execute-job', 'cancel-job-execution'],
};
</script>
