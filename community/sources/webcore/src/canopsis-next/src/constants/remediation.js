import { TRIGGERS_TYPES } from './common';

export const REMEDIATION_TABS = {
  instructions: 'instructions',
  configurations: 'configurations',
  jobs: 'jobs',
  statistics: 'statistics',
};

export const REMEDIATION_INSTRUCTION_TYPES = {
  manual: 0,
  auto: 1,
  simpleManual: 2,
};

export const REMEDIATION_INSTRUCTION_STATUSES = {
  approved: 0,
  createdAndApproveRequested: 1,
  updatedAndApproveRequested: 2,
  createdAndDismissed: 3,
  updatedAndDismissed: 4,
};

export const REMEDIATION_INSTRUCTION_APPROVAL_TYPES = {
  role: 0,
  user: 1,
};

export const REMEDIATION_INSTRUCTION_EXECUTION_STATUSES = {
  running: 0,
  paused: 1,
  completed: 2,
  aborted: 3,
  failed: 4,
};

export const REMEDIATION_JOB_EXECUTION_STATUSES = {
  running: 0,
  succeeded: 1,
  failed: 2,
  canceled: 3,
};

export const REMEDIATION_CONFIGURATION_JOBS_AUTH_TYPES_WITH_USERNAME = ['basic-auth'];

export const REMEDIATION_STATISTICS_CHART_DATA_TYPE = {
  percent: 'percent',
  remediation: 'remediation',
};

export const REMEDIATION_STATISTICS_BAR_PERCENTAGE = 0.5;

export const INSTRUCTION_EXECUTION_ICONS = {
  manualInProgress: 1,
  autoInProgress: 2,
  autoFailed: 3,
  manualFailed: 4,
  manualFailedWithInProgress: 5,
  autoFailedWithInProgress: 6,
  manualFailedWithManualAvailable: 7,
  autoFailedWithManualAvailable: 8,
  manualAvailable: 9,
  autoSuccessful: 10,
  manualSuccessful: 11,
  autoSuccessfulWithInProgress: 12,
  manualSuccessfulWithInProgress: 13,
  autoSuccessfulWithManualAvailable: 14,
  manualSuccessfulWithManualAvailable: 15,
};

export const REMEDIATION_AUTO_INSTRUCTION_TRIGGERS_TYPES = [
  TRIGGERS_TYPES.create,
  TRIGGERS_TYPES.statedec,
  TRIGGERS_TYPES.stateinc,
  TRIGGERS_TYPES.changestate,
  TRIGGERS_TYPES.unsnooze,
  TRIGGERS_TYPES.activate,
  TRIGGERS_TYPES.pbhenter,
  TRIGGERS_TYPES.pbhleave,
  TRIGGERS_TYPES.eventscount,
];

export const REMEDIATION_INSTRUCTION_EXECUTION_STEP_TYPES = {
  job: 0,
  manual: 1,
};

export const REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES = {
  completed: 1,
  failed: 2,
  aborted: 3,
  skipped: 5,
};
