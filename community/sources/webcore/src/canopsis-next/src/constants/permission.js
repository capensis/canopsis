import featuresService from '@/services/features';

import { ALARM_LIST_ACTIONS_TYPES } from './alarm';
import { CONTEXT_ACTIONS_TYPES, WEATHER_ACTIONS_TYPES, COUNTER_ACTIONS_TYPES } from './entity';

export const CRUD_ACTIONS = {
  can: 'can',
  create: 'create',
  read: 'read',
  update: 'update',
  delete: 'delete',
};

export const USERS_PERMISSIONS_TYPES = {
  crud: 'CRUD',
  rw: 'RW',
};

export const PERMISSIONS_TYPES_TO_ACTIONS = {
  [USERS_PERMISSIONS_TYPES.crud]: [
    CRUD_ACTIONS.create,
    CRUD_ACTIONS.read,
    CRUD_ACTIONS.update,
    CRUD_ACTIONS.delete,
  ],
  [USERS_PERMISSIONS_TYPES.rw]: [
    CRUD_ACTIONS.read,
    CRUD_ACTIONS.update,
    CRUD_ACTIONS.delete,
  ],
};

export const CANOPSIS_STACK = {
  go: 'go',
  python: 'python',
};

export const CANOPSIS_EDITION = {
  community: 'community',
  pro: 'pro',
};

export const EXPLOITATION_PAGES_RULES = {
  eventFilter: { stack: CANOPSIS_STACK.go },
  snmpRule: { edition: CANOPSIS_EDITION.pro },
  dynamicInfo: { edition: CANOPSIS_EDITION.pro },
  metaAlarmRule: { stack: CANOPSIS_STACK.go, edition: CANOPSIS_EDITION.pro },
  scenario: { stack: CANOPSIS_STACK.go },
};

export const ADMIN_PAGES_RULES = {
  remediation: { stack: CANOPSIS_STACK.go, edition: CANOPSIS_EDITION.pro },
  healthcheck: { stack: CANOPSIS_STACK.go, edition: CANOPSIS_EDITION.pro },
  kpi: { stack: CANOPSIS_STACK.go, edition: CANOPSIS_EDITION.pro },
};

export const NOTIFICATIONS_PAGES_RULES = {
  instructionStats: { stack: CANOPSIS_STACK.go, edition: CANOPSIS_EDITION.pro },
};

export const USER_PERMISSIONS_PREFIXES = {
  technical: {
    admin: 'models',
    exploitation: 'models_exploitation',
    notification: 'models_notification',
    profile: 'models_profile',
  },
  business: {
    common: 'common',
    alarmsList: 'listalarm',
    context: 'crudcontext',
    serviceWeather: 'serviceweather',
    testingWeather: 'testingweather',
    counter: 'counter',
    map: 'map',
  },
  api: 'api',
};

export const USERS_PERMISSIONS = {
  technical: {
    view: `${USER_PERMISSIONS_PREFIXES.technical.admin}_userview`,
    role: `${USER_PERMISSIONS_PREFIXES.technical.admin}_role`,
    permission: `${USER_PERMISSIONS_PREFIXES.technical.admin}_permission`,
    user: `${USER_PERMISSIONS_PREFIXES.technical.admin}_user`,
    parameters: `${USER_PERMISSIONS_PREFIXES.technical.admin}_parameters`,
    broadcastMessage: `${USER_PERMISSIONS_PREFIXES.technical.admin}_broadcastMessage`,
    playlist: `${USER_PERMISSIONS_PREFIXES.technical.admin}_playlist`,
    planning: `${USER_PERMISSIONS_PREFIXES.technical.admin}_planning`,
    planningType: `${USER_PERMISSIONS_PREFIXES.technical.admin}_planningType`,
    planningReason: `${USER_PERMISSIONS_PREFIXES.technical.admin}_planningReason`,
    planningExceptions: `${USER_PERMISSIONS_PREFIXES.technical.admin}_planningExceptions`,
    remediation: `${USER_PERMISSIONS_PREFIXES.technical.admin}_remediation`,
    remediationInstruction: `${USER_PERMISSIONS_PREFIXES.technical.admin}_remediationInstruction`,
    remediationJob: `${USER_PERMISSIONS_PREFIXES.technical.admin}_remediationJob`,
    remediationConfiguration: `${USER_PERMISSIONS_PREFIXES.technical.admin}_remediationConfiguration`,
    remediationStatistic: `${USER_PERMISSIONS_PREFIXES.technical.admin}_remediationStatistic`,
    healthcheck: `${USER_PERMISSIONS_PREFIXES.technical.admin}_healthcheck`,
    techmetrics: `${USER_PERMISSIONS_PREFIXES.technical.admin}_techmetrics`,
    engine: `${USER_PERMISSIONS_PREFIXES.technical.admin}_engine`,
    healthcheckStatus: `${USER_PERMISSIONS_PREFIXES.technical.admin}_healthcheckStatus`,
    kpi: `${USER_PERMISSIONS_PREFIXES.technical.admin}_kpi`,
    kpiFilters: `${USER_PERMISSIONS_PREFIXES.technical.admin}_kpiFilters`,
    kpiRatingSettings: `${USER_PERMISSIONS_PREFIXES.technical.admin}_kpiRatingSettings`,
    map: `${USER_PERMISSIONS_PREFIXES.technical.admin}_map`,
    shareToken: `${USER_PERMISSIONS_PREFIXES.technical.admin}_shareToken`,
    widgetTemplate: `${USER_PERMISSIONS_PREFIXES.technical.admin}_widgetTemplate`,
    exploitation: {
      eventFilter: `${USER_PERMISSIONS_PREFIXES.technical.exploitation}_eventFilter`,
      pbehavior: `${USER_PERMISSIONS_PREFIXES.technical.exploitation}_pbehavior`,
      snmpRule: `${USER_PERMISSIONS_PREFIXES.technical.exploitation}_snmpRule`,
      dynamicInfo: `${USER_PERMISSIONS_PREFIXES.technical.exploitation}_dynamicInfo`,
      metaAlarmRule: `${USER_PERMISSIONS_PREFIXES.technical.exploitation}_metaAlarmRule`,
      scenario: `${USER_PERMISSIONS_PREFIXES.technical.exploitation}_scenario`,
      idleRules: `${USER_PERMISSIONS_PREFIXES.technical.exploitation}_idleRules`,
      flappingRules: `${USER_PERMISSIONS_PREFIXES.technical.exploitation}_flappingRules`,
      resolveRules: `${USER_PERMISSIONS_PREFIXES.technical.exploitation}_resolveRules`,
    },
    notification: {
      instructionStats: `${USER_PERMISSIONS_PREFIXES.technical.notification}_instructionStats`,
    },
    profile: {
      corporatePattern: `${USER_PERMISSIONS_PREFIXES.technical.profile}_corporatePattern`,
    },
  },
  business: {
    alarmsList: {
      actions: {
        ack: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_ack`,
        fastAck: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_fastAck`,
        ackRemove: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_cancelAck`,
        pbehaviorAdd: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_pbehavior`,
        snooze: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_snoozeAlarm`,
        declareTicket: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_declareanIncident`,
        associateTicket: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_assignTicketNumber`,
        cancel: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_removeAlarm`,
        changeState: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_changeState`,
        history: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_history`,
        groupRequest: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_groupRequest`,
        manualMetaAlarmGroup: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_manualMetaAlarmGroup`,
        comment: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_comment`,

        listFilters: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_listFilters`,
        editFilter: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_editFilter`,
        addFilter: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_addFilter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_userFilter`,

        listRemediationInstructionsFilters:
          `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_listRemediationInstructionsFilters`,
        editRemediationInstructionsFilter:
          `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_editRemediationInstructionsFilter`,
        addRemediationInstructionsFilter:
          `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_addRemediationInstructionsFilter`,
        userRemediationInstructionsFilter:
          `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_userRemediationInstructionsFilter`,

        links: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_links`,

        correlation: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_correlation`,

        executeInstruction: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_executeInstruction`,

        category: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_category`,

        variablesHelp: `${USER_PERMISSIONS_PREFIXES.business.common}_variablesHelp`,

        exportAsCsv: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_exportAsCsv`,

        // TODO: USERS_RIGHTS Should be removed after create MR in all features
        ...featuresService.get('constants.USERS_RIGHTS.business.alarmsList.actions'),
        ...featuresService.get('constants.USERS_PERMISSIONS.business.alarmsList.actions'),
      },
    },
    context: {
      actions: {
        createEntity: `${USER_PERMISSIONS_PREFIXES.business.context}_createEntity`,
        editEntity: `${USER_PERMISSIONS_PREFIXES.business.context}_edit`,
        duplicateEntity: `${USER_PERMISSIONS_PREFIXES.business.context}_duplicate`,
        deleteEntity: `${USER_PERMISSIONS_PREFIXES.business.context}_delete`,
        pbehaviorAdd: `${USER_PERMISSIONS_PREFIXES.business.context}_pbehavior`,
        pbehaviorList: `${USER_PERMISSIONS_PREFIXES.business.context}_listPbehavior`,
        pbehaviorDelete: `${USER_PERMISSIONS_PREFIXES.business.context}_deletePbehavior`,
        massEnable: `${USER_PERMISSIONS_PREFIXES.business.context}_massEnable`,
        massDisable: `${USER_PERMISSIONS_PREFIXES.business.context}_massDisable`,

        listFilters: `${USER_PERMISSIONS_PREFIXES.business.context}_listFilters`,
        editFilter: `${USER_PERMISSIONS_PREFIXES.business.context}_editFilter`,
        addFilter: `${USER_PERMISSIONS_PREFIXES.business.context}_addFilter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.context}_userFilter`,

        category: `${USER_PERMISSIONS_PREFIXES.business.context}_category`,

        exportAsCsv: `${USER_PERMISSIONS_PREFIXES.business.context}_exportAsCsv`,
      },
    },
    serviceWeather: {
      actions: {
        entityAck: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_entityAck`,
        entityAssocTicket: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_entityAssocTicket`,
        entityComment: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_entityComment`,
        entityValidate: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_entityValidate`,
        entityInvalidate: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_entityInvalidate`,
        entityPause: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_entityPause`,
        entityPlay: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_entityPlay`,
        entityCancel: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_entityCancel`,
        entityManagePbehaviors: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_entityManagePbehaviors`,
        executeInstruction: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_executeInstruction`,

        entityLinks: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_entityLinks`,

        moreInfos: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_moreInfos`,
        alarmsList: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_alarmsList`,
        pbehaviorList: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_pbehaviorList`,

        listFilters: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_listFilters`,
        editFilter: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_editFilter`,
        addFilter: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_addFilter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_userFilter`,

        category: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_category`,
      },
    },
    counter: {
      actions: {
        alarmsList: `${USER_PERMISSIONS_PREFIXES.business.counter}_alarmsList`,

        variablesHelp: `${USER_PERMISSIONS_PREFIXES.business.common}_variablesHelp`,
      },
    },
    testingWeather: {
      actions: {
        alarmsList: `${USER_PERMISSIONS_PREFIXES.business.testingWeather}_alarmsList`,
      },
    },
    map: {
      actions: {
        alarmsList: `${USER_PERMISSIONS_PREFIXES.business.map}_alarmsList`,

        listFilters: `${USER_PERMISSIONS_PREFIXES.business.map}_listFilters`,
        editFilter: `${USER_PERMISSIONS_PREFIXES.business.map}_editFilter`,
        addFilter: `${USER_PERMISSIONS_PREFIXES.business.map}_addFilter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.map}_userFilter`,

        category: `${USER_PERMISSIONS_PREFIXES.business.map}_category`,
      },
    },
  },
  api: {
    general: {
      acl: `${USER_PERMISSIONS_PREFIXES.api}_acl`,
      appInfoRead: `${USER_PERMISSIONS_PREFIXES.api}_app_info_read`,
      alarmRead: `${USER_PERMISSIONS_PREFIXES.api}_alarm_read`,
      entity: `${USER_PERMISSIONS_PREFIXES.api}_entity`,
      entityservice: `${USER_PERMISSIONS_PREFIXES.api}_entityservice`,
      entitycategory: `${USER_PERMISSIONS_PREFIXES.api}_entitycategory`,
      event: `${USER_PERMISSIONS_PREFIXES.api}_event`,
      view: `${USER_PERMISSIONS_PREFIXES.api}_view`,
      viewgroup: `${USER_PERMISSIONS_PREFIXES.api}_viewgroup`,
      userInterfaceUpdate: `${USER_PERMISSIONS_PREFIXES.api}_user_interface_update`,
      userInterfaceDelete: `${USER_PERMISSIONS_PREFIXES.api}_user_interface_delete`,
      datastorageRead: `${USER_PERMISSIONS_PREFIXES.api}_datastorage_read`,
      datastorageUpdate: `${USER_PERMISSIONS_PREFIXES.api}_datastorage_update`,
      associativeTable: `${USER_PERMISSIONS_PREFIXES.api}_associative_table`,
      stateSettings: `${USER_PERMISSIONS_PREFIXES.api}_state_settings`,
      files: `${USER_PERMISSIONS_PREFIXES.api}_file`,
      healthcheck: `${USER_PERMISSIONS_PREFIXES.api}_healthcheck`,
      techmetrics: `${USER_PERMISSIONS_PREFIXES.api}_techmetrics`,
      contextgraph: `${USER_PERMISSIONS_PREFIXES.api}_contextgraph`,
      broadcastMessage: `${USER_PERMISSIONS_PREFIXES.api}_broadcast_message`,
      junit: `${USER_PERMISSIONS_PREFIXES.api}_junit`,
      notifications: `${USER_PERMISSIONS_PREFIXES.api}_notification`,
      metrics: `${USER_PERMISSIONS_PREFIXES.api}_metrics`,
      ratingSettings: `${USER_PERMISSIONS_PREFIXES.api}_rating_settings`,
      filter: `${USER_PERMISSIONS_PREFIXES.api}_kpi_filter`,
      corporatePattern: `${USER_PERMISSIONS_PREFIXES.api}_corporate_pattern`,
      api: `${USER_PERMISSIONS_PREFIXES.api}_api`,
      exportConfigurations: `${USER_PERMISSIONS_PREFIXES.api}_export_configurations`,
      map: `${USER_PERMISSIONS_PREFIXES.api}_map`,
    },
    rules: {
      action: `${USER_PERMISSIONS_PREFIXES.api}_action`,
      dynamicinfos: `${USER_PERMISSIONS_PREFIXES.api}_dynamicinfos`,
      eventFilter: `${USER_PERMISSIONS_PREFIXES.api}_eventfilter`,
      idleRule: `${USER_PERMISSIONS_PREFIXES.api}_idlerule`,
      metaalarmrule: `${USER_PERMISSIONS_PREFIXES.api}_metaalarmrule`,
      playlist: `${USER_PERMISSIONS_PREFIXES.api}_playlist`,
      flappingRule: `${USER_PERMISSIONS_PREFIXES.api}_flapping_rule`,
      resolveRule: `${USER_PERMISSIONS_PREFIXES.api}_resolve_rule`,
      snmpRule: `${USER_PERMISSIONS_PREFIXES.api}_snmprule`,
      snmpMib: `${USER_PERMISSIONS_PREFIXES.api}_snmpmib`,
      shareToken: `${USER_PERMISSIONS_PREFIXES.api}_share_token`,
    },
    remediation: {
      instruction: `${USER_PERMISSIONS_PREFIXES.api}_instruction`,
      jobConfig: `${USER_PERMISSIONS_PREFIXES.api}_job_config`,
      job: `${USER_PERMISSIONS_PREFIXES.api}_job`,
      execution: `${USER_PERMISSIONS_PREFIXES.api}_execution`,
      instructionApprove: `${USER_PERMISSIONS_PREFIXES.api}_instruction_approve`,
      messageRateStatsRead: `${USER_PERMISSIONS_PREFIXES.api}_message_rate_stats_read`,
    },
    pbehavior: {
      pbehavior: `${USER_PERMISSIONS_PREFIXES.api}_pbehavior`,
      pbehaviorException: `${USER_PERMISSIONS_PREFIXES.api}_pbehaviorexception`,
      pbehaviorReason: `${USER_PERMISSIONS_PREFIXES.api}_pbehaviorreason`,
      pbehaviorType: `${USER_PERMISSIONS_PREFIXES.api}_pbehaviortype`,
    },
  },
};

export const NOT_COMPLETED_USER_PERMISSIONS = [
  USERS_PERMISSIONS.business.alarmsList.actions.links,
  USERS_PERMISSIONS.business.serviceWeather.actions.entityLinks,
];

export const BUSINESS_USER_PERMISSIONS_ACTIONS_MAP = {
  alarmsList: {
    [ALARM_LIST_ACTIONS_TYPES.ack]: USERS_PERMISSIONS.business.alarmsList.actions.ack,
    [ALARM_LIST_ACTIONS_TYPES.fastAck]: USERS_PERMISSIONS.business.alarmsList.actions.fastAck,
    [ALARM_LIST_ACTIONS_TYPES.ackRemove]: USERS_PERMISSIONS.business.alarmsList.actions.ackRemove,
    [ALARM_LIST_ACTIONS_TYPES.pbehaviorAdd]: USERS_PERMISSIONS.business.alarmsList.actions.pbehaviorAdd,
    [ALARM_LIST_ACTIONS_TYPES.snooze]: USERS_PERMISSIONS.business.alarmsList.actions.snooze,
    [ALARM_LIST_ACTIONS_TYPES.declareTicket]: USERS_PERMISSIONS.business.alarmsList.actions.declareTicket,
    [ALARM_LIST_ACTIONS_TYPES.associateTicket]: USERS_PERMISSIONS.business.alarmsList.actions.associateTicket,
    [ALARM_LIST_ACTIONS_TYPES.cancel]: USERS_PERMISSIONS.business.alarmsList.actions.cancel,
    [ALARM_LIST_ACTIONS_TYPES.changeState]: USERS_PERMISSIONS.business.alarmsList.actions.changeState,
    [ALARM_LIST_ACTIONS_TYPES.history]: USERS_PERMISSIONS.business.alarmsList.actions.history,
    [ALARM_LIST_ACTIONS_TYPES.variablesHelp]: USERS_PERMISSIONS.business.alarmsList.actions.variablesHelp,
    [ALARM_LIST_ACTIONS_TYPES.comment]: USERS_PERMISSIONS.business.alarmsList.actions.comment,
    [ALARM_LIST_ACTIONS_TYPES.groupRequest]: USERS_PERMISSIONS.business.alarmsList.actions.groupRequest,
    [ALARM_LIST_ACTIONS_TYPES.manualMetaAlarmGroup]:
    USERS_PERMISSIONS.business.alarmsList.actions.manualMetaAlarmGroup,
    [ALARM_LIST_ACTIONS_TYPES.manualMetaAlarmUngroup]:
    USERS_PERMISSIONS.business.alarmsList.actions.manualMetaAlarmGroup,
    [ALARM_LIST_ACTIONS_TYPES.manualMetaAlarmUpdate]:
    USERS_PERMISSIONS.business.alarmsList.actions.manualMetaAlarmGroup,

    [ALARM_LIST_ACTIONS_TYPES.links]: USERS_PERMISSIONS.business.alarmsList.actions.links,
    [ALARM_LIST_ACTIONS_TYPES.correlation]: USERS_PERMISSIONS.business.alarmsList.actions.correlation,

    [ALARM_LIST_ACTIONS_TYPES.listFilters]: USERS_PERMISSIONS.business.alarmsList.actions.listFilters,
    [ALARM_LIST_ACTIONS_TYPES.editFilter]: USERS_PERMISSIONS.business.alarmsList.actions.editFilter,
    [ALARM_LIST_ACTIONS_TYPES.addFilter]: USERS_PERMISSIONS.business.alarmsList.actions.addFilter,
    [ALARM_LIST_ACTIONS_TYPES.userFilter]: USERS_PERMISSIONS.business.alarmsList.actions.userFilter,

    [ALARM_LIST_ACTIONS_TYPES.listRemediationInstructionsFilters]:
    USERS_PERMISSIONS.business.alarmsList.actions.listRemediationInstructionsFilters,
    [ALARM_LIST_ACTIONS_TYPES.editRemediationInstructionsFilter]:
    USERS_PERMISSIONS.business.alarmsList.actions.editRemediationInstructionsFilter,
    [ALARM_LIST_ACTIONS_TYPES.addRemediationInstructionsFilter]:
    USERS_PERMISSIONS.business.alarmsList.actions.addRemediationInstructionsFilter,

    [ALARM_LIST_ACTIONS_TYPES.executeInstruction]:
    USERS_PERMISSIONS.business.alarmsList.actions.executeInstruction,

    // TODO: Should be removed after create MR in all features
    ...featuresService.get('constants.BUSINESS_USER_RIGHTS_ACTIONS_MAP.alarmsList'),
    ...featuresService.get('constants.BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.alarmsList'),
  },

  context: {
    [CONTEXT_ACTIONS_TYPES.createEntity]: USERS_PERMISSIONS.business.context.actions.createEntity,
    [CONTEXT_ACTIONS_TYPES.editEntity]: USERS_PERMISSIONS.business.context.actions.editEntity,
    [CONTEXT_ACTIONS_TYPES.duplicateEntity]: USERS_PERMISSIONS.business.context.actions.duplicateEntity,
    [CONTEXT_ACTIONS_TYPES.deleteEntity]: USERS_PERMISSIONS.business.context.actions.deleteEntity,
    [CONTEXT_ACTIONS_TYPES.pbehaviorAdd]: USERS_PERMISSIONS.business.context.actions.pbehaviorAdd,
    [CONTEXT_ACTIONS_TYPES.pbehaviorList]: USERS_PERMISSIONS.business.context.actions.pbehaviorList,
    [CONTEXT_ACTIONS_TYPES.pbehaviorDelete]: USERS_PERMISSIONS.business.context.actions.pbehaviorDelete,
    [CONTEXT_ACTIONS_TYPES.pbehaviorDelete]: USERS_PERMISSIONS.business.context.actions.pbehaviorDelete,
    [CONTEXT_ACTIONS_TYPES.massEnable]: USERS_PERMISSIONS.business.context.actions.massEnable,
    [CONTEXT_ACTIONS_TYPES.massDisable]: USERS_PERMISSIONS.business.context.actions.massDisable,

    [CONTEXT_ACTIONS_TYPES.listFilters]: USERS_PERMISSIONS.business.context.actions.listFilters,
    [CONTEXT_ACTIONS_TYPES.editFilter]: USERS_PERMISSIONS.business.context.actions.editFilter,
    [CONTEXT_ACTIONS_TYPES.addFilter]: USERS_PERMISSIONS.business.context.actions.addFilter,
  },

  weather: {
    [WEATHER_ACTIONS_TYPES.entityAck]: USERS_PERMISSIONS.business.serviceWeather.actions.entityAck,
    [WEATHER_ACTIONS_TYPES.entityAssocTicket]:
      USERS_PERMISSIONS.business.serviceWeather.actions.entityAssocTicket,
    [WEATHER_ACTIONS_TYPES.entityValidate]: USERS_PERMISSIONS.business.serviceWeather.actions.entityValidate,
    [WEATHER_ACTIONS_TYPES.entityInvalidate]:
      USERS_PERMISSIONS.business.serviceWeather.actions.entityInvalidate,
    [WEATHER_ACTIONS_TYPES.entityPause]: USERS_PERMISSIONS.business.serviceWeather.actions.entityPause,
    [WEATHER_ACTIONS_TYPES.entityPlay]: USERS_PERMISSIONS.business.serviceWeather.actions.entityPlay,
    [WEATHER_ACTIONS_TYPES.entityCancel]: USERS_PERMISSIONS.business.serviceWeather.actions.entityCancel,
    [WEATHER_ACTIONS_TYPES.executeInstruction]:
      USERS_PERMISSIONS.business.serviceWeather.actions.executeInstruction,

    [WEATHER_ACTIONS_TYPES.entityLinks]: USERS_PERMISSIONS.business.serviceWeather.actions.entityLinks,

    [WEATHER_ACTIONS_TYPES.moreInfos]: USERS_PERMISSIONS.business.serviceWeather.actions.moreInfos,
    [WEATHER_ACTIONS_TYPES.alarmsList]: USERS_PERMISSIONS.business.serviceWeather.actions.alarmsList,
    [WEATHER_ACTIONS_TYPES.pbehaviorList]: USERS_PERMISSIONS.business.serviceWeather.actions.pbehaviorList,
    [WEATHER_ACTIONS_TYPES.entityComment]: USERS_PERMISSIONS.business.serviceWeather.actions.entityComment,
  },

  counter: {
    [COUNTER_ACTIONS_TYPES.alarmsList]: USERS_PERMISSIONS.business.counter.actions.alarmsList,
    [COUNTER_ACTIONS_TYPES.variablesHelp]: USERS_PERMISSIONS.business.counter.actions.variablesHelp,
  },
};

export const USER_PERMISSIONS_TO_PAGES_RULES = {
  /**
   * Admin pages
   */
  [USERS_PERMISSIONS.technical.remediation]: ADMIN_PAGES_RULES.remediation,
  [USERS_PERMISSIONS.technical.healthcheck]: ADMIN_PAGES_RULES.healthcheck,
  [USERS_PERMISSIONS.technical.kpi]: ADMIN_PAGES_RULES.kpi,

  /**
   * Exploitation pages
   */
  [USERS_PERMISSIONS.technical.exploitation.eventFilter]: EXPLOITATION_PAGES_RULES.eventFilter,
  [USERS_PERMISSIONS.technical.exploitation.snmpRule]: EXPLOITATION_PAGES_RULES.snmpRule,
  [USERS_PERMISSIONS.technical.exploitation.dynamicInfo]: EXPLOITATION_PAGES_RULES.dynamicInfo,
  [USERS_PERMISSIONS.technical.exploitation.metaAlarmRule]: EXPLOITATION_PAGES_RULES.metaAlarmRule,
  [USERS_PERMISSIONS.technical.exploitation.scenario]: EXPLOITATION_PAGES_RULES.scenario,

  /**
   * Notifications pages
   */
  [USERS_PERMISSIONS.technical.notification.instructionStats]: NOTIFICATIONS_PAGES_RULES.instructionStats,
};

export const DOCUMENTATION_LINKS = {
  /**
   * Exploitation
   */
  [USERS_PERMISSIONS.technical.exploitation.eventFilter]: 'guide-utilisation/menu-exploitation/filtres-evenements/',
  [USERS_PERMISSIONS.technical.exploitation.pbehavior]: 'guide-utilisation/cas-d-usage/comportements_periodiques/',
  [USERS_PERMISSIONS.technical.exploitation.snmpRule]: 'interconnexions/Supervision/SNMPtrap/',
  [USERS_PERMISSIONS.technical.exploitation.idleRules]: 'guide-utilisation/menu-exploitation/regles-inactivite/',
  [USERS_PERMISSIONS.technical.exploitation.resolveRules]: 'guide-utilisation/menu-exploitation/regles-resolution/',
  [USERS_PERMISSIONS.technical.exploitation.flappingRules]: 'guide-utilisation/menu-exploitation/regles-bagot/',
  [USERS_PERMISSIONS.technical.exploitation.dynamicInfo]: 'guide-utilisation/menu-exploitation/informations-dynamiques/',
  [USERS_PERMISSIONS.technical.exploitation.metaAlarmRule]: 'guide-utilisation/menu-exploitation/regles-metaalarme/',
  [USERS_PERMISSIONS.technical.exploitation.scenario]: 'guide-utilisation/menu-exploitation/scenarios/',

  /**
   * Admin
   */
  [USERS_PERMISSIONS.technical.broadcastMessage]: 'guide-utilisation/interface/broadcast-messages/',
  [USERS_PERMISSIONS.technical.playlist]: 'guide-utilisation/interface/playlists/',
  [USERS_PERMISSIONS.technical.planning]: 'guide-administration/moteurs/moteur-pbehavior/#administration-de-la-planification',
  [USERS_PERMISSIONS.technical.remediation]: 'guide-utilisation/remediation/',

  /**
   * Notifications
   */
  // [USERS_PERMISSIONS.technical.notification.instructionStats]: '', // TODO: TBD
};
