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
