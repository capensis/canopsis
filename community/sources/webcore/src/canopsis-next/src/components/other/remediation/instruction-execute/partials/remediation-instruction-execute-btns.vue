<template>
  <c-layout class="flex-grow-1" justify="end" gap="2">
    <v-btn
      :loading="cancelPending"
      :disabled="cancelPending"
      color="error"
      variant="outlined"
      @click="cancel"
    >
      {{ t('remediation.instructionExecute.cancelInstruction') }}
    </v-btn>
    <v-btn
      :loading="pausePending"
      :disabled="pausePending"
      color="secondary"
      @click="pause"
    >
      {{ t('remediation.instructionExecute.pauseInstruction') }}
    </v-btn>
  </c-layout>
</template>

<script>
import { useI18n } from '@/hooks/i18n';
import { usePendingHandler } from '@/hooks/query/pending';
import { useRemediationInstructionExecution } from '@/hooks/store/modules/remediation-instruction-execution';

export default {
  props: {
    instructionExecution: {
      type: Object,
      required: true,
    },
  },

  emits: ['close'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const closeModal = () => emit('close');

    const {
      cancelRemediationInstructionExecution,
      pauseRemediationInstructionExecution,
    } = useRemediationInstructionExecution();

    const {
      pending: cancelPending,
      handler: cancel,
    } = usePendingHandler(async () => {
      await cancelRemediationInstructionExecution({ id: props.instructionExecution._id });

      closeModal();
    });

    const {
      pending: pausePending,
      handler: pause,
    } = usePendingHandler(async () => {
      await pauseRemediationInstructionExecution({ id: props.instructionExecution._id });

      closeModal();
    });

    return {
      t,
      cancelPending,
      cancel,
      pausePending,
      pause,
    };
  },
};
</script>
