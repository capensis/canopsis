import { PBEHAVIOR_TYPE_TYPES, WIDGET_TYPES } from '@/constants';

export default {
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
      privateTitle: 'Create a private view',
    },
    edit: {
      title: 'Edit the view',
    },
    duplicate: {
      title: 'Duplicate the view - {viewTitle}',
      infoMessage: 'You\'re duplicating a view. All duplicated view\'s rows/widgets will be copied on the new view.',
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
      },
      color: {
        label: 'Use special color for this event?',
      },
    },
    success: {
      create: 'Pbehavior successfully created! You may need to wait 60 sec to see it in interface',
    },
    cancelConfirmation: 'Some data has been modified and will not be saved. Do you really want to close this menu?',
  },
  createPause: {
    title: 'Create Pause event',
  },
  createAckRemove: {
    title: 'Remove ack',
  },
  createUnCancel: {
    title: 'Create uncancel event',
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
  },
  service: {
    refreshEntities: 'Refresh entities list',
    editPbehaviors: 'Edit pbehaviors',
    massActionsDescription: 'You can choose entities to perform actions',
    actionInQueue: 'action in queue|actions in queue',
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
      [WIDGET_TYPES.barChart]: {
        title: 'Bar chart',
      },
      [WIDGET_TYPES.lineChart]: {
        title: 'Line chart',
      },
      [WIDGET_TYPES.pieChart]: {
        title: 'Pie chart',
      },
      [WIDGET_TYPES.numbers]: {
        title: 'Numbers',
      },
      [WIDGET_TYPES.userStatistics]: {
        title: 'User statistics',
      },
      [WIDGET_TYPES.alarmStatistics]: {
        title: 'Alarm statistics',
      },
      [WIDGET_TYPES.availability]: {
        title: 'Availability',
      },
      chart: {
        title: 'Chart',
      },
      report: {
        title: 'Report',
      },
    },
  },
  manageHistogramGroups: {
    title: {
      add: 'Add a group',
      edit: 'Edit a group',
    },
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
    duplicate: {
      title: 'Duplicate role',
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
    duplicate: {
      title: 'Duplicate SNMP rule',
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
    canonicalTypes: {
      [PBEHAVIOR_TYPE_TYPES.active]: 'Active',
      [PBEHAVIOR_TYPE_TYPES.inactive]: 'Inactive',
      [PBEHAVIOR_TYPE_TYPES.maintenance]: 'Maintenance',
      [PBEHAVIOR_TYPE_TYPES.pause]: 'Pause',
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
  linkToMetaAlarm: {
    title: 'Link to a meta alarm',
    noData: 'No meta alarm corresponding. Press <kbd>enter</kbd> to create a new one',
    fields: {
      metaAlarm: 'Select meta alarm or create a new one',
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
  createStateSetting: {
    create: {
      title: 'Create state compute method',
      success: 'State compute method created!',
    },
    edit: {
      title: 'Edit state compute method',
      success: 'State compute method edited!',
    },
    duplicate: {
      title: 'Duplicate state compute method',
      success: 'State compute method duplicated!',
    },
    remove: {
      success: 'State compute method deleted!',
    },
  },
  createJunitStateSetting: {
    edit: {
      title: 'JUnit test suite state settings',
      success: 'JUnit test suite state setting edited!',
    },
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
    dismissed: 'has dismissed your updates',
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
      title: 'Archive/delete entities. Are you sure ?',
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
  createServiceWeatherPattern: {
    create: {
      title: 'Create service weather filter',
    },
    edit: {
      title: 'Edit service weather filter',
    },
  },
  createCorporateServiceWeatherPattern: {
    create: {
      title: 'Create shared service weather filter',
    },
    edit: {
      title: 'Edit shared service weather filter',
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
  createWidgetTemplate: {
    create: {
      title: 'Create widget template',
    },
    edit: {
      title: 'Edit widget template',
    },
  },
  selectWidgetTemplateType: {
    title: 'Select widget template type',
  },
  entityDependenciesList: {
    title: '{name} impacted entities',
  },
  createDeclareTicketRule: {
    create: {
      title: 'Create a declare ticket rule',
    },
    edit: {
      title: 'Edit a declare ticket rule',
    },
    duplicate: {
      title: 'Duplicate a declare ticket rule',
    },
  },
  createDeclareTicketEvent: {
    title: 'Declare ticket',
  },
  executeDeclareTickets: {
    title: 'Ticket declaration status',
  },
  createAssociateTicketEvent: {
    title: 'Associate ticket number',
  },
  createAckEvent: {
    title: 'Ack',
  },
  createLinkRule: {
    create: {
      title: 'Create link generator',
    },
    edit: {
      title: 'Edit link generator',
    },
    duplicate: {
      title: 'Duplicate a link generator',
    },
  },
  createAlarmChart: {
    [WIDGET_TYPES.barChart]: {
      create: {
        title: 'Create bar chart',
      },
      edit: {
        title: 'Edit bar chart',
      },
    },
    [WIDGET_TYPES.lineChart]: {
      create: {
        title: 'Create line chart',
      },
      edit: {
        title: 'Edit line chart',
      },
    },
    [WIDGET_TYPES.numbers]: {
      create: {
        title: 'Create numbers chart',
      },
      edit: {
        title: 'Edit numbers chart',
      },
    },
  },
  importPbehaviorException: {
    title: 'Import exception dates',
  },
  createMaintenance: {
    enableMaintenance: 'Enable maintenance mode',
    setup: {
      title: 'Maintenance mode setup',
    },
    edit: {
      title: 'Edit maintenance mode',
    },
  },
  confirmationLeaveMaintenance: {
    title: 'Leave maintenance mode',
    text: 'Are you sure you want to leave the maintenance mode?\nAll users will be able to login to the system after leaving.',
  },
  confirmationCreateNewTicketForAlarm: {
    title: 'Confirm create tickets',
    text: 'This alarm already has tickets created.\nDo you want to create a new one?',
  },
  confirmationCreateNewTicketForAlarms: {
    title: 'Confirm create tickets',
    text: 'Some alarms already have tickets created.\nDo you want to create new tickets for them?',
  },
  createTag: {
    create: {
      title: 'Create a tag',
    },
    edit: {
      title: 'Edit a tag',
    },
    duplicate: {
      title: 'Duplicate a tag',
    },
  },
  createTheme: {
    create: {
      title: 'Create theme',
    },
    edit: {
      title: 'Edit theme',
    },
    duplicate: {
      title: 'Duplicate theme',
    },
  },
  archiveDisabledEntities: {
    text: 'Are you sure you want to archive disabled entities?\nThis action cannot be undone',
  },
  createIcon: {
    create: {
      title: 'Upload icon',
      success: 'Icon was uploaded',
    },
    remove: {
      success: 'Icon was removed',
    },
  },
  launchEventsRecording: {
    title: 'Launch events recording',
  },
  eventsRecording: {
    title: 'Events recording {date}',
    subtitle: '{count} events from RabbitMQ received',
    buttonTooltip: 'Delete received events',
  },
};
