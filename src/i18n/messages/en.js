import { ENTITIES_STATES, ENTITIES_STATUSES } from '@/constants';

export default {
  common: {
    watcher: 'Watcher',
    name: 'Name',
    description: 'Description',
    submit: 'Submit',
    enabled: 'Enabled',
    login: 'Login',
    yes: 'Yes',
    no: 'No',
    confirmation: 'Are you sure ?',
    by: 'By',
    date: 'Date',
    comment: 'Comment',
    end: 'End',
    username: 'Username',
    password: 'Password',
    title: 'Title',
    save: 'Save',
    label: 'Label',
    value: 'Value',
    add: 'Add',
    delete: 'Delete',
    edit: 'Edit',
    parse: 'Parse',
    home: 'Home',
    entries: 'entries',
    showing: 'showing',
    apply: 'Apply',
    to: 'to',
    of: 'of',
    actionsLabel: 'Actions',
    noResults: 'No results',
    actions: {
      close: 'Close',
      acknowledge: 'Acknowledge',
      acknowledgeAndReport: 'Acknowledge and report an incident',
      saveChanges: 'Save changes',
      reportIncident: 'Report an incident',
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
  },
  settings: {
    titles: {
      alarmListSettings: 'Alarm list settings',
      contextTableSettings: 'Context table settings',
    },
    widgetTitle: 'Widget title',
    columnName: 'Column name',
    defaultSortColumn: 'Default Sort Column',
    columnNames: 'Column names',
    periodicRefresh: 'Periodic refresh',
    defaultNumberOfElementsPerPage: 'Default number of elements/page',
    elementsPerPage: 'Elements per page',
    filterOnOpenResolved: 'Filter on Open/Resolved',
    open: 'Open',
    resolved: 'Resolved',
    filters: 'Filters',
    selectAFilter: 'Select a filter',
    infoPopup: 'Info popup',
    moreInfosModal: '"More Infos" Popup',
    contextTypeOfEntities: {
      title: 'Type of entities',
      fields: {
        component: 'Component',
        connector: 'Connector',
        resource: 'Resource',
        watcher: 'Watcher',
      },
    },
  },
  modals: {
    contextInfos: {
      title: 'Entities infos',
    },
    createEntity: {
      createTitle: 'Create an entity',
      editTitle: 'Edit an entity',
      fields: {
        type: 'Types',
        manageInfos: 'Manage Infos',
        form: 'Form',
        types: {
          connector: 'connector',
          component: 'component',
          resource: 'resource',
        },
      },
    },
    createAckEvent: {
      title: 'Add event: Ack',
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
      title: 'Add event: Snooze',
      fields: {
        duration: 'Duration',
      },
    },
    createCancelEvent: {
      title: 'Add event: Cancel',
      fields: {
        output: 'Note',
      },
    },
    createChangeStateEvent: {
      title: 'Add event: Change state',
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
      title: 'Put a pbehavior on these elements ?',
      fields: {
        name: 'Name',
        start: 'Start',
        stop: 'End',
        reason: 'Reason',
        type: 'Type',
        rRuleQuestion: 'Put a rrule on this pbehavior ?',
      },
    },
    createAckRemove: {
      title: 'Add event: Remove ack',
    },
    createDeclareTicket: {
      title: 'Add event: Declare ticket',
    },
    createAssociateTicket: {
      title: 'Add event: Associate ticket number',
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
    },
    moreInfos: {
      moreInfos: 'More infos',
      defineATemplate: 'To define a template for this window, go to the alarms list settings',
    },
  },
  tables: {
    contextList: {
      title: 'Context List',
      name: 'Name',
      id: 'Id',
      noDataText: 'Make a research',
    },
    alarmGeneral: {
      title: 'General',
      author: 'Author',
      connector: 'Connector',
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
      connector_name: 'Connector name',
      enabled: 'Is enabled',
      tstart: 'Begins',
      tstop: 'Ends',
      type_: 'Type',
      reason: 'Reason',
      rrule: 'Rrule',
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
  },
  success: {
    default: 'Done !',
  },
  mFilterEditor: {
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
      connector: 'Connector',
      connectorName: 'Connector name',
      component: 'Component',
      resource: 'Resource',
    },
  },
};
