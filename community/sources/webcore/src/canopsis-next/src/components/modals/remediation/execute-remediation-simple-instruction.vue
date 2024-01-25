<template>
  <modal-wrapper
    :close="closeModal"
    minimize
  >
    <template #title="">
      <span>{{ config.assignedInstruction.name }}</span>
    </template>
    <template #text="">
      <v-fade-transition>
        <v-layout
          v-if="pending"
          justify-center
        >
          <v-progress-circular
            color="primary"
            indeterminate
          />
        </v-layout>
        <remediation-instruction-simple-execute
          v-else
          :executed="executed"
          :jobs="jobs"
          @run:jobs="runJobs"
        />
      </v-fade-transition>
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="closeModal"
      >
        {{ $t('common.close') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { SOCKET_ROOMS } from '@/config';
import { MODALS, REMEDIATION_INSTRUCTION_EXECUTION_STATUSES } from '@/constants';

import Socket from '@/plugins/socket/services/socket';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesRemediationJobExecutionMixin } from '@/mixins/entities/remediation/job-execution';
import { entitiesRemediationInstructionMixin } from '@/mixins/entities/remediation/instruction';
import { entitiesRemediationInstructionExecutionMixin } from '@/mixins/entities/remediation/instruction-execution';

import RemediationInstructionSimpleExecute from '@/components/other/remediation/instruction-execute/remediation-instruction-simple-execute.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.executeRemediationSimpleInstruction,
  components: {
    RemediationInstructionSimpleExecute,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    entitiesRemediationJobExecutionMixin,
    entitiesRemediationInstructionExecutionMixin,
    entitiesRemediationInstructionMixin,
  ],
  data() {
    return {
      pending: true,
      executed: false,
      instruction: null,
      instructionExecution: null,
    };
  },
  computed: {
    instructionId() {
      return this.config.assignedInstruction?._id;
    },

    instructionExecutionId() {
      const { execution } = this.config.assignedInstruction;

      return execution?._id ?? this.instructionExecution?._id;
    },

    instructionJobs() {
      return this.instruction?.jobs?.map(({ job }) => job);
    },

    jobs() {
      return this.instructionExecution?.jobs ?? this.instructionJobs;
    },

    socketRoomName() {
      return `${SOCKET_ROOMS.execution}/${this.instructionExecutionId}`;
    },
  },
  async mounted() {
    this.fetchInstruction();

    if (this.config.assignedInstruction.execution) {
      await this.fetchInstructionExecution();

      this.joinToSocketRoom();
    }
  },
  beforeDestroy() {
    this.leaveFromSocketRoom();
  },
  methods: {
    async fetchInstruction() {
      this.pending = true;

      this.instruction = await this.fetchRemediationInstructionWithoutStore({ id: this.instructionId });

      this.pending = false;
    },

    async runJobs() {
      await this.createInstructionExecution();

      this.joinToSocketRoom();
    },

    setJobs(jobs) {
      this.instructionExecution.jobs = jobs;
    },

    /**
     * Join from execution room
     */
    joinToSocketRoom() {
      this.$socket
        .on(Socket.EVENTS_TYPES.customClose, this.socketCloseHandler)
        .on(Socket.EVENTS_TYPES.closeRoom, this.socketCloseRoomHandler)
        .join(this.socketRoomName)
        .addListener(this.setJobs);
    },

    /**
     * Leave from execution room
     */
    leaveFromSocketRoom() {
      this.$socket
        .off(Socket.EVENTS_TYPES.customClose, this.socketCloseHandler)
        .off(Socket.EVENTS_TYPES.closeRoom, this.socketCloseRoomHandler)
        .leave(this.socketRoomName)
        .removeListener(this.setJobs);
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
     * Fetch instruction execution method (resume or fetch if exists)
     *
     * @return {Promise<void>}
     */
    async fetchInstructionExecution() {
      const { execution } = this.config.assignedInstruction;

      try {
        if (execution) {
          if (execution.status === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.paused) {
            this.instructionExecution = await this.resumeRemediationInstructionExecution({
              id: this.instructionExecutionId,
            });
          } else {
            this.instructionExecution = await this.fetchRemediationInstructionExecutionWithoutStore({
              id: this.instructionExecutionId,
            });
          }

          this.executed = true;
        }
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: err.error || this.$t('errors.default') });

        this.closeModal();
      }
    },

    /**
     * Create instruction execution method
     *
     * @return {Promise<void>}
     */
    async createInstructionExecution() {
      try {
        this.instructionExecution = await this.createRemediationInstructionExecution({
          data: {
            alarm: this.config.alarmId,
            instruction: this.instructionId,
          },
        });

        this.executed = true;

        if (this.config.onExecute) {
          await this.config.onExecute();
        }
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: err.error || this.$t('errors.default') });

        this.closeModal();
      }
    },

    closeModal() {
      if (this.config.onClose) {
        this.config.onClose();
      }

      this.$modals.hide();
    },
  },
};
</script>
