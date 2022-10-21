import { merge } from 'lodash';

import {
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  EVENT_ENTITY_TYPES,
  STATS_TYPES,
  STATS_CRITICITY,
  QUICK_RANGES,
  TOURS,
  BROADCAST_MESSAGES_STATUSES,
  USER_PERMISSIONS_PREFIXES,
  PBEHAVIOR_RRULE_PERIODS_RANGES,
  WIDGET_TYPES,
  ACTION_TYPES,
  ENTITY_TYPES,
  TEST_SUITE_STATUSES,
  SIDE_BARS,
  STATE_SETTING_METHODS,
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES,
  REMEDIATION_INSTRUCTION_TYPES,
  IDLE_RULE_TYPES,
  IDLE_RULE_ALARM_CONDITIONS,
  USERS_PERMISSIONS,
  ALARMS_OPENED_VALUES,
  HEALTHCHECK_SERVICES_NAMES,
  HEALTHCHECK_ENGINES_NAMES,
  GROUPS_NAVIGATION_TYPES,
  ALARM_METRIC_PARAMETERS,
  USER_METRIC_PARAMETERS,
  EVENT_FILTER_TYPES,
  PATTERN_OPERATORS,
  PATTERN_TYPES,
  PATTERN_FIELD_TYPES,
  PBEHAVIOR_TYPE_TYPES,
  SCENARIO_TRIGGERS,
  WEATHER_ACTIONS_TYPES,
  MAP_TYPES,
  MERMAID_THEMES,
  EVENT_FILTER_EXTERNAL_DATA_TYPES,
  EVENT_FILTER_EXTERNAL_DATA_CONDITION_TYPES,
  EVENT_FILTER_PATTERN_FIELDS,
  SERVICE_WEATHER_STATE_COUNTERS,
} from '@/constants';

import featureService from '@/services/features';

export default merge({
  common: {
    ok: 'Ok',
    undefined: 'Not defined',
    entity: 'Entity | Entities',
    service: 'Service',
    widget: 'Widget',
    addWidget: 'Add widget',
    addTab: 'Add tab',
    shareLink: 'Create share link',
    addPbehavior: 'Add pbehavior',
    refresh: 'Refresh',
    toggleEditView: 'Toggle view edition mode',
    toggleEditViewSubtitle: 'If you want to save widget positions you should toggle off the editing mode for that',
    name: 'Name',
    description: 'Description',
    author: 'Author',
    submit: 'Submit',
    cancel: 'Cancel',
    continue: 'Continue',
    stop: 'Stop',
    options: 'Options',
    type: 'Type',
    quitEditing: 'Quit editing',
    enabled: 'Enabled',
    disabled: 'Disabled',
    login: 'Login',
    yes: 'Yes',
    no: 'No',
    default: 'Default',
    confirmation: 'Are you sure?',
    parameter: 'Parameter | Parameters',
    by: 'By',
    date: 'Date',
    comment: 'Comment | Comments',
    lastComment: 'Last comment',
    start: 'Start',
    end: 'End',
    message: 'Message',
    preview: 'Preview',
    recursive: 'Recursive',
    select: 'Select',
    states: 'Severities',
    state: 'Severity',
    sla: 'Sla',
    authors: 'Authors',
    stat: 'Stat',
    trend: 'Trend',
    user: 'User | Users',
    role: 'Role | Roles',
    import: 'Import',
    export: 'Export',
    profile: 'Profile',
    username: 'Username',
    password: 'Password',
    authKey: 'Auth. key',
    widgetId: 'Widget id',
    connect: 'Connect',
    optional: 'optional',
    logout: 'Logout',
    title: 'Title',
    save: 'Save',
    label: 'Label',
    field: 'Field',
    value: 'Value',
    limit: 'Limit',
    add: 'Add',
    create: 'Create',
    delete: 'Delete',
    show: 'Show',
    hide: 'Hide',
    edit: 'Edit',
    duplicate: 'Duplicate',
    play: 'Play',
    copyLink: 'Copy link',
    parse: 'Parse',
    home: 'Home',
    step: 'Step',
    paginationItems: 'showing {first} to {last} of {total} entries',
    apply: 'Apply',
    from: 'From',
    to: 'To',
    tags: 'tags',
    actionsLabel: 'Actions',
    noResults: 'No results',
    result: 'Result',
    exploitation: 'Exploitation',
    administration: 'Administration',
    forbidden: 'Forbidden',
    notFound: 'Not found',
    search: 'Search',
    filters: 'Filters',
    filter: 'Filter',
    emptyObject: 'Empty object',
    startDate: 'Start date',
    endDate: 'End date',
    link: 'Link | Links',
    stack: 'Stack',
    edition: 'Edition',
    icon: 'Icon',
    fullscreen: 'Fullscreen',
    interval: 'Interval',
    status: 'Status',
    unit: 'Unit',
    delay: 'Delay',
    begin: 'Begin',
    timezone: 'Timezone',
    reason: 'Reason',
    or: 'Or',
    and: 'And',
    priority: 'Priority',
    clear: 'Clear',
    deleteAll: 'Delete all',
    payload: 'Payload',
    note: 'Note',
    output: 'Output',
    displayName: 'Display name',
    created: 'Creation date',
    updated: 'Last update date',
    expired: 'Expired date',
    accessed: 'Accessed at',
    lastEventDate: 'Last event date',
    pattern: 'Pattern | Patterns',
    correlation: 'Correlation',
    periods: 'Periods',
    range: 'Range',
    duration: 'Duration',
    previous: 'Previous',
    next: 'Next',
    eventPatterns: 'Event patterns',
    alarmPatterns: 'Alarm patterns',
    entityPatterns: 'Entity patterns',
    pbehaviorPatterns: 'Pbehavior patterns',
    totalEntityPatterns: 'Total entity patterns',
    serviceWeatherPatterns: 'Service weather patterns',
    addFilter: 'Add filter',
    id: 'Id',
    reset: 'Reset',
    selectColor: 'Select color',
    triggers: 'Triggers',
    disableDuringPeriods: 'Disable during periods',
    retryDelay: 'Delay',
    retryUnit: 'Unit',
    retryCount: 'Repeat',
    ticket: 'Ticket',
    method: 'Method',
    url: 'URL',
    category: 'Category',
    infos: 'Infos',
    impactLevel: 'Impact level',
    impactState: 'Impact state',
    loadMore: 'Load more',
    download: 'Download',
    initiator: 'Initiator',
    percent: 'Percent | Percents',
    tests: 'Tests',
    total: 'Total',
    error: 'Error | Errors',
    failures: 'Failures',
    skipped: 'Skipped',
    current: 'Current',
    average: 'Average',
    information: 'Information | Informations',
    file: 'File',
    group: 'Group | Groups',
    view: 'View | Views',
    tab: 'Tab | Tabs',
    access: 'Access',
    communication: 'Communication | Communications',
    general: 'General',
    notification: 'Notification | Notifications',
    dismiss: 'Dismiss',
    approve: 'Approve',
    summary: 'Summary',
    recurrence: 'Recurrence',
    statistics: 'Statistics',
    action: 'Action | Actions',
    minimal: 'Minimal',
    optimal: 'Optimal',
    graph: 'Graph | Graphs',
    systemStatus: 'System status',
    downloadAsPng: 'Download as PNG',
    rating: 'Rating | Ratings',
    sampling: 'Sampling',
    parametersToDisplay: '{count} parameters to display',
    uptime: 'Uptime',
    maintenance: 'Maintenance',
    downtime: 'Downtime',
    toTheTop: 'To the top',
    time: 'Time',
    lastModifiedOn: 'Last modified on',
    lastModifiedBy: 'Last modified by',
    exportAsCsv: 'Export as csv',
    criteria: 'Criteria',
    ratingSettings: 'Rating settings',
    pbehavior: 'Pbehavior | Pbehaviors',
    activePbehavior: 'Active pbehavior | Active pbehaviors',
    searchBy: 'Search by',
    dictionary: 'Dictionary',
    condition: 'Condition | Conditions',
    template: 'Template',
    pbehaviorList: 'List periodic behaviors',
    canceled: 'Canceled',
    snooze: 'Snooze',
    snoozed: 'Snoozed',
    impact: 'Impact | Impacts',
    depend: 'Depend | Depends',
    componentInfo: 'Component info | Component infos',
    connector: 'Connector',
    connectorName: 'Connector name',
    component: 'Component',
    resource: 'Resource',
    extraDetail: 'Extra detail | Extra details',
    ack: 'Ack',
    acked: 'Acked',
    ackedAt: 'Acked at',
    ackedBy: 'Acked by',
    resolvedAt: 'Resolved at',
    extraInfo: 'Extra info | Extra infos',
    custom: 'Custom',
    eventType: 'Event type',
    sourceType: 'Source type',
    cycleDependency: 'Cycle dependency',
    checkPattern: 'Check pattern',
    itemFound: '{count} item found | {count} items found',
    canonicalType: 'Canonical type',
    map: 'Map | Maps',
    instructions: 'Instructions',
    playlist: 'Playlist | Playlists',
    ctrlZoom: 'Use ctrl + mouse wheel for zoom',
    calendar: 'Calendar',
    tag: 'Tag | Tags',
    sharedTokens: 'Shared tokens',
    notAvailable: 'N/a',
    addMore: 'Add more',
    attribute: 'Attribute',
    timeTaken: 'Time taken',
    enginesMetrics: 'Engines` metrics',
    failed: 'Failed',
    close: 'Close',
    timestamp: 'Timestamp',
    actions: {
      acknowledgeAndDeclareTicket: 'Acknowledge and declare ticket',
      acknowledgeAndAssociateTicket: 'Acknowledge and associate ticket',
      saveChanges: 'Save changes',
      reportIncident: 'Report an incident',
      [EVENT_ENTITY_TYPES.ack]: 'Acknowledge',
      [EVENT_ENTITY_TYPES.declareTicket]: 'Declare ticket',
      [EVENT_ENTITY_TYPES.validate]: 'Validate',
      [EVENT_ENTITY_TYPES.invalidate]: 'Invalidate',
      [EVENT_ENTITY_TYPES.pause]: 'Pause',
      [EVENT_ENTITY_TYPES.play]: 'Play',
      [EVENT_ENTITY_TYPES.cancel]: 'Cancel',
      [EVENT_ENTITY_TYPES.assocTicket]: 'Associate ticket',
      [EVENT_ENTITY_TYPES.comment]: 'Comment',
      [EVENT_ENTITY_TYPES.executeInstruction]: 'Execute instruction',
    },
    acknowledge: 'Acknowledge',
    acknowledgeAndDeclareTicket: 'Acknowledge and declare ticket',
    acknowledgeAndAssociateTicket: 'Acknowledge and associate ticket',
    saveChanges: 'Save changes',
    reportIncident: 'Report an incident',
    times: {
      second: 'second | seconds',
      minute: 'minute | minutes',
      hour: 'hour | hours',
      day: 'day | days',
      week: 'week | weeks',
      month: 'month | months',
      year: 'year | years',
    },
    timeFrequencies: {
      secondly: 'Secondly',
      minutely: 'Minutely',
      hourly: 'Hourly',
      daily: 'Daily',
      weekly: 'Weekly',
      monthly: 'Monthly',
      yearly: 'Yearly',
    },
    weekDays: {
      monday: 'Monday',
      tuesday: 'Tuesday',
      wednesday: 'Wednesday',
      thursday: 'Thursday',
      friday: 'Friday',
      saturday: 'Saturday',
      sunday: 'Sunday',
    },
    months: {
      january: 'January',
      february: 'February',
      march: 'March',
      april: 'April',
      may: 'May',
      june: 'June',
      july: 'July',
      august: 'August',
      september: 'September',
      october: 'October',
      november: 'November',
      december: 'December',
    },
    stateTypes: {
      [ENTITIES_STATES.ok]: 'Ok',
      [ENTITIES_STATES.minor]: 'Minor',
      [ENTITIES_STATES.major]: 'Major',
      [ENTITIES_STATES.critical]: 'Critical',
    },
    statusTypes: {
      [ENTITIES_STATUSES.closed]: 'Closed',
      [ENTITIES_STATUSES.ongoing]: 'Ongoing',
      [ENTITIES_STATUSES.flapping]: 'Flapping',
      [ENTITIES_STATUSES.stealthy]: 'Stealth',
      [ENTITIES_STATUSES.cancelled]: 'Canceled',
      [ENTITIES_STATUSES.noEvents]: 'No events',
    },
    operators: {
      [PATTERN_OPERATORS.equal]: 'Equal',
      [PATTERN_OPERATORS.contains]: 'Contains',
      [PATTERN_OPERATORS.notEqual]: 'Not equal',
      [PATTERN_OPERATORS.notContains]: 'Does not contain',
      [PATTERN_OPERATORS.beginsWith]: 'Begins with',
      [PATTERN_OPERATORS.notBeginWith]: 'Does not begin with',
      [PATTERN_OPERATORS.endsWith]: 'Ends with',
      [PATTERN_OPERATORS.notEndWith]: 'Does not end with',
      [PATTERN_OPERATORS.exist]: 'Exist',
      [PATTERN_OPERATORS.notExist]: 'Not exist',

      [PATTERN_OPERATORS.hasEvery]: 'Has every',
      [PATTERN_OPERATORS.hasOneOf]: 'Has one of',
      [PATTERN_OPERATORS.isOneOf]: 'Is one of',
      [PATTERN_OPERATORS.hasNot]: 'Has not',
      [PATTERN_OPERATORS.isNotOneOf]: 'Is not one of',
      [PATTERN_OPERATORS.isEmpty]: 'Is empty',
      [PATTERN_OPERATORS.isNotEmpty]: 'Is not empty',

      [PATTERN_OPERATORS.higher]: 'Higher than',
      [PATTERN_OPERATORS.lower]: 'Lower than',

      [PATTERN_OPERATORS.longer]: 'Longer',
      [PATTERN_OPERATORS.shorter]: 'Shorter',

      [PATTERN_OPERATORS.ticketAssociated]: 'Ticket is associated',
      [PATTERN_OPERATORS.ticketNotAssociated]: 'Ticket is not associated',

      [PATTERN_OPERATORS.canceled]: 'Canceled',
      [PATTERN_OPERATORS.notCanceled]: 'Not canceled',

      [PATTERN_OPERATORS.snoozed]: 'Snoozed',
      [PATTERN_OPERATORS.notSnoozed]: 'Not snoozed',

      [PATTERN_OPERATORS.acked]: 'Acked',
      [PATTERN_OPERATORS.notAcked]: 'Not acked',

      [PATTERN_OPERATORS.isGrey]: 'Gray tiles',
      [PATTERN_OPERATORS.isNotGrey]: 'Not gray tiles',

      [PATTERN_OPERATORS.with]: 'With',
      [PATTERN_OPERATORS.without]: 'Without',
    },
    entityEventTypes: {
      [EVENT_ENTITY_TYPES.ack]: 'Ack',
      [EVENT_ENTITY_TYPES.ackRemove]: 'Ack remove',
      [EVENT_ENTITY_TYPES.assocTicket]: 'Associate ticket',
      [EVENT_ENTITY_TYPES.declareTicket]: 'Declare ticket',
      [EVENT_ENTITY_TYPES.cancel]: 'Cancel',
      [EVENT_ENTITY_TYPES.uncancel]: 'Uncancel',
      [EVENT_ENTITY_TYPES.changeState]: 'Change state',
      [EVENT_ENTITY_TYPES.check]: 'Check',
      [EVENT_ENTITY_TYPES.comment]: 'Comment',
      [EVENT_ENTITY_TYPES.snooze]: 'Snooze',
    },
    scenarioTriggers: {
      [SCENARIO_TRIGGERS.create]: {
        text: 'Alarm creation',
      },
      [SCENARIO_TRIGGERS.statedec]: {
        text: 'Alarm state decrease',
      },
      [SCENARIO_TRIGGERS.changestate]: {
        text: 'Alarm state has been changed by "change state" action',
      },
      [SCENARIO_TRIGGERS.stateinc]: {
        text: 'Alarm state increase',
      },
      [SCENARIO_TRIGGERS.changestatus]: {
        text: 'Alarm status changes eg. flapping',
      },
      [SCENARIO_TRIGGERS.ack]: {
        text: 'Alarm has been acked',
      },
      [SCENARIO_TRIGGERS.ackremove]: {
        text: 'Alarm has been unacked',
      },
      [SCENARIO_TRIGGERS.cancel]: {
        text: 'Alarm has been cancelled',
      },
      [SCENARIO_TRIGGERS.uncancel]: {
        text: 'Alarm has been uncancelled',
        helpText: 'Probably legacy trigger, because there is no way to uncancel alarm when you cancel it in the UI, but it\'s possible to send an uncancel event via API',
      },
      [SCENARIO_TRIGGERS.comment]: {
        text: 'Alarm has been commented',
      },
      [SCENARIO_TRIGGERS.done]: {
        text: 'Alarm is "done"',
        helpText: 'Probably legacy, because there is no such action in the UI, but it\'s possible to send a done event via API',
      },
      [SCENARIO_TRIGGERS.declareticket]: {
        text: 'Ticket has been declared by the UI action',
      },
      [SCENARIO_TRIGGERS.declareticketwebhook]: {
        text: 'Ticket has been declared by the webhook',
      },
      [SCENARIO_TRIGGERS.assocticket]: {
        text: 'Ticket has been associated with an alarm',
      },
      [SCENARIO_TRIGGERS.snooze]: {
        text: 'Alarm has been snoozed',
      },
      [SCENARIO_TRIGGERS.unsnooze]: {
        text: 'Alarm has been unsnoozed',
      },
      [SCENARIO_TRIGGERS.resolve]: {
        text: 'Alarm has been resolved',
      },
      [SCENARIO_TRIGGERS.activate]: {
        text: 'Alarm has been activated',
      },
      [SCENARIO_TRIGGERS.pbhenter]: {
        text: 'Alarm enters a periodic behavior',
      },
      [SCENARIO_TRIGGERS.pbhleave]: {
        text: 'Alarm leaves a periodic behavior',
      },
      [SCENARIO_TRIGGERS.instructionfail]: {
        text: 'Manual instruction has failed',
      },
      [SCENARIO_TRIGGERS.autoinstructionfail]: {
        text: 'Auto instruction has failed',
      },
      [SCENARIO_TRIGGERS.instructionjobfail]: {
        text: 'Manual or auto instruction\'s job is failed',
      },
      [SCENARIO_TRIGGERS.instructionjobcomplete]: {
        text: 'Manual or auto instruction\'s job is completed',
      },
      [SCENARIO_TRIGGERS.instructioncomplete]: {
        text: 'Manual instruction is completed',
      },
      [SCENARIO_TRIGGERS.autoinstructioncomplete]: {
        text: 'Auto instruction is completed',
      },
    },
  },
  variableTypes: {
    string: 'String',
    number: 'Number',
    boolean: 'Boolean',
    null: 'Null',
    array: 'Array',
  },
  context: {
    impacts: 'Impacts',
    dependencies: 'Dependencies',
    noEventsFilter: 'No events filter',
    impactChain: 'Impact chain',
    impactDepends: 'Impact/Depends',
    treeOfDependencies: 'Tree of dependencies',
    infosSearchLabel: 'Search infos',
    eventStatisticsMessage: '{ok} OK events\n{ko} KO Events',
    eventStatistics: 'Event statistics',
    actions: {
      titles: {
        editEntity: 'Edit entity',
        duplicateEntity: 'Duplicate entity',
        deleteEntity: 'Delete entity',
        pbehavior: 'Periodical behavior',
        variablesHelp: 'List of available variables',
        massEnable: 'Enable entities',
        massDisable: 'Disable entities',
      },
    },
    fab: {
      common: 'Add a new entity',
      addService: 'Add a new service entity',
    },
    popups: {
      massDeleteWarning: 'The mass deletion cannot be applied for some of selected elements, so they won\'t be deleted.',
    },
  },
  search: {
    alarmAdvancedSearch: '<span>Help on the advanced research :</span>\n'
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
    contextAdvancedSearch: '<span>Help on the advanced research :</span>\n'
      + '<p>- [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt;</p> [ AND|OR [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt; ]\n'
      + '<p>The "-" before the research is required</p>\n'
      + '<p>Operators :\n'
      + '    <=, <,=, !=,>=, >, LIKE (For MongoDB regular expression)</p>\n'
      + '<p>Value\'s type : String between quote, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"</p>\n'
      + '<dl><dt>Examples :</dt><dt>- Name = "name_1"</dt>\n'
      + '    <dd>Entities whose names are "name_1"</dd><dt>- Name="name_1" AND Type="service"</dt>\n'
      + '    <dd>Entities whose names is "name_1" and the types is "service"</dd><dt>- infos.custom.value="Custom value" OR Type="resource"</dt>\n'
      + '    <dd>Entities whose infos.custom.value is "Custom value" or the type is "resource"</dd><dt>- infos.custom.value LIKE 1 OR infos.custom.value LIKE 2</dt>\n'
      + '    <dd>Entities whose infos.custom.value contains 1 or 2</dd><dt>- NOT Name = "name_1"</dt>\n'
      + '    <dd>Entities whose name isn\'t "name_1"</dd>\n'
      + '</dl>',
    dynamicInfoAdvancedSearch: '<span>Help on the advanced research :</span>\n'
      + '<p>- [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt;</p> [ AND|OR [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt; ]\n'
      + '<p>The "-" before the research is required</p>\n'
      + '<p>Operators :\n'
      + '    <=, <,=, !=,>=, >, LIKE (For MongoDB regular expression)</p>\n'
      + '<p>For querying patterns, use "pattern" keyword as the &lt;ColumnName&gt; alias</p>\n'
      + '<p>Value\'s type : String between quote, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"</p>\n'
      + '<dl><dt>Examples :</dt><dt>- description = "testdyninfo"</dt>\n'
      + '    <dd>Dynamic info rules descriptions are "testdyninfo"</dd><dt>- pattern = "SEARCHPATTERN1"</dt>\n'
      + '    <dd>Dynamic info rules whose one of its patterns is equal "SEARCHPATTERN1"</dd><dt>- pattern LIKE "SEARCHPATTERN2"</dt>\n'
      + '    <dd>Dynamic info rules whose one of its patterns contains "SEARCHPATTERN2"</dd>'
      + '</dl>',
    submit: 'Search',
    clear: 'Clear search input',
  },
  login: {
    base: 'Standard',
    LDAP: 'LDAP',
    loginWithCAS: 'Login with CAS',
    loginWithSAML: 'Login with SAML',
    documentation: 'Documentation',
    forum: 'Forum',
    website: 'Canopsis.com',
    connectionProtocols: 'Connection protocols',
    errors: {
      incorrectEmailOrPassword: 'Incorrect email or password',
    },
  },
  alarmList: {
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
        manualMetaAlarmGroup: 'Manual meta alarm management',
        manualMetaAlarmUngroup: 'Unlink alarm from manual meta alarm',
        comment: 'Comment',
        executeInstruction: 'Execute {instructionName}',
        resumeInstruction: 'Resume {instructionName}',
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
        rule: 'Rule | Rules',
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
        [EVENT_ENTITY_TYPES.declareTicket]: 'Ticket declared',
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
      [ALARM_METRIC_PARAMETERS.ticketActiveAlarms]: 'Number of active alarms with acks',
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
    },
  },
  weather: {
    moreInfos: 'More info',
  },
  pbehaviors: {
    connector: 'Connector Type',
    connectorName: 'Connector name',
    isEnabled: 'Is enabled',
    begins: 'Begins',
    ends: 'Ends',
    type: 'Type',
    reason: 'Reason',
    rrule: 'Recurrence',
    status: 'Status',
    created: 'Creation date',
    updated: 'Last update date',
    lastAlarmDate: 'Last alarm date',
    massRemove: 'Remove pbehaviors',
    massEnable: 'Enable pbehaviors',
    massDisable: 'Disable pbehaviors',
    searchHelp: '<span>Help on the advanced research :</span>\n'
      + '<p>- [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt;</p> [ AND|OR [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt; ]\n'
      + '<p>The "-" before the research is required</p>\n'
      + '<p>Operators : <=, <,=, !=,>=, >, LIKE (For MongoDB regular expression)</p>\n'
      + '<p>For querying patterns, use "pattern" keyword as the &lt;ColumnName&gt; alias</p>\n'
      + '<p>Value\'s type : String between quote, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"</p>\n'
      + '<dl>'
      + '  <dt>Examples :</dt>'
      + '  <dt>- name = "name_1"</dt>\n'
      + '  <dd>Pbehavior name are "name_1"</dd>\n'
      + '  <dt>- rrule = "rrule_1"</dt>\n'
      + '  <dd>Pbehavior rrule are "rrule_1"</dd>\n'
      + '  <dt>- filter = "filter_1"</dt>\n'
      + '  <dd>Pbehavior filter are "filter_1"</dd>\n'
      + '  <dt>- type.name = "type_name_1"</dt>\n'
      + '  <dd>Pbehavior type name are "type_name_1"</dd>\n'
      + '  <dt>- reason.name = "reason_name_1"</dt>\n'
      + '  <dd>Pbehavior reason name are "reason_name_1"</dd>'
      + '</dl>',
    tabs: {
      entities: 'Entities',
    },
  },
  settings: {
    titles: {
      [SIDE_BARS.alarmSettings]: 'Alarm list settings',
      [SIDE_BARS.contextSettings]: 'Context table settings',
      [SIDE_BARS.serviceWeatherSettings]: 'Service weather settings',
      [SIDE_BARS.statsCalendarSettings]: 'Stats calendar settings',
      [SIDE_BARS.textSettings]: 'Text settings',
      [SIDE_BARS.counterSettings]: 'Counter settings',
      [SIDE_BARS.testingWeatherSettings]: 'Testing weather',
      [SIDE_BARS.mapSettings]: 'Mapping widget settings',
    },
    openedTypes: {
      [ALARMS_OPENED_VALUES.opened]: 'Opened alarms',
      [ALARMS_OPENED_VALUES.resolved]: 'All resolved alarms',
      [ALARMS_OPENED_VALUES.all]: 'Opened and recent resolved alarms',
    },
    advancedSettings: 'Advanced settings',
    entityDisplaySettings: 'Entity display settings',
    entitiesUnderPbehaviorEnabled: 'Entities under PBh type inactive, Pause, Maintenance display',
    widgetTitle: 'Widget title',
    columnName: 'Column name',
    defaultSortColumn: 'Default sort column',
    sortColumnNoData: 'Press <kbd>enter</kbd> to create a new one',
    columnNames: 'Column names',
    exportColumnNames: 'Export column names',
    groupColumnNames: 'Column names for meta alarms',
    trackColumnNames: 'Track alarm source columns',
    treeOfDependenciesColumnNames: 'Column names for tree of dependencies',
    orderBy: 'Order by',
    periodicRefresh: 'Periodic refresh',
    defaultNumberOfElementsPerPage: 'Default number of elements/page',
    elementsPerPage: 'Elements per page',
    filterOnOpenResolved: 'Filter on Open/Resolved',
    open: 'Open',
    resolved: 'Resolved',
    filters: 'Filters',
    filterEditor: 'Filter',
    isAckNoteRequired: 'Note field required when ack?',
    isSnoozeNoteRequired: 'Note field required when snooze?',
    linksCategoriesAsList: 'Display links as a list?',
    linksCategoriesLimit: 'Number of category items',
    isMultiAckEnabled: 'Multiple ack',
    isMultiDeclareTicketEnabled: 'Multiple declare ticket',
    fastAckOutput: 'Fast-ack output',
    isHtmlEnabledOnTimeLine: 'HTML enabled on timeline?',
    isCorrelationEnabled: 'Is correlation enabled?',
    duration: 'Duration',
    tstop: 'End date',
    periodsNumber: 'Number of steps',
    yesNoMode: 'Yes/No mode',
    selectAFilter: 'Select a filter',
    lockedFilter: 'Filter locked in widget settings',
    exportAsCsv: 'Export widget as csv file',
    criticityLevels: 'Criticity levels',
    isPriorityEnabled: 'Show priority',
    clearFilterDisabled: 'Disable possibility to clear selected filter',
    alarmsColumns: 'Alarm list columns',
    entitiesColumns: 'Context explorer columns',
    entityInfoPopup: 'Entity info popup',
    modal: '(Modal)',
    exportCsv: {
      title: 'Export CSV',
      fields: {
        separator: 'Separator',
        datetimeFormat: 'Datetime format',
      },
    },
    colorsSelector: {
      title: 'Colors selector',
      statsCriticity: {
        [STATS_CRITICITY.ok]: 'ok',
        [STATS_CRITICITY.minor]: 'minor',
        [STATS_CRITICITY.major]: 'major',
        [STATS_CRITICITY.critical]: 'critical',
      },
    },
    infoPopup: {
      title: 'Info popup',
      fields: {
        column: 'Column',
      },
    },
    rowGridSize: {
      title: 'Widget\'s size',
      noData: 'No row corresponding. Press <kbd>enter</kbd> to create a new one',
      fields: {
        row: 'Row',
      },
    },
    moreInfosModal: '"More Infos" Popup',
    expandGridRangeSize: 'Expand card (more infos / timeline) width',
    weatherTemplate: 'Template - Weather item',
    modalTemplate: 'Template - Modal',
    entityTemplate: 'Template - Entities',
    blockTemplate: 'Template - Tile',
    columnMobile: 'Columns - Mobile',
    columnTablet: 'Columns - Tablet',
    columnDesktop: 'Columns - Desktop',
    limit: 'Limit',
    height: 'Height',
    margin: {
      title: 'Block margins',
      top: 'Margin - Top',
      right: 'Margin - Right',
      bottom: 'Margin - Bottom',
      left: 'Margin - Left',
    },
    contextTypeOfEntities: {
      title: 'Type of entities',
      fields: {
        [ENTITY_TYPES.component]: 'Component',
        [ENTITY_TYPES.connector]: 'Connector Type',
        [ENTITY_TYPES.resource]: 'Resource',
        [ENTITY_TYPES.service]: 'Service',
      },
    },
    considerPbehaviors: {
      title: 'Consider pbehaviors',
    },
    serviceWeatherModalTypes: {
      title: 'Type of the weather modal window',
      fields: {
        moreInfo: 'More info',
        alarmList: 'Alarm list',
        both: 'Both',
      },
    },
    templateEditor: 'Template',
    columns: {
      isHtml: 'Is it HTML?',
      withTemplate: 'Custom template',
      isState: 'Displayed as severity?',
    },
    liveReporting: {
      title: 'Live reporting',
    },
    counterLevels: {
      title: 'Levels',
      fields: {
        counter: 'Counter',
      },
    },
    counters: 'Counters',
    pbehaviorCounters: 'Pbehavior counters',
    entityStateCounters: 'Entity states counters',
    remediationInstructionsFilters: 'Instructions filters',
    colorIndicator: {
      title: 'Color indicator',
      fields: {
        displayAsSeverity: 'Display as severity',
        displayAsPriority: 'Display as priority',
      },
    },
    receiveByApi: 'Receive by the API',
    serverStorage: 'Server storage',
    filenameRecognition: 'Filename recognition',
    resultDirectory: 'Test results storage',
    screenshotDirectories: {
      title: 'Screenshots storage settings',
      helpText: 'Define where screenshots are stored',
    },
    screenshotMask: {
      title: 'Screenshots filename mask',
      helpText: '<dl>'
        + '<dt>Define the filename mask of which screenshots are created using the following variables:<dt>\n'
        + '<dd>- test case name %test_case%</dd>\n'
        + '<dd>- date (YYYY, MM, DD)</dd>\n'
        + '<dd>- time of execution (hh, mm, ss)</dd>'
        + '</dl>',
    },
    videoDirectories: {
      title: 'Video storage settings',
      helpText: 'Define where video are stored',
    },
    videoMask: {
      title: 'Videos filename mask',
      helpText: '<dl>'
        + '<dt>Define the filename mask of which videos are created using the following variables:<dt>\n'
        + '<dd>- test case name %test_case%</dd>\n'
        + '<dd>- date (YYYY, MM, DD)</dd>\n'
        + '<dd>- time of execution (hh, mm, ss)</dd>'
        + '</dl>',
    },
    stickyHeader: 'Sticky header',
    reportFileRegexp: {
      title: 'Report file mask',
      helpText: '<dl>'
        + '<dt>Define the filename regexp of which report:<dt>\n'
        + '<dd>For example:</dd>\n'
        + '<dd>"^(?P&lt;name&gt;\\\\w+)_(.+)\\\\.xml$"</dd>\n'
        + '</dl>',
    },
    density: {
      title: 'Default view',
      comfort: 'Comfort view',
      compact: 'Compact view',
    },
  },
  modals: {
    common: {
      titleButtons: {
        minimizeTooltip: 'You already have minimized modal window',
      },
    },
    contextInfos: {
      title: 'Entities infos',
    },
    createEntity: {
      create: {
        title: 'Create an entity',
      },
      edit: {
        title: 'Edit an entity',
      },
      duplicate: {
        title: 'Duplicate an entity',
      },
      success: {
        create: 'Entity successfully created!',
        edit: 'Entity successfully edited!',
        duplicate: 'Entity successfully duplicated!',
      },
    },
    createService: {
      create: {
        title: 'Create a service',
      },
      edit: {
        title: 'Edit a service',
      },
      duplicate: {
        title: 'Duplicate a service',
      },
      success: {
        create: 'Service successfully created!',
        edit: 'Service successfully edited!',
        duplicate: 'Service successfully duplicated!',
      },
    },
    createEntityInfo: {
      create: {
        title: 'Add an information',
      },
      edit: {
        title: 'Edit an information',
      },
    },
    view: {
      create: {
        title: 'Create a view',
      },
      edit: {
        title: 'Edit the view',
      },
      duplicate: {
        title: 'Duplicate the view - {viewTitle}',
        infoMessage: 'You\'re duplicating a view. All duplicated view\'s rows/widgets will be copied on the new view.',
      },
      noData: 'No group corresponding. Press <kbd>enter</kbd> to create a new one',
      fields: {
        periodicRefresh: 'Periodic refresh',
        groupIds: 'Choose a group, or create a new one',
        groupTags: 'Group tags',
      },
      success: {
        create: 'New view created!',
        edit: 'View successfully edited!',
        duplicate: 'View successfully duplicated!',
        delete: 'View successfully deleted!',
      },
      fail: {
        create: 'View creation failed...',
        edit: 'View edition failed...',
        duplicate: 'View duplication failed...',
        delete: 'View deletion failed...',
      },
    },
    createEvent: {
      fields: {
        output: 'Note',
      },
    },
    createAckEvent: {
      title: 'Ack',
      tooltips: {
        ackResources: 'Do you want to ack linked resources?',
      },
      fields: {
        ticket: 'Ticket number',
        output: 'Note',
        ackResources: 'Ack resources',
      },
    },
    confirmAckWithTicket: {
      continueAndAssociateTicket: 'Continue and associate ticket',
      infoMessage: `A ticket number has been specified.
        Maybe you wanted to associate this ticket number to the alarm.
        If so, click on "Continue and associate ticket" button.
        To continue the ack action without taking ticket number into account,
        click on "Continue" button.`,
    },
    createSnoozeEvent: {
      title: 'Snooze',
      fields: {
        duration: 'Duration',
      },
    },
    createCancelEvent: {
      title: 'Cancel',
    },
    createGroupRequestEvent: {
      title: 'Suggest group request for meta alarm',
    },
    createGroupEvent: {
      title: 'Create meta alarm',
    },
    createChangeStateEvent: {
      title: 'Change severity',
      states: {
        ok: 'Info',
        minor: 'Minor',
        major: 'Major',
        critical: 'Critical',
      },
      fields: {
        output: 'Note',
      },
    },
    createPbehavior: {
      create: {
        title: 'Create periodical behavior',
      },
      edit: {
        title: 'Edit periodic behavior',
      },
      duplicate: {
        title: 'Duplicate periodic behavior',
      },
      steps: {
        general: {
          title: 'General parameters',
          dates: 'Dates',
          fields: {
            enabled: 'Enabled',
            name: 'Name',
            reason: 'Reason',
            type: 'Type',
            start: 'Start',
            stop: 'End',
            fullDay: 'Whole day',
            noEnding: 'No ending',
            startOnTrigger: 'Start on trigger',
          },
        },
        filter: {
          title: 'Filter',
        },
        rrule: {
          title: 'Recurrence rule',
          exdate: 'Exclusion dates',
          buttons: {
            addExdate: 'Add an exclusion date',
          },
        },
        comments: {
          title: 'Comments',
          buttons: {
            addComment: 'Add comment',
          },
          fields: {
            message: 'Message',
          },
        },
      },
      errors: {
        invalid: 'Invalid',
      },
      success: {
        create: 'Pbehavior successfully created! You may need to wait 60 sec to see it in interface',
      },
      cancelConfirmation: 'Some data has been modified and will not be saved. Do you really want to close this menu?',
    },
    createPause: {
      title: 'Create Pause event',
      comment: 'Comment',
      reason: 'Reason',
    },
    createAckRemove: {
      title: 'Remove ack',
    },
    createDeclareTicket: {
      title: 'Declare ticket',
    },
    createAssociateTicket: {
      title: 'Associate ticket number',
      fields: {
        ticket: 'Number of the ticket',
      },
      alerts: {
        noAckItems: 'There is {count} item without ack. Ack event for the item will send before. | There is {count} items without ack. Ack events for items will send before.',
      },
    },
    liveReporting: {
      editLiveReporting: 'Live reporting',
      dateInterval: 'Date interval',
      today: 'Today',
      yesterday: 'Yesterday',
      last7Days: 'Last 7 days',
      last30Days: 'Last 30 days',
      thisMonth: 'This month',
      lastMonth: 'Last month',
      custom: 'Custom',
      tstart: 'Begins',
      tstop: 'Ends',
    },
    infoPopupSetting: {
      title: 'Info popup',
      add: 'Add',
      column: 'Column',
      addInfoPopup: {
        title: 'Add an info popup',
      },
    },
    variablesHelp: {
      variables: 'Variables',
      copyToClipboard: 'Copy to clipboard',
    },
    service: {
      refreshEntities: 'Refresh entities list',
      editPbehaviors: 'Edit pbehaviors',
      entity: {
        tabs: {
          info: 'Info',
          treeOfDependencies: 'Tree of dependencies',
        },
      },
    },
    createFilter: {
      create: {
        title: 'Create filter',
      },
      edit: {
        title: 'Edit filter',
      },
      duplicate: {
        title: 'Duplicate filter',
      },
      fields: {
        title: 'Title',
      },
      emptyFilters: 'No filters added yet',
    },
    colorPicker: {
      title: 'Color picker',
    },
    textEditor: {
      title: 'Text editor',
    },
    createWidget: {
      title: 'Select a widget',
      types: {
        [WIDGET_TYPES.alarmList]: {
          title: 'Alarm list',
        },
        [WIDGET_TYPES.context]: {
          title: 'Context explorer',
        },
        [WIDGET_TYPES.serviceWeather]: {
          title: 'Service weather',
        },
        [WIDGET_TYPES.statsCalendar]: {
          title: 'Stats calendar',
        },
        [WIDGET_TYPES.text]: {
          title: 'Text',
        },
        [WIDGET_TYPES.counter]: {
          title: 'Counter',
        },
        [WIDGET_TYPES.testingWeather]: {
          title: 'Junit scenarios',
        },
        [WIDGET_TYPES.map]: {
          title: 'Mapping',
        },
      },
    },
    manageHistogramGroups: {
      title: {
        add: 'Add a group',
        edit: 'Edit a group',
      },
    },
    addStat: {
      title: {
        add: 'Add a stat',
        edit: 'Edit a stat',
      },
      slaTooltip: 'The sla parameter should be a string of the form "<op> <value>", where <op> is <, >, <= or >= and <value> is a number',
    },
    group: {
      create: {
        title: 'Create group',
      },
      edit: {
        title: 'Edit group',
      },
      fields: {
        name: 'Name',
      },
      errors: {
        isNotEmpty: 'The group is not empty',
      },
    },
    alarmsList: {
      title: 'Alarm list',
      prefixTitle: '{prefix} - alarm list',
    },
    createUser: {
      create: {
        title: 'Create user',
      },
      edit: {
        title: 'Edit user',
      },
    },
    createRole: {
      create: {
        title: 'Create role',
      },
      edit: {
        title: 'Edit role',
      },
    },
    createEventFilter: {
      create: {
        title: 'Create event filter rule',
        success: 'Rule successfully created!',
      },
      duplicate: {
        title: 'Duplicate event filter rule',
        success: 'Rule successfully created!',
      },
      edit: {
        title: 'Edit an event filter rule',
        success: 'Rule successfully edited!',
      },
      remove: {
        success: 'Rule successfully removed!',
      },
    },
    metaAlarmRule: {
      create: {
        title: 'Create meta alarm rule',
        success: 'Rule successfully created!',
      },
      duplicate: {
        title: 'Duplicate meta alarm rule',
        success: 'Rule successfully created!',
      },
      edit: {
        title: 'Edit a meta alarm rule',
        success: 'Rule successfully edited!',
      },
      remove: {
        success: 'Rule successfully removed!',
      },
      editPattern: 'Edit pattern',
      actions: 'Actions',
    },
    viewTab: {
      create: {
        title: 'Create tab',
      },
      edit: {
        title: 'Edit tab',
      },
      duplicate: {
        title: 'Duplicate tab',
      },
      fields: {
        title: 'Title',
      },
    },
    createSnmpRule: {
      create: {
        title: 'Create SNMP rule',
      },
      edit: {
        title: 'Edit SNMP rule',
      },
    },
    selectView: {
      title: 'Select view',
    },
    selectViewTab: {
      title: 'Select tab',
    },
    createDynamicInfo: {
      alarmUpdate: 'The rule will update existing alarms!',
      create: {
        title: 'Create dynamic information',
        success: 'Dynamic information successfully created!',
      },
      edit: {
        title: 'Edit dynamic information',
        success: 'Dynamic information successfully edited!',
      },
      duplicate: {
        title: 'Duplicate dynamic information',
      },
      remove: {
        success: 'Dynamic information successfully removed!',
      },
      errors: {
        invalid: 'Invalid',
        emptyInfos: 'At least one info must be added.',
      },
      steps: {
        infos: {
          title: 'Informations',
        },
        patterns: {
          title: 'Patterns',
          alarmPatterns: 'Alarm patterns',
          entityPatterns: 'Entity patterns',
          validationError: 'At least one pattern must be set. Please add an alarm pattern and/or an entity pattern',
        },
      },
    },
    createDynamicInfoInformation: {
      create: {
        title: 'Add an information to the dynamic information rule',
      },
    },
    dynamicInfoTemplatesList: {
      title: 'Dynamic info templates',
    },
    createDynamicInfoTemplate: {
      create: {
        title: 'Create dynamic info template',
      },
      edit: {
        title: 'Edit dynamic info template',
      },
      fields: {
        names: 'Names',
      },
      buttons: {
        addName: 'Add new name',
      },
      errors: {
        noNames: 'You have to add at least 1 name',
      },
      emptyNames: 'No names added yet',
    },
    importExportViews: {
      title: 'Import/Export views',
      groups: 'Groups',
      views: 'Views',
    },
    createBroadcastMessage: {
      create: {
        title: 'Create broadcast message',
      },
      edit: {
        title: 'Edit broadcast message',
      },
      defaultMessage: 'Your message here',
    },
    createCommentEvent: {
      title: 'Add comment',
    },
    createPlaylist: {
      create: {
        title: 'Create playlist',
      },
      edit: {
        title: 'Edit playlist',
      },
      duplicate: {
        title: 'Duplicate playlist',
      },
      errors: {
        emptyTabs: 'You should add a tab',
      },
      fields: {
        interval: 'Interval',
        unit: 'Unit',
      },
      groups: 'Groups',
      manageTabs: 'Manage tabs',
    },
    pbehaviorPlanning: {
      title: 'Periodical behaviors',
    },
    selectExceptionsLists: {
      title: 'Choose list of exceptions',
    },
    createRrule: {
      title: 'Create recurrence rule',
    },
    createPbehaviorType: {
      title: 'Create type',
      iconNameHint: 'Enter a name of an icon from material.io',
      errors: {
        iconName: 'The name is invalid',
      },
      fields: {
        name: 'Name',
        description: 'Description',
        type: 'Type',
        priority: 'Priority',
        iconName: 'Icon name',
      },
    },
    pbehaviorRecurrentChangesConfirmation: {
      title: 'Modify',
      fields: {
        selected: 'Only selected period',
        all: 'All the periods',
      },
    },
    createPbehaviorReason: {
      title: 'Create reason',
      fields: {
        name: 'Name',
        description: 'Description',
      },
    },
    createPbehaviorException: {
      title: 'Create date of exception',
      addDate: 'Add date',
      fields: {
        name: 'Name',
        description: 'Description',
      },
      emptyExdates: 'No exdates added yet',
    },
    createManualMetaAlarm: {
      title: 'Manual meta alarm management',
      noData: 'No meta alarm corresponding. Press <kbd>enter</kbd> to create a new one',
      fields: {
        metaAlarm: 'Manual meta alarm',
      },
    },
    createRemediationInstruction: {
      create: {
        title: 'Create instruction',
        popups: {
          success: '{instructionName} has been successfully created',
        },
      },
      edit: {
        title: 'Modify instruction',
        popups: {
          success: '{instructionName} has been successfully modified',
        },
      },
      duplicate: {
        title: 'Duplicate instruction',
        popups: {
          success: '{instructionName} has been successfully duplicated',
        },
      },
    },
    createRemediationConfiguration: {
      create: {
        title: 'Create configuration',
        popups: {
          success: '{configurationName} has been successfully modified',
        },
      },
      edit: {
        title: 'Modify configuration',
        popups: {
          success: '{configurationName} has been successfully modified',
        },
      },
      duplicate: {
        title: 'Duplicate configuration',
        popups: {
          success: '{configurationName} has been successfully duplicated',
        },
      },
      fields: {
        host: 'Host',
        token: 'Authorization token',
      },
    },
    createRemediationJob: {
      create: {
        title: 'Create Job',
        popups: {
          success: '{jobName} has been successfully modified',
        },
      },
      edit: {
        title: 'Modify Job',
        popups: {
          success: '{jobName} has been successfully modified',
        },
      },
      duplicate: {
        title: 'Duplicate Job',
        popups: {
          success: '{jobName} has been successfully duplicated',
        },
      },
    },
    clickOutsideConfirmation: {
      title: 'Are you sure?',
      text: 'Changes will not be saved. Are you sure?',
      buttons: {
        save: 'Save',
        dontSave: 'Don\'t save',
        backToForm: 'Back to form',
      },
    },
    patterns: {
      title: 'Assign patterns',
    },
    rateInstruction: {
      title: 'Rate this instruction "{name}"',
      text: 'How useful was this instruction?',
    },
    createScenario: {
      create: {
        title: 'Create scenario',
        success: 'Scenario created!',
      },
      edit: {
        title: 'Modify scenario',
        success: 'Scenario modified!',
      },
      duplicate: {
        title: 'Duplicate scenario',
        success: 'Scenario duplicated!',
      },
      remove: {
        success: 'Scenario deleted!',
      },
    },
    serviceDependencies: {
      impacts: {
        title: 'Impacts for {name}',
      },
      dependencies: {
        title: 'Dependencies for {name}',
      },
    },
    stateSetting: {
      title: 'JUnit test suite state settings',
    },
    defineStorage: {
      title: 'Define result storage',
      field: {
        placeholder: 'Input the path to the result folder',
      },
    },
    defineXMLStorage: {
      title: 'Define XML storage',
      field: {
        placeholder: 'Input the path to the XML folder',
      },
    },
    defineScreenshotStorage: {
      title: 'Define screenshots storage',
      field: {
        placeholder: 'Input the path to the screenshots folder',
      },
    },
    defineVideoStorage: {
      title: 'Define video storage',
      field: {
        placeholder: 'Input the path to the video folder',
      },
    },
    remediationInstructionApproval: {
      title: 'Instruction approval',
      requested: 'requested for approval',
      tabs: {
        updated: 'Updated',
        original: 'Original',
      },
    },
    createAlarmIdleRule: {
      create: {
        title: 'Create alarm rule',
      },
      edit: {
        title: 'Edit alarm rule',
      },
      duplicate: {
        title: 'Duplicate alarm rule',
      },
    },
    createEntityIdleRule: {
      create: {
        title: 'Create entity rule',
      },
      edit: {
        title: 'Edit entity rule',
      },
      duplicate: {
        title: 'Duplicate entity rule',
      },
    },
    createAlarmStatusRule: {
      flapping: {
        create: {
          title: 'Create flapping rule',
        },
        edit: {
          title: 'Edit flapping rule',
        },
        duplicate: {
          title: 'Duplicate flapping rule',
        },
      },
      resolve: {
        create: {
          title: 'Create resolve rule',
        },
        edit: {
          title: 'Edit resolve rule',
        },
        duplicate: {
          title: 'Duplicate resolve rule',
        },
      },
    },
    webSocketError: {
      title: 'WebSocket connection error',
      text: '<p>Websockets are unavailable, so the following functionalities are restricted:</p>'
        + '<p>'
        + '<ul>'
        + '<li>Healthcheck header</li>'
        + '<li>Healthcheck network graph</li>'
        + '<li>Active broadcast messages</li>'
        + '<li>Active users sessions</li>'
        + '<li>Remediation execution</li>'
        + '</ul>'
        + '</p>'
        + '<p>Please check your server configuration.</p>',
      shortText: '<p>Websockets are unavailable, so the following functionalities are restricted:</p>'
        + '<p>'
        + '<ul>'
        + '<li>Active broadcast messages</li>'
        + '<li>Active users sessions</li>'
        + '</ul>'
        + '</p>'
        + '<p>Please check your server configuration.</p>',
    },
    confirmationPhrase: {
      phrase: 'Phrase',
      updateStorageSettings: {
        title: 'Updating storage policy. Are you sure ?',
        text: 'You are about to change the storage policy.\n'
          + '<strong>Associated operations, deleting data, won\'t be cancellable.</strong>',
        phraseText: 'Please, type the following to confirm:',
        phrase: 'update the storage policy',
      },
      cleanStorage: {
        title: 'Archive/delete disabled entities. Are you sure ?',
        text: 'You are about to archive and/or delete data.\n'
          + '<strong>Deletion operation won\'t be cancellable.</strong>',
        phraseText: 'Please, type the following to confirm:',
        phrase: 'archive or delete',
      },
    },
    pbehaviorsCalendar: {
      title: 'Periodic behaviors',
      entity: {
        title: 'Periodic behaviors - {name}',
      },
    },
    createAlarmPattern: {
      create: {
        title: 'Create alarm filter',
      },
      edit: {
        title: 'Edit alarm filter',
      },
    },
    createCorporateAlarmPattern: {
      create: {
        title: 'Create shared alarm filter',
      },
      edit: {
        title: 'Edit shared alarm filter',
      },
    },
    createEntityPattern: {
      create: {
        title: 'Create entity filter',
      },
      edit: {
        title: 'Edit entity filter',
      },
    },
    createCorporateEntityPattern: {
      create: {
        title: 'Create shared entity filter',
      },
      edit: {
        title: 'Edit shared entity filter',
      },
    },
    createPbehaviorPattern: {
      create: {
        title: 'Create pbehavior filter',
      },
      edit: {
        title: 'Edit pbehavior filter',
      },
    },
    createCorporatePbehaviorPattern: {
      create: {
        title: 'Create shared pbehavior filter',
      },
      edit: {
        title: 'Edit shared pbehavior filter',
      },
    },
    createMap: {
      title: 'Create a map',
    },
    createGeoMap: {
      create: {
        title: 'Create a geomap',
      },
      edit: {
        title: 'Edit a geomap',
      },
      duplicate: {
        title: 'Duplicate a geomap',
      },
    },
    createFlowchartMap: {
      create: {
        title: 'Create a flowchart',
      },
      edit: {
        title: 'Edit a flowchart',
      },
      duplicate: {
        title: 'Duplicate a flowchart',
      },
    },
    createMermaidMap: {
      create: {
        title: 'Create a mermaid diagram',
      },
      edit: {
        title: 'Edit a mermaid diagram',
      },
      duplicate: {
        title: 'Duplicate a mermaid diagram',
      },
    },
    createTreeOfDependenciesMap: {
      create: {
        title: 'Create a tree of dependencies diagram',
      },
      edit: {
        title: 'Edit a tree of dependencies diagram',
      },
      duplicate: {
        title: 'Duplicate a tree of dependencies diagram',
      },
      addEntity: 'Add entity',
      pinnedEntities: 'Pinned entities',
    },
    createShareToken: {
      create: {
        title: 'Create share token',
      },
    },
  },
  tables: {
    noData: 'No data',
    contextEntities: {
      columns: {
        name: 'Name',
        type: 'Type',
        _id: 'Id',
      },
    },
    noColumns: {
      message: 'You have to select at least 1 column',
    },
    broadcastMessages: {
      statuses: {
        [BROADCAST_MESSAGES_STATUSES.active]: 'Active',
        [BROADCAST_MESSAGES_STATUSES.pending]: 'Pending',
        [BROADCAST_MESSAGES_STATUSES.expired]: 'Expired',
      },
    },
  },
  recurrenceRule: {
    advancedHint: 'Separate numbers with a comma',
    freq: 'Frequency',
    until: 'Until',
    byweekday: 'By week day',
    count: 'Repeat',
    interval: 'Interval',
    wkst: 'Week start',
    bymonth: 'By month',
    bysetpos: 'By set position',
    bymonthday: 'By month day',
    byyearday: 'By year day',
    byweekno: 'By week n°',
    byhour: 'By hour',
    byminute: 'By minute',
    bysecond: 'By second',
    tabs: {
      simple: 'Simple',
      advanced: 'Advanced',
    },
    errors: {
      main: 'Please note that the recurrence rule you chose is not valid. We strongly advise you to modify it before saving changes.',
    },
    periodsRanges: {
      [PBEHAVIOR_RRULE_PERIODS_RANGES.thisWeek]: 'This week',
      [PBEHAVIOR_RRULE_PERIODS_RANGES.nextWeek]: 'Next week',
      [PBEHAVIOR_RRULE_PERIODS_RANGES.next2Weeks]: 'Next 2 weeks',
      [PBEHAVIOR_RRULE_PERIODS_RANGES.thisMonth]: 'This month',
      [PBEHAVIOR_RRULE_PERIODS_RANGES.nextMonth]: 'Next month',
    },
    tooltips: {
      bysetpos: 'If given, it must be one or many integers, positive or negative. Each given integer will specify an occurrence number, corresponding to the nth occurrence of the rule inside the frequency period. For example, a \'bysetpos\' of -1 if combined with a monthly frequency, and a \'byweekday\' of (Monday, Tuesday, Wednesday, Thursday, Friday), will result in the last work day of every month.',
      bymonthday: 'If given, it must be one or many integers, meaning the month days to apply the recurrence to.',
      byyearday: 'If given, it must be one or many integers, meaning the year days to apply the recurrence to.',
      byweekno: 'If given, it must be on or many integers, meaning the week numbers to apply the recurrence to. Week numbers have the meaning described in ISO8601, that is, the first week of the year is that containing at least four days of the new year.',
      byhour: 'If given, it must be one or many integers, meaning the hours to apply the recurrence to.',
      byminute: 'If given, it must be one or many integers, meaning the minutes to apply the recurrence to.',
      bysecond: 'If given, it must be one or many integers, meaning the seconds to apply the recurrence to.',
    },
  },
  errors: {
    default: 'Something went wrong...',
    lineNotEmpty: 'This line is not empty',
    JSONNotValid: 'Invalid JSON',
    versionNotFound: 'Unable to get application version',
    statsRequestProblem: 'An error occurred while retrieving stats data',
    statsWrongEditionError: "Stats widgets are not available with 'community' edition",
    socketConnectionProblem: 'Problem with connection to socket server',
    endDateLessOrEqualStartDate: 'End date should be after start date',
    unknownWidgetType: 'Unknown widget type: {type}',
    unique: 'Field must be unique',
    codeEditorProblem: 'Problem with code-editor',
  },
  warnings: {
    authTokenExpired: 'Authentication token was expired',
  },
  calendar: {
    today: 'Today',
    month: 'Month',
    week: 'Week',
    day: 'Day',
    pbehaviorPlanningLegend: {
      title: 'Legend',
      noData: 'There aren\'t any exception dates on calendar',
    },
  },
  success: {
    default: 'Done!',
    createEntity: 'Entity successfully created',
    editEntity: 'Entity successfully edited',
    pathCopied: 'Path copied to clipboard',
    linkCopied: 'Link copied to clipboard',
    authKeyCopied: 'Auth key copied to clipboard',
    widgetIdCopied: 'Widget id copied to clipboard',
  },
  filterEditor: {
    title: 'Filter editor',
    tabs: {
      visualEditor: 'Visual Editor',
      advancedEditor: 'Advanced Editor',
      results: 'Results',
    },
    buttons: {
      addRule: 'Add a rule',
      addGroup: 'Add a group',
      deleteGroup: 'Delete group',
    },
    hints: {
      alarm: {
        service: 'Service',
        connector: 'Connector Type',
        connectorName: 'Connector name',
        component: 'Component',
        resource: 'Resource',
      },
    },
    errors: {
      cantParseToVisualEditor: 'We can\'t parse this filter to Visual Editor',
      invalidJSON: 'Invalid JSON',
      required: 'You need to add at least one valid rule',
    },
  },
  filterSelector: {
    defaultFilter: 'Default filter',
    fields: {
      mixFilters: 'Mix filters',
    },
    buttons: {
      list: 'Manage filters',
    },
  },
  stats: {
    types: {
      [STATS_TYPES.alarmsCreated.value]: 'Alarms created',
      [STATS_TYPES.alarmsResolved.value]: 'Alarms resolved',
      [STATS_TYPES.alarmsCanceled.value]: 'Alarms canceled',
      [STATS_TYPES.alarmsAcknowledged.value]: 'Alarms acknowledged',
      [STATS_TYPES.ackTimeSla.value]: 'Ack time Sla',
      [STATS_TYPES.resolveTimeSla.value]: 'Resolve time Sla',
      [STATS_TYPES.timeInState.value]: 'Time in severity',
      [STATS_TYPES.stateRate.value]: 'Severity rate',
      [STATS_TYPES.mtbf.value]: 'MTBF',
      [STATS_TYPES.currentState.value]: 'Current severity',
      [STATS_TYPES.ongoingAlarms.value]: 'Ongoing alarms',
      [STATS_TYPES.currentOngoingAlarms.value]: 'Current ongoing alarms',
      [STATS_TYPES.currentOngoingAlarmsWithAck.value]: 'Current ongoing alarms with ack',
      [STATS_TYPES.currentOngoingAlarmsWithoutAck.value]: 'Current ongoing alarms without ack',
    },
  },
  eventFilter: {
    externalData: 'External data',
    actionsRequired: 'Please add at least one action',
    configRequired: 'No configuration defined. Please add at least one config parameter',
    idHelp: 'If no id is specified, a unique id will be generated automatically on rule creation',
    editPattern: 'Edit pattern',
    advanced: 'Advanced',
    addAField: 'Add a field',
    simpleEditor: 'Simple editor',
    field: 'Field',
    value: 'Value',
    advancedEditor: 'Advanced editor',
    comparisonRules: 'Comparison rules',
    editActions: 'Edit actions',
    addAction: 'Add an action',
    editAction: 'Edit an action',
    actions: 'Actions',
    onSuccess: 'On success',
    onFailure: 'On failure',
    configuration: 'Configuration',
    resource: 'Resource ID or template',
    component: 'Component ID or template',
    connector: 'Connector ID or template',
    connectorName: 'Connector name or template',
    duringPeriod: 'Applied during this period only',
    enrichmentOptions: 'Enrichment options',
    changeEntityOptions: 'Change entity options',
    noExternalData: 'No external data added yet',
    addExternalData: 'Add external data',
    reference: 'Reference',
    collection: 'Collection',
    externalDataTypes: {
      [EVENT_FILTER_EXTERNAL_DATA_TYPES.mongo]: 'MongoDB collection',
      [EVENT_FILTER_EXTERNAL_DATA_TYPES.api]: 'API',
    },
    externalDataConditionTypes: {
      [EVENT_FILTER_EXTERNAL_DATA_CONDITION_TYPES.select]: 'Select',
      [EVENT_FILTER_EXTERNAL_DATA_CONDITION_TYPES.regexp]: 'Regexp',
    },
    types: {
      [EVENT_FILTER_TYPES.drop]: 'Drop',
      [EVENT_FILTER_TYPES.break]: 'Break',
      [EVENT_FILTER_TYPES.enrichment]: 'Enrichment',
      [EVENT_FILTER_TYPES.changeEntity]: 'Change entity',
    },
    tooltips: {
      addValueRuleField: 'Add value rule field',
      editValueRuleField: 'Edit value rule field',
      addObjectRuleField: 'Add object rule field',
      editObjectRuleField: 'Edit object rule field',
      removeRuleField: 'Remove rule field',
      copyFromHelp: '<p>The accessible variables are: <strong>Event</strong></p>'
        + '<i>For example:</i> <span>"Event.ExtraInfos.datecustom"</span>',
      reference: 'Will be used in actions as <strong>.ExternalData.&lt;Reference&gt;</strong>',
    },
    actionsTypes: {
      [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy]: {
        text: 'Copy a value from a field of event to another',
        message: 'This action is used used to copy the value of a control in an event.',
        description: 'The parameters of the action are:\n- value: the name of the control whose value must be copied. It can be an event field, a subgroup of a regular expression, or an external data.\n- description (optional): the description.\n- name: the name of the event field into which the value must be copied.',
      },
      [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copyToEntityInfo]: {
        text: 'Copy a value from a field of an event to an info of an entity',
        message: 'This action is used to copy the field value of an event to the field of an entity.',
        description: 'The parameters of the action are:\n- description (optional): the description.\n- name: the name of the field of an entity.\n- value: the name of the control whose value must be copied. It can be an event field, a subgroup of a regular expression, or an external data.',
      },
      [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfo]: {
        text: 'Set an info of an entity to a constant',
        message: 'This action is used to set the dynamic information from an entity corresponding to the event.',
        description: 'The parameters of the action are:\n- description (optional): the description.\n- name: the name of the field.\n- value: the value of a field.',
      },
      [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate]: {
        text: 'Set a string info of an entity using a template',
        message: 'This action is used to modify the dynamic information from an entity corresponding to the event.',
        description: 'The parameters of the action are:\n- description (optional): the description\n- name: the name of the field.\n- value: the template used to determine the value of the data item.\nTemplates {{.Event.NomDuChamp}}, regular expressions or external data can be used.',
      },
      [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField]: {
        text: 'Set a field of an event to a constant',
        message: 'This action can be used to modify a field of the event.',
        description: 'The parameters of the action are:\n- description (optional): the description.\n- name: the name of the field.\n- value: the new value of the field.',
      },
      [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate]: {
        text: 'Set a string field of an event using a template',
        message: 'This action allows you to modify an event field from a template.',
        description: 'The parameters of the action are:\n- description (optional): the description.\n- name: the name of the field.\n- value: the template used to determine the value of the field.\n Templates {{.Event.NomDuChamp}}, regular expressions or external data can be used.',
      },
    },
    externalDataValues: {
      [EVENT_FILTER_PATTERN_FIELDS.component]: 'Component',
      [EVENT_FILTER_PATTERN_FIELDS.connector]: 'Connector',
      [EVENT_FILTER_PATTERN_FIELDS.connectorName]: 'Connector name',
      [EVENT_FILTER_PATTERN_FIELDS.resource]: 'Resource',
      [EVENT_FILTER_PATTERN_FIELDS.output]: 'Output',
      [EVENT_FILTER_PATTERN_FIELDS.extraInfos]: 'Extra infos',
    },
  },
  metaAlarmRule: {
    outputTemplate: 'Output template',
    thresholdType: 'Threshold type',
    thresholdRate: 'Threshold rate',
    thresholdCount: 'Threshold count',
    timeInterval: 'Time interval',
    valuePath: 'Value path | Value paths',
    autoResolve: 'Auto resolve',
    idHelp: 'If no id is specified, a unique id will be generated automatically on rule creation',
    corelId: 'Corel ID',
    corelIdHelp: '<p>The accessible variables are: <strong>.Alarm</strong> and <strong>.Entity</strong></p>'
      + '<i>For example:</i> <span>"{{ .Alarm.Value.Connector }}", "{{ .Entity.Component }}"</span>',
    corelStatus: 'Corel status',
    corelStatusHelp: '<p>The accessible variables are: <strong>.Alarm</strong> and <strong>.Entity</strong></p>'
      + '<i>For example:</i> <span>"{{ .Alarm.Value.Connector }}", "{{ .Entity.Component }}"</span>',
    corelParent: 'Corel parent',
    corelChild: 'Corel child',
    outputTemplateHelp: '<p>The accessible variables are:</p>'
      + '<p><strong>.Count</strong>: The number of consequence alarms attached to the meta alarm.</p>'
      + '<p><strong>.Children</strong>: The set of variables of the last consequence alarm attached to the meta alarm.</p>'
      + '<p><strong>.Rule</strong>: The administrative information of the meta alarm itself.</p>'
      + '<p>For example:</p>'
      + '<p>Count: <strong>{{ .Count }};</strong> Children: <strong>{{ .Children.Alarm.Value.State.Message }};</strong> Rule: <strong>{{ .Rule.Name }};</strong></p>'
      + '<p>A static informative message</p>'
      + '<p>Correlated by the rule <strong>{{ .Rule.Name }}</strong></p>',
    removeConfirmationText: 'When deleting a meta alarm rule, all corresponding meta alarms will be deleted as well.\n'
      + 'Are you sure to proceed with it?\n',
    errors: {
      noValuePaths: 'You have to add at least 1 value path',
    },
  },
  layout: {
    sideBar: {
      buttons: {
        edit: 'Toggle editing mode',
        create: 'Create view',
        settings: 'Settings',
      },
      loggedUsersCount: 'Active sessions',
      ordering: {
        popups: {
          success: 'The groups was reordered',
          error: 'Several groups wasn\'t reordered',
          periodicRefreshWasPaused: 'Periodic refresh was paused while you are editing the groups bar',
          periodicRefreshWasResumed: 'Periodic refresh was resumed',
        },
      },
    },
  },
  parameters: {
    tabs: {
      parameters: 'Parameters',
      importExportViews: 'Import/Export',
      stateSettings: 'State settings',
      storageSettings: 'Storage settings',
      notificationsSettings: 'Notifications settings',
    },
  },
  view: {
    errors: {
      emptyTabs: 'You should create a tab',
    },
    sharedViewUrl: 'Shared view url',
    shareView: 'Share view {name}',
    deleteRow: 'Delete row',
    deleteWidget: 'Delete widget',
    fullScreen: 'Full screen',
    fullScreenShortcut: 'Alt + Enter / Command + Enter',
    copyWidgetId: 'Copy widget ID',
    autoHeightButton: 'If this button is selected, height will be automatically calculated.',
  },
  validation: {
    messages: {
      _default: 'The value is not valid',
      after: 'The must be after {1}',
      after_with_inclusion: 'The must be after or equal to {1}',
      alpha: 'The field may only contain alphabetic characters',
      alpha_dash: 'The field may contain alpha-numeric characters as well as dashes and underscores',
      alpha_num: 'The field may only contain alpha-numeric characters',
      alpha_spaces: 'The field may only contain alphabetic characters as well as spaces',
      before: 'The must be before {1}',
      before_with_inclusion: 'The must be before or equal to {1}',
      between: 'The field must be between {1} and {2}',
      confirmed: 'The confirmation does not match',
      credit_card: 'The field is invalid',
      date_between: 'The must be between {1} and {2}',
      date_format: 'The must be in the format {1}',
      decimal: 'The field must be numeric and may contain {1} decimal points',
      digits: 'The field must be numeric and contains exactly {1} digits',
      dimensions: 'The field must be {1} pixels by {2} pixels',
      email: 'The field must be a valid email',
      excluded: 'The field must be a valid value',
      ext: 'The field must be a valid file',
      image: 'The field must be an image',
      included: 'The field must be a valid value',
      integer: 'The field must be an integer',
      ip: 'The field must be a valid ip address',
      ip_or_fqdn: 'The field must be a valid ip address or FQDN',
      length: 'The length must be {1}',
      max: 'The field may not be greater than {1} characters',
      max_value: 'The field must be {1} or less',
      mimes: 'The field must have a valid file type',
      min: 'The field must be at least {1} characters',
      min_value: 'The field must be {1} or more',
      numeric: 'The field may only contain numeric characters',
      regex: 'The field format is invalid',
      required: 'The field is required',
      required_if: 'The field is required when the {1} field has this value',
      size: 'The size must be less than {1}KB',
      url: 'The field is not a valid URL',
    },
    custom: {
      tstop: {
        after: 'End time should be later than {1}',
      },
      logo: {
        size: 'The {0} size must be less than {1} KB.',
      },
    },
  },
  home: {
    popups: {
      info: {
        noAccessToDefaultView: 'Access to default view forbidden. Redirecting to role default view.',
        notSelectedRoleDefaultView: 'No role default view selected.',
        noAccessToRoleDefaultView: 'Access to role default view forbidden.',
      },
    },
  },
  serviceWeather: {
    seeAlarms: 'See alarms',
    grey: 'Gray',
    primaryIcon: 'Primary icon',
    secondaryIcon: 'Secondary icon',
    massActions: 'Mass actions',
    cannotBeApplied: 'This action cannot be applied',
    actions: {
      [WEATHER_ACTIONS_TYPES.entityAck]: 'Acknowledge',
      [WEATHER_ACTIONS_TYPES.entityAckRemove]: 'Cancel acknowledge',
      [WEATHER_ACTIONS_TYPES.entityValidate]: 'Validate',
      [WEATHER_ACTIONS_TYPES.entityInvalidate]: 'Invalidate',
      [WEATHER_ACTIONS_TYPES.entityPause]: 'Pause',
      [WEATHER_ACTIONS_TYPES.entityPlay]: 'Play',
      [WEATHER_ACTIONS_TYPES.entityCancel]: 'Cancel',
      [WEATHER_ACTIONS_TYPES.entityAssocTicket]: 'Associate ticket',
      [WEATHER_ACTIONS_TYPES.entityComment]: 'Comment',
      [WEATHER_ACTIONS_TYPES.executeInstruction]: 'Execute instruction',
      [WEATHER_ACTIONS_TYPES.declareTicket]: 'Declare ticket',
    },
    iconTypes: {
      [PBEHAVIOR_TYPE_TYPES.inactive]: 'Inactive',
      [PBEHAVIOR_TYPE_TYPES.pause]: 'Pause',
      [PBEHAVIOR_TYPE_TYPES.maintenance]: 'Maintenance',
    },
    stateCounters: {
      [SERVICE_WEATHER_STATE_COUNTERS.all]: 'Number of alarms',
      [SERVICE_WEATHER_STATE_COUNTERS.active]: 'Number of active alarms',
      [SERVICE_WEATHER_STATE_COUNTERS.depends]: 'Number of dependencies',
      [SERVICE_WEATHER_STATE_COUNTERS.ok]: 'Ok',
      [SERVICE_WEATHER_STATE_COUNTERS.minor]: 'Minor',
      [SERVICE_WEATHER_STATE_COUNTERS.major]: 'Major',
      [SERVICE_WEATHER_STATE_COUNTERS.critical]: 'Critical',
      [SERVICE_WEATHER_STATE_COUNTERS.acked]: 'Acknowledged',
      [SERVICE_WEATHER_STATE_COUNTERS.unacked]: 'Not acknowledged',
      [SERVICE_WEATHER_STATE_COUNTERS.underPbehavior]: 'Under PBh',
      [SERVICE_WEATHER_STATE_COUNTERS.ackedUnderPbehavior]: 'Acknowledged under PBh',
    },
    stateCountersTooltips: {
      [SERVICE_WEATHER_STATE_COUNTERS.all]: 'alarms total',
      [SERVICE_WEATHER_STATE_COUNTERS.active]: 'active alarms',
      [SERVICE_WEATHER_STATE_COUNTERS.depends]: 'dependencies',
      [SERVICE_WEATHER_STATE_COUNTERS.ok]: 'OK state',
      [SERVICE_WEATHER_STATE_COUNTERS.minor]: 'minor alarms',
      [SERVICE_WEATHER_STATE_COUNTERS.major]: 'major alarms',
      [SERVICE_WEATHER_STATE_COUNTERS.critical]: 'critical alarms',
      [SERVICE_WEATHER_STATE_COUNTERS.acked]: 'alarms acked',
      [SERVICE_WEATHER_STATE_COUNTERS.unacked]: 'not acked',
      [SERVICE_WEATHER_STATE_COUNTERS.underPbehavior]: 'under PBh',
      [SERVICE_WEATHER_STATE_COUNTERS.ackedUnderPbehavior]: 'acked under PBh',
    },
  },
  contextGeneralTable: {
    addSelection: 'Add selection',
  },
  liveReporting: {
    button: 'Set a custom date range',
  },
  tours: {
    [TOURS.alarmsExpandPanel]: {
      step1: 'Details',
      step2: 'MoreInfos tab (Displayed only in case of existing confguration)',
      step3: 'Timeline tab',
    },
  },
  handlebars: {
    requestHelper: {
      errors: {
        timeout: 'Request timeout',
        unauthorized: 'Unauthorized',
        other: 'Error while fetching data',
      },
    },
  },
  importExportViews: {
    selectAll: 'Select all groups and views',
  },
  playlist: {
    player: {
      tooltips: {
        fullscreen: 'Actions are disabled in full screen mode',
      },
    },
  },

  permissions: {
    technical: {
      admin: 'Admin rights',
      exploitation: 'Exploitation rights',
      notification: 'Notification rights',
      profile: 'Profile rights',
    },
    business: {
      [USER_PERMISSIONS_PREFIXES.business.common]: 'Rights for common',
      [USER_PERMISSIONS_PREFIXES.business.alarmsList]: 'Rights for Alarms List',
      [USER_PERMISSIONS_PREFIXES.business.context]: 'Rights for Context Explorer',
      [USER_PERMISSIONS_PREFIXES.business.serviceWeather]: 'Rights for Service Weather',
      [USER_PERMISSIONS_PREFIXES.business.counter]: 'Rights for Counter',
      [USER_PERMISSIONS_PREFIXES.business.testingWeather]: 'Rights for Testing Weather',
      [USER_PERMISSIONS_PREFIXES.business.map]: 'Rights for Mapping',
    },
    api: {
      general: 'General',
      rules: 'Rules',
      remediation: 'Remediation',
      pbehavior: 'PBehavior',
    },
  },

  pbehavior: {
    periodsCalendar: 'Calendar with periods',
    buttons: {
      addFilter: 'Add filter',
      editFilter: 'Edit filter',
      addRRule: 'Add recurrence rule',
      editRrule: 'Edit recurrence rule',
    },
  },

  pbehaviorExceptions: {
    title: 'Exception dates',
    create: 'Add an exception date',
    choose: 'Choose list of exceptions',
    usingException: 'Cannot be deleted since it is in use',
    emptyExceptions: 'No exceptions added yet',
  },

  pbehaviorTypes: {
    usingType: 'Cannot be deleted since it is in use',
    defaultType: 'Type is default, because cannot be edited',
    types: {
      [PBEHAVIOR_TYPE_TYPES.active]: 'Active',
      [PBEHAVIOR_TYPE_TYPES.inactive]: 'Inactive',
      [PBEHAVIOR_TYPE_TYPES.pause]: 'Pause',
      [PBEHAVIOR_TYPE_TYPES.maintenance]: 'Maintenance',
    },
  },

  pbehaviorReasons: {
    usingReason: 'Cannot be deleted since it is in use',
  },

  planning: {
    tabs: {
      type: 'Type',
      reason: 'Reason',
      exceptions: 'Exception dates',
    },
  },

  healthcheck: {
    metricsUnavailable: 'Metrics are not collecting',
    notRunning: '{name} is unavailable',
    queueOverflow: 'Queue overflow',
    lackOfInstances: 'Lack of instances',
    diffInstancesConfig: 'Invalid instances configuration',
    queueLength: 'Queue length {queueLength}/{maxQueueLength}',
    instancesCount: 'Instances {instances}/{minInstances}',
    activeInstances: 'Only {instances} is active out of {minInstances}. The optimal number of instances is {optimalInstances}.',
    queueOverflowed: 'Queue is overflowed: {queueLength} messages out of {maxQueueLength}.\nPlease check the instances.',
    engineDown: '{name} is down, the system is not operational.\nPlease check the log or restart the service.',
    engineDownOrSlow: '{name} is down or responds too slow, the system is not operational.\nPlease check the log or restart the instance.',
    timescaleDown: '{name} is down, metrics and KPIs are not collecting.\nPlease check the log or restart the instance.',
    invalidEnginesOrder: 'Invalid engines configuration',
    invalidInstancesConfiguration: 'Invalid instances configuration: engine instances read or write to different queues.\nPlease check the instances.',
    chainConfigurationInvalid: 'Engines chain configuration is invalid.\nRefer below for the correct sequence of engines:',
    queueLimit: 'Queue length limit',
    defineQueueLimit: 'Define the engines queue length limit',
    notifyUsersQueueLimit: 'Users can be notified when the queue length limit is exceeded',
    numberOfInstances: 'Number of instances',
    notifyUsersNumberOfInstances: 'Users can be notified when the number of active instances is less than the minimal value. The optimal number of instances is shown when the engine state is unavailable.',
    messagesHistory: 'FIFO messages processing history',
    messagesLastHour: 'FIFO messages processing for the last hour',
    messagesPerHour: 'messages/hour',
    unknown: 'This system state is unavailable',
    systemStatusChipError: 'The system is not operational',
    systemStatusServerError: 'System configuration is invalid, please contact the administrator',
    systemsOperational: 'All systems are operational',
    validation: {
      max_value: 'The field must be equal or less than the optimal instance count',
      min_value: 'The field must be equal or more than the minimal instance count',
    },
    nodes: {
      [HEALTHCHECK_SERVICES_NAMES.mongo]: {
        name: 'MongoDB',
        edgeLabel: 'Status check',
      },

      [HEALTHCHECK_SERVICES_NAMES.rabbit]: {
        name: 'RabbitMQ',
        edgeLabel: 'Status check',
      },

      [HEALTHCHECK_SERVICES_NAMES.redis]: {
        name: 'Redis',
        edgeLabel: 'FIFO data\nRedis check',
      },

      [HEALTHCHECK_SERVICES_NAMES.events]: {
        name: 'Events',
      },

      [HEALTHCHECK_SERVICES_NAMES.api]: {
        name: 'Canopsis API',
      },

      [HEALTHCHECK_SERVICES_NAMES.enginesChain]: {
        name: 'Engines chain',
      },

      [HEALTHCHECK_SERVICES_NAMES.healthcheck]: {
        name: 'Healthcheck',
      },

      [HEALTHCHECK_ENGINES_NAMES.webhook]: {
        name: 'Webhook',
        description: 'Triggers the webhooks launch',
      },

      [HEALTHCHECK_ENGINES_NAMES.fifo]: {
        name: 'FIFO',
        edgeLabel: 'RabbitMQ status\nIncomming flow KPIs',
        description: 'Manages the queue of events and alarms',
      },

      [HEALTHCHECK_ENGINES_NAMES.axe]: {
        name: 'AXE',
        description: 'Creates alarms and performs actions with them',
      },

      [HEALTHCHECK_ENGINES_NAMES.che]: {
        name: 'CHE',
        description: 'Applies eventfilters and created entities',
      },

      [HEALTHCHECK_ENGINES_NAMES.pbehavior]: {
        name: 'Pbehavior',
        description: 'Checks if the alarm is under PBehvaior',
      },

      [HEALTHCHECK_ENGINES_NAMES.action]: {
        name: 'Action',
        description: 'Triggers the actions launch',
      },

      [HEALTHCHECK_ENGINES_NAMES.service]: {
        name: 'Service',
        description: 'Updates counters and generates service-events',
      },

      [HEALTHCHECK_ENGINES_NAMES.dynamicInfos]: {
        name: 'Dynamic infos',
        description: 'Adds dynamic infos to alarm',
      },

      [HEALTHCHECK_ENGINES_NAMES.correlation]: {
        name: 'Correlation',
        description: 'Adds dynamic infos to alarm',
      },

      [HEALTHCHECK_ENGINES_NAMES.remediation]: {
        name: 'Remediation',
        description: 'Triggers the instructions',
      },
    },
  },

  remediation: {
    tabs: {
      configurations: 'Configurations',
      jobs: 'Jobs',
      statistics: 'Remediation statistics',
    },
  },

  remediationInstructions: {
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
    needApprove: 'Approval is needed',
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
      runningInstruction: 'New changes cannot be applied to the instruction in progress. Would you like to cancel started instruction and apply new changes?',
      operationRequired: 'Please add at least one operation',
      stepRequired: 'Please add at least one step',
      jobRequired: 'Please add at least one job',
    },
  },

  remediationJobs: {
    addJobs: 'Add {count} job | Add {count} jobs',
    usingJob: 'Cannot be deleted since it is in use',
    table: {
      configuration: 'Configuration',
      jobId: 'Job ID',
    },
  },

  remediationConfigurations: {
    usingConfiguration: 'Cannot be deleted since it is in use',
    table: {
      host: 'Host',
    },
  },

  remediationInstructionExecute: {
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

  remediationInstructionsFilters: {
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

  remediationInstructionStats: {
    alarmsTimeline: 'Alarms timeline',
    alarmId: 'Alarm ID',
    executedAt: 'Executed at',
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
    actions: {
      needRate: 'Rate it!',
      rate: 'Rate',
    },
  },

  remediationPatterns: {
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

  remediationJob: {
    configuration: 'Configuration',
    jobId: 'Job ID',
    query: 'Query',
    multipleExecutions: 'Allow parallel execution',
    retryAmount: 'Retry amount',
    retryInterval: 'Retry interval',
    addPayload: 'Add payload',
    deletePayload: 'Delete payload',
    payloadHelp: '<p>The accessible variables are: <strong>.Alarm</strong> and <strong>.Entity</strong></p>'
      + '<i>For example:</i>'
      + '<pre>{\n  resource: "{{ .Alarm.Value.Resource }}",\n  entity: "{{ .Entity.ID }}"\n}</pre>',
    errors: {
      invalidJSON: 'Invalid JSON',
    },
  },

  remediationStatistic: {
    remediation: 'Remediation',
    fields: {
      all: 'All',
    },
    labels: {
      remediated: 'Remediated',
      notRemediated: 'Not remediated',
    },
    tooltips: {
      remediated: '{value} alarms remediated',
      assigned: '{value} alarms with instructions',
    },
  },

  scenario: {
    triggers: 'Triggers',
    emitTrigger: 'Emit trigger',
    withAuth: 'Do you need auth fields?',
    emptyResponse: 'Empty response',
    isRegexp: 'The value can be a RegExp',
    headerKey: 'Header key',
    headerValue: 'Header value',
    key: 'Key',
    skipVerify: 'Ignore HTTPS certificate verification',
    headers: 'Headers',
    declareTicket: 'Declare ticket',
    workflow: 'Workflow if this action didn’t match:',
    remainingAction: 'Continue with remaining actions',
    addAction: 'Add action',
    emptyActions: 'No actions added yet',
    output: 'Output Action Format',
    forwardAuthor: 'Forward author to the next step',
    urlHelp: '<p>The accessible variables are: <strong>.Alarm</strong>, <strong>.Entity</strong> and <strong>.Children</strong></p>'
      + '<i>For example:</i>'
      + '<pre>"https://exampleurl.com?resource={{ .Alarm.Value.Resource }}"</pre>'
      + '<pre>"https://exampleurl.com?entity_id={{ .Entity.ID }}"</pre>'
      + '<pre>"https://exampleurl.com?children_count={{ len .Children }}"</pre>'
      + '<pre>"https://exampleurl.com?children={{ range .Children }}{{ .ID }}{{ end }}"</pre>',
    outputHelp: '<p>The accessible variables are: <strong>.Alarm</strong> and <strong>.Entity</strong></p>'
      + '<i>For example:</i>'
      + '<pre>Resource - {{ .Alarm.Value.Resource }}. Entity - {{ .Entity.ID }}.</pre>',
    payloadHelp: '<p>The accessible variables are: <strong>.Alarm</strong>, <strong>.Entity</strong> and <strong>.Children</strong></p>'
      + '<i>For example:</i>'
      + '<pre>{\n'
      + '  resource: "{{ .Alarm.Value.Resource }}",\n'
      + '  entity: "{{ .Entity.ID }}",\n'
      + '  children_count: "{{ len .Children }}",\n'
      + '  children: {{ range .Children }}{{ .ID }}{{ end }}\n'
      + '}</pre>',
    actions: {
      [ACTION_TYPES.snooze]: 'Snooze',
      [ACTION_TYPES.pbehavior]: 'Pbehavior',
      [ACTION_TYPES.changeState]: 'Change state (Change and lock severity)',
      [ACTION_TYPES.ack]: 'Acknowledge',
      [ACTION_TYPES.ackremove]: 'Acknowledge remove',
      [ACTION_TYPES.assocticket]: 'Associate ticket',
      [ACTION_TYPES.cancel]: 'Cancel',
      [ACTION_TYPES.webhook]: 'Webhook',
    },
    tabs: {
      pattern: 'Pattern',
    },
    errors: {
      actionRequired: 'Please add at least one action',
      priorityExist: 'The priority of current scenario is already in use. Do you want to change the current scenario priority to {priority}?',
    },
  },

  mixedField: {
    types: {
      [PATTERN_FIELD_TYPES.string]: '@:variableTypes.string',
      [PATTERN_FIELD_TYPES.number]: '@:variableTypes.number',
      [PATTERN_FIELD_TYPES.boolean]: '@:variableTypes.boolean',
      [PATTERN_FIELD_TYPES.null]: '@:variableTypes.null',
      [PATTERN_FIELD_TYPES.stringArray]: '@:variableTypes.array',
    },
  },

  entity: {
    manageInfos: 'Manage Infos',
    form: 'Form',
    impact: 'Impact',
    depends: 'Depends',
    addInformation: 'Add Information',
    emptyInfos: 'No information',
    availabilityState: 'Hi availability state',
    types: {
      [ENTITY_TYPES.component]: 'Component',
      [ENTITY_TYPES.connector]: 'Connector',
      [ENTITY_TYPES.resource]: 'Resource',
      [ENTITY_TYPES.service]: 'Service',
    },
  },

  service: {
    outputTemplate: 'Output template',
    createCategory: 'Add new category',
    createCategoryHelp: 'Press <kbd>enter</kbd> to save',
    availabilityState: 'Hi availability state',
  },

  users: {
    seeProfile: 'See profile',
    selectDefaultView: 'Select default view',
    firstName: 'First name',
    lastName: 'Last name',
    email: 'Email',
    language: 'User interface language',
    auth: 'Auth',
    navigationType: 'Groups navigation type',
    active: 'Session active',
    activeConnects: 'Connections count',
    navigationTypes: {
      [GROUPS_NAVIGATION_TYPES.sideBar]: 'Side bar',
      [GROUPS_NAVIGATION_TYPES.topBar]: 'Top bar',
    },
    metrics: {
      [USER_METRIC_PARAMETERS.totalUserActivity]: 'Total activity time',
    },
  },

  role: {
    expirationSettings: 'Expiration settings',
    inactivityInterval: 'Inactivity interval',
    expirationInterval: 'Expiration interval',
    inactivityIntervalHelpText: 'Defines when the user is counted as inactive',
    expirationIntervalHelpText: 'Defines the inactivity time period after which the auth token is expired',
  },

  testSuite: {
    xmlFeed: 'XML feed',
    hostname: 'Host name',
    lastUpdate: 'Last update',
    totalTests: 'Total tests',
    disabledTests: 'Tests disabled',
    copyMessage: 'Copy system message',
    systemError: 'System error',
    systemErrorMessage: 'System error message',
    systemOut: 'System out',
    systemOutMessage: 'System out message',
    compareWithHistorical: 'Compare with historical data',
    className: 'Classname',
    line: 'Line',
    failureMessage: 'Failure message',
    noData: 'No system messages found in XML',
    tabs: {
      globalMessages: 'Global messages',
      gantt: 'Gantt',
      details: 'Details',
      screenshots: 'Screenshots',
      videos: 'Videos',
    },
    statuses: {
      [TEST_SUITE_STATUSES.passed]: 'Passed',
      [TEST_SUITE_STATUSES.skipped]: 'Skipped',
      [TEST_SUITE_STATUSES.error]: 'Error',
      [TEST_SUITE_STATUSES.failed]: 'Failed',
      [TEST_SUITE_STATUSES.total]: 'Total time taken',
    },
    popups: {
      systemMessageCopied: 'System message copied to clipboard',
    },
  },

  stateSetting: {
    worstLabel: 'The worst of:',
    worstHelpText: 'Canopsis counts the state for each criterion defined. The final state of JUnit test suite is taken as a worst of resulting states.',
    criterion: 'Criterion',
    serviceState: 'Service state',
    methods: {
      [STATE_SETTING_METHODS.worst]: 'Worst',
      [STATE_SETTING_METHODS.worstOfShare]: 'Worst of share',
    },
    states: {
      minor: 'Minor',
      major: 'Major',
      critical: 'Critical',
    },
  },

  storageSettings: {
    alarm: {
      title: 'Alarm data storage',
      titleHelp: 'When switched on, the resolved alarms data will be archived and/or deleted after the defined time period.',
      archiveAfter: 'Archive resolved alarms data after',
      deleteAfter: 'Delete resolved alarms data after',
    },
    junit: {
      title: 'JUnit data storage',
      deleteAfter: 'Delete test suites data after',
      deleteAfterHelpText: 'When switched on, the JUnit test suites data (XMLs, screenshots and videos) will be deleted after the defined time period.',
    },
    remediation: {
      title: 'Instructions data storage',
      accumulateAfter: 'Accumulate instructions statistics after',
      deleteAfter: 'Delete instructions data after',
      deleteAfterHelpText: 'When switched on, the instructions statistical data will be deleted after the defined time period.',
    },
    entity: {
      title: 'Entities data storage',
      titleHelp: 'All disabled entities with associated alarms can be archived (moved to the separate collection) and/or deleted forever.',
      archiveEntity: 'Archive disabled entities',
      deleteEntity: 'Delete disabled entities forever from archive',
      archiveDependencies: 'Remove the impacting and dependent entities as well',
      archiveDependenciesHelp: 'For connectors, all impacting and dependent components and resources will be archived or deleted forever. For components, all dependent resources will be archived or deleted forever as well.',
      cleanStorage: 'Clean storage',
    },
    pbehavior: {
      title: 'PBehavior data storage',
      deleteAfter: 'Delete PBehavior data after',
      deleteAfterHelpText: 'When switched on, inactive PBehaviors will be deleted after the defined time period from the last event.',
    },
    healthCheck: {
      title: 'Healthcheck data storage',
      deleteAfter: 'Delete FIFO incoming flow data after',
    },
    history: {
      scriptLaunched: 'Script launched at {launchedAt}.',
      alarm: {
        deletedCount: 'Alarms deleted: {count}.',
        archivedCount: 'Alarms archived: {count}.',
      },
      entity: {
        deletedCount: 'Entities deleted: {count}.',
        archivedCount: 'Entities archived: {count}.',
      },
    },
  },

  notificationSettings: {
    instruction: {
      header: 'Instructions',
      rate: '”Rate the instruction” notifications',
      rateFrequency: 'Frequency',
      duration: 'Time range',
    },
  },

  quickRanges: {
    title: 'Quick ranges',
    timeField: 'Time field',
    types: {
      [QUICK_RANGES.custom.value]: 'Custom',
      [QUICK_RANGES.last15Minutes.value]: 'Last 15 minutes',
      [QUICK_RANGES.last30Minutes.value]: 'Last 30 minutes',
      [QUICK_RANGES.last1Hour.value]: 'Last 1 hour',
      [QUICK_RANGES.last3Hour.value]: 'Last 3 hour',
      [QUICK_RANGES.last6Hour.value]: 'Last 6 hour',
      [QUICK_RANGES.last12Hour.value]: 'Last 12 hour',
      [QUICK_RANGES.last24Hour.value]: 'Last 24 hour',
      [QUICK_RANGES.last2Days.value]: 'Last 2 days',
      [QUICK_RANGES.last7Days.value]: 'Last 7 days',
      [QUICK_RANGES.last30Days.value]: 'Last 30 days',
      [QUICK_RANGES.last1Year.value]: 'Last 1 year',
      [QUICK_RANGES.yesterday.value]: 'Yesterday',
      [QUICK_RANGES.previousWeek.value]: 'Previous week',
      [QUICK_RANGES.previousMonth.value]: 'Previous month',
      [QUICK_RANGES.today.value]: 'Today',
      [QUICK_RANGES.todaySoFar.value]: 'Today so far',
      [QUICK_RANGES.thisWeek.value]: 'This week',
      [QUICK_RANGES.thisWeekSoFar.value]: 'This week so far',
      [QUICK_RANGES.thisMonth.value]: 'This month',
      [QUICK_RANGES.thisMonthSoFar.value]: 'This month so far',
    },
  },

  idleRules: {
    timeAwaiting: 'Time awaiting',
    timeRangeAwaiting: 'Time range awaiting',
    types: {
      [IDLE_RULE_TYPES.alarm]: 'Alarm rule',
      [IDLE_RULE_TYPES.entity]: 'Entity rule',
    },
    alarmConditions: {
      [IDLE_RULE_ALARM_CONDITIONS.lastEvent]: 'No events received',
      [IDLE_RULE_ALARM_CONDITIONS.lastUpdate]: 'No state changes',
    },
  },

  alarmStatusRules: {
    frequencyLimit: 'Frequency limit',
  },

  icons: {
    noEvents: 'No events received for {duration} by some of dependencies',
  },

  pageHeaders: {
    hideMessage: 'Got it! Hide',
    learnMore: 'Learn more on {link}',

    /**
     * Exploitation
     */
    [USERS_PERMISSIONS.technical.exploitation.eventFilter]: {
      title: 'Event filter',
      message: 'The event-filter is a feature of engine-che, allowing to define rules handling events.',
    },

    [USERS_PERMISSIONS.technical.exploitation.dynamicInfo]: {
      title: 'Dynamic informations',
      message: 'The Canopsis Dynamic infos are used to add information to the alarms. This information is defined with rules indicating under which conditions information must be presented on an alarm.',
    },

    [USERS_PERMISSIONS.technical.exploitation.metaAlarmRule]: {
      title: 'Meta alarm rule',
      message: 'Meta alarm rules can be used for grouping alarms by types and criteria (parent-child relationship, time interval, etc).',
    },

    [USERS_PERMISSIONS.technical.exploitation.idleRules]: {
      title: 'Idle rules',
      message: 'Idle rules for entities and alarms can be used in order to monitor events and alarm states in order to be aware when events are not receiving or alarm state is not changed for a long time because of errors or invalid configuration.',
    },

    [USERS_PERMISSIONS.technical.exploitation.flappingRules]: {
      title: 'Flapping rules',
      // message: '', // TODO: need to put description
    },

    [USERS_PERMISSIONS.technical.exploitation.resolveRules]: {
      title: 'Resolve rules',
      // message: '', // TODO: need to put description
    },

    [USERS_PERMISSIONS.technical.exploitation.pbehavior]: {
      title: 'PBehaviors',
      message: 'Canopsis periodical behaviors can be used in order to define a periods when the behavior has to be changed, e.g. for  maintenance or service range.',
    },

    [USERS_PERMISSIONS.technical.exploitation.scenario]: {
      title: 'Scenarios',
      message: 'The Canopsis scenarios can be used to conditionally trigger various types of actions on alarms.',
    },

    [USERS_PERMISSIONS.technical.exploitation.snmpRule]: {
      title: 'SNMP rules',
      message: 'The SNMP engine allows the processing of SNMP traps retrieved by the connector snmp2canopsis.',
    },

    /**
     * Admin access
     */
    [USERS_PERMISSIONS.technical.permission]: {
      title: 'Rights',
    },
    [USERS_PERMISSIONS.technical.role]: {
      title: 'Roles',
    },
    [USERS_PERMISSIONS.technical.user]: {
      title: 'Users',
    },

    /**
     * Admin communications
     */
    [USERS_PERMISSIONS.technical.broadcastMessage]: {
      title: 'Broadcast messages',
      message: 'The Canopsis broadcasting messages can be used for displaying banners and information messages that will appear in the Canopsis interface.',
    },
    [USERS_PERMISSIONS.technical.playlist]: {
      title: 'Playlists',
      message: 'Playlists can be used for the views customization which can be displayed one after another with an associated delay.',
    },
    [USERS_PERMISSIONS.technical.healthcheck]: {
      title: 'Healthcheck',
      message: 'The Healthcheck feature is the dashboard with states and errors indications of all systems included to the Canopsis.',
    },
    [USERS_PERMISSIONS.technical.engine]: {
      title: 'Engines',
      message: 'This page contains the information about the sequence and configuration of engines. To work properly, the chain of engines must be continuous.',
    },
    [USERS_PERMISSIONS.technical.kpi]: {
      title: 'KPI',
      message: '', // TODO: add correct message
    },
    [USERS_PERMISSIONS.technical.map]: {
      title: 'Maps',
      message: '', // TODO: add correct message
    },

    /**
     * Admin general
     */
    [USERS_PERMISSIONS.technical.parameters]: {
      title: 'Parameters',
    },
    [USERS_PERMISSIONS.technical.planning]: {
      title: 'Planning',
      message: 'The Canopsis Planning Administration functionality can be used for the periodic behavior types customization.',
    },
    [USERS_PERMISSIONS.technical.remediation]: {
      title: 'Instructions',
      message: 'The Canopsis Remediation feature is used for creation plans or instructions to correct situations.',
    },

    /**
     * Notifications
     */
    [USERS_PERMISSIONS.technical.notification.instructionStats]: {
      title: 'Instruction rating',
      message: 'This page contains the statistics on the instructions execution. Users can rate instructions based on their performance.',
    },
  },

  userInterface: {
    title: 'User interface',
    appTitle: 'App title',
    language: 'Default user interface language',
    footer: 'Login footer',
    description: 'Login page description',
    logo: 'Logo',
    infoPopupTimeout: 'Info popup timeout',
    errorPopupTimeout: 'Error popup timeout',
    allowChangeSeverityToInfo: 'Allow change severity to info',
    maxMatchedItems: 'Max matched items',
    checkCountRequestTimeout: 'Check max matched items request timeout (seconds)',
    tooltips: {
      maxMatchedItems: 'it need to warn user when number of items that match patterns is above this value',
      checkCountRequestTimeout: 'it need to define request timeout value for max matched items checking',
    },
  },

  kpi: {
    alarmMetrics: 'Alarm metrics',
    sli: 'SLI',
    metricsNotAvailable: 'TimescaleDB not running. Metrics are not available.',
    noData: 'No data available',
  },

  kpiMetrics: {
    parameter: 'Parameter to compare',
    tooltip: {
      [USER_METRIC_PARAMETERS.totalUserActivity]: '{value} total activity time',

      [ALARM_METRIC_PARAMETERS.createdAlarms]: '{value} created alarms',
      [ALARM_METRIC_PARAMETERS.activeAlarms]: '{value} active alarms',
      [ALARM_METRIC_PARAMETERS.nonDisplayedAlarms]: '{value} non-displayed alarms',
      [ALARM_METRIC_PARAMETERS.instructionAlarms]: '{value} alarms under auto remediation',
      [ALARM_METRIC_PARAMETERS.pbehaviorAlarms]: '{value} alarms under PBehavior',
      [ALARM_METRIC_PARAMETERS.correlationAlarms]: '{value} alarms with correlation',
      [ALARM_METRIC_PARAMETERS.ackAlarms]: '{value} alarms with acks',
      [ALARM_METRIC_PARAMETERS.ackActiveAlarms]: '{value} active alarms with acks',
      [ALARM_METRIC_PARAMETERS.cancelAckAlarms]: '{value} alarms with cancelled acks',
      [ALARM_METRIC_PARAMETERS.ticketActiveAlarms]: '{value} active alarms with acks',
      [ALARM_METRIC_PARAMETERS.withoutTicketActiveAlarms]: '{value} active alarms without tickets',
      [ALARM_METRIC_PARAMETERS.ratioCorrelation]: '{value}% of alarms with auto remediation',
      [ALARM_METRIC_PARAMETERS.ratioInstructions]: '{value}% alarms with instructions',
      [ALARM_METRIC_PARAMETERS.ratioTickets]: '{value}% of alarms with tickets created',
      [ALARM_METRIC_PARAMETERS.ratioNonDisplayed]: '{value}% of non-displayed alarms',
      [ALARM_METRIC_PARAMETERS.ratioRemediatedAlarms]: '{value}% of manually remediated alarms',
      [ALARM_METRIC_PARAMETERS.averageAck]: '{value} to ack alarms',
      [ALARM_METRIC_PARAMETERS.averageResolve]: '{value} to resolve alarms',
      [ALARM_METRIC_PARAMETERS.manualInstructionAssignedAlarms]: '{value} alarms with manual instructions',
      [ALARM_METRIC_PARAMETERS.manualInstructionExecutedAlarms]: '{value} manually remediated alarms',
    },
  },

  kpiFilters: {
    helpInformation: 'Here the filter patterns for additional slices of data for counters and ratings can be added.',
  },

  kpiRatingSettings: {
    helpInformation: 'The list of parameters to use for rating.',
  },

  snmpRule: {
    oid: 'oid',
    module: 'Select a mib module',
    output: 'output',
    resource: 'resource',
    component: 'component',
    connectorName: 'connector_name',
    toCustom: 'To custom',
    defineVar: 'Define matching snmp var',
    writeTemplate: 'Write template',
    state: 'severity',
    moduleMibObjects: 'Snmp vars match field',
    regex: 'Regex',
    formatter: 'Format (capture group with \\x)',
    uploadMib: 'Upload MIB',
    addSnmpRule: 'Add SNMP rule',
  },

  pattern: {
    patterns: 'Patterns',
    myPatterns: 'My patterns',
    corporatePatterns: 'Shared patterns',
    addRule: 'Add rule',
    addGroup: 'Add group',
    removeRule: 'Remove rule',
    advancedEditor: 'Advanced editor',
    simpleEditor: 'Simple editor',
    noData: 'No pattern set. Click \'@:pattern.addGroup\' button to start adding fields to the pattern',
    noDataDisabled: 'No pattern set.',
    discard: 'Discard pattern',
    types: {
      [PATTERN_TYPES.alarm]: 'Alarm pattern',
      [PATTERN_TYPES.entity]: 'Entity pattern',
      [PATTERN_TYPES.pbehavior]: 'Pbehavior pattern',
    },
    errors: {
      ruleRequired: 'Please add at least one rule',
      groupRequired: 'Please add at least one group',
      invalidPatterns: 'Patterns are invalid or there is a disabled pattern field',
      countOverLimit: 'The patterns you\'ve defined targets about {count} items. It can affect performance, are you sure ?',
      oldPattern: 'The current filter pattern is defined in old format. Please use the Advanced editor to view it. Filters in old format will be deprecated soon. Please create new patterns in our updated interface.',
      existExcluded: 'The rules include excluded rule.',
    },
  },

  filter: {
    oldPattern: 'Old pattern format',
  },

  map: {
    defineEntity: 'Define entity',
    addLink: 'Add link',
    addPoint: 'Add point',
    editPoint: 'Edit point',
    removePoint: 'Remove point',
    latitude: 'Latitude',
    longitude: 'Longitude',
    toggleAddingPointMode: 'Toggle adding point mode',
    usingMap: 'Map is linked',
    showAll: 'Show all ({count})',
    types: {
      [MAP_TYPES.geo]: 'Geo',
      [MAP_TYPES.flowchart]: 'Flowchart',
      [MAP_TYPES.mermaid]: 'Mermaid',
      [MAP_TYPES.treeOfDependencies]: 'Tree of dependencies',
    },
    layers: {
      openStreetMap: 'Open street map',
      points: 'Points',
    },
  },

  mermaid: {
    theme: 'Color theme',
    panzoom: {
      helpText: 'Useful shortcuts:\n'
        + 'Ctrl + mouse wheel - zoom in/out\n'
        + 'Shift + mouse wheel - horizontal scroll\n'
        + 'Alt + mouse wheel - vertical scroll\n'
        + 'Ctrl + Left mouse click + drag - pan the area',
    },
    themes: {
      [MERMAID_THEMES.default]: 'Default',
      [MERMAID_THEMES.base]: 'Base',
      [MERMAID_THEMES.dark]: 'Dark',
      [MERMAID_THEMES.forest]: 'Forest',
      [MERMAID_THEMES.neutral]: 'Neutral',
      [MERMAID_THEMES.canopsis]: 'Canopsis',
    },
    errors: {
      emptyMermaid: 'The diagram and points must be added',
    },
  },

  geomap: {
    layers: 'Layers',
    panzoom: {
      helpText: 'Useful shortcuts:\n'
        + 'Ctrl + mouse wheel - zoom in/out\n'
        + 'Left mouse click + drag - pan the area',
    },
    errors: {
      pointsRequired: 'The points must be added',
    },
  },

  flowchart: {
    shape: 'Shape | Shapes',
    icons: 'Icons',
    properties: 'Properties',
    color: 'Color',
    fill: 'Fill',
    stroke: 'Stroke',
    strokeWidth: 'Stroke width',
    strokeType: 'Stroke type',
    fontColor: 'Font color',
    fontSize: 'Font size',
    fontBackgroundColor: 'Font background color',
    lineType: 'Line type',
    backgroundColor: 'Background color',
    shapes: {
      rectangle: 'Rectangle',
      roundedRectangle: 'Rounded rectangle',
      square: 'Square',
      rhombus: 'Rhombus',
      circle: 'Circle',
      ellipse: 'Ellipse',
      parallelogram: 'Parallelogram',
      process: 'Process',
      document: 'Document',
      storage: 'Storage',
      curve: 'Curve',
      curveArrow: 'Curve arrow',
      bidirectionalCurve: 'Bidirectional curve',
      line: 'Line',
      arrowLine: 'Arrow line',
      bidirectionalArrowLine: 'Bidirectional arrow line',
      text: 'Text',
      textbox: 'Textbox',
      image: 'Image',
    },
    panzoom: {
      helpText: 'Useful shortcuts:\n'
        + 'Ctrl + mouse wheel - zoom in/out\n'
        + 'Ctrl + Left mouse click + drag - pan the area\n'
        + 'Middle mouse click + drag - pan the area\n'
        + 'Shift + mouse wheel - horizontal scroll\n'
        + 'Alt + mouse wheel - vertical scroll\n',
    },
    errors: {
      pointsRequired: 'The points must be added',
    },
  },
  treeOfDependencies: {
    panzoom: {
      helpText: 'Useful shortcuts:\n'
        + 'Ctrl + mouse wheel - zoom in/out\n'
        + 'Ctrl + Left mouse click + drag - pan the area\n',
    },
  },

  shareToken: {
    revokeToken: 'Revoke token',
    revokeSelectedTokens: 'Revoke selected tokens',
    tokenExpiration: 'Token expiration',
  },

  techMetric: {
    noDumps: 'No dumps available. Generate a new dump?',
    metricsDisabled: 'Engine\'s metrics are disabled',
    generateDump: 'Generate a new dump',
    downloadDump: 'Download dump',
  },
}, featureService.get('i18n.en'));
