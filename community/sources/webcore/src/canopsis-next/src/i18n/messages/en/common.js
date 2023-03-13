import { THEMES_NAMES } from '@/config';
import {
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  EVENT_ENTITY_TYPES,
  PATTERN_FIELD_TYPES,
  PATTERN_OPERATORS,
  TRIGGERS,
} from '@/constants';

export default {
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
  created: 'Creation date',
  updated: 'Last update date',
  expired: 'Expired date',
  accessed: 'Accessed at',
  lastEventDate: 'Last event date',
  lastPbehaviorDate: 'Last pbehavior date',
  activated: 'Activated',
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
  ack: 'Ack',
  acked: 'Acked',
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
  alarmId: 'Alarm ID',
  longOutput: 'Long output',
  timestamp: 'Timestamp',
  countOfMax: '{count} of {total}',
  trigger: 'Trigger | Triggers',
  column: 'Column | Columns',
  countOfTotal: '{count} of {total}',
  rule: 'Rule',
  initialLongOutput: 'Long initial output',
  totalStateChanges: 'Total state changes',
  noReceivedEvents: 'No events received for {duration} by some of dependencies',
  frequencyLimit: 'Frequency limit',
  clearSearch: 'Clear search input',
  noData: 'No data',
  noColumns: 'You have to select at least 1 column',
  theme: 'Theme | Themes',
  variableTypes: {
    string: 'String',
    number: 'Number',
    boolean: 'Boolean',
    null: 'Null',
    array: 'Array',
  },
  mixedField: {
    types: {
      [PATTERN_FIELD_TYPES.string]: '@:common.variableTypes.string',
      [PATTERN_FIELD_TYPES.number]: '@:common.variableTypes.number',
      [PATTERN_FIELD_TYPES.boolean]: '@:common.variableTypes.boolean',
      [PATTERN_FIELD_TYPES.null]: '@:common.variableTypes.null',
      [PATTERN_FIELD_TYPES.stringArray]: '@:common.variableTypes.array',
    },
  },
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

    [PATTERN_OPERATORS.activated]: 'Activated',
    [PATTERN_OPERATORS.inactive]: 'Inactive',

    [PATTERN_OPERATORS.regexp]: 'Regexp',
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

  triggers: {
    [TRIGGERS.create]: {
      text: 'Alarm creation',
    },
    [TRIGGERS.statedec]: {
      text: 'Alarm state decrease',
    },
    [TRIGGERS.changestate]: {
      text: 'Alarm state has been changed by "change state" action',
    },
    [TRIGGERS.stateinc]: {
      text: 'Alarm state increase',
    },
    [TRIGGERS.changestatus]: {
      text: 'Alarm status changes eg. flapping',
    },
    [TRIGGERS.ack]: {
      text: 'Alarm has been acked',
    },
    [TRIGGERS.ackremove]: {
      text: 'Alarm has been unacked',
    },
    [TRIGGERS.cancel]: {
      text: 'Alarm has been cancelled',
    },
    [TRIGGERS.uncancel]: {
      text: 'Alarm has been uncancelled',
      helpText: 'Probably legacy trigger, because there is no way to uncancel alarm when you cancel it in the UI, but it\'s possible to send an uncancel event via API',
    },
    [TRIGGERS.comment]: {
      text: 'Alarm has been commented',
    },
    [TRIGGERS.declareticket]: {
      text: 'Ticket has been declared by the UI action',
    },
    [TRIGGERS.declareticketwebhook]: {
      text: 'Ticket has been declared by the webhook',
    },
    [TRIGGERS.assocticket]: {
      text: 'Ticket has been associated with an alarm',
    },
    [TRIGGERS.snooze]: {
      text: 'Alarm has been snoozed',
    },
    [TRIGGERS.unsnooze]: {
      text: 'Alarm has been unsnoozed',
    },
    [TRIGGERS.resolve]: {
      text: 'Alarm has been resolved',
    },
    [TRIGGERS.activate]: {
      text: 'Alarm has been activated',
    },
    [TRIGGERS.pbhenter]: {
      text: 'Alarm enters a periodic behavior',
    },
    [TRIGGERS.pbhleave]: {
      text: 'Alarm leaves a periodic behavior',
    },
    [TRIGGERS.instructionfail]: {
      text: 'Manual instruction has failed',
    },
    [TRIGGERS.autoinstructionfail]: {
      text: 'Auto instruction has failed',
    },
    [TRIGGERS.instructionjobfail]: {
      text: 'Manual or auto instruction\'s job is failed',
    },
    [TRIGGERS.instructionjobcomplete]: {
      text: 'Manual or auto instruction\'s job is completed',
    },
    [TRIGGERS.instructioncomplete]: {
      text: 'Manual instruction is completed',
    },
    [TRIGGERS.autoinstructioncomplete]: {
      text: 'Auto instruction is completed',
    },
  },
  themes: {
    [THEMES_NAMES.canopsis]: 'Canopsis',
    [THEMES_NAMES.canopsisDark]: 'Canopsis dark',
    [THEMES_NAMES.colorBlind]: 'Color blind',
    [THEMES_NAMES.colorBlindDark]: 'Color blind dark',
  },
};
