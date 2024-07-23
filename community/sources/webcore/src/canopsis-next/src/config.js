import sanitizeHtml from 'sanitize-html';

import { removeTrailingSlashes } from '@/helpers/url';

export const {
  BASE_URL = '',
  VUE_APP_API_HOST = '',
  VUE_APP_LOCAL_STORAGE_ACCESS_TOKEN_KEY = '',
  VUE_APP_PAGINATION_LIMIT = '',
  VUE_APP_OPEN_STREET_LAYER_URL = '',
  VUE_APP_LINKIFY_PROTOCOLS = '',
} = process.env;

export const APP_HOST = removeTrailingSlashes(`${window.location.origin}${BASE_URL}`);

export const API_HOST = VUE_APP_API_HOST || window.location.origin;

export const SOCKET_HOST = API_HOST.replace(/^http(s?)/, 'wss');

export const SOCKET_ROUTE = '/api/v4/ws';

export const SOCKET_URL = removeTrailingSlashes(`${SOCKET_HOST}${SOCKET_ROUTE}`);

export const OPEN_STREET_LAYER_URL = VUE_APP_OPEN_STREET_LAYER_URL;

export const ROUTER_MODE = 'history';

export const ROUTER_ACCESS_TOKEN_KEY = 'access_token';

export const LOCAL_STORAGE_ACCESS_TOKEN_KEY = VUE_APP_LOCAL_STORAGE_ACCESS_TOKEN_KEY || 'accessToken';

export const PAGINATION_LIMIT = parseInt(VUE_APP_PAGINATION_LIMIT, 10);

export const PAGINATION_PER_PAGE_VALUES = [5, 10, 20, 50, 100];

export const PAGINATION_TOTAL_VISIBLE = 7;

export const DEFAULT_SANITIZE_OPTIONS = {
  allowedTags: sanitizeHtml.defaults.allowedTags.concat([
    'h1', 'h2', 'u', 'nl', 'font', 'img', 'video', 'audio', 'area', 'map', 'strike', 'button', 'span', 'address',
    'bdo', 'cite', 'q', 'dfn', 'var', 'dl', 'dt', 'dd', 'section', 'article', 'colgroup', 'col',

    /**
     * VUE COMPONENTS
     */
    'router-link', 'c-alarm-chip', 'c-alarm-tags-chips', 'c-entity-tags-chips', 'c-copy-wrapper', 'c-links-list',
    'service-entities-list', 'v-icon',
  ]),
  allowedAttributes: {
    '*': [
      'style', 'title', 'class', 'id', 'v-if', 'name', 'autoplay', 'colspan', 'controls', 'dir', 'align', 'width',
      'height', 'role',
    ],
    a: ['href', 'name', 'target'],
    img: ['src', 'alt'],
    font: ['color', 'size', 'face'],
    marquee: ['direction'],
    'router-link': ['href', 'name', 'target', 'to'],
    'c-alarm-chip': ['value'],
    'c-alarm-tags-chips': [':alarm', ':selected-tag', 'closable-active', 'inline-count', '@select', '@close'],
    'c-entity-tags-chips': [':entity', 'inline-count'],
    'c-copy-wrapper': ['value'],
    'c-links-list': [':links', ':category'],
    'v-icon': [
      'color', 'dark', 'dense', 'disabled', 'large', 'light', 'right', 'size', 'small', 'tag', 'x-large', 'x-small',
    ],
    'service-entities-list': [
      ':service', ':service-entities', ':widget-parameters', ':options', ':total-items', ':actions-requests',
      'entity-name-field', '@refresh', '@update:options', '@add:action',
    ],
  },
  allowedSchemes: sanitizeHtml.defaults.allowedSchemes.concat(['data']),
  disallowedTagsMode: 'escape',
  disallowedTagsAutoClose: false,
};

export const DEFAULT_LINKIFY_OPTIONS = {
  target: '_blank',
  ignoreTags: ['script', 'style'],
  validate: (str, type, token) => token?.hasProtocol?.(),
};

export const LINKIFY_PROTOCOLS = [
  ...VUE_APP_LINKIFY_PROTOCOLS.split(',').map(protocol => protocol.trim()).filter(Boolean),

  'ssh',
];

export const DEFAULT_MAX_MULTI_SORT_COLUMNS_COUNT = 3;

export const DEFAULT_WEATHER_LIMIT = 120;

export const DEFAULT_LOCALE = 'en';

export const LOCALES = {
  en: 'en',
  fr: 'fr',
};

export const DEFAULT_VIEW_STATS_INTERVAL = 120000;

export const LOCALE_PRIORITIES = {
  default: 1,
  global: 2,
  personal: 3,
};

export const POPUP_AUTO_CLOSE_DELAY = 3000;

export const ZOOM_OVERLAY_DELAY = 1500;

export const VUETIFY_ANIMATION_DELAY = 300;

export const POPUP_TICK_DELAY = 100;

export const SIDE_BAR_WIDTH = 250;

export const TOP_BAR_HEIGHT = 48;

export const EXPAND_DEFAULT_MAX_LETTERS = 50;

export const SCHEMA_EMBEDDED_KEY = '_embedded';

export const MEDIA_QUERIES_BREAKPOINTS = {
  m: 450,
  t: 900,
  l: 1400,
  xl: 'Infinity',
};

export const MAX_LOGO_SIZE_IN_KB = 16384;

export const DEFAULT_APP_TITLE = 'Canopsis';

export const EXPORT_VIEWS_AND_GROUPS_FILENAME_PREFIX = 'canopsis_groups_views-';

export const SOCKET_ROOMS = {
  healthcheck: 'healthcheck',
  healthcheckStatus: 'healthcheck-status',
  messageRates: 'message-rates',
  loggedUserCount: 'logged-user-count',
  broadcastMessages: 'broadcast-messages',
  execution: 'execution',
  declareticket: 'declareticket',
  alarms: 'alarms',
  alarmDetails: 'alarm-details',
  icons: 'icons',
  testscenario: 'testscenario',
};

export const API_ROUTES = {
  login: '/api/v4/login',
  logout: '/api/v4/logout',
  loggedUserCount: '/api/v4/logged-user-count',
  currentUser: '/api/v4/account/me',
  maintenance: '/api/v4/maintenance',
  alarms: {
    list: '/api/v4/alarms',
    bulkList: '/api/v4/bulk/alarms',
  },
  componentAlarms: '/api/v4/component-alarms',
  resolvedAlarms: '/api/v4/resolved-alarms',
  alarmDetails: '/api/v4/alarm-details',
  openAlarms: '/api/v4/open-alarms',
  alarmLinks: '/api/v4/alarm-links',
  alarmDisplayNames: '/api/v4/alarm-display-names',
  alarmExecutions: '/api/v4/cat/alarm-executions',
  entity: '/api/v4/entities',
  entityComments: '/api/v4/entity-comments',
  entityContextGraph: '/api/v4/entities/context-graph',
  entityCheckStateSetting: '/api/v4/entities/check-state-setting',
  entityStateSetting: '/api/v4/entities/state-setting',
  bulkEntitiesEnable: '/api/v4/bulk/entities/enable',
  bulkEntitiesDisable: '/api/v4/bulk/entities/disable',
  entityBasics: '/api/v4/entitybasics',
  service: '/api/v4/entityservices',
  serviceDependencies: '/api/v4/entityservice-dependencies',
  serviceImpacts: '/api/v4/entityservice-impacts',
  serviceAlarms: '/api/v4/entityservice-alarms',
  entityInfosDictionaryKeys: '/api/v4/entity-infos-dictionary/keys',
  weatherService: '/api/v4/weather-services',
  alarmListExport: '/api/v4/alarm-export',
  alarmTag: {
    list: '/api/v4/alarm-tags',
    bulkList: '/api/v4/bulk/alarm-tags',
  },
  contextExport: '/api/v4/entity-export',
  event: '/api/v4/event',
  userPreferences: '/api/v4/user-preferences',
  view: {
    list: '/api/v4/views',
    copy: '/api/v4/view-copy',
    export: '/api/v4/view-export',
    import: '/api/v4/view-import',
    positions: '/api/v4/view-positions',
    groups: '/api/v4/view-groups',
    tabs: '/api/v4/view-tabs',
    tabCopy: '/api/v4/view-tab-copy',
    tabPositions: '/api/v4/view-tab-positions',
  },
  widget: {
    list: '/api/v4/widgets',
    gridPositions: '/api/v4/widget-grid-positions',
    filters: '/api/v4/widget-filters',
    filterPositions: '/api/v4/widget-filter-positions',
  },
  widgetTemplate: '/api/v4/widget-templates',
  permissions: '/api/v4/permissions',
  users: '/api/v4/users',
  roles: {
    list: '/api/v4/roles',
    templates: '/api/v4/role-templates',
  },
  eventFilter: {
    list: '/api/v4/eventfilter',
    rules: '/api/v4/eventfilter/rules',
  },
  file: '/api/v4/file',
  fileAccess: '/api/v4/file-access',
  infos: {
    app: '/api/v4/app-info',
    userInterface: '/api/v4/internal/user_interface',
  },
  associativeTable: '/api/v4/associativetable',
  sessionsCount: '/api/v4/sessions-count',
  broadcastMessage: {
    list: '/api/v4/broadcast-message',
    activeList: '/api/v4/active-broadcast-message',
  },
  counter: '/api/v4/alarm-counters',
  playlist: '/api/v4/playlists',
  pbehavior: {
    timespan: '/api/v4/pbehavior-timespans',
    exceptions: '/api/v4/pbehavior-exceptions',
    exceptionImport: '/api/v4/pbehavior-exception-import',
    types: '/api/v4/pbehavior-types',
    nextTypesPriority: '/api/v4/pbehavior-types/next-priority',
    pbehaviors: '/api/v4/pbehaviors',
    bulkPbehaviors: '/api/v4/bulk/pbehaviors',
    pbehaviorComments: '/api/v4/pbehavior-comments',
    entities: '/api/v4/entities/pbehaviors',
    entitiesCalendar: '/api/v4/entities/pbehavior-calendar',
    reasons: '/api/v4/pbehavior-reasons',
    calendar: '/api/v4/pbehavior-calendar',
    bulkEntityPbehaviors: '/api/v4/bulk/entity-pbehaviors',
  },
  engineRunInfo: '/api/v4/engine-runinfo',
  cas: {
    login: '/api/v4/cas/login',
    loggedin: '/api/v4/cas/loggedin',
  },
  saml: {
    auth: '/api/v4/saml/auth',
  },
  oauth: {
    login: '/api/v4/oauth/',
  },
  scenario: {
    scenarios: '/api/v4/scenarios',
    testExecution: '/api/v4/cat/test-scenario-executions',
    testExecutionWebhooks: '/api/v4/cat/test-scenario-webhooks',
  },
  entityCategories: '/api/v4/entity-categories',
  stateSetting: '/api/v4/state-settings',
  dataStorage: '/api/v4/data-storage',
  notification: '/api/v4/notification',
  idleRules: '/api/v4/idle-rules',
  flappingRules: '/api/v4/flapping-rules',
  resolveRules: '/api/v4/resolve-rules',
  messageRateStats: '/api/v4/message-rate-stats',
  pattern: {
    list: '/api/v4/patterns',
    bulkList: '/api/v4/bulk/patterns',
    entitiesCount: '/api/v4/patterns-entities-count',
    alarmsCount: '/api/v4/patterns-alarms-count',
  },
  shareTokens: '/api/v4/share-tokens',
  techMetrics: '/api/v4/tech-metrics-export',
  templateVars: '/api/v4/template-vars',
  templateValidator: {
    declareTicketRules: '/api/v4/template-validator/declare-ticket-rules',
    scenarios: '/api/v4/template-validator/scenarios',
    eventFilterRules: '/api/v4/template-validator/event-filter-rules',
  },
  linkRule: '/api/v4/link-rules',
  bulkLinkRule: '/api/v4/bulk/link-rules',
  linkCategories: '/api/v4/link-categories',
  icons: '/api/v4/icons',
  themes: {
    list: '/api/v4/color-themes',
    bulkList: '/api/v4/bulk/color-themes',
  },
  healthcheck: {
    engines: '/api/v4/healthcheck',
    status: '/api/v4/healthcheck/status',
    enginesOrder: '/api/v4/healthcheck/engines-order',
    parameters: '/api/v4/healthcheck/parameters',
  },

  /**
   * Cat routes
   */
  filters: '/api/v4/cat/kpi-filters',
  ratingSettings: '/api/v4/cat/rating-settings',
  bulkRatingSettings: '/api/v4/cat/rating-settings/bulk',
  dynamicInfo: '/api/v4/cat/dynamic-infos',
  dynamicInfosDictionaryKeys: '/api/v4/cat/dynamic-infos-dictionary/keys',
  metaAlarmRule: '/api/v4/cat/metaalarmrules',
  remediation: {
    instructions: '/api/v4/cat/instructions',
    instructionStats: '/api/v4/cat/instruction-stats',
    jobs: '/api/v4/cat/jobs',
    configurations: '/api/v4/cat/job-configs',
    executions: '/api/v4/cat/executions',
    jobExecutions: '/api/v4/cat/job-executions',
    pausedExecutions: '/api/v4/cat/account/paused-executions',
    instructionComments: '/api/v4/cat/instruction-comments',
  },
  junit: {
    directory: '/api/v4/cat/junit/directory',
    testSuites: '/api/v4/cat/junit/test-suites',
    entityGantt: '/api/v4/cat/junit/test-suites-entity-gantt',
    history: '/api/v4/cat/junit/test-suites-history',
    widget: '/api/v4/cat/junit/test-suites-widget',
    file: '/api/v4/cat/junit/test-cases-file',
  },
  metrics: {
    settings: '/api/v4/cat/metrics-settings',
    alarm: '/api/v4/cat/metrics/alarm',
    exportAlarm: '/api/v4/cat/metrics-export/alarm',
    exportRating: '/api/v4/cat/metrics-export/rating',
    exportSli: '/api/v4/cat/metrics-export/sli',
    exportAggregate: '/api/v4/cat/metrics-export/aggregate',
    exportRemediation: '/api/v4/cat/metrics-export/remediation',
    exportMetric: '/api/v4/cat/metrics-export',
    exportAvailabilityByEntity: '/api/v4/cat/metrics-export/availability-by-entity',
    exportAvailability: '/api/v4/cat/metrics-export/availability',
    sli: '/api/v4/cat/metrics/sli',
    rating: '/api/v4/cat/metrics/rating',
    remediation: '/api/v4/cat/metrics/remediation',
    aggregate: '/api/v4/cat/metrics/aggregate',
    perfDataMetrics: '/api/v4/cat/perf-data-metrics',
    entityAlarmMetrics: '/api/v4/cat/entity-metrics/alarm',
    entityAggregateMetrics: '/api/v4/cat/entity-metrics/aggregate',
    entityAvailability: '/api/v4/cat/entity-metrics/availability',
    entityAggregateAvailability: '/api/v4/cat/entity-metrics/aggregate-availability',
    group: '/api/v4/cat/metrics/group',
    exportGroup: '/api/v4/cat/metrics-export/group',
    availability: '/api/v4/cat/metrics/availability',
  },
  maps: '/api/v4/cat/maps',
  bulkMaps: '/api/v4/cat/maps/bulk',
  mapState: '/api/v4/cat/map-state',
  metaAlarm: '/api/v4/cat/meta-alarms',
  snmpRule: '/api/v4/cat/snmprules',
  snmpMib: '/api/v4/cat/snmpmibs',
  declareTicket: {
    rules: '/api/v4/cat/declare-ticket-rules',
    bulkRules: '/api/v4/cat/bulk/declare-ticket-rules',
    alarmsAssigned: '/api/v4/cat/declare-ticket-assigned',
    testExecution: '/api/v4/cat/test-declare-ticket-executions',
    testExecutionWebhooks: '/api/v4/cat/test-declare-ticket-webhooks',
    declareTicketExecution: '/api/v4/cat/declare-ticket-executions',
    bulkDeclareTicket: '/api/v4/cat/bulk/declare-ticket-executions',
  },
  tags: '/api/v4/cat/tags',
  privateView: {
    groups: '/api/v4/cat/private-view-groups',
  },
  eventsRecord: '/api/v4/cat/event-records',
  eventsRecordCurrent: '/api/v4/cat/event-records-current',
};

export const COLORS = {
  primary: '#2fab63',
  secondary: '#2b3e4f',
  accent: '#82b1ff',
  error: '#ff5252',
  info: '#2196f3',
  success: '#4caf50',
  warning: '#fb8c00',
  background: '#ffffff',
  activeColor: '#000000',
  table: {
    background: '#FFFFFF',
    rowColor: '#FFFFFF',
    shiftRowColor: '#F5F5F5',
    hoverRowColor: '#F5F5F5',
  },
  state: {
    ok: '#00a65a',
    minor: '#fcdc00',
    major: '#ff9900',
    critical: '#f56954',
    pause: '#808080',
  },
  status: {
    closed: '#2FAB63',
    stealthy: '#6A6A6A',
    flapping: '#2B3E4F',
    ongoing: '#F9A825',
    cancelled: '#757575',
    noEvents: '#FF5252',
    unknown: '#000',
  },
  alarmSteps: {
    ack: '#35485D',
    ackRemove: '#35485D',
    declareTicket: '#6A6A6A',
    snooze: '#FF5252',
    comment: '#5B6E7F',
    metaalarmattach: '#8C1D18',
    junitTestCaseUpdate: '#979797',
    activate: '#01579B',
    resolve: '#2FAB63',
    instruction: '#6A6A6A',
  },
  alarmExtraDetails: {
    ack: '#9c27b0',
    canceled: '#607d8b',
    children: '#6d4c41',
    parent: '#6d4c41',
    comment: '#ba68c8',
    snooze: '#e91e63',
    ticket: '#2196F3',
  },
  entitiesEvents: {
    ack: '#9c27b0',
    ackRemove: '#9c27b0',
    declareTicket: '#2196f3',
    snooze: '#e91e63',
    pbhenter: '#18ffff',
    pbhleave: '#18ffff',
    comment: '#ba68c8',
    metaalarmattach: '#6d4c41',
  },
  testSuiteStatuses: {
    passed: '#2fab63',
    error: '#ffa800',
    failed: '#ff5252',
    skipped: '#5a6D80',
  },
  healthcheck: {
    error: '#FF5252',
    warning: '#FFA800',
    unknown: '#C4C4C4',
    edgeGray: '#979797',
    edgeBlack: '#000000',
  },
  impactState: [
    '#2FAB63',
    '#7CB342',
    '#CDDC39',
    '#DCE775',
    '#FFF59D',
    '#FFF176',
    '#FFEB3B',
    '#FDD835',
    '#FBC02D',
    '#FFB300',
    '#F9A825',
    '#FF8F00',
    '#F57C00',
    '#FF6F00',
    '#EF6C00',
    '#FF6A3C',
    '#FF5722',
    '#F4511E',
    '#F94C16',
    '#F44336',
    '#E64A19',
    '#D84315',
    '#E53935',
    '#D32F2F',
    '#C62828',
    '#BF360C',
    '#B71C1C',
    '#C51162',
    '#C2185B',
    '#AD1457',
    '#880E4F',
  ],
  kpi: {
    downtime: '#c4c4c4',
    maintenance: '#ffa800',
    uptime: '#5b6E7f',
  },
  metrics: {
    activeAlarms: '#5a6d7e',
    createdAlarms: '#fda701',
    nonDisplayedAlarms: '#fd693b',
    instructionAlarms: '#7bb242',
    pbehaviorAlarms: '#d64315',
    correlationAlarms: '#fdef75',
    ackAlarms: '#fd5252',
    cancelAckAlarms: '#9b27af',
    ackActiveAlarms: '#81b0fd',
    ticketActiveAlarms: '#2faa63',
    withoutTicketActiveAlarms: '#747474',
    ratioCorrelation: '#c31162',
    ratioInstructions: '#cbda39',
    ratioTickets: '#2195f1',
    ratioNonDisplayed: '#f9bf2d',
    averageAck: '#f5c6ab',
    averageResolve: '#1afd01',
    totalUserActivity: '#1fbbd1',
    manualInstructionExecutedAlarms: '#2FAB63',
    ratioRemediatedAlarms: '#FFA800',
    remediationStatisticAssignedRemediations: '#FFA800',
    remediationStatisticExecutedRemediations: '#5B6E7F',
    remediationStatisticRatioRemediations: '#5B6E7F',
    notAckedAlarms: '#dce775',
    notAckedInHourAlarms: '#1fbcd3',
    notAckedInFourHoursAlarms: '#afb42b',
    notAckedInDayAlarms: '#fff176',
  },
  mermaid: {
    primaryColor: '#bfe4ce',
    primaryBorderColor: '#2faa62',
    textColor: '#323232',
    lineColor: '#323232',
    noteBkgColor: '#75818c',
    noteTextColor: '#fdfdfd',
  },
  remediation: {
    executionStatus: {
      running: '#5A6D80',
      paused: '#323232',
      completed: '#2FAB63',
      aborted: '#979797',
      failed: '#FF5252',
    },
  },
  flowchart: {
    selection: '#5b6e7f',
    shapes: [
      '#fd897f',
      '#fd7faa',
      '#e87ffa',
      '#b287fd',
      '#8b9dfd',
      '#7fd6fd',
      '#83fdfd',
      '#a6fde9',
      '#b8f4c8',
      '#cafd8f',
      '#f2fd80',
      '#fdfd8c',
      '#fde37e',
      '#fdcf7f',
      '#fd9d7f',
      '#d5cac6',
      '#f3f3f3',
      '#cdd6da',
    ],
    border: [
      '#fd1743',
      '#f30056',
      '#d300f7',
      '#641ffd',
      '#2978fd',
      '#00affd',
      '#00b7d2',
      '#1de7b5',
      '#00e475',
      '#75fd03',
      '#c4fd00',
      '#fde800',
      '#fdc200',
      '#fd9000',
      '#fd3d00',
      '#8c6d62',
      '#9d9d9d',
      '#5f7c8a',
    ],
    text: [
      '#b61c1c',
      '#870e4e',
      '#49148b',
      '#311b91',
      '#0d46a0',
      '#01569a',
      '#005f63',
      '#004c3f',
      '#1b5d20',
      '#33681e',
      '#9d9c24',
      '#f37e17',
      '#fd6e00',
      '#e45000',
      '#be360c',
      '#4d342e',
      '#414141',
      '#37464e',
    ],
    background: [
      '#ffffff',
      '#eceff1',
      '#fafafa',
      '#efebe9',
      '#fbe9e7',
      '#f9fbe7',
      '#fff8e1',
      '#fffde7',
      '#fff3e0',
      '#f1f8e9',
      '#e0f2f1',
      '#e0f7fa',
      '#e8f5e9',
      '#f3e5f5',
      '#ede7f6',
      '#e8eaf6',
      '#e3f2fd',
      '#fce4ec',
      '#ffebee',
    ],
  },
};

export const CSS_COLORS_VARS = {
  primary: 'var(--v-primary-base)',
  secondary: 'var(--v-secondary-base)',
  error: 'var(--v-error-base)',

  state: {
    ok: 'var(--v-state-ok-base)',
    minor: 'var(--v-state-minor-base)',
    major: 'var(--v-state-major-base)',
    critical: 'var(--v-state-critical-base)',
    pause: 'var(--v-state-pause-base)',
  },

  status: {
    closed: 'var(--v-status-closed-base)',
    stealthy: 'var(--v-status-stealthy-base)',
    flapping: 'var(--v-status-flapping-base)',
    ongoing: 'var(--v-status-ongoing-base)',
    cancelled: 'var(--v-status-cancelled-base)',
    noEvents: 'var(--v-status-noEvents-base)',
  },
};

export const DEFAULT_THEME_COLORS = {
  primary: COLORS.primary,
  secondary: COLORS.secondary,
  accent: COLORS.accent,
  error: COLORS.error,
  info: COLORS.info,
  success: COLORS.success,
  warning: COLORS.warning,
  background: COLORS.background,
  active_color: COLORS.activeColor,

  state: COLORS.state,
  status: COLORS.status,
  table: {
    background: COLORS.table.background,
    row_color: COLORS.table.rowColor,
    hover_row_color: COLORS.table.hoverRowColor,
  },
};

export const FILE_BASE_URL = `${API_HOST}${API_ROUTES.file}`;

export const DOCUMENTATION_BASE_URL = 'https://doc.canopsis.net/';

export const EXPORT_FETCHING_INTERVAL = 2000;

export const MAX_PBEHAVIOR_DATES_DIFF_YEARS = 5;

export const HEALTHCHECK_HISTORY_FILENAME_PREFIX = 'healthcheck_history-';

export const HEALTHCHECK_LAST_HOUR_FILENAME_PREFIX = 'healthcheck_last_hour-';

export const KPI_ALARM_METRICS_FILENAME_PREFIX = 'kpi_alarm_metrics';

export const KPI_RATING_METRICS_FILENAME_PREFIX = 'kpi_rating_metrics';

export const KPI_SLI_METRICS_FILENAME_PREFIX = 'kpi_sli_metrics';

export const REMEDIATION_STATISTICS_FILENAME_PREFIX = 'remediation_statistics';

export const AVAILABILITY_FILENAME_PREFIX = 'availability';

export const APP_INFO_FETCHING_INTERVAL = 10000;

export const MIN_CLICK_OUTSIDE_DELAY_AFTER_REGISTERED = 100;
