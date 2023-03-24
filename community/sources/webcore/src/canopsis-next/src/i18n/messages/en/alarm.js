import { EVENT_ENTITY_TYPES, ALARM_METRIC_PARAMETERS } from '@/constants';

export default {
  eventsCount: 'Events count',
  liveReporting: 'Set a custom date range',
  ackAuthor: 'Ack author',
  lastCommentAuthor: 'Last comment author',
  lastCommentMessage: 'Last comment message',
  metaAlarm: 'Meta alarm',
  acknowledge: 'Acknowledge',
  ackResources: 'Ack resources',
  ackResourcesQuestion: 'Do you want to ack linked resources?',
  actionsRequired: 'Actions required',
  acknowledgeAndDeclareTicket: 'Acknowledge and declare ticket',
  acknowledgeAndAssociateTicket: 'Acknowledge and associate ticket',
  advancedSearch: '<span>Help on the advanced research :</span>\n'
    + '<p>- [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt;</p> [ AND|OR [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt; ]\n'
    + '<p>The "-" before the research is required</p>\n'
    + '<p>Operators :\n'
    + '    <=, <,=, !=,>=, >, LIKE (For MongoDB regular expression)</p>\n'
    + '<p>Value\'s type : String between quote, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"</p>\n'
    + '<dl><dt>Examples :</dt><dt>- Connector = "connector_1"</dt>\n'
    + '    <dd>Alarms whose connectors are "connector_1"</dd><dt>- Connector="connector_1" AND Resource="resource_3"</dt>\n'
    + '    <dd>Alarms whose connectors is "connector_1" and the resources is "resource_3"</dd><dt>- Connector="connector_1" OR Resource="resource_3"</dt>\n'
    + '    <dd>Alarms whose connectors is "connector_1" or the resources is "resource_3"</dd><dt>- Connector LIKE 1 OR Connector LIKE 2</dt>\n'
    + '    <dd>Alarms whose connectors contains 1 or 2</dd><dt>- NOT Connector = "connector_1"</dt>\n'
    + '    <dd>Alarms whose connectors isn\'t "connector_1"</dd>\n'
    + '</dl>',
  otherTickets: 'Other tickets are available in the expand panel',
  noAlarmFound: 'No alarms is found according to the patterns defined',
  associateTicketResources: 'Associate ticket for resources',
  followLink: 'Follow the link "{title}"',
  popups: {
    exportFailed: 'Failed to export alarms list in CSV format',
  },
  actions: {
    titles: {
      ack: 'Ack',
      fastAck: 'Fast ack',
      ackRemove: 'Cancel ack',
      pbehavior: 'Periodical behavior',
      snooze: 'Snooze alarm',
      declareTicket: 'Declare ticket',
      associateTicket: 'Associate ticket',
      cancel: 'Cancel alarm',
      changeState: 'Change and lock severity',
      variablesHelp: 'List of available variables',
      history: 'History',
      groupRequest: 'Suggest group request for meta alarm',
      createManualMetaAlarm: 'Manual meta alarm management',
      removeAlarmsFromManualMetaAlarm: 'Unlink alarm from manual meta alarm',
      comment: 'Comment',
    },
    iconsTitles: {
      ack: 'Ack',
      declareTicket: 'Declare ticket',
      canceled: 'Canceled',
      snooze: 'Snooze',
      pbehaviors: 'Periodic behaviors',
      grouping: 'Meta alarm',
      comment: 'Comment',
    },
    iconsFields: {
      ticketNumber: 'Ticket number',
      parents: 'Causes',
      children: 'Consequences',
    },
  },
  timeLine: {
    titlePaths: {
      by: 'by',
    },
    stateCounter: {
      header: 'Cropped Severity (since last change of status)',
      stateIncreased: 'State increased',
      stateDecreased: 'State decreases',
    },
    types: {
      [EVENT_ENTITY_TYPES.ack]: 'Ack',
      [EVENT_ENTITY_TYPES.ackRemove]: 'Ack removed',
      [EVENT_ENTITY_TYPES.stateinc]: 'State increased',
      [EVENT_ENTITY_TYPES.statedec]: 'State decreased',
      [EVENT_ENTITY_TYPES.statusinc]: 'Status increased',
      [EVENT_ENTITY_TYPES.statusdec]: 'Status decreased',
      [EVENT_ENTITY_TYPES.assocTicket]: 'Ticket associated',
      [EVENT_ENTITY_TYPES.declareTicket]: 'Ticket is declared successfully',
      [EVENT_ENTITY_TYPES.declareTicketFail]: 'Ticket declaration is failed',
      [EVENT_ENTITY_TYPES.webhookStart]: 'Webhook has started',
      [EVENT_ENTITY_TYPES.webhookComplete]: 'Webhook is executed successfully',
      [EVENT_ENTITY_TYPES.webhookFail]: 'Webhook is failed',
      [EVENT_ENTITY_TYPES.snooze]: 'Alarm snoozed',
      [EVENT_ENTITY_TYPES.unsooze]: 'Alarm unsnoozed',
      [EVENT_ENTITY_TYPES.changeState]: 'Change and lock severity',
      [EVENT_ENTITY_TYPES.pbhenter]: 'Periodic behavior enabled',
      [EVENT_ENTITY_TYPES.pbhleave]: 'Periodic behavior disabled',
      [EVENT_ENTITY_TYPES.cancel]: 'Alarm cancelled',
      [EVENT_ENTITY_TYPES.comment]: 'Alarm commented',
      [EVENT_ENTITY_TYPES.metaalarmattach]: 'Alarm linked to meta alarm',
      [EVENT_ENTITY_TYPES.instructionStart]: 'Instruction has been started',
      [EVENT_ENTITY_TYPES.instructionPause]: 'Instruction has been paused',
      [EVENT_ENTITY_TYPES.instructionResume]: 'Instruction has been resumed',
      [EVENT_ENTITY_TYPES.instructionComplete]: 'Instruction has been completed',
      [EVENT_ENTITY_TYPES.instructionAbort]: 'Instruction has been aborted',
      [EVENT_ENTITY_TYPES.instructionFail]: 'Instruction has been failed',
      [EVENT_ENTITY_TYPES.instructionJobStart]: 'Instruction job has been started',
      [EVENT_ENTITY_TYPES.instructionJobComplete]: 'Instruction job has been completed',
      [EVENT_ENTITY_TYPES.instructionJobAbort]: 'Instruction job has been aborted',
      [EVENT_ENTITY_TYPES.instructionJobFail]: 'Instruction job has been failed',
      [EVENT_ENTITY_TYPES.autoInstructionStart]: 'Instruction has been started automatically',
      [EVENT_ENTITY_TYPES.autoInstructionComplete]: 'Instruction has been completed automatically',
      [EVENT_ENTITY_TYPES.autoInstructionFail]: 'Instruction has been failed automatically',
      [EVENT_ENTITY_TYPES.autoInstructionAlreadyRunning]: 'Instruction has been started automatically for another alarm',
      [EVENT_ENTITY_TYPES.junitTestSuiteUpdate]: 'Test suite has been updated',
      [EVENT_ENTITY_TYPES.junitTestCaseUpdate]: 'Test case has been updated',
    },
  },
  tabs: {
    moreInfos: 'More infos',
    timeLine: 'Timeline',
    alarmsChildren: 'Alarms consequences',
    trackSource: 'Track source',
    impactChain: 'Impact chain',
    entityGantt: 'Gantt chart',
    ticketsDeclared: 'Tickets declared',
  },
  moreInfos: {
    defineATemplate: 'To define a template for this window, go to the alarms list settings',
  },
  infoPopup: 'Info popup',
  tooltips: {
    priority: 'The priority parameter is derived from the alarm severity multiplied by impact level of the entity on which the alarm is raised',
    runningManualInstructions: 'Manual instruction <strong>{title}</strong> in progress | Manual instructions <strong>{title}</strong> in progress',
    runningAutoInstructions: 'Automatic instruction <strong>{title}</strong> in progress | Automatic instructions <strong>{title}</strong> in progress',
    successfulManualInstructions: 'Manual instruction <strong>{title}</strong> is successful | Manual instructions <strong>{title}</strong> is successful',
    successfulAutoInstructions: 'Automatic instruction <strong>{title}</strong> is successful | Automatic instructions <strong>{title}</strong> is successful',
    failedManualInstructions: 'Manual instruction <strong>{title}</strong> is failed | Manual instructions <strong>{title}</strong> is failed',
    failedAutoInstructions: 'Automatic instruction <strong>{title}</strong> is failed | Automatic instructions <strong>{title}</strong> is failed',
    hasManualInstruction: 'There is a manual instruction for this type of an incident | There are a manual instructions for this type of an incident',
  },
  metrics: {
    [ALARM_METRIC_PARAMETERS.createdAlarms]: 'Number of created alarms',
    [ALARM_METRIC_PARAMETERS.activeAlarms]: 'Number of active alarms',
    [ALARM_METRIC_PARAMETERS.nonDisplayedAlarms]: 'Number of non-displayed alarms',
    [ALARM_METRIC_PARAMETERS.instructionAlarms]: 'Number of alarms under automatic remediation',
    [ALARM_METRIC_PARAMETERS.pbehaviorAlarms]: 'Number of alarms with PBehavior',
    [ALARM_METRIC_PARAMETERS.correlationAlarms]: 'Number of alarms with correlation',
    [ALARM_METRIC_PARAMETERS.ackAlarms]: 'Total number of acks',
    [ALARM_METRIC_PARAMETERS.ackActiveAlarms]: 'Number of active alarms with acks',
    [ALARM_METRIC_PARAMETERS.cancelAckAlarms]: 'Number of canceled acks',
    [ALARM_METRIC_PARAMETERS.ticketActiveAlarms]: 'Number of active alarms with tickets',
    [ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms]: 'Number of active alarms without tickets',
    [ALARM_METRIC_PARAMETERS.ratioCorrelation]: '% of correlated alarms',
    [ALARM_METRIC_PARAMETERS.ratioInstructions]: '% of alarms with auto remediation',
    [ALARM_METRIC_PARAMETERS.ratioTickets]: '% of alarms with tickets created',
    [ALARM_METRIC_PARAMETERS.ratioNonDisplayed]: '% of non-displayed alarms',
    [ALARM_METRIC_PARAMETERS.ratioRemediatedAlarms]: '% of manually remediated alarms',
    [ALARM_METRIC_PARAMETERS.averageAck]: 'Average time to ack alarms',
    [ALARM_METRIC_PARAMETERS.averageResolve]: 'Average time to resolve alarms',
    [ALARM_METRIC_PARAMETERS.manualInstructionExecutedAlarms]: 'Number of manually remediated alarms',
    [ALARM_METRIC_PARAMETERS.manualInstructionAssignedAlarms]: 'Number of alarms with manual instructions',
    [ALARM_METRIC_PARAMETERS.notAckedAlarms]: 'Number of not acked alarms',
    [ALARM_METRIC_PARAMETERS.notAckedInHourAlarms]: 'Number of not acked alarms with duration 1-4h',
    [ALARM_METRIC_PARAMETERS.notAckedInFourHoursAlarms]: 'Number of not acked alarms with duration 4-24h',
    [ALARM_METRIC_PARAMETERS.notAckedInDayAlarms]: 'Number of not acked alarms older than 24h',
  },
  fields: {
    displayName: 'Display name',
    initialOutput: 'Initial output',
    initialLongOutput: 'Initial long output',
    lastComment: 'Last comment',
    ackBy: 'Acked by',
    ackMessage: 'Acked message',
    ackInitiator: 'Acked initiator',
    stateMessage: 'State message',
    statusMessage: 'Status message',
    totalStateChanges: 'Total state changes',
    stateValue: 'State value',
    statusValue: 'Status value',
    ackAt: 'Acked at',
    stateAt: 'State changed at',
    statusAt: 'Status changed at',
    resolved: 'Resolved at',
    activationDate: 'Activation date',
    currentStateDuration: 'Current state duration',
    snoozeDuration: 'Snooze duration',
    pbhInactiveDuration: 'Pbehavior inactive duration',
    activeDuration: 'Active duration',
    eventsCount: 'Events count',
    extraDetails: 'Extra details',
    ticketAuthor: 'Ticket submitter',
    ticketId: 'Ticket id',
    ticketMessage: 'Ticket message',
    entityId: 'Entity ID',
    entityName: 'Entity name',
    entityCategoryName: 'Entity category name',
    entityType: 'Entity type',
    entityComponent: 'Entity component',
    entityConnector: 'Entity connector',
    entityImpactLevel: 'Entity impact level',
    entityKoEvents: 'Entity KO events',
    entityOkEvents: 'Entity OK events',
    entityInfos: 'Entity infos',
    entityComponentInfos: 'Entity component infos',
    entityLastPbehaviorDate: 'Entity last pbehavior date',
  },
};
