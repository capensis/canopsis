import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

export default {
  tabs: {
    configurations: 'Configurations',
    jobs: 'Jobs',
    statistics: 'Remediation statistics',
  },

  instruction: {
    usingInstruction: 'Cannot be deleted since it is in use',
    addStep: 'Add step',
    addOperation: 'Add operation',
    addJob: 'Add job',
    endpoint: 'Endpoint',
    job: 'Job | Jobs',
    listJobs: 'List of jobs',
    endpointAvatar: 'EP',
    workflow: 'Workflow if this step fails:',
    jobWorkflow: 'Workflow if this job fails:',
    remainingStep: 'Continue with remaining steps',
    remainingJob: 'Continue with remaining job',
    timeToComplete: 'Time to complete',
    emptySteps: 'No steps added yet',
    emptyOperations: 'No operations added yet',
    emptyJobs: 'No jobs added yet',
    timeoutAfterExecution: 'Timeout after instruction execution',
    requestApproval: 'Request for approval',
    type: 'Instruction type',
    approvalPending: 'Approval pending',
    approvalDismissed: 'Instruction is dismissed',
    needApprove: 'Approval is needed',
    executeInstruction: 'Execute {instructionName}',
    resumeInstruction: 'Resume {instructionName}',
    inProgressInstruction: '{instructionName} in progress...',
    types: {
      [REMEDIATION_INSTRUCTION_TYPES.simpleManual]: 'Manual simplified',
      [REMEDIATION_INSTRUCTION_TYPES.manual]: 'Manual',
      [REMEDIATION_INSTRUCTION_TYPES.auto]: 'Automatic',
    },
    tooltips: {
      endpoint: 'Endpoint should be in question in Yes/No format',
    },
    table: {
      rating: 'Rating',
      monthExecutions: '№ of executions\nthis month',
      lastExecutedOn: 'Last executed on',
    },
    errors: {
      operationRequired: 'Please add at least one operation',
      stepRequired: 'Please add at least one step',
      jobRequired: 'Please add at least one job',
    },
  },

  configuration: {
    host: 'Host',
    usingConfiguration: 'Cannot be deleted since it is in use',
  },

  instructionExecute: {
    timeToComplete: '{duration} to complete',
    completedAt: 'Completed at {time}',
    failedAt: 'Failed at {time}',
    startedAt: 'Started at {time}',
    closeConfirmationText: 'Would you like to resume this instruction later?',
    queueNumber: '{number} {name} jobs are in the queue',
    runJobs: 'Run jobs',
    popups: {
      success: '{instructionName} has been successfully completed',
      failed: '{instructionName} has been failed. Please escalate this problem further',
      connectionError: 'There is a problem with connection. Please click on refresh button or reload the page.',
      wasAborted: '{instructionName} has been aborted',
      wasPaused: 'The {instructionName} instruction on {alarmName} alarm was paused at {date}. You can resume it manually.',
      wasRemovedOrDisabled: 'The {instructionName} instruction was removed or disabled.',
    },
    jobs: {
      title: 'Jobs assigned:',
      startedAt: 'Started at',
      launchedAt: 'Launched at',
      completedAt: 'Completed at',
      waitAlert: 'External job executor is not responding, please contact your admin',
      skip: 'Skip job',
      await: 'Await',
      failedReason: 'Failed reason',
      output: 'Output',
      instructionFailed: 'Instruction failed',
      instructionComplete: 'Instruction complete',
      stopped: 'Stopped',
    },
  },

  instructionsFilter: {
    button: 'Create instructions filter',
    filterByInstructions: 'For alarms by instructions',
    with: 'Show alarms with selected instructions',
    without: 'Show alarms without selected instructions',
    selectAll: 'Select all',
    alarmsListDisplay: 'Alarms list display',
    allAlarms: 'Show all filtered alarms',
    showWithInProgress: 'Show filtered alarms with instructions in progress',
    showWithoutInProgress: 'Show filtered alarms without instructions in progress',
    hideWithInProgress: 'Hide filtered alarms with instructions in progress',
    hideWithoutInProgress: 'Hide filtered alarms without instructions in progress',
    selectedInstructions: 'Selected instructions',
    selectedInstructionsHelp: 'Instructions of selected type are excluded from the list',
    inProgress: 'In progress',
    chip: {
      with: 'WITH',
      without: 'WITHOUT',
      all: 'ALL',
    },
  },

  instructionStat: {
    alarmsTimeline: 'Alarms timeline',
    executedAt: 'End of execution at',
    lastExecutedOn: 'Last executed on',
    modifiedOn: 'Modified on',
    averageCompletionTime: 'Average time\nof completion',
    executionCount: 'Number of\nexecutions',
    totalExecutions: 'Total executions',
    successfulExecutions: 'Successful executions',
    alarmStates: 'Alarms affected by state',
    okAlarmStates: 'Number of resulting\nOK states',
    instructionChanged: 'The instruction has been changed',
    alarmResolvedDate: 'Alarm resolved date',
    showFailedExecutions: 'Show failed instruction executions',
    remediationDuration: 'Remediation duration',
    timeoutAfterExecution: 'Timeout after execution',
    actions: {
      needRate: 'Rate it!',
      rate: 'Rate',
    },
  },

  pattern: {
    tabs: {
      pbehaviorTypes: {
        title: 'Pbehavior types',
        fields: {
          activeOnTypes: 'Active on types',
          disabledOnTypes: 'Disabled on types',
        },
      },
    },
  },

  job: {
    configuration: 'Configuration',
    jobId: 'Job ID',
    addJobs: 'Add {count} job | Add {count} jobs',
    usingJob: 'Cannot be deleted since it is in use',
    query: 'Query',
    multipleExecutions: 'Allow parallel execution',
    jobWaitInterval: 'Job wait interval',
    addPayload: 'Add payload',
    deletePayload: 'Delete payload',
    payloadHelp: '<p>The accessible variables are: <strong>.Alarm</strong> and <strong>.Entity</strong></p>'
      + '<i>For example:</i>'
      + '<pre>{\n  resource: "{{ .Alarm.Value.Resource }}",\n  entity: "{{ .Entity.ID }}"\n}</pre>',
    errors: {
      invalidJSON: 'Invalid JSON',
    },
  },

  statistic: {
    remediation: 'Remediation',
    allInstructions: 'All instructions',
    manualInstructions: 'Manual instructions',
    autoInstructions: 'Automatic instructions',
    labels: {
      remediated: 'Remediated',
      withAssignedRemediations: 'Remediable (With assigned instructions)',
    },
    tooltips: {
      remediated: '{value} alarms remediated',
      assigned: '{value} alarms with instructions',
    },
  },
};
