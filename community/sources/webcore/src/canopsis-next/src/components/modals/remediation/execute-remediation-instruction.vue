<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(:close="close", minimize)
      template(#title="")
        span {{ config.assignedInstruction.name }}
      template(#text="")
        v-fade-transition
          remediation-instruction-execute(
            v-if="instructionExecution",
            :instruction-execution="instructionExecution",
            :next-pending="nextPending",
            :previous-pending="previousPending",
            @next-step="nextStep",
            @next-operation="nextOperation",
            @previous-operation="previousOperation",
            @execute-job="executeJob",
            @cancel-job-execution="cancelJobExecution"
          )
          v-layout(v-else, justify-center)
            v-progress-circular(color="primary", indeterminate)
</template>

<script>
import { pick } from 'lodash';

import { SOCKET_ROOMS } from '@/config';
import { MODALS, RESPONSE_STATUSES } from '@/constants';

import Socket from '@/plugins/socket/services/socket';

import { getEmptyRemediationJobExecution } from '@/helpers/entities/remediation/job/form';
import {
  isInstructionExecutionAborted,
  isInstructionExecutionCompleted,
  isInstructionExecutionFailed,
  isInstructionExecutionPaused,
  isInstructionExecutionRunning,
} from '@/helpers/entities/remediation/instruction-execution/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesRemediationJobExecutionMixin } from '@/mixins/entities/remediation/job-execution';
import { entitiesRemediationInstructionExecutionMixin } from '@/mixins/entities/remediation/instruction-execution';

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
    entitiesRemediationJobExecutionMixin,
    entitiesRemediationInstructionExecutionMixin,
  ],
  data() {
    return {
      pending: true,
      nextPending: false,
      previousPending: false,
      instructionExecution: null,
    };
  },
  computed: {
    instruction() {
      return this.config.assignedInstruction;
    },

    instructionId() {
      return this.instruction?._id;
    },

    instructionExecutionId() {
      const { execution } = this.instruction;

      return execution?._id ?? this.instructionExecution?._id;
    },

    socketRoomName() {
      return `${SOCKET_ROOMS.execution}/${this.instructionExecutionId}`;
    },
  },
  watch: {
    async instructionExecution(instructionExecution) {
      if (!isInstructionExecutionRunning(instructionExecution)) {
        const isFailedExecution = isInstructionExecutionFailed(instructionExecution)
          || isInstructionExecutionAborted(instructionExecution);

        const type = isFailedExecution ? 'failed' : 'success';
        const text = this.$t(`remediation.instructionExecute.popups.${type}`, {
          instructionName: instructionExecution.name,
        });

        if (isFailedExecution) {
          this.$popups.error({ text });
        } else {
          this.$popups.success({ text });
        }

        if (this.config.onComplete) {
          await this.config.onComplete(instructionExecution);
        }

        await this.closeModal();
      }
    },
  },
  async mounted() {
    await this.fetchInstructionExecution();

    this.joinToSocketRoom();
  },
  beforeDestroy() {
    this.leaveFromSocketRoom();
  },
  methods: {
    /**
     * Join from execution room
     */
    joinToSocketRoom() {
      if (
        !this.instructionExecutionId
        || isInstructionExecutionFailed(this.instructionExecution)
        || isInstructionExecutionAborted(this.instructionExecution)
        || isInstructionExecutionCompleted(this.instructionExecution)
      ) {
        return;
      }

      this.$socket
        .on(Socket.EVENTS_TYPES.customClose, this.socketCloseHandler)
        .on(Socket.EVENTS_TYPES.closeRoom, this.socketCloseRoomHandler)
        .join(this.socketRoomName)
        .addListener(this.setOperation);
    },

    /**
     * Leave from execution room
     */
    leaveFromSocketRoom() {
      if (this.instructionExecutionId) {
        this.$socket
          .off(Socket.EVENTS_TYPES.customClose, this.socketCloseHandler)
          .off(Socket.EVENTS_TYPES.closeRoom, this.socketCloseRoomHandler)
          .leave(this.socketRoomName)
          .removeListener(this.setOperation);
      }
    },

    async tryToCallWithHandleNotFound(func) {
      try {
        await func();
      } catch (err) {
        const message = err?.status === RESPONSE_STATUSES.notFound
          ? this.$t('remediation.instructionExecute.popups.wasRemovedOrDisabled', {
            instructionName: this.instructionExecution?.name,
          })
          : err?.data?.error;

        console.error(err);

        this.$popups.error({ text: message || this.$t('errors.default') });

        this.closeModal();
      }
    },

    async callFunctionWithPreviousPending(func) {
      this.previousPending = true;

      await this.tryToCallWithHandleNotFound(func);

      this.previousPending = false;
    },

    async callFunctionWithNextPending(func) {
      this.nextPending = true;

      await this.tryToCallWithHandleNotFound(func);

      this.nextPending = false;
    },

    /**
     * Goto next step
     *
     * @param {boolean} success
     * @return {Promise<void>}
     */
    nextStep(success = false) {
      return this.callFunctionWithNextPending(async () => {
        this.instructionExecution = await this.nextStepRemediationInstructionExecution({
          id: this.instructionExecutionId,
          data: { failed: !success },
        });
      });
    },

    /**
     * Goto next operation
     *
     * @return {Promise<void>}
     */
    async nextOperation() {
      return this.callFunctionWithNextPending(async () => {
        this.instructionExecution = await this.nextOperationRemediationInstructionExecution({
          id: this.instructionExecutionId,
        });
      });
    },

    /**
     * Goto previous operation
     *
     * @return {Promise<void>}
     */
    previousOperation() {
      return this.callFunctionWithPreviousPending(async () => {
        this.instructionExecution = await this.previousOperationRemediationInstructionExecution({
          id: this.instructionExecutionId,
        });
      });
    },

    /**
     * Create remediation instruction execution with error handler
     *
     * @return {Promise<void>}
     */
    createExecution() {
      return this.tryToCallWithHandleNotFound(async () => {
        this.instructionExecution = await this.createRemediationInstructionExecution({
          data: {
            alarm: this.config.alarmId,
            instruction: this.instructionId,
          },
        });
      });
    },

    /**
     * Fetch remediation instruction execution with error handler
     *
     * @return {Promise<void>}
     */
    fetchExecution() {
      return this.tryToCallWithHandleNotFound(async () => {
        this.instructionExecution = await this.fetchRemediationInstructionExecutionWithoutStore({
          id: this.instructionExecutionId,
        });
      });
    },

    /**
     * Resume remediation instruction execution with error handler
     *
     * @return {Promise<void>}
     */
    resumeExecution() {
      return this.tryToCallWithHandleNotFound(async () => {
        this.instructionExecution = await this.resumeRemediationInstructionExecution({
          id: this.instructionExecutionId,
        });
      });
    },

    /**
     * Cancel remediation instruction execution with error handler
     *
     * @return {Promise<void>}
     */
    cancelExecution() {
      return this.tryToCallWithHandleNotFound(
        () => this.cancelRemediationInstructionExecution({ id: this.instructionExecutionId }),
      );
    },

    /**
     * Pause remediation instruction execution with error handler
     *
     * @return {Promise<void>}
     */
    pauseExecution() {
      return this.tryToCallWithHandleNotFound(
        () => this.pauseRemediationInstructionExecution({ id: this.instructionExecutionId }),
      );
    },

    /**
     * Execute special job by operation
     *
     * @param {RemediationJobExecution} job
     * @param {RemediationInstructionStepOperation} operation
     * @return {Promise<void>}
     */
    async executeJob({ job, operation }) {
      try {
        const updatedJob = await this.createRemediationJobExecution({
          data: {
            execution: this.instructionExecutionId,
            job: job.job_id,
            operation: operation.operation_id,
          },
        });

        this.setJob(updatedJob, operation);
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: err.error || this.$t('errors.default') });
      }
    },

    /**
     * Cancel special job by operation
     *
     * @param {RemediationJobExecution} job
     * @param {RemediationInstructionStepOperation} operation
     * @return {Promise<void>}
     */
    async cancelJobExecution({ job, operation }) {
      try {
        await this.cancelRemediationJobExecution({ id: job._id });

        const updatedJob = {
          ...getEmptyRemediationJobExecution(),
          ...pick(job, ['_id', 'job_id', 'name', 'payload', 'query']),
        };

        this.setJob(updatedJob, operation);
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: err.error || this.$t('errors.default') });
      }
    },

    /**
     * Socket customClose event handler (we need to use for connection checking)
     */
    socketCloseHandler() {
      if (!this.$socket.isConnectionOpen) {
        this.closeModal();
        this.$popups.error({
          text: this.$t('remediation.instructionExecute.popups.connectionError'),
          autoClose: false,
        });
      }
    },

    /**
     * Socket closeRoom event handler
     */
    socketCloseRoomHandler() {
      this.closeModal();
      this.$popups.error({
        text: this.$t('remediation.instructionExecute.popups.wasAborted', {
          instructionName: this.instructionExecution?.name,
        }),
        autoClose: false,
      });
    },

    /**
     * Set job into special operation into current instructionExecution
     *
     * @param {RemediationJobExecution} job
     * @param {RemediationInstructionStepOperation} operation
     */
    setJob(job, operation) {
      this.setOperation({
        ...operation,

        jobs: operation.jobs.map(operationJob => (
          operationJob.job_id === job.job_id
            ? job
            : operationJob
        )),
      });
    },

    /**
     * Set operation into current instructionExecution
     *
     * @param {RemediationInstructionStepOperation} operation
     */
    setOperation(operation) {
      for (const step of this.instructionExecution.steps) {
        const operationIndex = step.operations
          .findIndex(({ operation_id: operationId }) => operationId === operation.operation_id);

        if (operationIndex !== -1) {
          this.$set(step.operations, operationIndex, operation);

          return;
        }
      }
    },

    /**
     * Confirmation modal hide method
     *
     * @return {Promise<void>}
     */
    async closeModal() {
      if (this.config.onClose) {
        await this.config.onClose();
      }

      this.$modals.hide();
    },

    /**
     * Fetch instruction execution method (create if not exists, resume or fetch if exists)
     *
     * @return {Promise<void>}
     */
    async fetchInstructionExecution() {
      const { execution } = this.config.assignedInstruction;

      this.pending = true;

      if (!execution) {
        await this.createExecution();
      } else if (isInstructionExecutionPaused(execution)) {
        await this.resumeExecution();
      } else {
        await this.fetchExecution();
      }

      if (this.config.onExecute) {
        await this.config.onExecute();
      }

      this.pending = false;
    },

    /**
     * Close handler
     */
    close() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          hideTitle: true,
          text: this.$t('remediation.instructionExecute.closeConfirmationText'),
          action: async () => {
            await this.pauseExecution();

            await this.closeModal();
          },
          cancel: async (cancelled) => {
            if (!cancelled) {
              return;
            }

            await this.cancelExecution();

            await this.closeModal();
          },
        },
      });
    },
  },
};
</script>
