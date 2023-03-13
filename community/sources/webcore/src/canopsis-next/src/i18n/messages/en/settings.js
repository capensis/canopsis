import {
  STATS_CRITICITY,
  ENTITY_TYPES,
  SIDE_BARS,
  ALARMS_OPENED_VALUES,
} from '@/constants';

export default {
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
  filters: 'Filters',
  filterEditor: 'Filter',
  isAckNoteRequired: 'Note field required when ack?',
  isSnoozeNoteRequired: 'Note field required when snooze?',
  inlineLinksCount: 'Inline links count',
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
  resolvedAlarmsColumns: 'Column names for resolved alarms',
  activeAlarmsColumns: 'Column names for active alarms',
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
  columns: {
    customLabel: 'Custom label',
    isHtml: 'Is it HTML?',
    withTemplate: 'Custom template',
    isState: 'Displayed as severity?',
    onlyIcon: 'Show only links icons',
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
    ultraCompact: 'Ultra compact view',
  },
};
