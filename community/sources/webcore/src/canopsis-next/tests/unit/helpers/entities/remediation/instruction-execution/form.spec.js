import { INSTRUCTION_EXECUTION_ICONS, REMEDIATION_INSTRUCTION_EXECUTION_STATUSES } from '@/constants';

import {
  hasInstructionWithoutAnyExecution,
  isInstructionExecutionIconManualInProgress,
  isInstructionExecutionIconAutoInProgress,
  isInstructionExecutionIconInProgress,
  isInstructionExecutionIconFailed,
  isInstructionExecutionIconSuccess,
  isInstructionExecutionManual,
  isInstructionExecutionCompleted,
  isInstructionExecutionFailed,
  isInstructionExecutionAborted,
  isInstructionExecutionPaused,
  isInstructionExecutionRunning,
} from '@/helpers/entities/remediation/instruction-execution/form';

describe('helpers/entities/remediation/instruction-execution/form', () => {
  describe('hasInstructionWithoutAnyExecution', () => {
    it('Returns true for withoutAnyExecution icon', () => {
      expect(hasInstructionWithoutAnyExecution(INSTRUCTION_EXECUTION_ICONS.withoutAnyExecution)).toBe(true);
    });

    it('Returns false for other icons', () => {
      expect(hasInstructionWithoutAnyExecution(INSTRUCTION_EXECUTION_ICONS.manualInProgress)).toBe(false);
      expect(hasInstructionWithoutAnyExecution(INSTRUCTION_EXECUTION_ICONS.autoSuccessful)).toBe(false);
      expect(hasInstructionWithoutAnyExecution(INSTRUCTION_EXECUTION_ICONS.manualAvailable)).toBe(false);
    });
  });

  describe('isInstructionExecutionIconManualInProgress', () => {
    it('Returns true for manual in progress icons', () => {
      expect(isInstructionExecutionIconManualInProgress(INSTRUCTION_EXECUTION_ICONS.manualInProgress)).toBe(true);
      expect(isInstructionExecutionIconManualInProgress(INSTRUCTION_EXECUTION_ICONS.manualFailedWithInProgress))
        .toBe(true);
      expect(isInstructionExecutionIconManualInProgress(INSTRUCTION_EXECUTION_ICONS.manualSuccessfulWithInProgress))
        .toBe(true);
    });

    it('Returns false for non-manual in progress icons', () => {
      expect(isInstructionExecutionIconManualInProgress(INSTRUCTION_EXECUTION_ICONS.autoInProgress)).toBe(false);
      expect(isInstructionExecutionIconManualInProgress(INSTRUCTION_EXECUTION_ICONS.manualFailed)).toBe(false);
      expect(isInstructionExecutionIconManualInProgress(INSTRUCTION_EXECUTION_ICONS.manualSuccessful)).toBe(false);
      expect(isInstructionExecutionIconManualInProgress(INSTRUCTION_EXECUTION_ICONS.withoutAnyExecution)).toBe(false);
    });
  });

  describe('isInstructionExecutionIconAutoInProgress', () => {
    it('Returns true for auto in progress icons', () => {
      expect(isInstructionExecutionIconAutoInProgress(INSTRUCTION_EXECUTION_ICONS.autoInProgress)).toBe(true);
      expect(isInstructionExecutionIconAutoInProgress(INSTRUCTION_EXECUTION_ICONS.autoFailedWithInProgress)).toBe(true);
      expect(isInstructionExecutionIconAutoInProgress(INSTRUCTION_EXECUTION_ICONS.autoSuccessfulWithInProgress))
        .toBe(true);
    });

    it('Returns false for non-auto in progress icons', () => {
      expect(isInstructionExecutionIconAutoInProgress(INSTRUCTION_EXECUTION_ICONS.manualInProgress)).toBe(false);
      expect(isInstructionExecutionIconAutoInProgress(INSTRUCTION_EXECUTION_ICONS.autoFailed)).toBe(false);
      expect(isInstructionExecutionIconAutoInProgress(INSTRUCTION_EXECUTION_ICONS.autoSuccessful)).toBe(false);
      expect(isInstructionExecutionIconAutoInProgress(INSTRUCTION_EXECUTION_ICONS.withoutAnyExecution)).toBe(false);
    });
  });

  describe('isInstructionExecutionIconInProgress', () => {
    it('Returns true for any in progress icons', () => {
      expect(isInstructionExecutionIconInProgress(INSTRUCTION_EXECUTION_ICONS.manualInProgress)).toBe(true);
      expect(isInstructionExecutionIconInProgress(INSTRUCTION_EXECUTION_ICONS.autoInProgress)).toBe(true);
      expect(isInstructionExecutionIconInProgress(INSTRUCTION_EXECUTION_ICONS.manualFailedWithInProgress)).toBe(true);
      expect(isInstructionExecutionIconInProgress(INSTRUCTION_EXECUTION_ICONS.autoFailedWithInProgress)).toBe(true);
      expect(isInstructionExecutionIconInProgress(INSTRUCTION_EXECUTION_ICONS.manualSuccessfulWithInProgress))
        .toBe(true);
      expect(isInstructionExecutionIconInProgress(INSTRUCTION_EXECUTION_ICONS.autoSuccessfulWithInProgress)).toBe(true);
    });

    it('Returns false for non-in progress icons', () => {
      expect(isInstructionExecutionIconInProgress(INSTRUCTION_EXECUTION_ICONS.manualFailed)).toBe(false);
      expect(isInstructionExecutionIconInProgress(INSTRUCTION_EXECUTION_ICONS.autoFailed)).toBe(false);
      expect(isInstructionExecutionIconInProgress(INSTRUCTION_EXECUTION_ICONS.manualSuccessful)).toBe(false);
      expect(isInstructionExecutionIconInProgress(INSTRUCTION_EXECUTION_ICONS.autoSuccessful)).toBe(false);
      expect(isInstructionExecutionIconInProgress(INSTRUCTION_EXECUTION_ICONS.withoutAnyExecution)).toBe(false);
    });
  });

  describe('isInstructionExecutionIconFailed', () => {
    it('Returns true for failed icons', () => {
      expect(isInstructionExecutionIconFailed(INSTRUCTION_EXECUTION_ICONS.autoFailed)).toBe(true);
      expect(isInstructionExecutionIconFailed(INSTRUCTION_EXECUTION_ICONS.manualFailed)).toBe(true);
      expect(isInstructionExecutionIconFailed(INSTRUCTION_EXECUTION_ICONS.manualFailedWithInProgress)).toBe(true);
      expect(isInstructionExecutionIconFailed(INSTRUCTION_EXECUTION_ICONS.autoFailedWithInProgress)).toBe(true);
      expect(isInstructionExecutionIconFailed(INSTRUCTION_EXECUTION_ICONS.autoFailedWithManualAvailable)).toBe(true);
      expect(isInstructionExecutionIconFailed(INSTRUCTION_EXECUTION_ICONS.manualFailedWithManualAvailable)).toBe(true);
    });

    it('Returns false for non-failed icons', () => {
      expect(isInstructionExecutionIconFailed(INSTRUCTION_EXECUTION_ICONS.manualInProgress)).toBe(false);
      expect(isInstructionExecutionIconFailed(INSTRUCTION_EXECUTION_ICONS.autoInProgress)).toBe(false);
      expect(isInstructionExecutionIconFailed(INSTRUCTION_EXECUTION_ICONS.manualSuccessful)).toBe(false);
      expect(isInstructionExecutionIconFailed(INSTRUCTION_EXECUTION_ICONS.autoSuccessful)).toBe(false);
      expect(isInstructionExecutionIconFailed(INSTRUCTION_EXECUTION_ICONS.withoutAnyExecution)).toBe(false);
    });
  });

  describe('isInstructionExecutionIconSuccess', () => {
    it('Returns true for successful icons', () => {
      expect(isInstructionExecutionIconSuccess(INSTRUCTION_EXECUTION_ICONS.autoSuccessful)).toBe(true);
      expect(isInstructionExecutionIconSuccess(INSTRUCTION_EXECUTION_ICONS.manualSuccessful)).toBe(true);
      expect(isInstructionExecutionIconSuccess(INSTRUCTION_EXECUTION_ICONS.manualSuccessfulWithInProgress)).toBe(true);
      expect(isInstructionExecutionIconSuccess(INSTRUCTION_EXECUTION_ICONS.autoSuccessfulWithInProgress)).toBe(true);
      expect(isInstructionExecutionIconSuccess(INSTRUCTION_EXECUTION_ICONS.autoSuccessfulWithManualAvailable))
        .toBe(true);
      expect(isInstructionExecutionIconSuccess(INSTRUCTION_EXECUTION_ICONS.manualSuccessfulWithManualAvailable))
        .toBe(true);
    });

    it('Returns false for non-successful icons', () => {
      expect(isInstructionExecutionIconSuccess(INSTRUCTION_EXECUTION_ICONS.manualInProgress)).toBe(false);
      expect(isInstructionExecutionIconSuccess(INSTRUCTION_EXECUTION_ICONS.autoInProgress)).toBe(false);
      expect(isInstructionExecutionIconSuccess(INSTRUCTION_EXECUTION_ICONS.manualFailed)).toBe(false);
      expect(isInstructionExecutionIconSuccess(INSTRUCTION_EXECUTION_ICONS.autoFailed)).toBe(false);
      expect(isInstructionExecutionIconSuccess(INSTRUCTION_EXECUTION_ICONS.withoutAnyExecution)).toBe(false);
    });
  });

  describe('isInstructionExecutionManual', () => {
    it('Returns true for manual icons', () => {
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.manualInProgress)).toBe(true);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.manualFailed)).toBe(true);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.manualFailedWithInProgress)).toBe(true);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.manualFailedWithManualAvailable)).toBe(true);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.manualAvailable)).toBe(true);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.manualSuccessful)).toBe(true);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.manualSuccessfulWithInProgress)).toBe(true);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.manualSuccessfulWithManualAvailable)).toBe(true);
    });

    it('Returns false for non-manual icons', () => {
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.autoInProgress)).toBe(false);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.autoFailed)).toBe(false);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.autoSuccessful)).toBe(false);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.autoFailedWithInProgress)).toBe(false);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.autoFailedWithManualAvailable)).toBe(false);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.autoSuccessfulWithInProgress)).toBe(false);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.autoSuccessfulWithManualAvailable)).toBe(false);
      expect(isInstructionExecutionManual(INSTRUCTION_EXECUTION_ICONS.withoutAnyExecution)).toBe(false);
    });
  });

  describe('isInstructionExecutionCompleted', () => {
    it('Returns true for completed status', () => {
      expect(isInstructionExecutionCompleted({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.completed,
      })).toBe(true);
    });

    it('Returns false for non-completed statuses', () => {
      expect(isInstructionExecutionCompleted({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running,
      })).toBe(false);
      expect(isInstructionExecutionCompleted({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.failed,
      })).toBe(false);
      expect(isInstructionExecutionCompleted({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.aborted,
      })).toBe(false);
      expect(isInstructionExecutionCompleted({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.paused,
      })).toBe(false);
    });

    it('Returns false for undefined object', () => {
      expect(isInstructionExecutionCompleted()).toBe(false);
      expect(isInstructionExecutionCompleted({})).toBe(false);
    });
  });

  describe('isInstructionExecutionFailed', () => {
    it('Returns true for failed status', () => {
      expect(isInstructionExecutionFailed({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.failed,
      })).toBe(true);
    });

    it('Returns false for non-failed statuses', () => {
      expect(isInstructionExecutionFailed({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running,
      })).toBe(false);
      expect(isInstructionExecutionFailed({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.completed,
      })).toBe(false);
      expect(isInstructionExecutionFailed({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.aborted,
      })).toBe(false);
      expect(isInstructionExecutionFailed({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.paused,
      })).toBe(false);
    });

    it('Returns false for undefined object', () => {
      expect(isInstructionExecutionFailed()).toBe(false);
      expect(isInstructionExecutionFailed({})).toBe(false);
    });
  });

  describe('isInstructionExecutionAborted', () => {
    it('Returns true for aborted status', () => {
      expect(isInstructionExecutionAborted({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.aborted,
      })).toBe(true);
    });

    it('Returns false for non-aborted statuses', () => {
      expect(isInstructionExecutionAborted({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running,
      })).toBe(false);
      expect(isInstructionExecutionAborted({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.completed,
      })).toBe(false);
      expect(isInstructionExecutionAborted({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.failed,
      })).toBe(false);
      expect(isInstructionExecutionAborted({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.paused,
      })).toBe(false);
    });

    it('Returns false for undefined object', () => {
      expect(isInstructionExecutionAborted()).toBe(false);
      expect(isInstructionExecutionAborted({})).toBe(false);
    });
  });

  describe('isInstructionExecutionPaused', () => {
    it('Returns true for paused status', () => {
      expect(isInstructionExecutionPaused({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.paused,
      })).toBe(true);
    });

    it('Returns false for non-paused statuses', () => {
      expect(isInstructionExecutionPaused({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running,
      })).toBe(false);
      expect(isInstructionExecutionPaused({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.completed,
      })).toBe(false);
      expect(isInstructionExecutionPaused({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.failed,
      })).toBe(false);
      expect(isInstructionExecutionPaused({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.aborted,
      })).toBe(false);
    });

    it('Returns false for undefined object', () => {
      expect(isInstructionExecutionPaused()).toBe(false);
      expect(isInstructionExecutionPaused({})).toBe(false);
    });
  });

  describe('isInstructionExecutionRunning', () => {
    it('Returns true for running status', () => {
      expect(isInstructionExecutionRunning({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running,
      })).toBe(true);
    });

    it('Returns false for non-running statuses', () => {
      expect(isInstructionExecutionRunning({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.paused,
      })).toBe(false);
      expect(isInstructionExecutionRunning({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.completed,
      })).toBe(false);
      expect(isInstructionExecutionRunning({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.failed,
      })).toBe(false);
      expect(isInstructionExecutionRunning({
        status: REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.aborted,
      })).toBe(false);
    });

    it('Returns false for undefined object', () => {
      expect(isInstructionExecutionRunning()).toBe(false);
      expect(isInstructionExecutionRunning({})).toBe(false);
    });
  });
});
