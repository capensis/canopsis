import { COLORS, MEDIA_QUERIES_BREAKPOINTS } from '@/config';

import featuresService from '@/services/features';

export const CRUD_ACTIONS = {
  create: 'create',
  read: 'read',
  update: 'update',
  delete: 'delete',
};

export const ENTITIES_TYPES = {
  alarm: 'alarm',
  entity: 'entity',
  watcher: 'watcher',
  watcherEntity: 'watcherEntity',
  weather: 'weather',
  pbehavior: 'pbehavior',
  event: 'event',
  userPreference: 'userPreference',
  group: 'group',
  view: 'view',
  viewTab: 'viewTab',
  widget: 'widget',
  stat: 'stat',
  user: 'user',
  role: 'role',
  eventFilterRule: 'eventFilterRule',
  metaAlarmRule: 'metaAlarmRule',
  webhook: 'webhook',
  snmpRule: 'snmpRule',
  action: 'action',
  heartbeat: 'heartbeat',
  dynamicInfo: 'dynamicInfo',
  broadcastMessage: 'broadcastMessage',
  playlist: 'playlist',
  pbehaviorExceptions: 'pbehaviorExceptions',
  pbehaviorTypes: 'pbehaviorTypes',
  pbehaviorReasons: 'pbehaviorReasons',
  remediationInstruction: 'remediationInstruction',
  remediationJob: 'remediationJob',
  remediationConfiguration: 'remediationConfiguration',
  remediationInstructionExecution: 'remediationInstructionExecution',
};

export const MODALS = {
  createEvent: 'create-event',
  createAckEvent: 'create-ack-event',
  confirmAckWithTicket: 'confirm-ack-with-ticket',
  createAssociateTicketEvent: 'create-associate-ticket-event',
  createCommentEvent: 'create-comment-event',
  createChangeStateEvent: 'create-change-state-event',
  createDeclareTicketEvent: 'create-declare-ticket-event',
  createSnoozeEvent: 'create-snooze-event',
  variablesHelp: 'variables-help',
  createEntity: 'create-entity',
  createWatcher: 'create-watcher',
  addEntityInfo: 'add-entity-info',
  watcher: 'watcher',
  createWatcherPauseEvent: 'create-watcher-pause-event',
  pbehaviorList: 'pbehavior-list',
  editLiveReporting: 'edit-live-reporting',
  infoPopupSetting: 'info-popup-setting',
  addInfoPopup: 'add-info-popup',
  confirmation: 'confirmation',
  clickOutsideConfirmation: 'click-outside-confirmation',
  createWidget: 'create-widget',
  createFilter: 'create-filter',
  alarmsList: 'alarms-list',
  addStat: 'add-stat',
  statsDateInterval: 'stats-date-interval',
  statsDisplayMode: 'stats-display-mode',
  colorPicker: 'color-picker',
  textEditor: 'text-editor',
  textFieldEditor: 'text-field-editor',
  selectView: 'select-view',
  createView: 'create-view',
  createGroup: 'create-group',
  createUser: 'create-user',
  createRole: 'create-role',
  createRight: 'create-right',
  createBroadcastMessage: 'create-broadcast-message',
  createEventFilterRule: 'create-event-filter-rule',
  createMetaAlarmRule: 'create-meta-alarm-rule',
  createPattern: 'create-pattern',
  createPatternRule: 'create-pattern-rule',
  addEventFilterRuleToPattern: 'add-event-filter-rule-to-pattern',
  eventFilterRuleActions: 'event-filter-rule-actions',
  eventFilterRuleExternalData: 'event-filter-rule-external-data',
  eventFilterRuleCreateAction: 'event-filter-rule-create-action',
  filtersList: 'filters-list',
  createWebhook: 'create-webhook',
  createSnmpRule: 'create-snmp-rule',
  selectViewTab: 'select-view-tab',
  createAction: 'create-action',
  createHeartbeat: 'create-heartbeat',
  createDynamicInfo: 'create-dynamic-info',
  createDynamicInfoInformation: 'create-dynamic-info-information',
  importExportViews: 'import-groups-and-views',
  dynamicInfoTemplatesList: 'dynamic-info-templates-list',
  createDynamicInfoTemplate: 'create-dynamic-info-template',
  createPlaylist: 'create-playlist',
  managePlaylistTabs: 'manage-playlist-tabs',
  pbehaviorPlanning: 'pbehavior-planning',
  createRRule: 'create-r-rule',
  selectExceptionsLists: 'select-exceptions-lists',
  pbehaviorRecurrentChangesConfirmation: 'pbehavior-recurrent-changes-confirmation',
  createPbehavior: 'create-pbehavior',
  createPbehaviorType: 'create-pbehavior-type',
  createPbehaviorReason: 'create-pbehavior-reason',
  createPbehaviorException: 'create-pbehavior-exception',
  createManualMetaAlarm: 'create-manual-meta-alarm',
  createRemediationInstruction: 'create-remediation-instruction',
  createRemediationConfiguration: 'create-remediation-configuration',
  createRemediationJob: 'create-remediation-job',
  createRemediationInstructionsFilter: 'create-remediation-instructions-filter',
  executeRemediationInstruction: 'execute-remediation-instruction',
  imageViewer: 'image-viewer',
  patterns: 'patterns',
  rate: 'rate',
};

export const EVENT_ENTITY_TYPES = {
  ack: 'ack',
  fastAck: 'fastAck',
  ackRemove: 'ackremove',
  pbehaviorAdd: 'pbehaviorAdd',
  pbehaviorList: 'pbehaviorList',
  assocTicket: 'assocticket',
  cancel: 'cancel',
  delete: 'delete',
  changeState: 'changestate',
  declareTicket: 'declareticket',
  snooze: 'snooze',
  done: 'done',
  validate: 'validate',
  invalidate: 'invalidate',
  pause: 'pause',
  play: 'play',
  groupRequest: 'groupRequest',
  group: 'group',
  pbhenter: 'pbhenter',
  pbhleave: 'pbhleave',
  comment: 'comment',
  manualMetaAlarmGroup: 'manual_metaalarm_group',
  manualMetaAlarmUngroup: 'manual_metaalarm_ungroup',
  manualMetaAlarmUpdate: 'manual_metaalarm_update',
  stateinc: 'stateinc',
  statedec: 'statedec',
  statusinc: 'statusinc',
  statusdec: 'statusdec',
  unsooze: 'unsooze',
  metaalarmattach: 'metaalarmattach',
  executeInstruction: 'executeInstruction',
  instructionStart: 'instructionstart',
  instructionPause: 'instructionpause',
  instructionResume: 'instructionresume',
  instructionComplete: 'instructioncomplete',
  instructionAbort: 'instructionabort',
  instructionFail: 'instructionfail',
  instructionJobStart: 'instructionjobstart',
  instructionJobComplete: 'instructionjobcomplete',
  instructionJobAbort: 'instructionjobabort',
  instructionJobFail: 'instructionjobfail',
};

export const ENTITY_INFOS_TYPE = {
  state: 'state',
  status: 'status',
  action: 'action',
};

export const ENTITIES_STATES_KEYS = {
  ok: 'ok',
  minor: 'minor',
  major: 'major',
  critical: 'critical',
};

export const ENTITIES_STATES = {
  ok: 0,
  minor: 1,
  major: 2,
  critical: 3,
};

export const ENTITIES_STATUSES = {
  off: 0,
  ongoing: 1,
  stealthy: 2,
  flapping: 3,
  cancelled: 4,
};

export const ENTITIES_STATES_STYLES = {
  [ENTITIES_STATES.ok]: {
    color: COLORS.state.ok,
    text: 'ok',
    icon: 'assistant_photo',
  },
  [ENTITIES_STATES.minor]: {
    color: COLORS.state.minor,
    text: 'minor',
    icon: 'assistant_photo',
  },
  [ENTITIES_STATES.major]: {
    color: COLORS.state.major,
    text: 'major',
    icon: 'assistant_photo',
  },
  [ENTITIES_STATES.critical]: {
    color: COLORS.state.critical,
    text: 'critical',
    icon: 'assistant_photo',
  },
};

export const WATCHER_STATES = {
  ok: 'ok',
  minor: 'minor',
  major: 'major',
  critical: 'critical',
  pause: 'pause',
};

export const WATCHER_STATES_COLORS = {
  [WATCHER_STATES.ok]: ENTITIES_STATES_STYLES[ENTITIES_STATES.ok].color,
  [WATCHER_STATES.minor]: ENTITIES_STATES_STYLES[ENTITIES_STATES.minor].color,
  [WATCHER_STATES.major]: ENTITIES_STATES_STYLES[ENTITIES_STATES.major].color,
  [WATCHER_STATES.critical]: ENTITIES_STATES_STYLES[ENTITIES_STATES.critical].color,
  [WATCHER_STATES.pause]: COLORS.state.pause,
};

export const COUNTER_STATES_ICONS = {
  [ENTITIES_STATES_KEYS.ok]: 'wb_sunny',
  [ENTITIES_STATES_KEYS.minor]: 'person',
  [ENTITIES_STATES_KEYS.major]: 'person',
  [ENTITIES_STATES_KEYS.critical]: 'wb_cloudy',
};

export const PBEHAVIOR_TYPE_TYPES = {
  active: 'active',
  inactive: 'inactive',
  maintenance: 'maintenance',
  pause: 'pause',
};

export const WEATHER_ICONS = {
  [WATCHER_STATES.ok]: 'wb_sunny',
  [WATCHER_STATES.minor]: 'person',
  [WATCHER_STATES.major]: 'person',
  [WATCHER_STATES.critical]: 'wb_cloudy',
  [PBEHAVIOR_TYPE_TYPES.maintenance]: 'build',
  [PBEHAVIOR_TYPE_TYPES.inactive]: 'brightness_3',
  [PBEHAVIOR_TYPE_TYPES.pause]: 'pause',
};

export const ENTITY_STATUS_STYLES = {
  [ENTITIES_STATUSES.off]: {
    color: COLORS.status.off,
    text: 'off',
    icon: 'keyboard_arrow_up',
  },
  [ENTITIES_STATUSES.ongoing]: {
    color: COLORS.status.ongoing,
    text: 'ongoing',
    icon: 'keyboard_arrow_up',
  },
  [ENTITIES_STATUSES.stealthy]: {
    color: COLORS.status.stealthy,
    text: 'stealthy',
    icon: 'keyboard_arrow_up',
  },
  [ENTITIES_STATUSES.flapping]: {
    color: COLORS.status.flapping,
    text: 'flapping',
    icon: 'keyboard_arrow_up',
  },
  [ENTITIES_STATUSES.cancelled]: {
    color: COLORS.status.cancelled,
    text: 'cancelled',
    icon: 'keyboard_arrow_up',
  },
};

export const WIDGET_TYPES = {
  alarmList: 'AlarmsList',
  context: 'Context',
  weather: 'ServiceWeather',
  statsHistogram: 'StatsHistogram',
  statsCurves: 'StatsCurves',
  statsTable: 'StatsTable',
  statsCalendar: 'StatsCalendar',
  statsNumber: 'StatsNumber',
  statsPareto: 'StatsPareto',
  text: 'Text',
  counter: 'Counter',
};

export const WIDGET_ICONS = {
  [WIDGET_TYPES.alarmList]: 'view_list',
  [WIDGET_TYPES.context]: 'view_list',
  [WIDGET_TYPES.weather]: 'view_module',
  [WIDGET_TYPES.statsHistogram]: 'bar_chart',
  [WIDGET_TYPES.statsCurves]: 'show_chart',
  [WIDGET_TYPES.statsTable]: 'table_chart',
  [WIDGET_TYPES.statsCalendar]: 'calendar_today',
  [WIDGET_TYPES.statsNumber]: 'table_chart',
  [WIDGET_TYPES.statsPareto]: 'multiline_chart',
  [WIDGET_TYPES.text]: 'view_headline',
  [WIDGET_TYPES.counter]: 'view_module',
};

export const SIDE_BARS = {
  alarmSettings: 'alarm-settings',
  contextSettings: 'context-settings',
  weatherSettings: 'weather-settings',
  statsHistogramSettings: 'stats-histogram-settings',
  statsCurvesSettings: 'stats-curves-settings',
  statsTableSettings: 'stats-table-settings',
  statsCalendarSettings: 'stats-calendar-settings',
  statsNumberSettings: 'stats-number-settings',
  statsParetoSettings: 'stats-pareto-settings',
  textSettings: 'text-settings',
  counterSettings: 'counter-settings',
};

export const SIDE_BARS_BY_WIDGET_TYPES = {
  [WIDGET_TYPES.alarmList]: SIDE_BARS.alarmSettings,
  [WIDGET_TYPES.context]: SIDE_BARS.contextSettings,
  [WIDGET_TYPES.weather]: SIDE_BARS.weatherSettings,
  [WIDGET_TYPES.statsTable]: SIDE_BARS.statsTableSettings,
  [WIDGET_TYPES.statsCalendar]: SIDE_BARS.statsCalendarSettings,
  [WIDGET_TYPES.statsNumber]: SIDE_BARS.statsNumberSettings,
  [WIDGET_TYPES.statsHistogram]: SIDE_BARS.statsHistogramSettings,
  [WIDGET_TYPES.statsCurves]: SIDE_BARS.statsCurvesSettings,
  [WIDGET_TYPES.statsPareto]: SIDE_BARS.statsParetoSettings,
  [WIDGET_TYPES.text]: SIDE_BARS.textSettings,
  [WIDGET_TYPES.counter]: SIDE_BARS.counterSettings,
};

export const EVENT_ENTITY_STYLE = {
  [EVENT_ENTITY_TYPES.ack]: {
    color: COLORS.entitiesEvents.ack,
    icon: 'playlist_add_check',
  },
  [EVENT_ENTITY_TYPES.fastAck]: {
    icon: 'check',
  },
  [EVENT_ENTITY_TYPES.pbehaviorAdd]: {
    icon: 'pause',
  },
  [EVENT_ENTITY_TYPES.pbehaviorList]: {
    icon: 'list',
  },
  [EVENT_ENTITY_TYPES.ackRemove]: {
    color: COLORS.entitiesEvents.ackRemove,
    icon: 'not_interested',
  },
  [EVENT_ENTITY_TYPES.declareTicket]: {
    color: COLORS.entitiesEvents.declareTicket,
    icon: 'report_problem',
  },
  [EVENT_ENTITY_TYPES.assocTicket]: {
    icon: 'local_play',
  },
  [EVENT_ENTITY_TYPES.delete]: {
    icon: 'delete',
  },
  [EVENT_ENTITY_TYPES.changeState]: {
    icon: 'thumbs_up_down',
  },
  [EVENT_ENTITY_TYPES.snooze]: {
    color: COLORS.entitiesEvents.snooze,
    icon: 'alarm',
  },
  [EVENT_ENTITY_TYPES.done]: {
    color: COLORS.entitiesEvents.done,
    icon: 'assignment_turned_in',
  },
  [EVENT_ENTITY_TYPES.validate]: {
    icon: 'thumb_up',
  },
  [EVENT_ENTITY_TYPES.invalidate]: {
    icon: 'thumb_down',
  },
  [EVENT_ENTITY_TYPES.pause]: {
    icon: 'pause',
  },
  [EVENT_ENTITY_TYPES.play]: {
    icon: 'play_arrow',
  },
  [EVENT_ENTITY_TYPES.groupRequest]: {
    icon: 'note_add',
  },
  [EVENT_ENTITY_TYPES.pbhenter]: {
    color: COLORS.entitiesEvents.pbhenter,
    icon: 'pause',
  },
  [EVENT_ENTITY_TYPES.pbhleave]: {
    color: COLORS.entitiesEvents.pbhleave,
    icon: 'play_arrow',
  },
  groupConsequences: {
    icon: 'center_focus_strong',
  },
  groupCauses: {
    icon: 'center_focus_weak',
  },
  [EVENT_ENTITY_TYPES.comment]: {
    color: COLORS.entitiesEvents.comment,
    icon: 'comment',
  },
  [EVENT_ENTITY_TYPES.manualMetaAlarmGroup]: {
    icon: 'center_focus_strong',
  },
  [EVENT_ENTITY_TYPES.manualMetaAlarmUngroup]: {
    icon: 'link_off',
  },
  [EVENT_ENTITY_TYPES.metaalarmattach]: {
    color: COLORS.entitiesEvents.metaalarmattach,
    icon: 'center_focus_weak',
  },
  [EVENT_ENTITY_TYPES.executeInstruction]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionStart]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionPause]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionResume]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionComplete]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionAbort]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionFail]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionJobStart]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionJobComplete]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionJobAbort]: {
    icon: 'assignment',
  },
  [EVENT_ENTITY_TYPES.instructionJobFail]: {
    icon: 'assignment',
  },
};

export const UNKNOWN_VALUE_STYLE = {
  color: COLORS.status.unknown,
  text: 'Invalid val',
  icon: 'clear',
};

export const FILTER_OPERATORS = {
  equal: 'equal',
  notEqual: 'not equal',
  in: 'in',
  notIn: 'not in',
  beginsWith: 'begins with',
  doesntBeginWith: 'doesn\'t begin with',
  contains: 'contains',
  doesntContains: 'doesn\'t contain',
  endsWith: 'ends with',
  doesntEndWith: 'doesn\'t end with',
  isEmpty: 'is empty',
  isEmptyArray: 'is empty array',
  isNotEmpty: 'is not empty',
  isNull: 'is null',
  isNotNull: 'is not null',
};

export const FILTER_OPERATORS_FOR_ARRAY = [FILTER_OPERATORS.in, FILTER_OPERATORS.notIn];

export const FILTER_MONGO_OPERATORS = {
  and: '$and',
  or: '$or',
  equal: '$eq',
  notEqual: '$ne',
  in: '$in',
  notIn: '$nin',
  regex: '$regex',
};

export const FILTER_INPUT_TYPES = {
  string: 'string',
  number: 'number',
  boolean: 'boolean',
  null: 'null',
};

export const FILTER_DEFAULT_VALUES = {
  condition: FILTER_MONGO_OPERATORS.and,
  rule: {
    field: '',
    operator: '',
    input: '',
    inputType: FILTER_INPUT_TYPES.string,
  },
  group: {
    condition: FILTER_MONGO_OPERATORS.and,
    groups: {},
    rules: {},
  },
};

export const DATETIME_FORMATS = {
  long: 'DD/MM/YYYY H:mm:ss',
  medium: 'DD/MM H:mm',
  short: 'DD/MM/YYYY',
  time: 'H:mm:ss',
  dateTimePicker: 'DD/MM/YYYY HH:mm',
  dateTimePickerWithSeconds: 'DD/MM/YYYY HH:mm:ss',
  datePicker: 'DD/MM/YYYY',
  timePicker: 'HH:mm',
  timePickerWithSeconds: 'HH:mm:ss',
  veeValidateDateTimeFormat: 'dd/MM/yyyy HH:mm',
  refreshFieldFormat: 'Y __ D __ H _ m _ s _',
};

export const STATS_OPTIONS = {
  recursive: 'recursive',
  states: 'states',
  authors: 'authors',
  sla: 'sla',
};

export const STATS_TYPES = {
  alarmsCreated: {
    value: 'alarms_created',
    options: [STATS_OPTIONS.recursive, STATS_OPTIONS.states, STATS_OPTIONS.authors],
  },
  alarmsResolved: {
    value: 'alarms_resolved',
    options: [STATS_OPTIONS.recursive, STATS_OPTIONS.states, STATS_OPTIONS.authors],
  },
  alarmsCanceled: {
    value: 'alarms_canceled',
    options: [STATS_OPTIONS.recursive, STATS_OPTIONS.states, STATS_OPTIONS.authors],
  },
  alarmsAcknowledged: {
    value: 'alarms_acknowledged',
    options: [STATS_OPTIONS.recursive, STATS_OPTIONS.states, STATS_OPTIONS.authors],
  },
  ackTimeSla: {
    value: 'ack_time_sla',
    options: [STATS_OPTIONS.recursive, STATS_OPTIONS.states, STATS_OPTIONS.authors, STATS_OPTIONS.sla],
  },
  resolveTimeSla: {
    value: 'resolve_time_sla',
    options: [STATS_OPTIONS.recursive, STATS_OPTIONS.states, STATS_OPTIONS.authors, STATS_OPTIONS.sla],
  },
  timeInState: {
    value: 'time_in_state',
    options: [STATS_OPTIONS.states],
  },
  stateRate: {
    value: 'state_rate',
    options: [STATS_OPTIONS.states],
  },
  mtbf: {
    value: 'mtbf',
    options: [STATS_OPTIONS.recursive],
  },
  currentState: {
    value: 'current_state',
    options: [],
  },
  ongoingAlarms: {
    value: 'ongoing_alarms',
    options: [STATS_OPTIONS.states],
  },
  currentOngoingAlarms: {
    value: 'current_ongoing_alarms',
    options: [STATS_OPTIONS.states],
  },
  currentOngoingAlarmsWithAck: {
    value: 'current_ongoing_alarms_with_ack',
    options: [STATS_OPTIONS.states],
  },
  currentOngoingAlarmsWithoutAck: {
    value: 'current_ongoing_alarms_without_ack',
    options: [STATS_OPTIONS.states],
  },
};

export const STATS_DURATION_UNITS = {
  hour: 'h',
  day: 'd',
  week: 'w',
  month: 'm',
};

export const STATS_CRITICITY = {
  ok: 'ok',
  minor: 'minor',
  major: 'major',
  critical: 'critical',
};

export const STATS_QUICK_RANGES = {
  custom: {
    value: 'custom',
  },
  last2Days: {
    value: 'last2Days',
    start: 'now-2d',
    stop: 'now',
  },
  last7Days: {
    value: 'last7Days',
    start: 'now-7d',
    stop: 'now',
  },
  last30Days: {
    value: 'last30Days',
    start: 'now-30d',
    stop: 'now',
  },
  last1Year: {
    value: 'last1Year',
    start: 'now-1y',
    stop: 'now',
  },
  yesterday: {
    value: 'yesterday',
    start: 'now-1d/d',
    stop: 'now-1d/d',
  },
  previousWeek: {
    value: 'previousWeek',
    start: 'now-1w/w',
    stop: 'now-1w/w',
  },
  previousMonth: {
    value: 'previousMonth',
    start: 'now-1m/m',
    stop: 'now-1m/m',
  },
  today: {
    value: 'today',
    start: 'now/d',
    stop: 'now/d',
  },
  todaySoFar: {
    value: 'todaySoFar',
    start: 'now/d',
    stop: 'now',
  },
  thisWeek: {
    value: 'thisWeek',
    start: 'now/w',
    stop: 'now/w',
  },
  thisWeekSoFar: {
    value: 'thisWeekSoFar',
    start: 'now/w',
    stop: 'now',
  },
  thisMonth: {
    value: 'thisMonth',
    start: 'now/m',
    stop: 'now/m',
  },
  thisMonthSoFar: {
    value: 'thisMonthSoFar',
    start: 'now/m',
    stop: 'now',
  },
  last1Hour: {
    value: 'last1Hour',
    start: 'now-1h',
    stop: 'now',
  },
  last3Hour: {
    value: 'last3Hour',
    start: 'now-3h',
    stop: 'now',
  },
  last6Hour: {
    value: 'last6Hour',
    start: 'now-6h',
    stop: 'now',
  },
  last12Hour: {
    value: 'last12Hour',
    start: 'now-12h',
    stop: 'now',
  },
  last24Hour: {
    value: 'last24Hour',
    start: 'now-24h',
    stop: 'now',
  },
};

export const STATS_DEFAULT_COLOR = COLORS.statsDefault;

export const STATS_DISPLAY_MODE = {
  value: 'value',
  criticity: 'criticity',
};

export const STATS_DISPLAY_MODE_PARAMETERS = {
  criticityLevels: {
    ok: 0,
    minor: 10,
    major: 20,
    critical: 30,
  },
  colors: {
    ok: COLORS.state.ok,
    minor: COLORS.state.minor,
    major: COLORS.state.major,
    critical: COLORS.state.critical,
  },
};

export const STATS_CURVES_POINTS_STYLES = {
  circle: 'circle',
  cross: 'cross',
  crossRot: 'crossRot',
  dash: 'dash',
  line: 'line',
  rect: 'rect',
  rectRounded: 'rectRounded',
  rectRot: 'rectRot',
  star: 'star',
  triangle: 'triangle',
};

export const WIDGET_MAX_SIZE = 12;

export const WIDGET_MIN_SIZE = 3;

export const RETRY_DEFAULT_DELAY = 180;

export const STATS_CALENDAR_COLORS = {
  alarm: {
    ok: COLORS.state.ok,
    minor: COLORS.state.minor,
    major: COLORS.state.major,
    critical: COLORS.state.critical,
  },
};

export const USERS_RIGHTS_MASKS = {
  default: 1,
  create: 8,
  read: 4,
  update: 2,
  delete: 1,
};

export const USERS_RIGHTS_TYPES = {
  crud: 'CRUD',
  rw: 'RW',
};

export const USERS_RIGHTS_TYPES_TO_MASKS = {
  [USERS_RIGHTS_TYPES.crud]: [
    USERS_RIGHTS_MASKS.create,
    USERS_RIGHTS_MASKS.read,
    USERS_RIGHTS_MASKS.update,
    USERS_RIGHTS_MASKS.delete,
  ],
  [USERS_RIGHTS_TYPES.rw]: [
    USERS_RIGHTS_MASKS.read,
    USERS_RIGHTS_MASKS.update,
    USERS_RIGHTS_MASKS.delete,
  ],
};

export const USER_RIGHTS_PREFIXES = {
  technical: {
    admin: 'models',
    exploitation: 'models_exploitation',
  },
  business: {
    common: 'common',
    alarmsList: 'listalarm',
    context: 'crudcontext',
    weather: 'serviceweather',
    counter: 'counter',
  },
  api: 'api',
};

export const USERS_RIGHTS = {
  technical: {
    view: `${USER_RIGHTS_PREFIXES.technical.admin}_userview`,
    role: `${USER_RIGHTS_PREFIXES.technical.admin}_role`,
    action: `${USER_RIGHTS_PREFIXES.technical.admin}_action`,
    user: `${USER_RIGHTS_PREFIXES.technical.admin}_user`,
    parameters: `${USER_RIGHTS_PREFIXES.technical.admin}_parameters`,
    broadcastMessage: `${USER_RIGHTS_PREFIXES.technical.admin}_broadcastMessage`,
    playlist: `${USER_RIGHTS_PREFIXES.technical.admin}_playlist`,
    planning: `${USER_RIGHTS_PREFIXES.technical.admin}_planning`,
    planningType: `${USER_RIGHTS_PREFIXES.technical.admin}_planningType`,
    planningReason: `${USER_RIGHTS_PREFIXES.technical.admin}_planningReason`,
    planningExceptions: `${USER_RIGHTS_PREFIXES.technical.admin}_planningExceptions`,
    remediation: `${USER_RIGHTS_PREFIXES.technical.admin}_remediation`,
    remediationInstruction: `${USER_RIGHTS_PREFIXES.technical.admin}_remediationInstruction`,
    remediationJob: `${USER_RIGHTS_PREFIXES.technical.admin}_remediationJob`,
    remediationConfiguration: `${USER_RIGHTS_PREFIXES.technical.admin}_remediationConfiguration`,
    engine: `${USER_RIGHTS_PREFIXES.technical.admin}_engine`,
    exploitation: {
      eventFilter: `${USER_RIGHTS_PREFIXES.technical.exploitation}_eventFilter`,
      pbehavior: `${USER_RIGHTS_PREFIXES.technical.exploitation}_pbehavior`,
      webhook: `${USER_RIGHTS_PREFIXES.technical.exploitation}_webhook`,
      snmpRule: `${USER_RIGHTS_PREFIXES.technical.exploitation}_snmpRule`,
      action: `${USER_RIGHTS_PREFIXES.technical.exploitation}_action`,
      heartbeat: `${USER_RIGHTS_PREFIXES.technical.exploitation}_heartbeat`,
      dynamicInfo: `${USER_RIGHTS_PREFIXES.technical.exploitation}_dynamicInfo`,
      metaAlarmRule: `${USER_RIGHTS_PREFIXES.technical.exploitation}_metaAlarmRule`,
    },
  },
  business: {
    alarmsList: {
      actions: {
        ack: `${USER_RIGHTS_PREFIXES.business.alarmsList}_ack`,
        fastAck: `${USER_RIGHTS_PREFIXES.business.alarmsList}_fastAck`,
        ackRemove: `${USER_RIGHTS_PREFIXES.business.alarmsList}_cancelAck`,
        pbehaviorAdd: `${USER_RIGHTS_PREFIXES.business.alarmsList}_pbehavior`,
        snooze: `${USER_RIGHTS_PREFIXES.business.alarmsList}_snoozeAlarm`,
        pbehaviorList: `${USER_RIGHTS_PREFIXES.business.alarmsList}_listPbehavior`,
        declareTicket: `${USER_RIGHTS_PREFIXES.business.alarmsList}_declareanIncident`,
        associateTicket: `${USER_RIGHTS_PREFIXES.business.alarmsList}_assignTicketNumber`,
        cancel: `${USER_RIGHTS_PREFIXES.business.alarmsList}_removeAlarm`,
        changeState: `${USER_RIGHTS_PREFIXES.business.alarmsList}_changeState`,
        history: `${USER_RIGHTS_PREFIXES.business.alarmsList}_history`,
        groupRequest: `${USER_RIGHTS_PREFIXES.business.alarmsList}_groupRequest`,
        manualMetaAlarmGroup: `${USER_RIGHTS_PREFIXES.business.alarmsList}_manualMetaAlarmGroup`,
        comment: `${USER_RIGHTS_PREFIXES.business.alarmsList}_comment`,

        listFilters: `${USER_RIGHTS_PREFIXES.business.alarmsList}_listFilters`,
        editFilter: `${USER_RIGHTS_PREFIXES.business.alarmsList}_editFilter`,
        addFilter: `${USER_RIGHTS_PREFIXES.business.alarmsList}_addFilter`,
        userFilter: `${USER_RIGHTS_PREFIXES.business.alarmsList}_userFilter`,

        listRemediationInstructionsFilters:
          `${USER_RIGHTS_PREFIXES.business.alarmsList}_listRemediationInstructionsFilters`,
        editRemediationInstructionsFilter:
          `${USER_RIGHTS_PREFIXES.business.alarmsList}_editRemediationInstructionsFilter`,
        addRemediationInstructionsFilter:
          `${USER_RIGHTS_PREFIXES.business.alarmsList}_addRemediationInstructionsFilter`,
        userRemediationInstructionsFilter:
          `${USER_RIGHTS_PREFIXES.business.alarmsList}_userRemediationInstructionsFilter`,

        links: `${USER_RIGHTS_PREFIXES.business.alarmsList}_links`,

        correlation: `${USER_RIGHTS_PREFIXES.business.alarmsList}_correlation`,

        executeInstruction: `${USER_RIGHTS_PREFIXES.business.alarmsList}_executeInstruction`,

        variablesHelp: `${USER_RIGHTS_PREFIXES.business.common}_variablesHelp`,

        ...featuresService.get('constants.USERS_RIGHTS.business.alarmsList.actions'),
      },
    },
    context: {
      actions: {
        createEntity: `${USER_RIGHTS_PREFIXES.business.context}_createEntity`,
        editEntity: `${USER_RIGHTS_PREFIXES.business.context}_edit`,
        duplicateEntity: `${USER_RIGHTS_PREFIXES.business.context}_duplicate`,
        deleteEntity: `${USER_RIGHTS_PREFIXES.business.context}_delete`,
        pbehaviorAdd: `${USER_RIGHTS_PREFIXES.business.context}_pbehavior`,
        pbehaviorList: `${USER_RIGHTS_PREFIXES.business.context}_listPbehavior`,
        pbehaviorDelete: `${USER_RIGHTS_PREFIXES.business.context}_deletePbehavior`,

        listFilters: `${USER_RIGHTS_PREFIXES.business.context}_listFilters`,
        editFilter: `${USER_RIGHTS_PREFIXES.business.context}_editFilter`,
        addFilter: `${USER_RIGHTS_PREFIXES.business.context}_addFilter`,
        userFilter: `${USER_RIGHTS_PREFIXES.business.context}_userFilter`,
      },
    },
    weather: {
      actions: {
        entityAck: `${USER_RIGHTS_PREFIXES.business.weather}_entityAck`,
        entityAssocTicket: `${USER_RIGHTS_PREFIXES.business.weather}_entityAssocTicket`,
        entityComment: `${USER_RIGHTS_PREFIXES.business.weather}_entityComment`,
        entityValidate: `${USER_RIGHTS_PREFIXES.business.weather}_entityValidate`,
        entityInvalidate: `${USER_RIGHTS_PREFIXES.business.weather}_entityInvalidate`,
        entityPause: `${USER_RIGHTS_PREFIXES.business.weather}_entityPause`,
        entityPlay: `${USER_RIGHTS_PREFIXES.business.weather}_entityPlay`,
        entityCancel: `${USER_RIGHTS_PREFIXES.business.weather}_entityCancel`,
        entityManagePbehaviors: `${USER_RIGHTS_PREFIXES.business.weather}_entityManagePbehaviors`,

        entityLinks: `${USER_RIGHTS_PREFIXES.business.weather}_entityLinks`,

        moreInfos: `${USER_RIGHTS_PREFIXES.business.weather}_moreInfos`,
        alarmsList: `${USER_RIGHTS_PREFIXES.business.weather}_alarmsList`,
        pbehaviorList: `${USER_RIGHTS_PREFIXES.business.weather}_pbehaviorList`,

        variablesHelp: `${USER_RIGHTS_PREFIXES.business.common}_variablesHelp`,
      },
    },
    counter: {
      actions: {
        alarmsList: `${USER_RIGHTS_PREFIXES.business.counter}_alarmsList`,

        variablesHelp: `${USER_RIGHTS_PREFIXES.business.common}_variablesHelp`,
      },
    },
  },
  api: {
    alarmUpdate: `${USER_RIGHTS_PREFIXES.api}_alarm_update`,
    alarmDelete: `${USER_RIGHTS_PREFIXES.api}_alarm_delete`,
    alarmFilter: `${USER_RIGHTS_PREFIXES.api}_alarmfilter`,
    idleRule: `${USER_RIGHTS_PREFIXES.api}_idlerule`,
    eventFilter: `${USER_RIGHTS_PREFIXES.api}_eventfilter`,
    action: `${USER_RIGHTS_PREFIXES.api}_action`,
    webhook: `${USER_RIGHTS_PREFIXES.api}_webhook`,
    metaalarmrule: `${USER_RIGHTS_PREFIXES.api}_metaalarmrule`,
    playlist: `${USER_RIGHTS_PREFIXES.api}_playlist`,
    dynamicinfos: `${USER_RIGHTS_PREFIXES.api}_dynamicinfos`,
    heartbeat: `${USER_RIGHTS_PREFIXES.api}_heartbeat`,
    watcher: `${USER_RIGHTS_PREFIXES.api}_watcher`,
    viewgroup: `${USER_RIGHTS_PREFIXES.api}_viewgroup`,
    view: `${USER_RIGHTS_PREFIXES.api}_view`,
    pbehavior: `${USER_RIGHTS_PREFIXES.api}_pbehavior`,
    pbehaviorType: `${USER_RIGHTS_PREFIXES.api}_pbehaviortype`,
    pbehaviorReason: `${USER_RIGHTS_PREFIXES.api}_pbehaviorreason`,
    pbehaviorException: `${USER_RIGHTS_PREFIXES.api}_pbehaviorexception`,
    event: `${USER_RIGHTS_PREFIXES.api}_event`,
    engine: `${USER_RIGHTS_PREFIXES.api}_engine`,
    entityRead: `${USER_RIGHTS_PREFIXES.api}_entity_read`,
    entityUpdate: `${USER_RIGHTS_PREFIXES.api}_entity_update`,
    entityDelete: `${USER_RIGHTS_PREFIXES.api}_entity_delete`,
  },
};

export const NOT_COMPLETED_USER_RIGHTS = [
  USERS_RIGHTS.business.alarmsList.actions.links,
  USERS_RIGHTS.business.weather.actions.entityLinks,
];

export const WIDGETS_ACTIONS_TYPES = {
  alarmsList: {
    ack: 'ack',
    fastAck: 'fastAck',
    ackRemove: 'ackRemove',
    pbehaviorAdd: 'pbehaviorAdd',
    pbehaviorList: 'pbehaviorList',
    moreInfos: 'moreInfos',
    snooze: 'snooze',
    declareTicket: 'declareTicket',
    associateTicket: 'associateTicket',
    cancel: 'cancel',
    changeState: 'changeState',
    variablesHelp: 'variablesHelp',
    history: 'history',
    groupRequest: 'groupRequest',
    manualMetaAlarmGroup: 'manualMetaAlarmGroup',
    manualMetaAlarmUngroup: 'manualMetaAlarmUngroup',
    manualMetaAlarmUpdate: 'manualMetaAlarmUpdate',
    comment: 'comment',

    ...featuresService.get('constants.WIDGETS_ACTIONS_TYPES.alarmsList'),

    links: 'links',

    correlation: 'correlation',

    listFilters: 'listFilters',
    editFilter: 'editFilter',
    addFilter: 'addFilter',
    userFilter: 'userFilter',

    listRemediationInstructionsFilters: 'listRemediationInstructionsFilters',
    editRemediationInstructionsFilter: 'editRemediationInstructionsFilter',
    addRemediationInstructionsFilter: 'addRemediationInstructionsFilter',
    userRemediationInstructionsFilter: 'userRemediationInstructionsFilter',

    executeInstruction: 'executeInstruction',
  },
  context: {
    createEntity: 'createEntity',
    editEntity: 'editEntity',
    duplicateEntity: 'duplicateEntity',
    deleteEntity: 'deleteEntity',
    pbehaviorAdd: 'pbehaviorAdd',
    pbehaviorList: 'pbehaviorList',
    pbehaviorDelete: 'pbehaviorDelete',
    variablesHelp: 'variablesHelp',

    listFilters: 'listFilters',
    editFilter: 'editFilter',
    addFilter: 'addFilter',
  },
  weather: {
    entityAck: 'entityAck',
    entityAssocTicket: 'entityAssocTicket',
    entityValidate: 'entityValidate',
    entityInvalidate: 'entityInvalidate',
    entityPause: 'entityPause',
    entityPlay: 'entityPlay',
    entityCancel: 'entityCancel',
    entityComment: 'entityComment',

    entityLinks: 'entityLinks',

    moreInfos: 'moreInfos',
    alarmsList: 'alarmsList',
    pbehaviorList: 'pbehaviorList',
    variablesHelp: 'variablesHelp',
  },
  counter: {
    alarmsList: 'alarmsList',
    variablesHelp: 'variablesHelp',
  },
};

export const BUSINESS_USER_RIGHTS_ACTIONS_MAP = {
  alarmsList: {
    [WIDGETS_ACTIONS_TYPES.alarmsList.ack]: USERS_RIGHTS.business.alarmsList.actions.ack,
    [WIDGETS_ACTIONS_TYPES.alarmsList.fastAck]: USERS_RIGHTS.business.alarmsList.actions.fastAck,
    [WIDGETS_ACTIONS_TYPES.alarmsList.ackRemove]: USERS_RIGHTS.business.alarmsList.actions.ackRemove,
    [WIDGETS_ACTIONS_TYPES.alarmsList.pbehaviorAdd]: USERS_RIGHTS.business.alarmsList.actions.pbehaviorAdd,
    [WIDGETS_ACTIONS_TYPES.alarmsList.pbehaviorList]: USERS_RIGHTS.business.alarmsList.actions.pbehaviorList,
    [WIDGETS_ACTIONS_TYPES.alarmsList.snooze]: USERS_RIGHTS.business.alarmsList.actions.snooze,
    [WIDGETS_ACTIONS_TYPES.alarmsList.declareTicket]: USERS_RIGHTS.business.alarmsList.actions.declareTicket,
    [WIDGETS_ACTIONS_TYPES.alarmsList.associateTicket]: USERS_RIGHTS.business.alarmsList.actions.associateTicket,
    [WIDGETS_ACTIONS_TYPES.alarmsList.cancel]: USERS_RIGHTS.business.alarmsList.actions.cancel,
    [WIDGETS_ACTIONS_TYPES.alarmsList.changeState]: USERS_RIGHTS.business.alarmsList.actions.changeState,
    [WIDGETS_ACTIONS_TYPES.alarmsList.history]: USERS_RIGHTS.business.alarmsList.actions.history,
    [WIDGETS_ACTIONS_TYPES.alarmsList.variablesHelp]: USERS_RIGHTS.business.alarmsList.actions.variablesHelp,
    [WIDGETS_ACTIONS_TYPES.alarmsList.comment]: USERS_RIGHTS.business.alarmsList.actions.comment,
    [WIDGETS_ACTIONS_TYPES.alarmsList.groupRequest]: USERS_RIGHTS.business.alarmsList.actions.groupRequest,
    [WIDGETS_ACTIONS_TYPES.alarmsList.manualMetaAlarmGroup]:
      USERS_RIGHTS.business.alarmsList.actions.manualMetaAlarmGroup,
    [WIDGETS_ACTIONS_TYPES.alarmsList.manualMetaAlarmUngroup]:
      USERS_RIGHTS.business.alarmsList.actions.manualMetaAlarmGroup,
    [WIDGETS_ACTIONS_TYPES.alarmsList.manualMetaAlarmUpdate]:
      USERS_RIGHTS.business.alarmsList.actions.manualMetaAlarmGroup,

    [WIDGETS_ACTIONS_TYPES.alarmsList.links]: USERS_RIGHTS.business.alarmsList.actions.links,
    [WIDGETS_ACTIONS_TYPES.alarmsList.correlation]: USERS_RIGHTS.business.alarmsList.actions.correlation,

    [WIDGETS_ACTIONS_TYPES.alarmsList.listFilters]: USERS_RIGHTS.business.alarmsList.actions.listFilters,
    [WIDGETS_ACTIONS_TYPES.alarmsList.editFilter]: USERS_RIGHTS.business.alarmsList.actions.editFilter,
    [WIDGETS_ACTIONS_TYPES.alarmsList.addFilter]: USERS_RIGHTS.business.alarmsList.actions.addFilter,
    [WIDGETS_ACTIONS_TYPES.alarmsList.userFilter]: USERS_RIGHTS.business.alarmsList.actions.userFilter,

    [WIDGETS_ACTIONS_TYPES.alarmsList.listRemediationInstructionsFilters]:
      USERS_RIGHTS.business.alarmsList.actions.listRemediationInstructionsFilters,
    [WIDGETS_ACTIONS_TYPES.alarmsList.editRemediationInstructionsFilter]:
      USERS_RIGHTS.business.alarmsList.actions.editRemediationInstructionsFilter,
    [WIDGETS_ACTIONS_TYPES.alarmsList.addRemediationInstructionsFilter]:
      USERS_RIGHTS.business.alarmsList.actions.addRemediationInstructionsFilter,

    [WIDGETS_ACTIONS_TYPES.alarmsList.executeInstruction]: USERS_RIGHTS.business.alarmsList.actions.executeInstruction,

    ...featuresService.get('constants.BUSINESS_USER_RIGHTS_ACTIONS_MAP.alarmsList'),
  },

  context: {
    [WIDGETS_ACTIONS_TYPES.context.createEntity]: USERS_RIGHTS.business.context.actions.createEntity,
    [WIDGETS_ACTIONS_TYPES.context.editEntity]: USERS_RIGHTS.business.context.actions.editEntity,
    [WIDGETS_ACTIONS_TYPES.context.duplicateEntity]: USERS_RIGHTS.business.context.actions.duplicateEntity,
    [WIDGETS_ACTIONS_TYPES.context.deleteEntity]: USERS_RIGHTS.business.context.actions.deleteEntity,
    [WIDGETS_ACTIONS_TYPES.context.pbehaviorAdd]: USERS_RIGHTS.business.context.actions.pbehaviorAdd,
    [WIDGETS_ACTIONS_TYPES.context.pbehaviorList]: USERS_RIGHTS.business.context.actions.pbehaviorList,
    [WIDGETS_ACTIONS_TYPES.context.pbehaviorDelete]: USERS_RIGHTS.business.context.actions.pbehaviorDelete,

    [WIDGETS_ACTIONS_TYPES.context.listFilters]: USERS_RIGHTS.business.context.actions.listFilters,
    [WIDGETS_ACTIONS_TYPES.context.editFilter]: USERS_RIGHTS.business.context.actions.editFilter,
    [WIDGETS_ACTIONS_TYPES.context.addFilter]: USERS_RIGHTS.business.context.actions.addFilter,
  },

  weather: {
    [WIDGETS_ACTIONS_TYPES.weather.entityAck]: USERS_RIGHTS.business.weather.actions.entityAck,
    [WIDGETS_ACTIONS_TYPES.weather.entityAssocTicket]: USERS_RIGHTS.business.weather.actions.entityAssocTicket,
    [WIDGETS_ACTIONS_TYPES.weather.entityValidate]: USERS_RIGHTS.business.weather.actions.entityValidate,
    [WIDGETS_ACTIONS_TYPES.weather.entityInvalidate]: USERS_RIGHTS.business.weather.actions.entityInvalidate,
    [WIDGETS_ACTIONS_TYPES.weather.entityPause]: USERS_RIGHTS.business.weather.actions.entityPause,
    [WIDGETS_ACTIONS_TYPES.weather.entityPlay]: USERS_RIGHTS.business.weather.actions.entityPlay,
    [WIDGETS_ACTIONS_TYPES.weather.entityCancel]: USERS_RIGHTS.business.weather.actions.entityCancel,
    [WIDGETS_ACTIONS_TYPES.weather.entityManagePbehaviors]:
      USERS_RIGHTS.business.weather.actions.entityManagePbehaviors,

    [WIDGETS_ACTIONS_TYPES.weather.entityLinks]: USERS_RIGHTS.business.weather.actions.entityLinks,

    [WIDGETS_ACTIONS_TYPES.weather.moreInfos]: USERS_RIGHTS.business.weather.actions.moreInfos,
    [WIDGETS_ACTIONS_TYPES.weather.alarmsList]: USERS_RIGHTS.business.weather.actions.alarmsList,
    [WIDGETS_ACTIONS_TYPES.weather.pbehaviorList]: USERS_RIGHTS.business.weather.actions.pbehaviorList,
    [WIDGETS_ACTIONS_TYPES.weather.variablesHelp]: USERS_RIGHTS.business.weather.actions.variablesHelp,
    [WIDGETS_ACTIONS_TYPES.weather.entityComment]: USERS_RIGHTS.business.weather.actions.entityComment,
  },

  counter: {
    [WIDGETS_ACTIONS_TYPES.counter.alarmsList]: USERS_RIGHTS.business.counter.actions.alarmsList,
    [WIDGETS_ACTIONS_TYPES.counter.variablesHelp]: USERS_RIGHTS.business.counter.actions.variablesHelp,
  },
};

export const GROUPS_NAVIGATION_TYPES = {
  sideBar: 'side-bar',
  topBar: 'top-bar',
};

export const EVENT_FILTER_RULE_TYPES = {
  drop: 'drop',
  break: 'break',
  enrichment: 'enrichment',
};

export const META_ALARMS_RULE_TYPES = {
  relation: 'relation',
  timebased: 'timebased',
  attribute: 'attribute',
  complex: 'complex',
  valuegroup: 'valuegroup',

  /**
   * Manual group type doesn't using in the form
   * We are using it only inside alarms list widget
   */
  manualgroup: 'manualgroup',
};

export const META_ALARMS_THRESHOLD_TYPES = {
  thresholdRate: 'thresholdRate',
  thresholdCount: 'thresholdCount',
};

export const EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES = {
  pass: 'pass',
  break: 'break',
  drop: 'drop',
};

export const EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES = {
  setField: {
    value: 'set_field',
    options: {
      name: {
        text: 'Name',
        value: 'name',
        required: true,
      },
      value: {
        text: 'Value',
        value: 'value',
        required: true,
      },
    },
  },
  setFieldFromTemplate: {
    value: 'set_field_from_template',
    options: {
      name: {
        text: 'Name',
        value: 'name',
        required: true,
      },
      value: {
        text: 'Value',
        value: 'value',
        required: true,
      },
    },
  },
  setEntityInfoFromTemplate: {
    value: 'set_entity_info_from_template',
    options: {
      name: {
        text: 'Name',
        value: 'name',
        required: true,
      },
      description: {
        text: 'Description',
        value: 'description',
        required: false,
      },
      value: {
        text: 'Value',
        value: 'value',
        required: true,
      },
    },
  },
  copy: {
    value: 'copy',
    options: {
      from: {
        text: 'From',
        value: 'from',
        required: true,
      },
      to: {
        text: 'To',
        value: 'to',
        required: true,
      },
    },
  },
};

export const EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES_MAP = {
  [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField.value]: EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField,
  [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate.value]:
    EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate,

  [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate.value]:
    EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate,

  [EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy.value]: EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy,
};

export const SERVICE_WEATHER_WIDGET_MODAL_TYPES = {
  moreInfo: 'more-info',
  alarmList: 'alarm-list',
  both: 'both',
};

export const WEATHER_EVENT_DEFAULT_ENTITY = 'engine';

export const WEATHER_AUTOREMOVE_BYPAUSE_OUTPUT = 'MDS_AUTOREMOVE_BYPAUSE';

export const WEATHER_ACK_EVENT_OUTPUT = {
  ack: 'MDS_ACKNOWLEDGE',
  validateOk: 'MDS_VALIDATEOK',
  validateCancel: 'MDS_VALIDATECANCEL',
};

export const EVENT_DEFAULT_ORIGIN = 'canopsis';

export const SORT_ORDERS = {
  asc: 'ASC',
  desc: 'DESC',
};

export const WEBHOOK_TRIGGERS = {
  create: 'create',
  stateinc: 'stateinc',
  statedec: 'statedec',
  statusinc: 'statusinc',
  statusdec: 'statusdec',
  changestate: 'changestate',
  ack: 'ack',
  ackremove: 'ackremove',
  cancel: 'cancel',
  uncancel: 'uncancel',
  comment: 'comment',
  done: 'done',
  declareticket: 'declareticket',
  declareticketwebhook: 'declareticketwebhook',
  assocticket: 'assocticket',
  snooze: 'snooze',
  unsnooze: 'unsnooze',
  resolve: 'resolve',
  activate: 'activate',
};

export const EVENT_FILTER_RULE_OPERATORS = ['>=', '>', '<', '<=', 'regex_match'];

export const WEBHOOK_EVENT_FILTER_RULE_OPERATORS = ['>=', '>', '<', '<=', 'regex_match'];

export const SNMP_STATE_TYPES = {
  simple: 'simple',
  template: 'template',
};

export const ACTION_TYPES = {
  snooze: 'snooze',
  pbehavior: 'pbehavior',
  changeState: 'changestate',
  ack: 'ack',
  ackremove: 'ackremove',
  assocticket: 'assocticket',
  declareticket: 'declareticket',
  cancel: 'cancel',
};

export const ACTION_AUTHOR = 'engine-action';

export const ACTION_FORM_FIELDS_MAP_BY_TYPE = {
  [ACTION_TYPES.snooze]: 'snoozeParameters',
  [ACTION_TYPES.pbehavior]: 'pbehaviorParameters',
  [ACTION_TYPES.changeState]: 'changeStateParameters',
  [ACTION_TYPES.ack]: 'ackParameters',
  [ACTION_TYPES.ackremove]: 'ackremoveParameters',
  [ACTION_TYPES.assocticket]: 'assocticketParameters',
  [ACTION_TYPES.declareticket]: 'declareticketParameters',
  [ACTION_TYPES.cancel]: 'cancelParameters',
};

export const CANOPSIS_STACK = {
  go: 'go',
  python: 'python',
};

export const CANOPSIS_EDITION = {
  core: 'core',
  cat: 'cat',
};

export const CANOPSIS_DOCUMENTATION = 'https://doc.canopsis.net';

export const CANOPSIS_WEBSITE = 'https://www.capensis.fr/canopsis/';

export const CANOPSIS_FORUM = 'https://community.capensis.org/';

export const ALARMS_LIST_TIME_LINE_SYSTEM_AUTHOR = 'canopsis.engine';

export const HEARTBEAT_DURATION_UNITS = {
  minute: 'm',
  hour: 'h',
};

export const TIME_UNITS = {
  second: 's',
  minute: 'm',
  hour: 'h',
  day: 'd',
  week: 'w',
  month: 'M',
  year: 'y',
};

export const AVAILABLE_SORTED_TIME_UNITS = [
  TIME_UNITS.year,
  TIME_UNITS.month,
  TIME_UNITS.day,
  TIME_UNITS.hour,
  TIME_UNITS.minute,
  TIME_UNITS.second,
];

export const DEFAULT_DURATION_FORMAT = 'D __ H _ m _ s _';

export const AVAILABLE_TIME_UNITS = {
  second: {
    text: 'common.times.second',
    value: TIME_UNITS.second,
  },
  minute: {
    text: 'common.times.minute',
    value: TIME_UNITS.minute,
  },
  hour: {
    text: 'common.times.hour',
    value: TIME_UNITS.hour,
  },
  day: {
    text: 'common.times.day',
    value: TIME_UNITS.day,
  },
  week: {
    text: 'common.times.week',
    value: TIME_UNITS.week,
  },
  month: {
    text: 'common.times.month',
    value: TIME_UNITS.month,
  },
  year: {
    text: 'common.times.year',
    value: TIME_UNITS.year,
  },
};

export const DURATION_UNITS = {
  minute: AVAILABLE_TIME_UNITS.minute,
  hour: AVAILABLE_TIME_UNITS.hour,
  day: AVAILABLE_TIME_UNITS.day,
  week: AVAILABLE_TIME_UNITS.week,
  month: AVAILABLE_TIME_UNITS.month,
  year: AVAILABLE_TIME_UNITS.month,
};

export const SNOOZE_DURATION_UNITS = {
  second: AVAILABLE_TIME_UNITS.second,
  ...DURATION_UNITS,
};

export const ALARM_ENTITY_FIELDS = {
  connector: 'v.connector',
  connectorName: 'v.connector_name',
  component: 'v.component',
  resource: 'v.resource',
  output: 'v.output',
  extraDetails: 'extra_details',
  state: 'v.state.val',
  status: 'v.status.val',
};

export const DEFAULT_ALARMS_WIDGET_COLUMNS = [
  {
    labelKey: 'tables.alarmGeneral.connector',
    value: ALARM_ENTITY_FIELDS.connector,
  },
  {
    labelKey: 'tables.alarmGeneral.connectorName',
    value: ALARM_ENTITY_FIELDS.connectorName,
  },
  {
    labelKey: 'tables.alarmGeneral.component',
    value: ALARM_ENTITY_FIELDS.component,
  },
  {
    labelKey: 'tables.alarmGeneral.resource',
    value: ALARM_ENTITY_FIELDS.resource,
  },
  {
    labelKey: 'tables.alarmGeneral.output',
    value: ALARM_ENTITY_FIELDS.output,
  },
  {
    labelKey: 'tables.alarmGeneral.extraDetails',
    value: ALARM_ENTITY_FIELDS.extraDetails,
  },
  {
    labelKey: 'tables.alarmGeneral.state',
    value: ALARM_ENTITY_FIELDS.state,
  },
  {
    labelKey: 'tables.alarmGeneral.status',
    value: ALARM_ENTITY_FIELDS.status,
  },
];

export const DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS = [
  {
    labelKey: 'tables.alarmGeneral.connector',
    value: ALARM_ENTITY_FIELDS.connector,
  },
  {
    labelKey: 'tables.alarmGeneral.connectorName',
    value: ALARM_ENTITY_FIELDS.connectorName,
  },
  {
    labelKey: 'tables.alarmGeneral.resource',
    value: ALARM_ENTITY_FIELDS.resource,
  },
  {
    labelKey: 'tables.alarmGeneral.output',
    value: ALARM_ENTITY_FIELDS.output,
  },
  {
    labelKey: 'tables.alarmGeneral.extraDetails',
    value: ALARM_ENTITY_FIELDS.extraDetails,
  },
  {
    labelKey: 'tables.alarmGeneral.state',
    value: ALARM_ENTITY_FIELDS.state,
  },
  {
    labelKey: 'tables.alarmGeneral.status',
    value: ALARM_ENTITY_FIELDS.status,
  },
];

export const PERIODIC_REFRESH_UNITS = {
  second: AVAILABLE_TIME_UNITS.second,
  minute: AVAILABLE_TIME_UNITS.minute,
  hour: AVAILABLE_TIME_UNITS.hour,
};

export const DEFAULT_PERIODIC_REFRESH = {
  interval: 60,
  enabled: false,
  unit: TIME_UNITS.second,
};

export const DEFAULT_TIME_INTERVAL = {
  interval: 60,
  unit: TIME_UNITS.second,
};

export const ADMIN_PAGES_RULES = {
  remediation: { stack: CANOPSIS_STACK.go, edition: CANOPSIS_EDITION.cat },
};

export const EXPLOITATION_PAGES_RULES = {
  eventFilter: { stack: CANOPSIS_STACK.go },
  webhooks: { stack: CANOPSIS_STACK.go, edition: CANOPSIS_EDITION.cat },
  snmpRule: { edition: CANOPSIS_EDITION.cat },
  heartbeat: { stack: CANOPSIS_STACK.go },
  action: { stack: CANOPSIS_STACK.go },
  dynamicInfo: { edition: CANOPSIS_EDITION.cat },
  metaAlarmRule: { stack: CANOPSIS_STACK.go, edition: CANOPSIS_EDITION.cat },
};

export const USER_RIGHTS_TO_PAGES_RULES = {
  /**
   * Admin pages
   */
  [USERS_RIGHTS.technical.remediation]: ADMIN_PAGES_RULES.remediation,

  /**
   * Exploitation pages
   */
  [USERS_RIGHTS.technical.exploitation.eventFilter]: EXPLOITATION_PAGES_RULES.eventFilter,
  [USERS_RIGHTS.technical.exploitation.webhook]: EXPLOITATION_PAGES_RULES.webhooks,
  [USERS_RIGHTS.technical.exploitation.snmpRule]: EXPLOITATION_PAGES_RULES.snmpRule,
  [USERS_RIGHTS.technical.exploitation.heartbeat]: EXPLOITATION_PAGES_RULES.heartbeat,
  [USERS_RIGHTS.technical.exploitation.action]: EXPLOITATION_PAGES_RULES.action,
  [USERS_RIGHTS.technical.exploitation.dynamicInfo]: EXPLOITATION_PAGES_RULES.dynamicInfo,
  [USERS_RIGHTS.technical.exploitation.metaAlarmRule]: EXPLOITATION_PAGES_RULES.metaAlarmRule,
};

export const WIDGET_TYPES_RULES = {
  [WIDGET_TYPES.statsHistogram]: { edition: CANOPSIS_EDITION.cat },
  [WIDGET_TYPES.statsCurves]: { edition: CANOPSIS_EDITION.cat },
  [WIDGET_TYPES.statsTable]: { edition: CANOPSIS_EDITION.cat },
  [WIDGET_TYPES.statsCalendar]: { edition: CANOPSIS_EDITION.cat },
  [WIDGET_TYPES.statsNumber]: { edition: CANOPSIS_EDITION.cat },
  [WIDGET_TYPES.statsPareto]: { edition: CANOPSIS_EDITION.cat },
};

export const POPUP_TYPES = {
  success: 'success',
  info: 'info',
  warning: 'warning',
  error: 'error',
};

export const ALARMS_GROUP_PREFIX = {
  CAUSES: 'causes_',
  CONSEQUENCES: 'consequences_',
};

export const GRID_SIZES = {
  min: 0,
  max: 12,
  step: 1,
};

export const TOURS = {
  alarmsExpandPanel: 'alarmsExpandPanel',
};

export const DEFAULT_BROADCAST_MESSAGE_COLOR = '#e75e40';

export const BROADCAST_MESSAGES_STATUSES = {
  active: 0,
  pending: 1,
  expired: 2,
};

export const AVAILABLE_COUNTERS = {
  total: 'total',
  total_active: 'total_active',
  snooze: 'snooze',
  ack: 'ack',
  ticket: 'ticket',
  pbehavior_active: 'pbehavior_active',
};

export const DEFAULT_COUNTER_BLOCK_TEMPLATE = `<h2 style="text-align: justify;">
  <strong>{{ counter.filter.title }}</strong></h2><br>
  <center><strong><span style="font-size: 18px;">{{ counter.total_active }} alarmes actives</span></strong></center>
  <br>Seuil mineur à {{ levels.values.minor }}, seuil critique à {{ levels.values.critical }}
  <p style="text-align: justify;">{{ counter.ack }} acquittées, {{ counter.ticket}} avec ticket</p>`;

export const PORTALS_NAMES = {
  additionalTopBarItems: 'additional-top-bar-items',
};

export const WIDGET_GRID_SIZES_KEYS = {
  mobile: 'mobile',
  tablet: 'tablet',
  desktop: 'desktop',
};

export const MQ_KEYS_TO_WIDGET_GRID_SIZES_KEYS_MAP = {
  m: WIDGET_GRID_SIZES_KEYS.mobile,
  t: WIDGET_GRID_SIZES_KEYS.tablet,
  l: WIDGET_GRID_SIZES_KEYS.desktop,
  xl: WIDGET_GRID_SIZES_KEYS.desktop,
};

export const WIDGET_LAYOUT_MAX_WIDTHS = {
  [WIDGET_GRID_SIZES_KEYS.desktop]: '100%',
  [WIDGET_GRID_SIZES_KEYS.tablet]: `${MEDIA_QUERIES_BREAKPOINTS.t}px`,
  [WIDGET_GRID_SIZES_KEYS.mobile]: `${MEDIA_QUERIES_BREAKPOINTS.m}px`,
};

export const WIDGET_GRID_SIZES_STYLES = {
  [WIDGET_GRID_SIZES_KEYS.mobile]: {
    value: WIDGET_GRID_SIZES_KEYS.mobile,
    icon: 'stay_primary_portrait',
  },
  [WIDGET_GRID_SIZES_KEYS.tablet]: {
    value: WIDGET_GRID_SIZES_KEYS.tablet,
    icon: 'tablet_mac',
  },
  [WIDGET_GRID_SIZES_KEYS.desktop]: {
    value: WIDGET_GRID_SIZES_KEYS.desktop,
    icon: 'desktop_windows',
  },
};

export const WIDGET_GRID_ROW_HEIGHT = 20;

export const WIDGET_GRID_COLUMNS_COUNT = 12;

export const META_ALARM_EVENT_DEFAULT_FIELDS = {
  component: 'metaalarm',
  connector: 'engine',
  connector_name: 'correlation',
  source_type: 'metaalarm',
};

export const PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES = {
  selected: 0,
  all: 1,
};

export const DEFAULT_TIMEZONE = 'Europe/Paris';

export const PLANNING_TABS = {
  types: 'types',
  reasons: 'reasons',
  exceptions: 'exceptions',
};

export const COUNTERS_LIMIT = 3;

export const PBEHAVIOR_RRULE_PERIODS_RANGES = {
  thisWeek: 'thisWeek',
  nextWeek: 'nextWeek',
  next2Weeks: 'next2Weeks',
  thisMonth: 'thisMonth',
  nextMonth: 'nextMonth',
};

export const REMEDIATION_TABS = {
  instructions: 'instructions',
  configurations: 'configurations',
  jobs: 'jobs',
};

export const REMEDIATION_WORKFLOW_TYPES = {
  stop: true,
  continue: false,
};

export const MAX_LIMIT = 10000;

export const REMEDIATION_CONFIGURATION_TYPES = {
  rundeck: 'rundeck',
  awx: 'awx',
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

export const REMEDIATION_INSTRUCTION_FILTER_ALL = 'all';

/**
 * 19/01/2038 @ 3:14am (UTC) in unix timestamp
 *
 * @type {number}
 */
export const MAX_PBEHAVIOR_DEFAULT_TSTOP = 2147483647;

export const ENGINES_NAMES = {
  event: 'event',
  webhook: 'engine-webhook',
  fifo: 'engine-fifo',
  axe: 'engine-axe',
  che: 'engine-che',
  pbehavior: 'engine-pbehavior',
  action: 'engine-action',
  watcher: 'engine-watcher',
  dynamicInfo: 'engine-dynamic-info',
  correlation: 'engine-correlation',
  heartbeat: 'engine-heartbeat',
};

export const ENGINES_QUEUE_NAMES = {
  webhook: 'Engine_webhook',
  fifo: 'Engine_fifo',
  axe: 'Engine_axe',
  che: 'Engine_che',
  pbehavior: 'Engine_pbehavior',
  action: 'Engine_action',
  watcher: 'Engine_watcher',
  dynamicInfo: 'Engine_dynamic_infos',
  correlation: 'Engine_correlation',
  heartbeat: 'Engine_heartbeat',
};

export const ENGINES_NAMES_TO_QUEUE_NAMES = {
  [ENGINES_QUEUE_NAMES.webhook]: ENGINES_NAMES.webhook,
  [ENGINES_QUEUE_NAMES.fifo]: ENGINES_NAMES.fifo,
  [ENGINES_QUEUE_NAMES.axe]: ENGINES_NAMES.axe,
  [ENGINES_QUEUE_NAMES.che]: ENGINES_NAMES.che,
  [ENGINES_QUEUE_NAMES.pbehavior]: ENGINES_NAMES.pbehavior,
  [ENGINES_QUEUE_NAMES.action]: ENGINES_NAMES.action,
  [ENGINES_QUEUE_NAMES.watcher]: ENGINES_NAMES.watcher,
  [ENGINES_QUEUE_NAMES.dynamicInfo]: ENGINES_NAMES.dynamicInfo,
  [ENGINES_QUEUE_NAMES.correlation]: ENGINES_NAMES.correlation,
  [ENGINES_QUEUE_NAMES.heartbeat]: ENGINES_NAMES.heartbeat,
};

export const CAT_ENGINES = [
  ENGINES_NAMES.correlation,
  ENGINES_NAMES.dynamicInfo,
  ENGINES_NAMES.webhook,
];

export const REQUEST_METHODS = {
  post: 'POST',
  get: 'GET',
  put: 'PUT',
  patch: 'PATCH',
  delete: 'DELETE',
  head: 'HEAD',
  connect: 'CONNECT',
  options: 'OPTIONS',
  trace: 'TRACE',
};
