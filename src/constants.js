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
  row: 'row',
  widget: 'widget',
  stat: 'stat',
  user: 'user',
  role: 'role',
  eventFilterRule: 'eventFilterRule',
};

export const MODALS = {
  createAckEvent: 'create-ack-event',
  createAssociateTicketEvent: 'create-associate-ticket-event',
  createCancelEvent: 'create-cancel-event',
  createChangeStateEvent: 'create-change-state-event',
  createDeclareTicketEvent: 'create-declare-ticket-event',
  createSnoozeEvent: 'create-snooze-event',
  createPbehavior: 'create-pbehavior',
  createEntity: 'create-entity',
  createWatcher: 'create-watcher',
  watcher: 'watcher',
  pbehaviorList: 'pbehavior-list',
  editLiveReporting: 'edit-live-reporting',
  moreInfos: 'more-infos',
  confirmation: 'confirmation',
  createWidget: 'create-widget',
  createFilter: 'create-filter',
  manageHistogramGroups: 'manage-histogram-groups',
  calendarAlarmsList: 'calendar-alarms-list',
  addStat: 'add-stat',
  colorPicker: 'color-picker',
  textEditor: 'text-editor',
  selectView: 'select-view',
  createView: 'create-view',
  createGroup: 'create-group',
  createUser: 'create-user',
  createRole: 'create-role',
  createRight: 'create-right',
  createEventFilterRule: 'create-event-filter-rule',
  createEventFilterRulePattern: 'create-event-filter-rule-pattern',
  eventFilterRuleActions: 'event-filter-rule-actions',
  eventFilterRuleExternalData: 'event-filter-rule-external-data',
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
};

export const ENTITY_INFOS_TYPE = {
  state: 'state',
  status: 'status',
  action: 'action',
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
    color: 'green',
    text: 'ok',
    icon: 'assistant_photo',
  },
  [ENTITIES_STATES.minor]: {
    color: 'gold',
    text: 'minor',
    icon: 'assistant_photo',
  },
  [ENTITIES_STATES.major]: {
    color: 'orange',
    text: 'major',
    icon: 'assistant_photo',
  },
  [ENTITIES_STATES.critical]: {
    color: 'red',
    text: 'critical',
    icon: 'assistant_photo',
  },
};

export const WATCHER_STATES_COLORS = {
  [ENTITIES_STATES.ok]: 'green darken-2',
  [ENTITIES_STATES.minor]: 'yellow darken-2',
  [ENTITIES_STATES.major]: 'orange darken-2',
  [ENTITIES_STATES.critical]: 'red darken-2',
};

export const PBEHAVIOR_TYPES = {
  maintenance: 'Maintenance',
  outOfSurveillance: 'Hors plage horaire de surveillance',
};

export const WEATHER_ICONS = {
  [ENTITIES_STATES.ok]: 'wb_sunny',
  [ENTITIES_STATES.minor]: 'person',
  [ENTITIES_STATES.major]: 'person',
  [ENTITIES_STATES.critical]: 'wb_cloudy',
  maintenance: 'build',
  outOfSurveillance: 'brightness_3',
  pause: 'pause',
};

export const WATCHER_PBEHAVIOR_COLOR = 'grey lighten-1';

export const ENTITY_STATUS_STYLES = {
  [ENTITIES_STATUSES.off]: {
    color: 'black',
    text: 'off',
    icon: 'keyboard_arrow_up',
  },
  [ENTITIES_STATUSES.ongoing]: {
    color: 'grey',
    text: 'ongoing',
    icon: 'keyboard_arrow_up',
  },
  [ENTITIES_STATUSES.stealthy]: {
    color: 'gold',
    text: 'stealthy',
    icon: 'keyboard_arrow_up',
  },
  [ENTITIES_STATUSES.flapping]: {
    color: 'orange',
    text: 'flapping',
    icon: 'keyboard_arrow_up',
  },
  [ENTITIES_STATUSES.cancelled]: {
    color: 'red',
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
};

export const EVENT_ENTITY_STYLE = {
  [EVENT_ENTITY_TYPES.ack]: {
    color: '#9c27b0',
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
    color: '#9c27b0',
    icon: 'not_interested',
  },
  [EVENT_ENTITY_TYPES.declareTicket]: {
    color: '#2196f3',
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
    color: '#e91e63',
    icon: 'alarm',
  },
  [EVENT_ENTITY_TYPES.done]: {
    color: 'green',
    icon: 'assignment_turned_in',
  },
};

export const UNKNOWN_VALUE_STYLE = {
  color: 'black',
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
  isNotEmpty: 'is not empty',
  isNull: 'is null',
  isNotNull: 'is not null',
};

export const FILTER_DEFAULT_VALUES = {
  condition: '$and',
  rule: {
    field: '',
    operator: '',
    input: '',
  },
  group: {
    condition: '$and',
    groups: {},
    rules: {},
  },
};

export const DATETIME_FORMATS = {
  long: 'DD/MM/YYYY H:mm:ss',
  short: 'DD/MM/YYYY',
  time: 'H:mm:ss',
};

export const STATS_TYPES = {
  alarmsCreated: {
    value: 'alarms_created',
    options: ['recursive', 'states', 'authors'],
  },
  alarmsResolved: {
    value: 'alarms_resolved',
    options: ['recursive', 'states', 'authors'],
  },
  alarmsCanceled: {
    value: 'alarms_canceled',
    options: ['recursive', 'states', 'authors'],
  },
  ackTimeSla: {
    value: 'ack_time_sla',
    options: ['recursive', 'states', 'authors', 'sla'],
  },
  resolveTimeSla: {
    value: 'resolve_time_sla',
    options: ['recursive', 'states', 'authors', 'sla'],
  },
  timeInState: {
    value: 'time_in_state',
    options: ['states'],
  },
  stateRate: {
    value: 'state_rate',
    options: ['states'],
  },
  mtbf: {
    value: 'mtbf',
    options: [],
  },
  currentState: {
    value: 'current_state',
    options: [],
  },
  ongoingAlarms: {
    value: 'ongoing_alarms',
    options: ['states'],
  },
  currentOngoingAlarms: {
    value: 'current_ongoing_alarms',
    options: ['states'],
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

export const WIDGET_MAX_SIZE = 12;

export const WIDGET_MIN_SIZE = 3;

export const STATS_CALENDAR_COLORS = {
  alarm: {
    ok: '#66BB6A',
    minor: '#FFEE58',
    major: '#FFA726',
    critical: '#FF7043',
  },
};

export const LIVE_REPORTING_INTERVALS = {
  today: 'today',
  yesterday: 'yesterday',
  last7Days: 'last7Days',
  last30Days: 'last30Days',
  thisMonth: 'thisMonth',
  lastMonth: 'lastMonth',
  custom: 'custom',
};

export const USERS_RIGHTS_MASKS = {
  create: 8,
  read: 4,
  update: 2,
  delete: 1,
};

export const USERS_RIGHTS_TYPES = {
  crud: 'CRUD',
  rw: 'RW',
};

export const USERS_RIGHTS = {
  technical: {
    view: 'models_userview',
    selector: 'models_selector',
    role: 'models_role',
    action: 'models_action',
    filter: 'models_filter',
    curve: 'models_curve',
    account: 'models_account',
    user: 'models_user',
    serie: 'models_serie',
    topology: 'models_topology',
    profile: 'models_profile',
    derogation: 'models_derogation',
    event: 'models_event',
    eventlog: 'models_eventlog',
    job: 'models_job',
    taskmail: 'models_taskmail',
    taskdataclean: 'models_taskdataclean',
    tasklinklist: 'models_tasklinklist',
    taskfile: 'models_taskfile',
    taskshell: 'models_taskshell',
    taskackcentreon: 'models_taskackcentreon',
    taskreporting: 'models_taskreporting',
  },
  business: {
    alarmList: {
      actions: {
        ack: 'listalarm_ack',
        fastAck: 'listalarm_fastAck',
        ackRemove: 'listalarm_cancelAck',
        pbehavior: 'listalarm_pbehavior',
        snooze: 'listalarm_snoozeAlarm',
        pbehaviorList: 'listalarm_listPbehavior',
        declareTicket: 'listalarm_declareanIncident',
        associateTicket: 'listalarm_assignTicketNumber',
        cancel: 'listalarm_removeAlarm',
        changeState: 'listalarm_changeState',
        listFilters: 'listalarm_listFilters',
        editFilter: 'listalarm_editFilter',
        addFilter: 'listalarm_addFilter',
      },
    },
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
