export const HEALTHCHECK_SERVICES_NAMES = {
  mongo: 'MongoDB',
  redis: 'Redis',
  rabbit: 'RabbitMQ',
  timescaleDB: 'TimescaleDB',
  events: 'Events',
  api: 'API',
  healthcheck: 'healthcheck',
  enginesChain: 'engines-chain',
};

export const HEALTHCHECK_SERVICES_RENDERED_POSITIONS_DIFF_FACTORS = {
  [HEALTHCHECK_SERVICES_NAMES.events]: { x: -1, y: 2 },
  [HEALTHCHECK_SERVICES_NAMES.mongo]: { x: -3, y: 0 },
  [HEALTHCHECK_SERVICES_NAMES.timescaleDB]: { x: -3, y: 1 },
  [HEALTHCHECK_SERVICES_NAMES.api]: { x: -2, y: 0 },
  [HEALTHCHECK_SERVICES_NAMES.rabbit]: { x: -1, y: 1 },
  [HEALTHCHECK_SERVICES_NAMES.healthcheck]: { x: -2, y: -0.5 },
  [HEALTHCHECK_SERVICES_NAMES.redis]: { x: -1, y: 0 },
  [HEALTHCHECK_SERVICES_NAMES.enginesChain]: { x: 0, y: -0.5 },
};

export const HEALTHCHECK_ENGINES_NAMES = {
  webhook: 'engine-webhook',
  fifo: 'engine-fifo',
  axe: 'engine-axe',
  che: 'engine-che',
  pbehavior: 'engine-pbehavior',
  action: 'engine-action',
  service: 'engine-service',
  dynamicInfos: 'engine-dynamic-infos',
  correlation: 'engine-correlation',
  remediation: 'engine-remediation',
};

export const HEALTHCHECK_ENGINES_REFERENCE_EDGES = [
  {
    from: HEALTHCHECK_ENGINES_NAMES.fifo,
    to: HEALTHCHECK_ENGINES_NAMES.che,
  },
  {
    from: HEALTHCHECK_ENGINES_NAMES.che,
    to: HEALTHCHECK_ENGINES_NAMES.pbehavior,
  },
  {
    from: HEALTHCHECK_ENGINES_NAMES.pbehavior,
    to: HEALTHCHECK_ENGINES_NAMES.axe,
  },
  {
    from: HEALTHCHECK_ENGINES_NAMES.axe,
    to: HEALTHCHECK_ENGINES_NAMES.remediation,
  },
  {
    from: HEALTHCHECK_ENGINES_NAMES.axe,
    to: HEALTHCHECK_ENGINES_NAMES.service,
  },
  {
    from: HEALTHCHECK_ENGINES_NAMES.service,
    to: HEALTHCHECK_ENGINES_NAMES.action,
  },
];

export const HEALTHCHECK_ENGINES_PRO_REFERENCE_EDGES = [
  {
    from: HEALTHCHECK_ENGINES_NAMES.fifo,
    to: HEALTHCHECK_ENGINES_NAMES.che,
  },
  {
    from: HEALTHCHECK_ENGINES_NAMES.che,
    to: HEALTHCHECK_ENGINES_NAMES.pbehavior,
  },
  {
    from: HEALTHCHECK_ENGINES_NAMES.pbehavior,
    to: HEALTHCHECK_ENGINES_NAMES.axe,
  },
  {
    from: HEALTHCHECK_ENGINES_NAMES.axe,
    to: HEALTHCHECK_ENGINES_NAMES.remediation,
  },
  {
    from: HEALTHCHECK_ENGINES_NAMES.axe,
    to: HEALTHCHECK_ENGINES_NAMES.correlation,
  },
  {
    from: HEALTHCHECK_ENGINES_NAMES.correlation,
    to: HEALTHCHECK_ENGINES_NAMES.service,
  },
  {
    from: HEALTHCHECK_ENGINES_NAMES.service,
    to: HEALTHCHECK_ENGINES_NAMES.dynamicInfos,
  },
  {
    from: HEALTHCHECK_ENGINES_NAMES.dynamicInfos,
    to: HEALTHCHECK_ENGINES_NAMES.action,
  },
  {
    from: HEALTHCHECK_ENGINES_NAMES.action,
    to: HEALTHCHECK_ENGINES_NAMES.webhook,
  },
];

export const ENGINES_QUEUE_NAMES = {
  webhook: 'Engine_webhook',
  fifo: 'Engine_fifo',
  axe: 'Engine_axe',
  che: 'Engine_che',
  pbehavior: 'Engine_pbehavior',
  action: 'Engine_action',
  service: 'Engine_service',
  dynamicInfo: 'Engine_dynamic_infos',
  correlation: 'Engine_correlation',
};

export const ENGINES_NAMES_TO_QUEUE_NAMES = {
  [ENGINES_QUEUE_NAMES.webhook]: HEALTHCHECK_ENGINES_NAMES.webhook,
  [ENGINES_QUEUE_NAMES.fifo]: HEALTHCHECK_ENGINES_NAMES.fifo,
  [ENGINES_QUEUE_NAMES.axe]: HEALTHCHECK_ENGINES_NAMES.axe,
  [ENGINES_QUEUE_NAMES.che]: HEALTHCHECK_ENGINES_NAMES.che,
  [ENGINES_QUEUE_NAMES.pbehavior]: HEALTHCHECK_ENGINES_NAMES.pbehavior,
  [ENGINES_QUEUE_NAMES.action]: HEALTHCHECK_ENGINES_NAMES.action,
  [ENGINES_QUEUE_NAMES.service]: HEALTHCHECK_ENGINES_NAMES.service,
  [ENGINES_QUEUE_NAMES.dynamicInfo]: HEALTHCHECK_ENGINES_NAMES.dynamicInfos,
  [ENGINES_QUEUE_NAMES.correlation]: HEALTHCHECK_ENGINES_NAMES.correlation,
};

export const PRO_ENGINES = [
  HEALTHCHECK_ENGINES_NAMES.correlation,
  HEALTHCHECK_ENGINES_NAMES.dynamicInfos,
  HEALTHCHECK_ENGINES_NAMES.webhook,
];

export const HEALTHCHECK_NETWORK_GRAPH_OPTIONS = {
  nodeSpace: 110, // Magic variable. Was calculated by imperative method
  spacingFactor: 2,
  fitPadding: 15,
  wheelSensitivity: 0.5,
  minZoom: 0.05,
  maxZoom: 1.5,
  nodeSize: 60,
};

export const MESSAGE_STATS_INTERVALS = {
  hour: 'hour',
  minute: 'minute',
};

export const HEALTHCHECK_HISTORY_GRAPH_RECEIVED_FACTOR = 1.2;

export const TECH_METRICS_EXPORT_STATUSES = {
  none: 0,
  running: 1,
  success: 2,
  failed: 3,
  disabled: 4,
};
