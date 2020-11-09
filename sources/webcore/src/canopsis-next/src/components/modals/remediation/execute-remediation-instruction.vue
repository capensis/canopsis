<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ config.assignedInstruction.name }}
      template(slot="text")
        remediation-instruction-execute(
          v-if="executionInstruction",
          :execution-instruction="executionInstruction"
        )
        v-layout(v-else, justify-center)
          v-progress-circular(indeterminate, color="primary")
</template>

<script>
import { MODALS, REMEDIATION_INSTRUCTION_EXECUTION_STATUSES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import authMixin from '@/mixins/auth';
import entitiesRemediationInstructionExecutionMixin from '@/mixins/entities/remediation/executions';

import RemediationInstructionExecute from '@/components/other/remediation/instruction-execute/remediation-instruction-execute.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.executeRemediationInstruction,
  components: {
    ModalWrapper,
    RemediationInstructionExecute,
  },
  mixins: [
    authMixin,
    modalInnerMixin,
    entitiesRemediationInstructionExecutionMixin,
  ],
  data() {
    const { execution } = this.modal.config.assignedInstruction;

    return {
      executionInstructionId: execution && execution._id,
    };
  },
  computed: {
    executionInstruction() {
      return this.getRemediationInstructionExecution(this.executionInstructionId);
    },
  },
  watch: {
    executionInstruction(executionInstruction) {
      if (executionInstruction.status !== REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running) {
        const isFailedExecution = [
          REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.failed,
          REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.aborted,
        ].includes(executionInstruction.status);
        const type = isFailedExecution ? 'failed' : 'success';
        const text = this.$t(`remediationInstructionExecute.popups.${type}`, {
          instructionName: executionInstruction.name,
        });

        if (isFailedExecution) {
          this.$popups.error({ text });
        } else {
          this.$popups.success({ text });
        }

        if (this.config.onFinished) {
          this.config.onFinished();
        }

        this.$modals.hide();
      }
    },
  },
  async mounted() {
    await this.fetchInstructionExecution();
  },
  methods: {
    async createInstructionExecution() {
      const { _id: instructionId } = this.config.assignedInstruction;

      const instructionExecution = await this.createRemediationInstructionExecution({
        data: {
          alarm: this.config.alarm._id,
          instruction: instructionId,
        },
      });

      this.executionInstructionId = instructionExecution._id;

      if (this.config.onCreate) {
        await this.config.onCreate();
      }
    },

    async fetchInstructionExecution() {
      const { execution } = this.config.assignedInstruction;

      try {
        if (!execution) {
          await this.createInstructionExecution();
        } else if (execution.status === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.paused) {
          await this.resumeRemediationInstructionExecution({ id: execution._id });
        } else {
          await this.fetchRemediationInstructionExecution({ id: execution._id });
        }
      } catch (err) {
        this.$popups.error({ text: err.error || this.$t('errors.default') });
        this.$modals.hide();
      }
    },
  },
};
</script>
