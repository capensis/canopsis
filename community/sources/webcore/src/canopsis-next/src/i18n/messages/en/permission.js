import { USER_PERMISSIONS_PREFIXES, USERS_PERMISSIONS } from '@/constants';

export default {
  technical: {
    admin: 'Admin rights',
    exploitation: 'Exploitation rights',
    notification: 'Notification rights',
    profile: 'Profile rights',
  },
  business: {
    [USER_PERMISSIONS_PREFIXES.business.common]: 'Rights for common',
    [USER_PERMISSIONS_PREFIXES.business.alarmsList]: 'Rights for Alarms List',
    [USER_PERMISSIONS_PREFIXES.business.context]: 'Rights for Context Explorer',
    [USER_PERMISSIONS_PREFIXES.business.serviceWeather]: 'Rights for Service Weather',
    [USER_PERMISSIONS_PREFIXES.business.counter]: 'Rights for Counter',
    [USER_PERMISSIONS_PREFIXES.business.testingWeather]: 'Rights for Testing Weather',
    [USER_PERMISSIONS_PREFIXES.business.map]: 'Rights for Mapping',
    [USER_PERMISSIONS_PREFIXES.business.barChart]: 'Rights for Bar Chart',
    [USER_PERMISSIONS_PREFIXES.business.lineChart]: 'Rights for Line Chart',
    [USER_PERMISSIONS_PREFIXES.business.pieChart]: 'Rights for Pie Chart',
    [USER_PERMISSIONS_PREFIXES.business.numbers]: 'Rights for Numbers',
  },
  api: {
    general: 'General',
    rules: 'Rules',
    remediation: 'Remediation',
    pbehavior: 'PBehavior',
  },
  permissions: {
    /**
     * Technical General Permissions
     */
    [USERS_PERMISSIONS.technical.view]: {
      name: 'Views',
      description: '',
    },
    [USERS_PERMISSIONS.technical.role]: {
      name: 'Roles',
      description: '',
    },
    [USERS_PERMISSIONS.technical.permission]: {
      name: 'Permission',
      description: '',
    },
    [USERS_PERMISSIONS.technical.user]: {
      name: 'Users',
      description: '',
    },
    [USERS_PERMISSIONS.technical.parameters]: {
      name: 'Parameters',
      description: '',
    },
    [USERS_PERMISSIONS.technical.broadcastMessage]: {
      name: 'Broadcast messages',
      description: '',
    },
    [USERS_PERMISSIONS.technical.playlist]: {
      name: 'Playlists',
      description: '',
    },
    [USERS_PERMISSIONS.technical.planning]: {
      name: 'Planning',
      description: '',
    },
    [USERS_PERMISSIONS.technical.planningType]: {
      name: 'Planning types',
      description: '',
    },
    [USERS_PERMISSIONS.technical.planningReason]: {
      name: 'Planning reasons',
      description: '',
    },
    [USERS_PERMISSIONS.technical.planningExceptions]: {
      name: 'Planning exceptions',
      description: '',
    },
    [USERS_PERMISSIONS.technical.remediation]: {
      name: 'Remediation',
      description: '',
    },
    [USERS_PERMISSIONS.technical.remediationInstruction]: {
      name: 'Remediation instructions',
      description: '',
    },
    [USERS_PERMISSIONS.technical.remediationJob]: {
      name: 'Remediation jobs',
      description: '',
    },
    [USERS_PERMISSIONS.technical.remediationConfiguration]: {
      name: 'Remediation configurations',
      description: '',
    },
    [USERS_PERMISSIONS.technical.remediationStatistic]: {
      name: 'Remediation statistics',
      description: '',
    },
    [USERS_PERMISSIONS.technical.healthcheck]: {
      name: 'Healthcheck',
      description: '',
    },
    [USERS_PERMISSIONS.technical.techmetrics]: {
      name: 'Tech metrics',
      description: '',
    },
    [USERS_PERMISSIONS.technical.engine]: {
      name: 'Engines',
      description: '',
    },
    [USERS_PERMISSIONS.technical.healthcheckStatus]: {
      name: 'Healthcheck status',
      description: '',
    },
    [USERS_PERMISSIONS.technical.kpi]: {
      name: 'KPI',
      description: '',
    },
    [USERS_PERMISSIONS.technical.kpiFilters]: {
      name: 'KPI filters',
      description: '',
    },
    [USERS_PERMISSIONS.technical.kpiRatingSettings]: {
      name: 'KPI rating settings',
      description: '',
    },
    [USERS_PERMISSIONS.technical.kpiCollectionSettings]: {
      name: 'KPI collection settings',
      description: '',
    },
    [USERS_PERMISSIONS.technical.map]: {
      name: 'Maps',
      description: '',
    },
    [USERS_PERMISSIONS.technical.shareToken]: {
      name: 'Share tokens',
      description: '',
    },
    [USERS_PERMISSIONS.technical.widgetTemplate]: {
      name: 'Widget templates',
      description: '',
    },

    /**
     * Technical Exploitation Permissions
     */
    [USERS_PERMISSIONS.technical.exploitation.eventFilter]: {
      name: 'Event filters',
      description: '',
    },
    [USERS_PERMISSIONS.technical.exploitation.pbehavior]: {
      name: 'PBehaviors',
      description: '',
    },
    [USERS_PERMISSIONS.technical.exploitation.snmpRule]: {
      name: 'SNMP Rules',
      description: '',
    },
    [USERS_PERMISSIONS.technical.exploitation.dynamicInfo]: {
      name: 'Dynamic info rules',
      description: '',
    },
    [USERS_PERMISSIONS.technical.exploitation.metaAlarmRule]: {
      name: 'Meta alarm rules',
      description: '',
    },
    [USERS_PERMISSIONS.technical.exploitation.scenario]: {
      name: 'Scenarios',
      description: '',
    },
    [USERS_PERMISSIONS.technical.exploitation.idleRules]: {
      name: 'IDLE Rules',
      description: '',
    },
    [USERS_PERMISSIONS.technical.exploitation.flappingRules]: {
      name: 'Flapping rules',
      description: '',
    },
    [USERS_PERMISSIONS.technical.exploitation.resolveRules]: {
      name: 'Resolve rules',
      description: '',
    },
    [USERS_PERMISSIONS.technical.exploitation.declareTicketRule]: {
      name: 'Declare ticket rules',
      description: '',
    },
    [USERS_PERMISSIONS.technical.exploitation.linkRule]: {
      name: 'Link rules',
      description: '',
    },

    /**
     * Technical Notification Permissions
     */
    [USERS_PERMISSIONS.technical.notification.instructionStats]: {
      name: 'Instruction statistics',
      description: '',
    },

    /**
     * Technical Notification Permissions
     */
    [USERS_PERMISSIONS.technical.notification.instructionStats]: {
      name: 'Instruction statistics',
      description: '',
    },

    /**
     * Technical Profile Permissions
     */
    [USERS_PERMISSIONS.technical.profile.corporatePattern]: {
      name: 'Corporate patterns',
      description: '',
    },

    /**
     * API Permissions
     */
    [USERS_PERMISSIONS.api.general.acl]: {
      name: 'ACL',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.appInfoRead]: {
      name: 'Read app information',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.alarmRead]: {
      name: 'Read alarms',
      description: 'Users with this right can see the list of alarms',
    },
    [USERS_PERMISSIONS.api.general.alarmUpdate]: {
      name: 'Update alarms',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.entity]: {
      name: 'Entities',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.entityservice]: {
      name: 'Services',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.entitycategory]: {
      name: 'Entity categories',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.event]: {
      name: 'Event',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.view]: {
      name: 'Views',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.viewgroup]: {
      name: 'View groups',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.userInterfaceUpdate]: {
      name: 'Update user interface',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.userInterfaceDelete]: {
      name: 'Delete user interface',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.datastorageRead]: {
      name: 'Read data storage',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.datastorageUpdate]: {
      name: 'Update data storage',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.associativeTable]: {
      name: 'Associative table',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.stateSettings]: {
      name: 'State settings',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.files]: {
      name: 'File',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.healthcheck]: {
      name: 'Healthcheck',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.techmetrics]: {
      name: 'Tech metrics',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.contextgraph]: {
      name: 'Context graph',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.broadcastMessage]: {
      name: 'Broadcast message',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.junit]: {
      name: 'JUnit',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.notifications]: {
      name: 'Notifications',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.metrics]: {
      name: 'Metrics',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.metricsSettings]: {
      name: 'Metrics settings',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.ratingSettings]: {
      name: 'Rating settings',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.filter]: {
      name: 'Filters',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.corporatePattern]: {
      name: 'Corporate patters',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.exportConfigurations]: {
      name: 'Export configurations',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.map]: {
      name: 'Maps',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.shareToken]: {
      name: 'Share token',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.declareTicketExecution]: {
      name: 'Run ticket declaration',
      description: '',
    },
    [USERS_PERMISSIONS.api.general.widgetTemplate]: {
      name: 'Widget templates',
      description: '',
    },

    [USERS_PERMISSIONS.api.rules.action]: {
      name: 'Actions',
      description: '',
    },
    [USERS_PERMISSIONS.api.rules.dynamicinfos]: {
      name: 'Dynamic infos',
      description: '',
    },
    [USERS_PERMISSIONS.api.rules.eventFilter]: {
      name: 'Event filter',
      description: '',
    },
    [USERS_PERMISSIONS.api.rules.idleRule]: {
      name: 'Idle rules',
      description: '',
    },
    [USERS_PERMISSIONS.api.rules.metaalarmrule]: {
      name: 'Meta alarm rules',
      description: '',
    },
    [USERS_PERMISSIONS.api.rules.playlist]: {
      name: 'Playlists',
      description: '',
    },
    [USERS_PERMISSIONS.api.rules.flappingRule]: {
      name: 'Flapping rules',
      description: '',
    },
    [USERS_PERMISSIONS.api.rules.resolveRule]: {
      name: 'Resolve rules',
      description: '',
    },
    [USERS_PERMISSIONS.api.rules.snmpRule]: {
      name: 'SNMP rules',
      description: '',
    },
    [USERS_PERMISSIONS.api.rules.snmpMib]: {
      name: 'SNMP mib',
      description: '',
    },
    [USERS_PERMISSIONS.api.rules.declareTicketRule]: {
      name: 'Ticket declaration rules',
      description: '',
    },
    [USERS_PERMISSIONS.api.rules.linkRule]: {
      name: 'Link rules',
      description: '',
    },

    [USERS_PERMISSIONS.api.remediation.instruction]: {
      name: 'Instructions',
      description: '',
    },
    [USERS_PERMISSIONS.api.remediation.jobConfig]: {
      name: 'Job configs',
      description: '',
    },
    [USERS_PERMISSIONS.api.remediation.job]: {
      name: 'Jobs',
      description: '',
    },
    [USERS_PERMISSIONS.api.remediation.execution]: {
      name: 'Run instructions',
      description: '',
    },
    [USERS_PERMISSIONS.api.remediation.instructionApprove]: {
      name: 'Approve instruction',
      description: '',
    },
    [USERS_PERMISSIONS.api.remediation.messageRateStatsRead]: {
      name: 'Read message rate stats',
      description: '',
    },

    [USERS_PERMISSIONS.api.pbehavior.pbehavior]: {
      name: 'PBehaviors',
      description: '',
    },
    [USERS_PERMISSIONS.api.pbehavior.pbehaviorException]: {
      name: 'PBehavior exceptions',
      description: '',
    },
    [USERS_PERMISSIONS.api.pbehavior.pbehaviorReason]: {
      name: 'PBehavior reasons',
      description: '',
    },
    [USERS_PERMISSIONS.api.pbehavior.pbehaviorType]: {
      name: 'PBehavior types',
      description: '',
    },
  },
};
