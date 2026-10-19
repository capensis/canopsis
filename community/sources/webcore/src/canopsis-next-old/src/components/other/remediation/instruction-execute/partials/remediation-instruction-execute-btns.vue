<template>
  <v-layout class="gap-2" justify-end>
    <v-btn
      :loading="cancelPending"
      :disabled="cancelPending"
      color="error"
      outlined
      @click="cancel"
    >
      {{ $t('remediation.instructionExecute.cancelInstruction') }}
    </v-btn>
    <v-btn
      :loading="pausePending"
      :disabled="pausePending"
      color="secondary"
      @click="pause"
    >
      {{ $t('remediation.instructionExecute.pauseInstruction') }}
    </v-btn>
  </v-layout>
</template>

<script>
import { usePendingHandler } from '@/hooks/query/pending';
import { useRemediationInstructionExecution } from '@/hooks/store/modules/remediation-instruction-execution';

export default {
  props: {
    instructionExecution: {
      type: Object,
      required: true,
    },
  },
  setup(props, { emit }) {
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

    /**
     * Pauses the remediation instruction execution by ID and closes the modal.
     */
    const {
      pending: pausePending,
      handler: pause,
    } = usePendingHandler(async () => {
      await pauseRemediationInstructionExecution({ id: props.instructionExecution._id });

      closeModal();
    });

    return {
      cancelPending,
      cancel,
      pausePending,
      pause,
    };
  },
};
</script>
