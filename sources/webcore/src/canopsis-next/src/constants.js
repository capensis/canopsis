import { COLORS } from '@/config';
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
  viewRow: 'viewRow',
  viewTab: 'viewTab',
  widget: 'widget',
  stat: 'stat',
  user: 'user',
  role: 'role',
  eventFilterRule: 'eventFilterRule',
  webhook: 'webhook',
  snmpRule: 'snmpRule',
  action: 'action',
  heartbeat: 'heartbeat',
};

export const MODALS = {
  createAckEvent: 'create-ack-event',
  confirmAckWithTicket: 'confirm-ack-with-ticket',
  createAssociateTicketEvent: 'create-associate-ticket-event',
  createCancelEvent: 'create-cancel-event',
  createChangeStateEvent: 'create-change-state-event',
  createDeclareTicketEvent: 'create-declare-ticket-event',
  createSnoozeEvent: 'create-snooze-event',
  variablesHelp: 'variables-help',
  createPbehavior: 'create-pbehavior',
  createEntity: 'create-entity',
  createWatcher: 'create-watcher',
  addEntityInfo: 'add-entity-info',
  watcher: 'watcher',
  createWatcherAssocTicketEvent: 'create-watcher-assoc-ticket-event',
  createWatcherPauseEvent: 'create-watcher-pause-event',
  pbehaviorList: 'pbehavior-list',
  editLiveReporting: 'edit-live-reporting',
  moreInfos: 'more-infos',
  infoPopupSetting: 'info-popup-setting',
  addInfoPopup: 'add-info-popup',
  confirmation: 'confirmation',
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
  createViewTab: 'create-view-tab',
  createGroup: 'create-group',
  createUser: 'create-user',
  createRole: 'create-role',
  createRight: 'create-right',
  createEventFilterRule: 'create-event-filter-rule',
  createEventFilterRulePattern: 'create-event-filter-rule-pattern',
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

export const PBEHAVIOR_TYPES = {
  maintenance: 'Maintenance',
  unmonitored: 'Hors plage horaire de surveillance',
  pause: 'pause',
};

export const PAUSE_REASONS = {
  authorisationProblem: 'Problème d\'habilitation',
  robotProblem: 'Problème Robot',
  scenarioProblem: 'Problème Scénario',
  flashFunctionnalProblem: 'Problème Flash Fonctionnel',
  other: 'Autre',
};

export const WEATHER_ICONS = {
  [WATCHER_STATES.ok]: 'wb_sunny',
  [WATCHER_STATES.minor]: 'person',
  [WATCHER_STATES.major]: 'person',
  [WATCHER_STATES.critical]: 'wb_cloudy',
  maintenance: 'build',
  unmonitored: 'brightness_3',
  [WATCHER_STATES.pause]: 'pause',
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
  isNotEmpty: 'is not empty',
  isNull: 'is null',
  isNotNull: 'is not null',
};

export const FILTER_INPUT_TYPES = {
  string: 'string',
  number: 'number',
  boolean: 'boolean',
  null: 'null',
};

export const FILTER_DEFAULT_VALUES = {
  condition: '$and',
  rule: {
    field: '',
    operator: '',
    input: '',
    inputType: FILTER_INPUT_TYPES.string,
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
  dateTimePicker: 'DD/MM/YYYY HH:mm',
  dateTimePickerWithSeconds: 'DD/MM/YYYY HH:mm:ss',
  datePicker: 'DD/MM/YYYY',
  timePicker: 'HH:mm',
  timePickerWithSeconds: 'HH:mm:ss',
  veeValidateDateTimeFormat: 'dd/MM/yyyy HH:mm',
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

export const STATS_CALENDAR_COLORS = {
  alarm: {
    ok: COLORS.state.ok,
    minor: COLORS.state.minor,
    major: COLORS.state.major,
    critical: COLORS.state.critical,
  },
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
    role: 'models_role',
    action: 'models_action',
    user: 'models_user',
    parameters: 'models_parameters',
    exploitation: {
      eventFilter: 'models_exploitation_eventFilter',
      pbehavior: 'models_exploitation_pbehavior',
      webhook: 'models_exploitation_webhook',
      snmpRule: 'models_exploitation_snmpRule',
      action: 'models_exploitation_action',
      heartbeat: 'models_exploitation_heartbeat',
    },
  },
  business: {
    alarmsList: {
      actions: {
        ack: 'listalarm_ack',
        fastAck: 'listalarm_fastAck',
        ackRemove: 'listalarm_cancelAck',
        pbehaviorAdd: 'listalarm_pbehavior',
        snooze: 'listalarm_snoozeAlarm',
        pbehaviorList: 'listalarm_listPbehavior',
        declareTicket: 'listalarm_declareanIncident',
        associateTicket: 'listalarm_assignTicketNumber',
        cancel: 'listalarm_removeAlarm',
        changeState: 'listalarm_changeState',
        listFilters: 'listalarm_listFilters',
        editFilter: 'listalarm_editFilter',
        addFilter: 'listalarm_addFilter',
        userFilter: 'listalarm_userFilter',

        links: 'listalarm_links',

        ...featuresService.get('constants.USERS_RIGHTS.business.alarmsList.actions'),
      },
    },
    context: {
      actions: {
        createEntity: 'crudcontext_createEntity',
        editEntity: 'crudcontext_edit',
        duplicateEntity: 'crudcontext_duplicate',
        deleteEntity: 'crudcontext_delete',
        pbehaviorAdd: 'crudcontext_pbehavior',
        pbehaviorList: 'crudcontext_listPbehavior',
        pbehaviorDelete: 'crudcontext_deletePbehavior',

        listFilters: 'crudcontext_listFilters',
        editFilter: 'crudcontext_editFilter',
        addFilter: 'crudcontext_addFilter',
        userFilter: 'crudcontext_userFilter',
      },
    },
    weather: {
      actions: {
        entityAck: 'serviceweather_entityAck',
        entityAssocTicket: 'serviceweather_entityAssocTicket',
        entityValidate: 'serviceweather_entityValidate',
        entityInvalidate: 'serviceweather_entityInvalidate',
        entityPause: 'serviceweather_entityPause',
        entityPlay: 'serviceweather_entityPlay',
        entityCancel: 'serviceweather_entityCancel',
        entityManagePbehaviors: 'serviceweather_entityManagePbehaviors',

        entityLinks: 'serviceweather_entityLinks',

        moreInfos: 'serviceweather_moreInfos',
        alarmsList: 'serviceweather_alarmsList',
      },
    },
  },
};

export const NOT_COMPLETED_USER_RIGHTS_KEYS = [
  'business.alarmsList.actions.links',
  'business.weather.actions.entityLinks',
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


    ...featuresService.get('constants.WIDGETS_ACTIONS_TYPES.alarmsList'),

    links: 'links',

    listFilters: 'listFilters',
    editFilter: 'editFilter',
    addFilter: 'addFilter',
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

    entityLinks: 'entityLinks',

    moreInfos: 'moreInfos',
    alarmsList: 'alarmsList',
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

    [WIDGETS_ACTIONS_TYPES.alarmsList.links]: USERS_RIGHTS.business.alarmsList.actions.links,

    [WIDGETS_ACTIONS_TYPES.alarmsList.listFilters]: USERS_RIGHTS.business.alarmsList.actions.listFilters,
    [WIDGETS_ACTIONS_TYPES.alarmsList.editFilter]: USERS_RIGHTS.business.alarmsList.actions.editFilter,
    [WIDGETS_ACTIONS_TYPES.alarmsList.addFilter]: USERS_RIGHTS.business.alarmsList.actions.addFilter,

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
  assocticket: 'assocticket',
  snooze: 'snooze',
  unsnooze: 'unsnooze',
  resolve: 'resolve',
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
};

export const ACTION_AUTHOR = 'engine-action';

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

export const DURATION_UNITS = {
  minute: {
    text: 'common.times.minute',
    value: 'm',
  },
  hour: {
    text: 'common.times.hour',
    value: 'h',
  },
  day: {
    text: 'common.times.day',
    value: 'd',
  },
  week: {
    text: 'common.times.week',
    value: 'w',
  },
  month: {
    text: 'common.times.month',
    value: 'M',
  },
  year: {
    text: 'common.times.year',
    value: 'y',
  },
};

export const EXPLOITATION_PAGES_RULES = {
  eventFilter: { stack: CANOPSIS_STACK.go },
  webhooks: { stack: CANOPSIS_STACK.go, edition: CANOPSIS_EDITION.cat },
  snmpRule: { stack: CANOPSIS_STACK.go, edition: CANOPSIS_EDITION.cat },
  heartbeat: { stack: CANOPSIS_STACK.go },
  action: { stack: CANOPSIS_STACK.go },
};

export const USER_RIGHTS_TO_EXPLOITATION_PAGES_RULES = {
  [USERS_RIGHTS.technical.exploitation.eventFilter]: EXPLOITATION_PAGES_RULES.eventFilter,
  [USERS_RIGHTS.technical.exploitation.webhook]: EXPLOITATION_PAGES_RULES.webhooks,
  [USERS_RIGHTS.technical.exploitation.snmpRule]: EXPLOITATION_PAGES_RULES.snmpRule,
  [USERS_RIGHTS.technical.exploitation.heartbeat]: EXPLOITATION_PAGES_RULES.heartbeat,
  [USERS_RIGHTS.technical.exploitation.action]: EXPLOITATION_PAGES_RULES.action,
};
