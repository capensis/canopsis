<template>
  <div>
    <c-text-field
      :model-value="instructionExecution.description"
      :label="t('common.description')"
      variant="filled"
      readonly
    />
    <c-text-field
      :model-value="instructionExecution.alarm.display_name"
      :label="t('alarm.alarmDisplayName')"
      variant="filled"
      readonly
      disabled
    />
    <remediation-instruction-execute-steps
      :steps="instructionExecution.steps"
      :execution-id="instructionExecution._id"
      :next-pending="nextPending"
      :previous-pending="previousPending"
      @next-step="$emit('next-step', $event)"
      @next-operation="$emit('next-operation')"
      @previous-operation="$emit('previous-operation')"
      @execute-job="$emit('execute-job', $event)"
      @cancel-job-execution="$emit('cancel-job-execution', $event)"
    />
  </div>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

import RemediationInstructionExecuteSteps from './remediation-instruction-execute-steps.vue';

export default {
  components: { RemediationInstructionExecuteSteps },

  props: {
    instructionExecution: {
      type: Object,
      required: true,
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

  setup() {
    const { t } = useI18n();

    return { t };
  },
};
</script>
