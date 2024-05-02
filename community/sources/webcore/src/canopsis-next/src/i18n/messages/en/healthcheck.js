import { HEALTHCHECK_ENGINES_NAMES, HEALTHCHECK_SERVICES_NAMES } from '@/constants';

export default {
  metricsUnavailable: 'Metrics are not collecting',
  notRunning: '{name} is unavailable',
  queueOverflow: 'Queue overflow',
  lackOfInstances: 'Lack of instances',
  diffInstancesConfig: 'Invalid instances configuration',
  queueLength: 'Queue length {queueLength}/{maxQueueLength}',
  instancesCount: 'Instances {instances}/{minInstances}',
  activeInstances: 'Only {instances} is active out of {minInstances}. The optimal number of instances is {optimalInstances}.',
  queueOverflowed: 'Queue is overflowed: {queueLength} messages out of {maxQueueLength}.\nPlease check the instances.',
  engineDown: '{name} is down, the system is not operational.\nPlease check the log or restart the service.',
  engineDownOrSlow: '{name} is down or responds too slow, the system is not operational.\nPlease check the log or restart the instance.',
  timescaleDown: '{name} is down, metrics and KPIs are not collecting.\nPlease check the log or restart the instance.',
  invalidEnginesOrder: 'Invalid engines configuration',
  invalidInstancesConfiguration: 'Invalid instances configuration: engine instances read or write to different queues.\nPlease check the instances.',
  chainConfigurationInvalid: 'Engines chain configuration is invalid.\nRefer below for the correct sequence of engines:',
  queueLimit: 'Queue length limit',
  defineQueueLimit: 'Define the engines queue length limit',
  notifyUsersQueueLimit: 'Users can be notified when the queue length limit is exceeded',
  messagesLimit: 'Messages processing limit',
  defineMessageLimit: 'Define the FIFO messages processing limit (per minute)',
  notifyUsersMessagesLimit: 'Users can be notified when the FIFO messages processing limit is exceeded',
  numberOfInstances: 'Number of instances',
  notifyUsersNumberOfInstances: 'Users can be notified when the number of active instances is less than the minimal value. The optimal number of instances is shown when the engine state is unavailable.',
  messagesHistory: 'FIFO messages processing history',
  messagesLastHour: 'FIFO messages processing for the last hour',
  messagesPerHour: 'messages/hour',
  messagesPerMinute: 'messages/minute',
  unknown: 'This system state is unavailable',
  systemStatusChipError: 'The system is not operational',
  systemStatusServerError: 'System configuration is invalid, please contact the administrator',
  systemsOperational: 'All systems are operational',
  validation: {
    max_value: 'The field must be equal or less than the optimal instance count',
    min_value: 'The field must be equal or more than the minimal instance count',
  },
  nodes: {
    [HEALTHCHECK_SERVICES_NAMES.mongo]: {
      name: 'MongoDB',
      edgeLabel: 'Status check',
    },

    [HEALTHCHECK_SERVICES_NAMES.rabbit]: {
      name: 'RabbitMQ',
      edgeLabel: 'Status check',
    },

    [HEALTHCHECK_SERVICES_NAMES.redis]: {
      name: 'Redis',
      edgeLabel: 'FIFO data\nRedis check',
    },

    [HEALTHCHECK_SERVICES_NAMES.events]: {
      name: 'Events',
    },

    [HEALTHCHECK_SERVICES_NAMES.api]: {
      name: 'Canopsis API',
    },

    [HEALTHCHECK_SERVICES_NAMES.enginesChain]: {
      name: 'Engines chain',
    },

    [HEALTHCHECK_SERVICES_NAMES.healthcheck]: {
      name: 'Healthcheck',
    },

    [HEALTHCHECK_ENGINES_NAMES.snmp]: {
      name: 'SNMP',
    },

    [HEALTHCHECK_ENGINES_NAMES.webhook]: {
      name: 'Webhook',
      description: 'Triggers the webhooks launch',
    },

    [HEALTHCHECK_ENGINES_NAMES.fifo]: {
      name: 'FIFO',
      edgeLabel: 'RabbitMQ status\nIncomming flow KPIs',
      description: 'Manages the queue of events and alarms',
    },

    [HEALTHCHECK_ENGINES_NAMES.axe]: {
      name: 'AXE',
      description: 'Creates alarms and performs actions with them',
    },

    [HEALTHCHECK_ENGINES_NAMES.che]: {
      name: 'CHE',
      description: 'Applies eventfilters and created entities',
    },

    [HEALTHCHECK_ENGINES_NAMES.pbehavior]: {
      name: 'Pbehavior',
      description: 'Checks if the alarm is under PBehvaior',
    },

    [HEALTHCHECK_ENGINES_NAMES.action]: {
      name: 'Action',
      description: 'Triggers the actions launch',
    },

    [HEALTHCHECK_ENGINES_NAMES.dynamicInfos]: {
      name: 'Dynamic infos',
      description: 'Adds dynamic infos to alarm',
    },

    [HEALTHCHECK_ENGINES_NAMES.correlation]: {
      name: 'Correlation',
      description: 'Adds dynamic infos to alarm',
    },

    [HEALTHCHECK_ENGINES_NAMES.remediation]: {
      name: 'Remediation',
      description: 'Triggers the instructions',
    },
  },
};
