<template lang="pug">
  div
    v-text-field(
      :value="executionInstruction.description",
      :label="$t('common.description')",
      readonly,
      box
    )
    remediation-instruction-execute-steps(:steps="executionInstruction.steps")
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
      method: 'refresh',
      delay: INSTRUCTION_EXECUTE_FETCHING_INTERVAL,
    }),
  ],
  props: {
    executionInstructionId: {
      type: [String, Number],
      required: true,
    },
  },
  computed: {
    executionInstruction() {
      return this.getRemediationInstructionExecution(this.executionInstructionId);
    },
  },
  methods: {
    async refresh() {
      await this.fetchRemediationInstructionExecution({ id: this.executionInstructionId });
    },
  },
};
</script>
