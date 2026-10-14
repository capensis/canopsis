import { mapActions } from 'pinia';

import { SOCKET_ROOMS } from '@/config';
import {
  MAX_LIMIT,
  REMEDIATION_INSTRUCTION_EXECUTION_STATUSES,
  USER_PERMISSIONS,
} from '@/constants';

import { useRemediationInstructionExecutionStore } from '@/stores/modules/remediation-instruction-execution';

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
  [REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.waitingResult]: {
    messageKey: 'remediation.instructionExecute.popups.wasWaitingResult',
    type: 'info',
  },
};

export const appRemediationInstructionExecutionsPopupsMixin = {
  beforeUnmount() {
    this.leaveFromSimpleManualExecutions();
  },

  computed: {
    hasExecuteInstructionAccess() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.executeInstruction);
    },
  },

  methods: {
    ...mapActions(useRemediationInstructionExecutionStore, {
      fetchPausedExecutionsWithoutStore: 'fetchPausedListWithoutStore',
    }),

    runExecutionsPopups() {
      this.joinToSimpleManualExecutions();
      this.showPausedExecutionsPopup();
    },

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

      this.$popups[popupConfig.type]({
        text,
        context: { execution },
        autoClose: false,
      });
    },

    joinToSimpleManualExecutions() {
      if (!this.hasExecuteInstructionAccess) {
        return;
      }

      this.$socket
        .join(SOCKET_ROOMS.simplifiedManualExecutions)
        .addListener(this.showInstructionPopup);
    },

    leaveFromSimpleManualExecutions() {
      if (!this.hasExecuteInstructionAccess) {
        return;
      }

      this.$socket
        .leave(SOCKET_ROOMS.simplifiedManualExecutions)
        .removeListener(this.showInstructionPopup);
    },

    async showPausedExecutionsPopup() {
      if (!this.hasExecuteInstructionAccess) {
        return;
      }

      const pausedExecutions = await this.fetchPausedExecutionsWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      if (!pausedExecutions?.length) {
        return;
      }

      pausedExecutions.forEach(this.showInstructionPopup);
    },
  },
};

export default appRemediationInstructionExecutionsPopupsMixin;
