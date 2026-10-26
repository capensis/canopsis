import { removeTrailingSlashes } from '@/helpers/url';

export const {
  BASE_URL = '',
  VITE_API_HOST = '',
  VITE_LOCAL_STORAGE_ACCESS_TOKEN_KEY = '',
  VITE_PAGINATION_LIMIT = '',
  VITE_OPEN_STREET_LAYER_URL = '',
  VITE_LINKIFY_PROTOCOLS = '',
} = import.meta.env;

export const APP_HOST = removeTrailingSlashes(`${window.location.origin}${BASE_URL}`);

export const API_HOST = VITE_API_HOST || window.location.origin;

export const SOCKET_HOST = API_HOST.replace(/^http(s?)/, 'wss');

export const SOCKET_ROUTE = '/internal/ws';

export const SOCKET_URL = removeTrailingSlashes(`${SOCKET_HOST}${SOCKET_ROUTE}`);

export const OPEN_STREET_LAYER_URL = VITE_OPEN_STREET_LAYER_URL;

export const ROUTER_MODE = 'history';

export const ROUTER_ACCESS_TOKEN_KEY = 'access_token';

export const LOCAL_STORAGE_ACCESS_TOKEN_KEY = VITE_LOCAL_STORAGE_ACCESS_TOKEN_KEY || 'accessToken';

export const LOCAL_STORAGE_WARNING_POPUP_KEY = 'warningPopup';

export const PAGINATION_LIMIT = parseInt(VITE_PAGINATION_LIMIT, 10) || 10;

export const PAGINATION_PER_PAGE_VALUES = [5, 10, 20, 50, 100];

export const PAGINATION_TOTAL_VISIBLE = 7;

export const DEFAULT_SANITIZE_OPTIONS = {
  ADD_TAGS: [
    'h1',
    'h2',
    'u',
    'nl',
    'font',
    'img',
    'video',
    'audio',
    'area',
    'map',
    'strike',
    'button',
    'span',
    'address',
    'bdo',
    'cite',
    'q',
    'dfn',
    'var',
    'dl',
    'dt',
    'dd',
    'section',
    'article',
    'colgroup',
    'col',
    'center',

    'router-link',
    'c-alarm-chip',
    'c-alarm-tags-chips',
    'c-entity-tags-chips',
    'c-copy-wrapper',
    'c-links-list',
    'service-entities-list',
    'v-icon',
    'v-row',
    'v-chip',
    'c-remediation-instruction-execution-see-details',
    'c-request-helper',
    'c-now',

    'div',
  ],

  ADD_ATTR: [
    'style',
    'title',
    'class',
    'id',
    'v-if',
    'name',
    'autoplay',
    'colspan',
    'controls',
    'dir',
    'align',
    'width',
    'height',
    'role',

    'href',
    'target',
    'src',
    'alt',
    'color',
    'size',
    'face',
    'direction',
    'to',

    'data-id',
    'data-c',
    'format',
    'timezone',
    'realtime',
    'offset',

    ':alarm',
    ':selected-tags',
    ':entity',
    ':links',
    ':category',
    ':service',
    ':service-entities',
    ':widget-parameters',
    ':options',
    ':total-items',
    ':actions-requests',
    ':execution',

    '@select',
    '@close',
    '@refresh',
    '@update:options',
    '@add:action',

    'value',
    'name-filter',
    'regex-filter',
    'regex-filter-flags',
    'closable-active',
    'inline-count',
    'entity-name-field',
    'helper-id',
    'key',

    'dark',
    'dense',
    'disabled',
    'large',
    'light',
    'right',
    'small',
    'tag',
    'x-large',
    'x-small',
    'text-color',
  ],

  ALLOW_DATA_ATTR: true,

  ADD_DATA_URI_TAGS: ['img', 'video', 'audio'],
};

export const DEFAULT_LINKIFY_OPTIONS = {
  target: '_blank',
  ignoreTags: ['script', 'style'],
  validate: (str, type, token) => token?.hasProtocol?.(),
};

export const LINKIFY_PROTOCOLS = [
  ...import.meta.env.VITE_LINKIFY_PROTOCOLS.split(',')
    .map(protocol => protocol.trim())
    .filter(Boolean),

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

/** Default width (px) for plugin sidebars when `config.width` is omitted */
export const DEFAULT_SIDEBAR_DRAWER_WIDTH = 450;

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
  notifications: 'notifications',
  executions: 'executions',
  simplifiedManualExecutions: 'simplified-manual-executions',
  pbehaviorPatterns: 'pbehavior-patterns',
  llmChat: 'llmchat',
  anomalyMonitoredConnectorStates: 'anomaly-monitored-connector-states',
};

export const API_ROUTES = {
  login: '/internal/login',
  logout: '/internal/logout',
  loggedUserCount: '/internal/logged-user-count',
  currentUser: '/internal/account/me',
  maintenance: '/internal/maintenance',
  alarms: {
    list: '/internal/alarms',
    bulkList: '/internal/bulk/alarms',
  },
  componentAlarms: '/internal/component-alarms',
  resolvedAlarms: '/internal/resolved-alarms',
  alarmDetails: '/internal/alarm-details',
  openAlarms: '/internal/open-alarms',
  alarmLinks: '/internal/alarm-links',
  alarmDisplayNames: '/internal/alarm-display-names',
  alarmExecutions: '/internal/cat/alarm-executions',
  entity: '/internal/entities',
  entityComments: '/internal/entity-comments',
  entityContextGraph: '/internal/entities/context-graph',
  entityCheckStateSetting: '/internal/entities/check-state-setting',
  entityStateSetting: '/internal/entities/state-setting',
  entityPause: '/internal/entities/pause',
  entityResume: '/internal/entities/resume',
  bulkEntitiesEnable: '/internal/bulk/entities/enable',
  bulkEntitiesDisable: '/internal/bulk/entities/disable',
  bulkEntitiesPause: '/internal/bulk/entities/pause',
  bulkEntitiesResume: '/internal/bulk/entities/resume',
  entityBasics: '/internal/entitybasics',
  service: '/internal/entityservices',
  bulkService: '/internal/bulk/entityservices',
  serviceDependencies: '/internal/entityservice-dependencies',
  serviceImpacts: '/internal/entityservice-impacts',
  serviceAlarms: '/internal/entityservice-alarms',
  entityUpstream: '/internal/entity-upstream',
  entityDownstreams: '/internal/entity-downstreams',
  entityInfosDictionaryKeys: '/internal/entity-infos-dictionary/keys',
  entityInfosProperties: '/internal/entity-infos-properties',
  bulkEntityInfosProperties: '/internal/bulk/entity-infos-properties',
  weatherService: '/internal/weather-services',
  alarmListExport: '/internal/alarm-export',
  alarmTimelineExport: '/internal/alarm-timeline-export',
  alarmTag: {
    list: '/internal/alarm-tags',
    bulkList: '/internal/bulk/alarm-tags',
    label: '/internal/alarm-tag-labels',
  },
  contextExport: '/internal/entity-export',
  event: '/internal/event',
  userPreferences: '/internal/user-preferences',
  view: {
    list: '/internal/views',
    copy: '/internal/view-copy',
    export: '/internal/view-export',
    import: '/internal/view-import',
    positions: '/internal/view-positions',
    groups: '/internal/view-groups',
    tabs: '/internal/view-tabs',
    tabCopy: '/internal/view-tab-copy',
    tabPositions: '/internal/view-tab-positions',
  },
  widget: {
    list: '/internal/widgets',
    gridPositions: '/internal/widget-grid-positions',
    filters: '/internal/widget-filters',
    filterPositions: '/internal/widget-filter-positions',
  },
  widgetTemplate: '/internal/widget-templates',
  permissions: '/internal/permissions',
  users: '/internal/users',
  topUserByActiveConnects: '/internal/users/top/active-connects',
  bulkUsersEnable: '/internal/bulk/users/enable',
  bulkUsersDisable: '/internal/bulk/users/disable',
  bulkUsers: '/internal/bulk/users',
  roles: {
    list: '/internal/roles',
    templates: '/internal/role-templates',
    bulk: '/internal/bulk/roles',
    bulkPermissions: '/internal/bulk/role-permissions',
  },
  eventFilter: {
    list: '/internal/eventfilter',
    rules: '/internal/eventfilter/rules',
    bulk: '/internal/bulk/eventfilters',
    bulkEnable: '/internal/bulk/eventfilters/enable',
    bulkDisable: '/internal/bulk/eventfilters/disable',
  },
  file: '/internal/file',
  fileAccess: '/internal/file-access',
  infos: {
    app: '/internal/app-info',
    userInterface: '/internal/internal/user_interface',
  },
  associativeTable: '/internal/associativetable',
  sessionsCount: '/internal/sessions-count',
  broadcastMessage: {
    list: '/internal/broadcast-message',
    activeList: '/internal/active-broadcast-message',
  },
  commentTemplate: '/internal/comment-templates',
  counter: '/internal/alarm-counters',
  playlist: '/internal/playlists',
  bulkPlaylist: '/internal/bulk/playlists',
  bulkPlaylistEnable: '/internal/bulk/playlists/enable',
  bulkPlaylistDisable: '/internal/bulk/playlists/disable',
  pbehavior: {
    timespan: '/internal/pbehavior-timespans',
    exceptions: '/internal/pbehavior-exceptions',
    exceptionImport: '/internal/pbehavior-exception-import',
    bulkExceptions: '/internal/bulk/pbehavior-exceptions',
    bulkExceptionsHide: '/internal/bulk/pbehavior-exceptions/hide',
    bulkExceptionsUnhide: '/internal/bulk/pbehavior-exceptions/unhide',
    types: '/internal/pbehavior-types',
    nextTypesPriority: '/internal/pbehavior-types/next-priority',
    bulkTypes: '/internal/bulk/pbehavior-types',
    bulkTypesHide: '/internal/bulk/pbehavior-types/hide',
    bulkTypesUnhide: '/internal/bulk/pbehavior-types/unhide',
    pbehaviors: '/internal/pbehaviors',
    bulkPbehaviors: '/internal/bulk/pbehaviors',
    bulkPbehaviorsEnable: '/internal/bulk/pbehaviors/enable',
    bulkPbehaviorsDisable: '/internal/bulk/pbehaviors/disable',
    pbehaviorComments: '/internal/pbehavior-comments',
    entities: '/internal/entities/pbehaviors',
    entitiesCalendar: '/internal/entities/pbehavior-calendar',
    reasons: '/internal/pbehavior-reasons',
    bulkReasons: '/internal/bulk/pbehavior-reasons',
    bulkReasonsHide: '/internal/bulk/pbehavior-reasons/hide',
    bulkReasonsUnhide: '/internal/bulk/pbehavior-reasons/unhide',
    calendar: '/internal/pbehavior-calendar',
    bulkEntityPbehaviors: '/internal/bulk/entity-pbehaviors',
    patterns: '/internal/pbehavior-patterns',
    allPatterns: '/internal/all-pbehavior-patterns',
  },
  engineRunInfo: '/internal/engine-runinfo',
  cas: {
    login: '/internal/cas/login',
    loggedin: '/internal/cas/loggedin',
  },
  saml: {
    auth: '/internal/saml/auth',
  },
  oauth: {
    login: '/internal/oauth/',
  },
  scenario: {
    scenarios: '/internal/scenarios',
    testExecution: '/internal/cat/test-scenario-executions',
    testExecutionWebhooks: '/internal/cat/test-scenario-webhooks',
    bulk: '/internal/bulk/scenarios',
    bulkEnable: '/internal/bulk/scenarios/enable',
    bulkDisable: '/internal/bulk/scenarios/disable',
  },
  entityCategories: '/internal/entity-categories',
  stateSetting: '/internal/state-settings',
  dataStorage: '/internal/data-storage',
  notificationSettings: '/internal/notification-settings',
  notifications: '/internal/notifications',
  idleRules: '/internal/idle-rules',
  bulkIdleRules: '/internal/bulk/idle-rules',
  bulkIdleRulesEnable: '/internal/bulk/idle-rules/enable',
  bulkIdleRulesDisable: '/internal/bulk/idle-rules/disable',
  flappingRules: '/internal/flapping-rules',
  bulkFlappingRules: '/internal/bulk/flapping-rules',
  bulkFlappingRulesEnable: '/internal/bulk/flapping-rules/enable',
  bulkFlappingRulesDisable: '/internal/bulk/flapping-rules/disable',
  resolveRules: '/internal/resolve-rules',
  bulkResolveRules: '/internal/bulk/resolve-rules',
  bulkResolveRulesEnable: '/internal/bulk/resolve-rules/enable',
  bulkResolveRulesDisable: '/internal/bulk/resolve-rules/disable',
  messageRateStats: '/internal/message-rate-stats',
  pattern: {
    list: '/internal/patterns',
    bulkList: '/internal/bulk/patterns',
    entitiesCount: '/internal/patterns-entities-count',
    alarmsCount: '/internal/patterns-alarms-count',
    entitiesOptimize: '/internal/patterns-entities-optimize',
  },
  shareTokens: '/internal/share-tokens',
  apiTokens: '/internal/api-tokens',
  techMetrics: '/internal/tech-metrics-export',
  techMetricsSettings: '/internal/tech-metrics-settings',
  templateVars: '/internal/template-vars',
  templateData: '/internal/template-data',
  templateTests: '/internal/template-test',
  templateVarsSpecific: {
    entityServices: '/internal/entityservice-template-vars',
    eventFilters: '/internal/eventfilter-template-vars',
    scenarios: '/internal/scenario-template-vars',
    linkRules: '/internal/link-rule-template-vars',
    widgets: '/internal/widget-template-vars',
  },
  templateValidation: {
    entityServices: '/internal/entityservice-template-validate',
    eventFilters: '/internal/eventfilter-template-validate',
    scenarios: '/internal/scenario-template-validate',
    linkRules: '/internal/link-rule-template-validate',
    widgets: '/internal/widget-template-validate',
  },
  copyVars: {
    eventFilters: '/internal/eventfilter-copy-vars',
  },
  linkRule: '/internal/link-rules',
  bulkLinkRule: '/internal/bulk/link-rules',
  bulkLinkRuleEnable: '/internal/bulk/link-rules/enable',
  bulkLinkRuleDisable: '/internal/bulk/link-rules/disable',
  linkCategories: '/internal/link-categories',
  icons: '/internal/icons',
  externalDataTables: '/internal/external-data-tables',
  bulkExternalDataTables: '/internal/bulk/external-data-tables',
  externalDataImport: '/internal/external-data-import',
  externalDataExport: '/internal/external-data-export',
  themes: {
    list: '/internal/color-themes',
    bulkList: '/internal/bulk/color-themes',
  },
  healthcheck: {
    engines: '/internal/healthcheck',
    status: '/internal/healthcheck/status',
    enginesOrder: '/internal/healthcheck/engines-order',
    parameters: '/internal/healthcheck/parameters',
  },
  dbExport: {
    pbehaviors: '/internal/pbehaviors-db-export',
    eventFilters: '/internal/eventfilter-db-export',
    linkRules: '/internal/link-rules-db-export',
    idleRules: '/internal/idle-rules-db-export',
    flappingRules: '/internal/flapping-rules-db-export',
    scenarios: '/internal/scenarios-db-export',

    /**
     * Cat routes for export
     */
    dynamicInfos: '/internal/cat/dynamic-infos-db-export',
    declareTickets: '/internal/cat/declare-ticket-rules-db-export',
    metaAlarmRules: '/internal/cat/metaalarmrules-db-export',
    instructions: '/internal/cat/instructions-db-export',
    jobs: '/internal/cat/jobs-db-export',
    jobConfigs: '/internal/cat/job-configs-db-export',
  },

  patternFields: {
    flappingRule: '/internal/flapping-rule-pattern-fields',
    idleRule: '/internal/idle-rule-pattern-fields',
    linkRule: '/internal/link-rule-pattern-fields',
    resolveRule: '/internal/resolve-rule-pattern-fields',
    pbehavior: '/internal/pbehavior-pattern-fields',
    alarmTag: '/internal/alarm-tag-pattern-fields',
    widgetFilter: '/internal/widget-filter-pattern-fields',
    service: '/internal/entityservice-pattern-fields',
    stateSetting: '/internal/state-setting-pattern-fields',
    eventFilter: '/internal/eventfilter-pattern-fields',
    scenario: '/internal/scenario-pattern-fields',

    /**
     * Cat routes for pattern fields
     */
    metaalarmrule: '/internal/cat/metaalarmrule-pattern-fields',
    declareTicketRule: '/internal/cat/declare-ticket-rule-pattern-fields',
    instruction: '/internal/cat/instruction-pattern-fields',
    kpiFilter: '/internal/cat/kpi-filter-pattern-fields',
    dynamicInfos: '/internal/cat/dynamic-infos-pattern-fields',
    eventRecord: '/internal/cat/event-record-pattern-fields',
  },

  /**
   * Cat routes
   */
  filters: '/internal/cat/kpi-filters',
  ratingSettings: '/internal/cat/rating-settings',
  bulkRatingSettings: '/internal/cat/rating-settings/bulk',
  dynamicInfo: '/internal/cat/dynamic-infos',
  bulkDynamicInfo: '/internal/cat/bulk/dynamic-infos',
  bulkDynamicInfoEnable: '/internal/cat/bulk/dynamic-infos/enable',
  bulkDynamicInfoDisable: '/internal/cat/bulk/dynamic-infos/disable',
  dynamicInfosDictionaryKeys: '/internal/cat/dynamic-infos-dictionary/keys',
  metaAlarmRule: '/internal/cat/metaalarmrules',
  metaAlarmRuleAlarmCount: '/internal/cat/metaalarmrule-alarm-count',
  bulkMetaAlarmRule: '/internal/cat/bulk/metaalarmrules',
  bulkMetaAlarmRuleEnable: '/internal/cat/bulk/metaalarmrules/enable',
  bulkMetaAlarmRuleDisable: '/internal/cat/bulk/metaalarmrules/disable',
  remediation: {
    instructions: '/internal/cat/instructions',
    bulkInstructions: '/internal/cat/bulk/instructions',
    bulkInstructionsEnable: '/internal/cat/bulk/instructions/enable',
    bulkInstructionsDisable: '/internal/cat/bulk/instructions/disable',
    instructionStats: '/internal/cat/instruction-stats',
    jobs: '/internal/cat/jobs',
    configurations: '/internal/cat/job-configs',
    executions: '/internal/cat/executions',
    jobExecutions: '/internal/cat/job-executions',
    alarmExecutions: '/internal/cat/alarm-executions',
    pausedExecutions: '/internal/cat/account/executions',
    instructionComments: '/internal/cat/instruction-comments',
    executionStatuses: '/internal/cat/execution-statuses',
    readExecution: '/internal/cat/executions',
  },
  junit: {
    directory: '/internal/cat/junit/directory',
    testSuites: '/internal/cat/junit/test-suites',
    entityGantt: '/internal/cat/junit/test-suites-entity-gantt',
    history: '/internal/cat/junit/test-suites-history',
    widget: '/internal/cat/junit/test-suites-widget',
    file: '/internal/cat/junit/test-cases-file',
  },
  metrics: {
    settings: '/internal/cat/metrics-settings',
    alarm: '/internal/cat/metrics/alarm',
    exportAlarm: '/internal/cat/metrics-export/alarm',
    exportRating: '/internal/cat/metrics-export/rating',
    exportSli: '/internal/cat/metrics-export/sli',
    exportAggregate: '/internal/cat/metrics-export/aggregate',
    exportRemediation: '/internal/cat/metrics-export/remediation',
    exportMetric: '/internal/cat/metrics-export',
    exportAvailabilityByEntity: '/internal/cat/metrics-export/availability-by-entity',
    exportAvailability: '/internal/cat/metrics-export/availability',
    sli: '/internal/cat/metrics/sli',
    rating: '/internal/cat/metrics/rating',
    remediation: '/internal/cat/metrics/remediation',
    aggregate: '/internal/cat/metrics/aggregate',
    perfDataMetrics: '/internal/cat/perf-data-metrics',
    entityAlarmMetrics: '/internal/cat/entity-metrics/alarm',
    entityAggregateMetrics: '/internal/cat/entity-metrics/aggregate',
    entityAvailability: '/internal/cat/entity-metrics/availability',
    entityAggregateAvailability: '/internal/cat/entity-metrics/aggregate-availability',
    group: '/internal/cat/metrics/group',
    exportGroup: '/internal/cat/metrics-export/group',
    availability: '/internal/cat/metrics/availability',
  },
  maps: '/internal/cat/maps',
  bulkMaps: '/internal/cat/bulk/maps',
  mapState: '/internal/cat/map-state',
  metaAlarm: '/internal/cat/meta-alarms',
  snmpRule: '/internal/cat/snmprules',
  bulkSnmpRule: '/internal/cat/bulk/snmprules',
  bulkSnmpRuleEnable: '/internal/cat/bulk/snmprules/enable',
  bulkSnmpRuleDisable: '/internal/cat/bulk/snmprules/disable',
  snmpMib: '/internal/cat/snmpmibs',
  webhookTokenRule: '/internal/cat/webhook-token-rules',
  anomalyMonitoredConnector: '/internal/cat/anomalies/monitored-connectors',
  anomalyMonitoredConnectorStates: '/internal/cat/anomalies/monitored-connectors/states',
  bulkAnomalyMonitoredConnector: '/internal/cat/bulk/anomalies/monitored-connectors',
  ticketStatusJobs: '/internal/cat/job-management/ticket-status',
  bulkTicketStatusJobs: {
    play: '/internal/cat/bulk/job-management/ticket-status/sync',
    pause: '/internal/cat/bulk/job-management/ticket-status/pause',
  },
  remediationJobs: '/internal/cat/job-management/remediation',
  bulkRemediationJobs: {
    stop: '/internal/cat/bulk/job-management/remediation/stop',
  },
  webhookJobs: '/internal/cat/job-management/webhook',
  bulkWebhookJobs: {
    stop: '/internal/cat/bulk/job-management/webhook/stop',
  },
  webhookTokenJobs: '/internal/cat/job-management/webhook-token',
  bulkWebhookTokenJobs: {
    stop: '/internal/cat/bulk/job-management/webhook-token/stop',
  },
  declareTicket: {
    rules: '/internal/cat/declare-ticket-rules',
    bulkRules: '/internal/cat/bulk/declare-ticket-rules',
    bulkRulesEnable: '/internal/cat/bulk/declare-ticket-rules/enable',
    bulkRulesDisable: '/internal/cat/bulk/declare-ticket-rules/disable',
    alarmsAssigned: '/internal/cat/declare-ticket-assigned',
    testExecution: '/internal/cat/test-declare-ticket-executions',
    testExecutionWebhooks: '/internal/cat/test-declare-ticket-webhooks',
    declareTicketExecution: '/internal/cat/declare-ticket-executions',
    bulkDeclareTicket: '/internal/cat/bulk/declare-ticket-executions',
  },
  tags: '/internal/cat/tags',
  privateView: {
    groups: '/internal/cat/private-view-groups',
  },
  eventsRecord: {
    list: '/internal/cat/event-records',
    export: '/internal/cat/event-records/exports',
    current: '/internal/cat/event-records-current',
    event: '/internal/cat/event-records/events',
    bulkEvent: '/internal/cat/event-records/events/bulk',
  },
  entityInfosLog: '/internal/cat/entity-infos-log',
  templateVarsSpecificCat: {
    declareTicketRules: '/internal/cat/declare-ticket-rule-template-vars',
    dynamicInfos: '/internal/cat/dynamic-infos-template-vars',
    instructions: '/internal/cat/instruction-template-vars',
    jobs: '/internal/cat/job-template-vars',
    metaAlarmRules: '/internal/cat/metaalarmrule-template-vars',
    webhookTokenRules: '/internal/cat/webhook-token-rule-template-vars',
  },
  templateValidationCat: {
    declareTicketRules: '/internal/cat/declare-ticket-rule-template-validate',
    dynamicInfos: '/internal/cat/dynamic-infos-template-validate',
    instructions: '/internal/cat/instruction-template-validate',
    jobs: '/internal/cat/job-template-validate',
    metaAlarmRules: '/internal/cat/metaalarmrule-template-validate',
    webhookTokenRules: '/internal/cat/webhook-token-rule-template-validate',
  },
  copyVarsCat: {
    dynamicInfos: '/internal/cat/dynamic-infos-copy-vars',
  },
  llms: {
    list: '/internal/cat/llm-configs',
    bulk: '/internal/cat/bulk/llm-configs',
    bulkEnable: '/internal/cat/bulk/llm-configs/enable',
    bulkDisable: '/internal/cat/bulk/llm-configs/disable',
    bulkHistoryLink: '/internal/cat/bulk/llm-history/link',
    models: '/internal/cat/llm-config-models',
  },
};

export const COLORS = {
  primary: '#2DAC61',
  secondary: '#2B3E4F',
  accent: '#678CCA',
  error: '#D22E2E',
  errorBackground: '#FFCABE',
  warning: '#FB8C00',
  warningBackground: '#FFD27A',
  success: '#22874C',
  successBackground: '#90DFB1',
  info: '#0A6EBD',
  infoBackground: '#A6D9FF',
  background: '#ffffff',
  activeColor: '#000000',
  table: {
    background: '#FFFFFF',
    rowColor: '#FFFFFF',
    shiftRowColor: '#F5F5F5',
    shiftRowDarkColor: '#616161',
    hoverRowColor: '#F5F5F5',
    hoverRowDarkColor: '#616161',
  },
  widget: {
    titleBackground: '#FFFFFF',
    contentBackground: '#FFFFFF',
    borderColor: '#C7C7C7',
  },
  state: {
    ok: '#24834C',
    minor: '#FFF176',
    major: '#FFA800',
    critical: '#C62828',
    pause: '#808080',
  },
  status: {
    closed: '#2FAB63',
    stealthy: '#6A6A6A',
    flapping: '#2B3E4F',
    ongoing: '#F9A825',
    cancelled: '#757575',
    noEvents: '#FF5252',
    unknown: '#808080',
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
    ack: '#9C27B0',
    canceled: '#607D8B',
    children: '#6D4C41',
    parent: '#6D4C41',
    comment: '#BA68C8',
    snooze: '#E91E63',
    tickets: '#2196F3',
    cancel_delay: '#808080',
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
  declareTicketRuleTicketStatusChip: {
    unknown: '#E0E0E0',
    open: '#BBDEFB',
    assigned: '#FFC107',
    inProgress: '#FFECB3',
    closed: '#4CAF50',
  },
  aiChat: {
    chipBackground: '#E0E0E0',
    chipText: '#616161',
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
  impactStateGrouped: ['#568822', '#b0e37b', '#fff176', '#ffb7b8', '#c50711'],
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
      skipped: '#BDBDBD',
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
  background: 'var(--v-application-background-base)',
  primary: 'var(--v-primary-base)',
  secondary: 'var(--v-secondary-base)',
  error: 'var(--v-error-base)',
  success: 'var(--v-success-base)',

  state: {
    ok: 'var(--v-state-ok-base)',
    minor: 'var(--v-state-minor-base)',
    major: 'var(--v-state-major-base)',
    critical: 'var(--v-state-critical-base)',
    pause: 'var(--v-state-pause-base)',
  },

  status: {
    closed: 'var(--v-success-base)',
    stealthy: 'var(--v-status-stealthy-base)',
    flapping: 'var(--v-status-flapping-base)',
    ongoing: 'var(--v-status-ongoing-base)',
    cancelled: 'var(--v-status-cancelled-base)',
    noEvents: 'var(--v-error-base)',
    unknown: 'var(--v-status-unknown-base)',
  },
};

export const DEFAULT_THEME_COLORS = {
  primary: COLORS.primary,
  secondary: COLORS.secondary,
  accent: COLORS.accent,
  error: COLORS.error,
  error_background: COLORS.errorBackground,
  warning: COLORS.warning,
  warning_background: COLORS.warningBackground,
  success: COLORS.success,
  success_background: COLORS.successBackground,
  info: COLORS.info,
  info_background: COLORS.infoBackground,
  background: COLORS.background,
  active_color: COLORS.activeColor,

  state: COLORS.state,
  status: COLORS.status,
  icons: COLORS.alarmExtraDetails,
  table: {
    background: COLORS.table.background,
    row_color: COLORS.table.rowColor,
    hover_row_color: COLORS.table.hoverRowColor,
  },
  widget: {
    title_background: COLORS.widget.titleBackground,
    content_background: COLORS.widget.contentBackground,
    border_color: COLORS.widget.borderColor,
  },
};

export const FILE_BASE_URL = `${API_HOST}${API_ROUTES.file}`;

export const DOCUMENTATION_BASE_URL = 'https://doc.canopsis.net/';

export const EXPORT_FETCHING_INTERVAL = 2000;

export const IMPORT_FETCHING_INTERVAL = 500;

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
