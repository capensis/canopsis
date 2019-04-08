import {
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  EVENT_ENTITY_TYPES,
  STATS_TYPES,
  STATS_CRITICITY,
  STATS_QUICK_RANGES,
} from '@/constants';

export default {
  common: {
    ok: 'Ok',
    undefined: 'Not defined',
    entity: 'Entity',
    watcher: 'Watcher',
    pbehaviors: 'PBehaviors',
    widget: 'Widget',
    addWidget: 'Add widget',
    addTab: 'Add tab',
    addPbehavior: 'Add pbehavior',
    refresh: 'Refresh',
    toggleEditView: 'Toggle view edition mode',
    name: 'Name',
    description: 'Description',
    author: 'Author',
    submit: 'Submit',
    cancel: 'Cancel',
    options: 'Options',
    type: 'Type',
    quitEditing: 'Quit editing',
    enabled: 'Enabled',
    disabled: 'Disabled',
    login: 'Login',
    yes: 'Yes',
    no: 'No',
    default: 'Default',
    confirmation: 'Are you sure ?',
    parameters: 'Parameters',
    by: 'By',
    date: 'Date',
    comment: 'Comment | Comments',
    end: 'End',
    recursive: 'Recursive',
    select: 'Select',
    states: 'States',
    sla: 'Sla',
    authors: 'Authors',
    stat: 'Stat',
    trend: 'Trend',
    users: 'Users',
    roles: 'Roles',
    rights: 'Rights',
    username: 'Username',
    password: 'Password',
    authKey: 'Auth. key',
    connect: 'Connect',
    optionnal: 'optionnal',
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
    edit: 'Edit',
    parse: 'Parse',
    home: 'Home',
    step: 'Step',
    entries: 'entries',
    showing: 'showing',
    apply: 'Apply',
    to: 'to',
    of: 'of',
    tags: 'tags',
    actionsLabel: 'Actions',
    noResults: 'No results',
    exploitation: 'Exploitation',
    administration: 'Administration',
    forbidden: 'Forbidden',
    search: 'Search',
    filters: 'Filters',
    webhooks: 'Webhooks',
    emptyObject: 'Empty object',
    startDate: 'Start date',
    endDate: 'End date',
    links: 'Links',
    filter: 'Filter',
    actions: {
      close: 'Close',
      acknowledgeAndReport: 'Acknowledge and report an incident',
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
    },
    times: {
      second: 'second | seconds',
      minute: 'minute | minutes',
      hour: 'hour | hours',
      day: 'day | days',
      week: 'week | weeks',
      month: 'month | months',
      year: 'year | years',
    },
  },
  user: {
    role: 'Role',
    defaultView: 'Default view',
  },
  context: {
    impacts: 'Impacts',
    dependencies: 'Dependencies',
    moreInfos: {
      infos: 'Informations',
      type: 'Type',
      enabled: 'Enabled',
      disabled: 'Disabled',
      lastActiveDate: 'Last Active Date',
      infosSearchLabel: 'Search infos',
      tabs: {
        main: 'Main',
        pbehaviors: 'Pbehaviors',
        impactDepends: 'Impact/Depends',
        infos: 'Infos',
      },
    },
    actions: {
      titles: {
        editEntity: 'Edit entity',
        duplicateEntity: 'Duplicate entity',
        deleteEntity: 'Delete entity',
        pbehavior: 'Periodical behavior',
        variablesHelp: 'List of available variables',
      },
    },
  },
  search: {
    advancedSearch: '<span>Help on the advanced research :</span>\n' +
      '<p>- [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt;</p> [ AND|OR [ NOT ] &lt;ColumnName&gt; &lt;Operator&gt; &lt;Value&gt; ]\n' +
      '<p>The "-" before the research is required</p>\n' +
      '<p>Operators :\n' +
      '    <=, <,=, !=,>=, >, LIKE (For MongoDB regular expression)</p>\n' +
      '<p>Value\'s type : String between quote, Boolean ("TRUE", "FALSE"), Integer, Float, "NULL"</p>\n' +
      '<dl><dt>Examples :</dt><dt>- Connector = "connector_1"</dt>\n' +
      '    <dd>Alarms whose connectors are "connector_1"</dd><dt>- Connector="connector_1" AND Resource="resource_3"</dt>\n' +
      '    <dd>Alarms whose connectors is "connector_1" and the ressources is "resource_3"</dd><dt>- Connector="connector_1" OR Resource="resource_3"</dt>\n' +
      '    <dd>Alarms whose connectors is "connector_1" or the ressources is "resource_3"</dd><dt>- Connector LIKE 1 OR Connector LIKE 2</dt>\n' +
      '    <dd>Alarms whose connectors contains 1 or 2</dd><dt>- NOT Connector = "connector_1"</dt>\n' +
      '    <dd>Alarms whose connectors isn\'t "connector_1"</dd>\n' +
      '</dl>',
  },
  entities: {
    watcher: 'watcher',
    entities: 'entities',
  },
  login: {
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
        pbehaviorList: 'List periodic behaviors',
        declareTicket: 'Declare ticket',
        associateTicket: 'Associate ticket',
        cancel: 'Cancel alarm',
        changeState: 'Change criticity',
        moreInfos: 'More infos',
        variablesHelp: 'List of available variables',
      },
      iconsTitles: {
        ack: 'Ack',
        declareTicket: 'Declare ticket',
        canceled: 'Canceled',
        snooze: 'Snooze',
        pbehaviors: 'Periodic behaviors',
      },
      iconsFields: {
        ticketNumber: 'Ticket number',
      },
    },
    infoPopup: 'Info popup',
  },
  pbehaviors: {
    connector: 'Connector',
    connectorName: 'Connector name',
    isEnabled: 'Is enabled',
    begins: 'Begins',
    ends: 'Ends',
    type: 'Type',
    reason: 'Reason',
    rrule: 'Rrule',
    tabs: {
      filter: 'Filter',
      eids: 'Eids',
      comments: 'Comments',
    },
  },
  settings: {
    titles: {
      alarmListSettings: 'Alarm list settings',
      contextTableSettings: 'Context table settings',
      weatherSettings: 'Service weather settings',
      statsHistogramSettings: 'Histogram settings',
      statsCurvesSettings: 'Curve settings',
      statsTableSettings: 'Stats table settings',
      statsCalendarSettings: 'Stats calendar settings',
      statsNumberSettings: 'Stats number settings',
      textSettings: 'Text settings',
    },
    advancedSettings: 'Advanced settings',
    widgetTitle: 'Widget title',
    columnName: 'Column name',
    defaultSortColumn: 'Default sort column',
    columnNames: 'Column names',
    periodicRefresh: 'Periodic refresh',
    defaultNumberOfElementsPerPage: 'Default number of elements/page',
    elementsPerPage: 'Elements per page',
    filterOnOpenResolved: 'Filter on Open/Resolved',
    open: 'Open',
    resolved: 'Resolved',
    filters: 'Filters',
    filterEditor: 'Filter',
    isAckNoteRequired: 'Note field required when ack ?',
    duration: 'Duration',
    tstop: 'End date',
    periodsNumber: 'Number of steps',
    statName: 'Stat name',
    statsSelect: {
      title: 'Stats select',
      required: 'Select at least 1 stat',
    },
    yesNoMode: 'Yes/No mode',
    selectAFilter: 'Select a filter',
    criticityLevels: 'Criticity levels',
    colorsSelector: {
      title: 'Colors selector',
      statsCriticity: {
        [STATS_CRITICITY.ok]: 'ok',
        [STATS_CRITICITY.minor]: 'minor',
        [STATS_CRITICITY.major]: 'major',
        [STATS_CRITICITY.critical]: 'critical',
      },
    },
    statsDateInterval: {
      monthPeriodInfo: "If you select a 'monthly' period, start and end date will be rounded to the first day of the month, at 00:00 UTC",
      quickRanges: {
        [STATS_QUICK_RANGES.custom.value]: 'Custom',
        [STATS_QUICK_RANGES.last2Days.value]: 'Last 2 days',
        [STATS_QUICK_RANGES.last7Days.value]: 'Last 7 days',
        [STATS_QUICK_RANGES.last30Days.value]: 'Last 30 days',
        [STATS_QUICK_RANGES.last1Year.value]: 'Last 1 year',
        [STATS_QUICK_RANGES.yesterday.value]: 'Yesterday',
        [STATS_QUICK_RANGES.previousWeek.value]: 'Previous week',
        [STATS_QUICK_RANGES.previousMonth.value]: 'Previous month',
        [STATS_QUICK_RANGES.today.value]: 'Today',
        [STATS_QUICK_RANGES.todaySoFar.value]: 'Today so far',
        [STATS_QUICK_RANGES.thisWeek.value]: 'This week',
        [STATS_QUICK_RANGES.thisWeekSoFar.value]: 'This week so far',
        [STATS_QUICK_RANGES.thisMonth.value]: 'This month',
        [STATS_QUICK_RANGES.thisMonthSoFar.value]: 'This month so far',
        [STATS_QUICK_RANGES.last1Hour.value]: 'Last 1 hour',
        [STATS_QUICK_RANGES.last3Hour.value]: 'Last 3 hour',
        [STATS_QUICK_RANGES.last6Hour.value]: 'Last 6 hour',
        [STATS_QUICK_RANGES.last12Hour.value]: 'Last 12 hour',
        [STATS_QUICK_RANGES.last24Hour.value]: 'Last 24 hour',
      },
    },
    statsNumbers: {
      title: 'Stats numbers',
      yesNoMode: 'Yes/No mode',
      defaultStat: 'Default: Alarms created',
      sortOrder: 'Sort order',
      displayMode: 'Display Mode',
      selectAColor: 'Select a color',
    },
    infoPopup: {
      title: 'Info popup',
      fields: {
        column: 'Column',
        template: 'Template',
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
    weatherTemplate: 'Template - Weather item',
    modalTemplate: 'Template - Modal',
    entityTemplate: 'Template - Entities',
    columnSM: 'Columns - Small',
    columnMD: 'Columns - Medium',
    columnLG: 'Columns - Large',
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
        component: 'Component',
        connector: 'Connector',
        resource: 'Resource',
        watcher: 'Watcher',
      },
    },
    statSelector: {
      error: {
        alreadyExist: 'Stat with this name already exists',
      },
    },
    statsGroups: {
      title: 'Stats groups',
      manageGroups: 'Add a group',
      required: 'Create at least 1 group',
    },
    statsColor: {
      title: 'Stats color',
      pickColor: 'Pick a color',
    },
    considerPbehaviors: {
      title: 'Consider pbehaviors',
    },
    serviceWeatherModalTypes: {
      title: 'Type of the weather modal window',
      fields: {
        moreInfo: 'More info',
        alarmList: 'Alarm list',
      },
    },
    templateEditor: 'Template',
  },
  modals: {
    contextInfos: {
      title: 'Entities infos',
    },
    createEntity: {
      createTitle: 'Create an entity',
      editTitle: 'Edit an entity',
      duplicateTitle: 'Duplicate an entity',
      manageInfos: {
        infosList: 'Informations',
        addInfo: 'Add Information',
        noInfos: 'No information',
      },
      fields: {
        type: 'Type',
        manageInfos: 'Manage Infos',
        form: 'Form',
        impact: 'Impact',
        depends: 'Depends',
        types: {
          connector: 'connector',
          component: 'component',
          resource: 'resource',
        },
      },
      success: {
        create: 'Entity successfully created !',
        edit: 'Entity successfully edited !',
        duplicate: 'Entity successfully duplicated !',
      },
    },
    createWatcher: {
      createTitle: 'Create a watcher',
      editTitle: 'Edit a watcher',
      duplicateTitle: 'Duplicate a watcher',
      displayName: 'Name',
      success: {
        create: 'Watcher successfully created !',
        edit: 'Watcher successfully edited !',
        duplicate: 'Watcher successfully duplicated !',
      },
    },
    addEntityInfo: {
      addTitle: 'Add an information',
      editTitle: 'Edit an information',
    },
    view: {
      select: {
        title: 'Select a view',
      },
      create: {
        title: 'Create a view',
      },
      edit: {
        title: 'Edit the view',
      },
      duplicate: {
        title: 'Duplicate the view',
        infoMessage: 'You\'re duplicating a view. All duplicated view\'s rows/widgets will be copied on the new view.',
      },
      noData: 'No group corresponding. Press <kbd>enter</kbd> to create a new one',
      fields: {
        groupIds: 'Choose a group, or create a new one',
        groupTags: 'Group tags',
      },
      success: {
        create: 'New view created !',
        edit: 'View successfully edited !',
        delete: 'View successfully deleted !',
      },
      fail: {
        create: 'View creation failed...',
        edit: 'View edition failed...',
        delete: 'View deletion failed...',
      },
      errors: {
        rightCreating: 'Error on right creating',
        rightRemoving: 'Error on right removing',
      },
    },
    createAckEvent: {
      title: 'Ack',
      tooltips: {
        ackResources: 'Do you want to ack linked resources ?',
      },
      fields: {
        ticket: 'Ticket number',
        output: 'Note',
        ackResources: 'Ack resources',
      },
    },
    createSnoozeEvent: {
      title: 'Snooze',
      fields: {
        duration: 'Duration',
      },
    },
    createCancelEvent: {
      title: 'Cancel',
      fields: {
        output: 'Note',
      },
    },
    createChangeStateEvent: {
      title: 'Change state',
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
      title: 'Create periodical behavior',
      fields: {
        name: 'Name',
        start: 'Start',
        stop: 'End',
        reason: 'Reason',
        type: 'Type',
        message: 'Message',
        rRuleQuestion: 'Put a rrule on this pbehavior ?',
      },
      buttons: {
        addComment: 'Add comment',
      },
      success: {
        create: 'Pbehavior successfully created ! You may need to wait 60sec to see it in interface',
      },
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
    moreInfos: {
      moreInfos: 'More infos',
      defineATemplate: 'To define a template for this window, go to the alarms list settings',
    },
    infoPopupSetting: {
      title: 'Info popup',
      add: 'Add',
      column: 'Column',
      template: 'Template',
      addInfoPopup: {
        title: 'Add an info popup',
      },
    },
    variablesHelp: {
      variables: 'Variables',
      copyToClipboard: 'Copy to clipboard',
    },
    watcher: {
      criticity: 'Criticity',
      organization: 'Organization',
      numberOk: 'Number Ok',
      numberKo: 'Number Ko',
      state: 'State',
      name: 'Name',
      org: 'Org',
      noData: 'No data',
      ticketing: 'Ticketing',
      application_crit_label: 'Criticality',
      product_line: 'Product line',
      service_period: 'Plage surveillanc',
      isInCarat: 'Cartographic repository',
      application_label: 'Description',
      target_platform: 'Environment',
      scenario_label: 'Label',
      scenario_probe_name: 'Sonde',
      scenario_calendar: 'Range of execution',
      actionPending: 'action(s) pending',
    },
    filter: {
      create: {
        title: 'Create filter',
      },
      edit: {
        title: 'Edit filter',
      },
      fields: {
        title: 'Title',
      },
    },
    colorPicker: {
      title: 'Color picker',
    },
    textEditor: {
      title: 'Text editor',
    },
    widgetCreation: {
      title: 'Select a widget',
      types: {
        alarmList: {
          title: 'Alarm list',
        },
        context: {
          title: 'Context explorer',
        },
        weather: {
          title: 'Service weather',
        },
        statsHistogram: {
          title: 'Stats histogram',
        },
        statsCurves: {
          title: 'Stats curves',
        },
        statsTable: {
          title: 'Stats table',
        },
        statsCalendar: {
          title: 'Stats calendar',
        },
        statsNumber: {
          title: 'Stats number',
        },
        text: {
          title: 'Text',
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
      title: 'Alarms list',
    },
    createUser: {
      title: 'Create user',
      fields: {
        username: 'Username',
        firstName: 'First name',
        lastName: 'Last name',
        email: 'Email',
        password: 'Password',
        language: 'User interface language',
        role: 'Role',
        enabled: 'Enabled',
      },
    },
    editUser: {
      title: 'Edit user',
    },
    createRole: {
      title: 'Create role',
    },
    editRole: {
      title: 'Edit role',
    },
    createRight: {
      title: 'Create right',
      fields: {
        id: 'ID',
        description: 'Description',
        type: 'Type',
      },
    },
    eventFilterRule: {
      create: {
        title: 'Create event filter rule',
        success: 'Rule successfully created !',
      },
      duplicate: {
        title: 'Duplicate event filter rule',
        success: 'Rule successfully created !',
      },
      edit: {
        title: 'Edit an event filter rule',
        success: 'Rule successfully edited !',
      },
      remove: {
        success: 'Rule successfully removed !',
      },
      priority: 'Priority',
      editPattern: 'Edit pattern',
      advanced: 'Advanced',
      addAField: 'Add a field',
      simpleEditor: 'Simple editor',
      field: 'Field',
      value: 'Value',
      advancedEditor: 'Advanced editor',
      comparisonRules: 'Comparison rules',
      enrichmentOptions: 'Enrichment options',
      editActions: 'Edit actions',
      addAction: 'Add an action',
      actions: 'Actions',
      from: 'From',
      to: 'To',
      externalData: 'External data',
      onSuccess: 'On success',
      onFailure: 'On failure',
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
    createWebhook: {
      create: {
        title: 'Create webhook',
        success: 'Webhook successfully created !',
      },
      edit: {
        title: 'Edit webhook',
        success: 'Webhook successfully edited !',
      },
      remove: {
        success: 'Webhook successfully removed !',
      },
      fields: {
        id: 'ID',
      },
      tooltips: {
        id: 'This field is optional, if no ID is entered, an ID will be auto-generated.',
      },
    },
    statsDateInterval: {
      title: 'Stats - Date interval',
      fields: {
        periodValue: 'Period value',
        periodUnit: 'Period unit',
      },
      errors: {
        endDateLessOrEqualStartDate: 'End date should be after start date',
      },
    },
  },
  tables: {
    noData: 'No data',
    contextList: {
      title: 'Context List',
      name: 'Name',
      type: 'Type',
      id: 'Id',
    },
    alarmGeneral: {
      title: 'General',
      author: 'Author',
      connector: 'Connector',
      connectorName: 'Connector name',
      component: 'Component',
      resource: 'Resource',
      output: 'Output',
      lastUpdateDate: 'Last update date',
      creationDate: 'Creation date',
      duration: 'Duration',
      state: 'State',
      status: 'Status',
      extraDetails: 'Extra details',
    },
    /**
     * This object for pbehavior fields from database
     */
    pbehaviorList: {
      name: 'Name',
      author: 'Author',
      connector: 'Connector',
      connectorName: 'Connector name',
      enabled: 'Is enabled',
      tstart: 'Begins',
      tstop: 'Ends',
      type_: 'Type',
      reason: 'Reason',
      rrule: 'Rrule',
    },
    rolesList: {
      name: 'Name',
      actions: 'Actions',
    },
    alarmStatus: {
      [ENTITIES_STATUSES.off]: 'Off',
      [ENTITIES_STATUSES.ongoing]: 'Ongoing',
      [ENTITIES_STATUSES.flapping]: 'Flapping',
      [ENTITIES_STATUSES.stealthy]: 'Stealthy',
      [ENTITIES_STATUSES.cancelled]: 'Canceled',
    },
    alarmStates: {
      [ENTITIES_STATES.ok]: 'Info',
      [ENTITIES_STATES.minor]: 'Minor',
      [ENTITIES_STATES.major]: 'Major',
      [ENTITIES_STATES.critical]: 'Critical',
    },
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
    admin: {
      users: {
        columns: {
          username: 'Username',
          role: 'Role',
          enabled: 'Enabled',
        },
      },
    },
  },
  rRule: {
    advancedHint: 'Separate numbers with a comma',
    textLabel: 'Rrule',
    stringLabel: 'Summary',
    tabs: {
      simple: 'Simple',
      advanced: 'Advanced',
    },
    errors: {
      main: 'Please note that the Rrule you choose is not valid. We strongly advise you to modify it before saving changes.',
    },
    fields: {
      freq: 'Frequency',
      until: 'Until',
      byweekday: 'By week day',
      count: 'Repeat',
      interval: 'Interval',
      wkst: 'Week start',
      bymonth: 'By month',
      bysetpos: {
        label: 'By set position',
        tooltip: 'If given, it must be one or many integers, positive or negative. Each given integer will specify an occurrence number, corresponding to the nth occurrence of the rule inside the frequency period. For example, a \'bysetpos\' of -1 if combined with a monthly frequency, and a \'byweekday\' of (Monday, Tuesday, Wednesday, Thursday, Friday), will result in the last work day of every month.',
      },
      bymonthday: {
        label: 'By month day',
        tooltip: 'If given, it must be one or many integers, meaning the month days to apply the recurrence to.',
      },
      byyearday: {
        label: 'By year day',
        tooltip: 'If given, it must be one or many integers, meaning the year days to apply the recurrence to.',
      },
      byweekno: {
        label: 'By week n°',
        tooltip: 'If given, it must be on or many integers, meaning the week numbers to apply the recurrence to. Week numbers have the meaning described in ISO8601, that is, the first week of the year is that containing at least four days of the new year.',
      },
      byhour: {
        label: 'By hour',
        tooltip: 'If given, it must be one or many integers, meaning the hours to apply the recurrence to.',
      },
      byminute: {
        label: 'By minute',
        tooltip: 'If given, it must be one or many integers, meaning the minutes to apply the recurrence to.',
      },
      bysecond: {
        label: 'By second',
        tooltip: 'If given, it must be one or many integers, meaning the seconds to apply the recurrence to.',
      },
    },
  },
  errors: {
    default: 'Something went wrong...',
    lineNotEmpty: 'This line is not empty',
    JSONNotValid: 'Invalid JSON',
    versionNotFound: 'Unable to get application version',
  },
  calendar: {
    today: 'Today',
    month: 'Month',
    week: 'Week',
    day: 'Day',
  },
  success: {
    default: 'Done !',
    createEntity: 'Entity successfully created',
    editEntity: 'Entity successfully edited',
    pathCopied: 'Path copied to clipboard',
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
    resultsTableHeaders: {
      alarm: {
        connector: 'Connector',
        connectorName: 'Connector name',
        component: 'Component',
        resource: 'Resource',
      },
      entity: {
        id: 'ID',
        name: 'Name',
        type: 'Type',
      },
    },
    errors: {
      cantParseToVisualEditor: 'We can\'t parse this filter to Visual Editor',
      invalidJSON: 'Invalid JSON',
      required: 'You need to add at least one valid rule',
    },
  },
  filterSelector: {
    fields: {
      mixFilters: 'Mix filters',
    },
  },
  validator: {
    unique: 'Field must be unique',
  },
  stats: {
    types: {
      [STATS_TYPES.alarmsCreated.value]: 'Alarms created',
      [STATS_TYPES.alarmsResolved.value]: 'Alarms resolved',
      [STATS_TYPES.alarmsCanceled.value]: 'Alarms canceled',
      [STATS_TYPES.ackTimeSla.value]: 'Ack time Sla',
      [STATS_TYPES.resolveTimeSla.value]: 'Resolve time Sla',
      [STATS_TYPES.timeInState.value]: 'Time in state',
      [STATS_TYPES.stateRate.value]: 'State rate',
      [STATS_TYPES.mtbf.value]: 'MTBF',
      [STATS_TYPES.currentState.value]: 'Current state',
      [STATS_TYPES.ongoingAlarms.value]: 'Ongoing alarms',
      [STATS_TYPES.currentOngoingAlarms.value]: 'Current ongoing alarms',
    },
  },
  eventFilter: {
    title: 'Event filter',
    type: 'Type',
    pattern: 'Pattern',
    priority: 'Priority',
    enabled: 'Enabled',
    actions: 'Actions',
    externalDatas: 'External datas',
    actionsRequired: 'Please add at least one action',
    id: 'Id',
    idHelp: 'If no id is specified, an unique id will be generated automatically on rule creation',
  },
  layout: {
    sideBar: {
      buttons: {
        edit: 'Toggle editing mode',
        create: 'Create view',
        settings: 'Settings',
      },
    },
  },
  parameters: {
    interfaceLanguage: 'Interface language',
    groupsNavigationType: {
      title: 'Groups navigation type',
      items: {
        sideBar: 'Side bar',
        topBar: 'Top bar',
      },
    },
  },
  view: {
    errors: {
      emptyTabs: 'You should create a tab',
    },
    deleteRow: 'Delete row',
    deleteWidget: 'Delete widget',
    fullScreen: 'Full screen',
    fullScreenShortcut: 'Alt + Enter / Command + Enter',
  },
  patternsList: {
    noData: 'No pattern set. Click \'Add\' button to start adding fields to the pattern',
    noDataDisabled: 'No pattern set.',
  },
  webhook: {
    title: 'Webhooks',
    table: {
      headers: {
        id: 'ID',
        requestMethod: 'Request method',
        requestUrl: 'Request URL',
      },
    },
    tabs: {
      hook: {
        title: 'Hook',
        fields: {
          triggers: 'Triggers',
          eventPatterns: 'Event patterns',
          alarmPatterns: 'Alarm patterns',
          entityPatterns: 'Entity patterns',
        },
      },
      request: {
        title: 'Request',
        fields: {
          method: 'Method',
          url: 'URL',
          authSwitch: 'Do you need auth fields?',
          auth: 'Auth',
          username: 'Username',
          password: 'Password',
          headers: 'Headers',
          headerKey: 'Header key',
          headerValue: 'Header value',
          payload: 'Payload',
        },
      },
      declareTicket: {
        title: 'Declare ticket',
        fields: {
          text: 'Key',
          value: 'Value',
        },
      },
    },
  },
  validation: {
    custom: {
      tstop: {
        after: 'The {0} should be after than {1}',
      },
    },
  },
};
