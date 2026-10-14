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
};

export const HEALTHCHECK_ENGINES_NAMES = {
  snmp: 'engine-snmp',
  webhook: 'engine-webhook',
  fifo: 'engine-fifo',
  axe: 'engine-axe',
  che: 'engine-che',
  pbehavior: 'engine-pbehavior',
  action: 'engine-action',
  dynamicInfos: 'engine-dynamic-infos',
  correlation: 'engine-correlation',
  remediation: 'engine-remediation',
};

export const PRO_ENGINES = [
  HEALTHCHECK_ENGINES_NAMES.correlation,
  HEALTHCHECK_ENGINES_NAMES.dynamicInfos,
  HEALTHCHECK_ENGINES_NAMES.webhook,
  HEALTHCHECK_ENGINES_NAMES.remediation,
  HEALTHCHECK_ENGINES_NAMES.snmp,
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
