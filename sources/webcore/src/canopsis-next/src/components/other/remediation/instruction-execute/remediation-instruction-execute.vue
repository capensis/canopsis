<template lang="pug">
  div
    v-text-field(
      :value="executionInstruction.description",
      :label="$t('common.description')",
      readonly,
      box
    )
    remediation-instruction-execute-steps(
      :steps="executionInstruction.steps",
      :execution-id="executionInstruction._id",
      @next-step="nextStep"
    )
</template>

<script>
import { INSTRUCTION_EXECUTE_FETCHING_INTERVAL } from '@/config';

import pollingMixin from '@/mixins/polling';
import entitiesRemediationInstructionExecutionMixin from '@/mixins/entities/remediation/executions';

import RemediationInstructionExecuteSteps from './remediation-instruction-execute-steps.vue';

export default {
  components: { RemediationInstructionExecuteSteps },
  mixins: [
    entitiesRemediationInstructionExecutionMixin,
    pollingMixin({
      method: 'fetchExecution',
      delay: INSTRUCTION_EXECUTE_FETCHING_INTERVAL,
    }),
  ],
  props: {
    executionInstruction: {
      type: Object,
      required: true,
    },
  },
  methods: {
    async fetchExecution() {
      await this.fetchRemediationInstructionExecution({ id: this.executionInstruction._id });
    },

    nextStep(success) {
      this.nextStepRemediationInstructionExecution({
        id: this.executionInstruction._id,
        data: { failed: !success },
      });
    },
  },
};
</script>
