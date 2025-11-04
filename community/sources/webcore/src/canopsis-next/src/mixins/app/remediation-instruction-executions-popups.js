import { createNamespacedHelpers } from 'vuex';

import { SOCKET_ROOMS } from '@/config';
import { MAX_LIMIT, REMEDIATION_INSTRUCTION_EXECUTION_STATUSES, USER_PERMISSIONS } from '@/constants';

const EXECUTION_STATUSES_TO_POPUPS = {
  [REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.completed]: {
    messageKey: 'remediation.instructionExecute.popups.wasFinished',
    type: 'success',
  },
  [REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.failed]: {
    messageKey: 'remediation.instructionExecute.popups.wasFailed',
    type: 'error',
  },
  [REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.paused]: {
    messageKey: 'remediation.instructionExecute.popups.wasPaused',
    type: 'info',
  },
};

const { mapActions } = createNamespacedHelpers('remediationInstructionExecution');

export const appRemediationInstructionExecutionsPopupsMixin = {
  beforeDestroy() {
    this.leaveFromSimpleManualExecutions();
  },
  computed: {
    hasExecuteInstructionAccess() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.executeInstruction);
    },
  },
  methods: {
    ...mapActions({
      fetchPausedExecutionsWithoutStore: 'fetchPausedListWithoutStore',
    }),

    /**
     * Runs execution popups by joining manual executions
     * and showing popups for paused executions
     */
    runExecutionsPopups() {
      this.joinToSimpleManualExecutions();
      this.showPausedExecutionsPopup();
    },

    /**
     * Shows an instruction popup based on the execution status
     *
     * @param {Object} execution - The instruction execution object
     * @param {string} execution.status - The execution status
     * @param {string} execution.name - The instruction name
     */
    showInstructionPopup(execution) {
      const popupConfig = EXECUTION_STATUSES_TO_POPUPS[execution.status];

      if (!popupConfig) {
        console.warn(`Unknown execution status for popup: ${execution.status}`);

        return;
      }

      const prefix = this.$t(popupConfig.messageKey, {
        instructionName: execution.name,
      });

      const text = `${prefix} <c-remediation-instruction-execution-see-details :execution="execution" />`;

      this.$popups[popupConfig.type]({ text, context: { execution }, autoClose: false });
    },

    /**
     * Joins the socket room for simplified manual executions
     * and adds a listener to show instruction popups
     */
    joinToSimpleManualExecutions() {
      if (!this.hasExecuteInstructionAccess) {
        return;
      }

      this.$socket.join(SOCKET_ROOMS.simplifiedManualExecutions)
        .addListener(this.showInstructionPopup);
    },

    /**
     * Leaves the socket room for simplified manual executions
     * and removes the instruction popup listener
     */
    leaveFromSimpleManualExecutions() {
      if (!this.hasExecuteInstructionAccess) {
        return;
      }

      this.$socket.leave(SOCKET_ROOMS.simplifiedManualExecutions)
        .removeListener(this.showInstructionPopup);
    },

    /**
     * Fetches paused executions and shows a popup for each one
     */
    async showPausedExecutionsPopup() {
      const pausedExecutions = await this.fetchPausedExecutionsWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      if (!pausedExecutions || !pausedExecutions.length) {
        return;
      }

      pausedExecutions.forEach(this.showInstructionPopup);
    },
  },
};
