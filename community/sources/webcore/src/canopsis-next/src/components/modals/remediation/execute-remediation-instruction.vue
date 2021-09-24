<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(:close="close", minimize)
      template(slot="title")
        span {{ config.assignedInstruction.name }}
      template(slot="text")
        remediation-instruction-execute(
          v-if="executionInstruction",
          :execution-instruction="executionInstruction"
        )
        v-layout(v-else, justify-center)
          v-progress-circular(color="primary", indeterminate)
</template>

<script>
import { MODALS, REMEDIATION_INSTRUCTION_EXECUTION_STATUSES } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { authMixin } from '@/mixins/auth';
import { pollingMixinCreator } from '@/mixins/polling';
import { entitiesInfoMixin } from '@/mixins/entities/info';
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
    modalInnerMixin,
    authMixin,
    entitiesRemediationInstructionExecutionMixin,
    entitiesInfoMixin,
    pollingMixinCreator({ method: 'pingInstructionExecution' }),
  ],
  data() {
    const { execution } = this.modal.config.assignedInstruction;

    return {
      executionInstructionId: execution && execution._id,
      pending: true,
    };
  },
  computed: {
    executionInstruction() {
      return this.getRemediationInstructionExecution(this.executionInstructionId);
    },

    pollingDelay() {
      return (this.remediationPauseManualInstructionIntervalSeconds - 1) * 1000;
    },
  },
  watch: {
    async executionInstruction(executionInstruction) {
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

        this.stopPolling();

        if (this.config.onComplete) {
          await this.config.onComplete(executionInstruction);
        }

        this.$modals.hide();
      }
    },
  },
  async mounted() {
    this.pending = true;

    await this.fetchInstructionExecution();

    this.pending = false;
  },
  methods: {
    async pingInstructionExecution() {
      try {
        if (!this.executionInstruction || this.pending) {
          return;
        }

        await this.pingRemediationInstructionExecution({ id: this.executionInstruction._id });
      } catch (err) {
        this.$modals.hide();
        this.$popups.error({
          text: this.$t('remediationInstructionExecute.popups.connectionError'),
          autoClose: false,
        });
      }
    },

    async createInstructionExecution() {
      const { _id: instructionId } = this.config.assignedInstruction;

      const instructionExecution = await this.createRemediationInstructionExecution({
        data: {
          alarm: this.config.alarm._id,
          instruction: instructionId,
        },
      });

      this.executionInstructionId = instructionExecution._id;
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

        if (this.config.onOpen) {
          await this.config.onOpen();
        }
      } catch (err) {
        this.$popups.error({ text: err.error || this.$t('errors.default') });
        this.$modals.hide();
      }
    },

    close() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          hideTitle: true,
          text: this.$t('remediationInstructionExecute.closeConfirmationText'),
          action: async () => {
            await this.pauseRemediationInstructionExecution({ id: this.executionInstruction._id });

            if (this.config.onClose) {
              await this.config.onClose();
            }

            this.$modals.hide();
          },
          cancel: async (cancelled) => {
            if (!cancelled) {
              return;
            }

            await this.cancelRemediationInstructionExecution({ id: this.executionInstruction._id });

            if (this.config.onClose) {
              await this.config.onClose();
            }

            this.$modals.hide();
          },
        },
      });
    },
  },
};
</script>
