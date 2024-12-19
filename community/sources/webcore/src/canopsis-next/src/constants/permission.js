import featuresService from '@/services/features';

import { groupedPermissionToPermission } from '@/helpers/permission';

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
  declareTicketRule: { edition: CANOPSIS_EDITION.pro },
};

export const ADMIN_PAGES_RULES = {
  remediation: { stack: CANOPSIS_STACK.go, edition: CANOPSIS_EDITION.pro },
  healthcheck: { stack: CANOPSIS_STACK.go },
  kpi: { stack: CANOPSIS_STACK.go, edition: CANOPSIS_EDITION.pro },
  tag: { stack: CANOPSIS_STACK.go, edition: CANOPSIS_EDITION.pro },
  map: { edition: CANOPSIS_EDITION.pro },
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
    barChart: 'barchart',
    lineChart: 'linechart',
    pieChart: 'piechart',
    numbers: 'numbers',
    userStatistics: 'userStatistics',
    alarmStatistics: 'alarmStatistics',
    availability: 'availability',
  },
  api: 'api',
};

export const USER_VIEWS_PERMISSIONS = {
  viewActions: 'view_actions',
  viewGeneral: 'view_general',
};

export const USER_PERMISSIONS_GROUPS = {
  commonviews: 'commonviews',
  commonviewsPlaylist: 'commonviews_playlist',

  widgets: 'widgets',

  widgetsCommon: 'widgets_common',

  widgetsAlarmstatistics: 'widgets_alarmstatistics',
  widgetsAlarmstatisticsWidgetsettings: 'widgets_alarmstatistics_widgetsettings',
  widgetsAlarmstatisticsViewsettings: 'widgets_alarmstatistics_viewsettings',

  widgetsAvailability: 'widgets_availability',
  widgetsAvailabilityActions: 'widgets_availability_actions',
  widgetsAvailabilityWidgetsettings: 'widgets_availability_widgetsettings',
  widgetsAvailabilityViewsettings: 'widgets_availability_viewsettings',

  widgetsBarchart: 'widgets_barchart',
  widgetsBarchartWidgetsettings: 'widgets_barchart_widgetsettings',
  widgetsBarchartViewsettings: 'widgets_barchart_viewsettings',

  widgetsCounter: 'widgets_counter',

  widgetsContext: 'widgets_context',
  widgetsContextViewsettings: 'widgets_context_viewsettings',
  widgetsContextEntityactions: 'widgets_context_entityactions',
  widgetsContextActions: 'widgets_context_actions',
  widgetsContextWidgetsettings: 'widgets_context_widgetsettings',

  widgetsLinechart: 'widgets_linechart',
  widgetsLinechartWidgetsettings: 'widgets_linechart_widgetsettings',
  widgetsLinechartViewsettings: 'widgets_linechart_viewsettings',

  widgetsAlarmslist: 'widgets_alarmslist',
  widgetsAlarmslistAlarmactions: 'widgets_alarmslist_alarmactions',
  widgetsAlarmslistViewsettings: 'widgets_alarmslist_viewsettings',
  widgetsAlarmslistActions: 'widgets_alarmslist_actions',
  widgetsAlarmslistWidgetsettings: 'widgets_alarmslist_widgetsettings',

  widgetsMap: 'widgets_map',
  widgetsMapViewsettings: 'widgets_map_viewsettings',
  widgetsMapWidgetsettings: 'widgets_map_widgetsettings',

  widgetsNumbers: 'widgets_numbers',
  widgetsNumbersWidgetsettings: 'widgets_numbers_widgetsettings',
  widgetsNumbersViewsettings: 'widgets_numbers_viewsettings',

  widgetsPiechart: 'widgets_piechart',
  widgetsPiechartWidgetsettings: 'widgets_piechart_widgetsettings',
  widgetsPiechartViewsettings: 'widgets_piechart_viewsettings',

  widgetsServiceweather: 'widgets_serviceweather',
  widgetsServiceweatherAlarmactions: 'widgets_serviceweather_alarmactions',
  widgetsServiceweatherViewsettings: 'widgets_serviceweather_viewsettings',
  widgetsServiceweatherWidgetsettings: 'widgets_serviceweather_widgetsettings',

  widgetsTestingweather: 'widgets_testingweather',

  widgetsUserstatistics: 'widgets_userstatistics',
  widgetsUserstatisticsWidgetsettings: 'widgets_userstatistics_widgetsettings',
  widgetsUserstatisticsViewsettings: 'widgets_userstatistics_viewsettings',

  api: 'api',

  apiGeneral: 'api_general',
  apiRules: 'api_rules',
  apiRemediation: 'api_remediation',
  apiPlanning: 'api_planning',

  technical: 'technical',

  technicalAdmin: 'technical_admin',
  technicalAdminCommunication: 'technical_admin_communication',
  technicalAdminGeneral: 'technical_admin_general',
  technicalAdminAccess: 'technical_admin_access',

  technicalExploitation: 'technical_exploitation',
  technicalNotification: 'technical_notification',
  technicalViewsandwidgets: 'technical_viewsandwidgets',
  technicalProfile: 'technical_profile',
  technicalToken: 'technical_token',
};

export const API_USER_PERMISSIONS_ROOT_GROUPS = [USER_PERMISSIONS_GROUPS.api];

export const VIEW_USER_PERMISSIONS_NAMES = {
  general: 'view_general',
  actions: 'view_actions',
};

export const USER_PERMISSIONS = {
  technical: {
    view: `${USER_PERMISSIONS_PREFIXES.technical.admin}_userview`,
    privateView: `${USER_PERMISSIONS_PREFIXES.technical.admin}_privateView`,
    role: `${USER_PERMISSIONS_PREFIXES.technical.admin}_role`,
    permission: `${USER_PERMISSIONS_PREFIXES.technical.admin}_permission`,
    user: `${USER_PERMISSIONS_PREFIXES.technical.admin}_user`,
    parameters: `${USER_PERMISSIONS_PREFIXES.technical.admin}_parameters`,
    broadcastMessage: `${USER_PERMISSIONS_PREFIXES.technical.admin}_broadcastMessage`,
    playlist: `${USER_PERMISSIONS_PREFIXES.technical.admin}_playlist`,
    planningType: `${USER_PERMISSIONS_PREFIXES.technical.admin}_planningType`,
    planningReason: `${USER_PERMISSIONS_PREFIXES.technical.admin}_planningReason`,
    planningExceptions: `${USER_PERMISSIONS_PREFIXES.technical.admin}_planningExceptions`,
    remediationInstruction: `${USER_PERMISSIONS_PREFIXES.technical.admin}_remediationInstruction`,
    remediationJob: `${USER_PERMISSIONS_PREFIXES.technical.admin}_remediationJob`,
    remediationConfiguration: `${USER_PERMISSIONS_PREFIXES.technical.admin}_remediationConfiguration`,
    remediationStatistic: `${USER_PERMISSIONS_PREFIXES.technical.admin}_remediationStatistic`,
    healthcheck: `${USER_PERMISSIONS_PREFIXES.technical.admin}_healthcheck`,
    techmetrics: `${USER_PERMISSIONS_PREFIXES.technical.admin}_techmetrics`,
    healthcheckStatus: `${USER_PERMISSIONS_PREFIXES.technical.admin}_healthcheckStatus`,
    kpi: `${USER_PERMISSIONS_PREFIXES.technical.admin}_kpi`,
    kpiFilters: `${USER_PERMISSIONS_PREFIXES.technical.admin}_kpiFilters`,
    kpiRatingSettings: `${USER_PERMISSIONS_PREFIXES.technical.admin}_kpiRatingSettings`,
    kpiCollectionSettings: `${USER_PERMISSIONS_PREFIXES.technical.admin}_kpiCollectionSettings`,
    map: `${USER_PERMISSIONS_PREFIXES.technical.admin}_map`,
    shareToken: `${USER_PERMISSIONS_PREFIXES.technical.admin}_shareToken`,
    maintenance: `${USER_PERMISSIONS_PREFIXES.technical.admin}_maintenance`,
    widgetTemplate: `${USER_PERMISSIONS_PREFIXES.technical.admin}_widgetTemplate`,
    stateSetting: `${USER_PERMISSIONS_PREFIXES.technical.admin}_stateSetting`,
    tag: `${USER_PERMISSIONS_PREFIXES.technical.admin}_tag`,
    storageSettings: `${USER_PERMISSIONS_PREFIXES.technical.admin}_storageSettings`,
    icon: `${USER_PERMISSIONS_PREFIXES.technical.admin}_icon`,
    eventsRecord: `${USER_PERMISSIONS_PREFIXES.technical.admin}_eventsRecord`,
    viewImportExport: `${USER_PERMISSIONS_PREFIXES.technical.admin}_view_import_export`,
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
      declareTicketRule: `${USER_PERMISSIONS_PREFIXES.technical.exploitation}_declareTicketRule`,
      linkRule: `${USER_PERMISSIONS_PREFIXES.technical.exploitation}_linkRule`,
    },
    notification: {
      common: USER_PERMISSIONS_PREFIXES.technical.notification,
      instructionStats: `${USER_PERMISSIONS_PREFIXES.technical.notification}_instructionStats`,
    },
    profile: {
      corporatePattern: `${USER_PERMISSIONS_PREFIXES.technical.profile}_corporatePattern`,
      theme: `${USER_PERMISSIONS_PREFIXES.technical.profile}_color_theme`,
    },
  },
  business: {
    alarmsList: {
      actions: {
        ack: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_ack`,
        fastAck: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_fastAck`,
        ackRemove: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_cancelAck`,
        pbehaviorAdd: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_pbehavior`,
        fastPbehaviorAdd: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_fastPbehavior`,
        snooze: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_snoozeAlarm`,
        declareTicket: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_declareanIncident`,
        associateTicket: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_assignTicketNumber`,
        cancel: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_removeAlarm`,
        unCancel: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_unCancel`,
        fastCancel: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_fastRemoveAlarm`,
        changeState: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_changeState`,
        history: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_history`,
        manualMetaAlarmGroup: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_manualMetaAlarmGroup`,
        metaAlarmGroup: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_metaAlarmGroup`,
        comment: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_comment`,
        exportPdf: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_exportPdf`,

        filter: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_filter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_userFilter`,

        remediationInstructionsFilter:
          `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_remediationInstructionsFilter`,
        userRemediationInstructionsFilter:
          `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_userRemediationInstructionsFilter`,

        links: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_links`,

        correlation: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_correlation`,

        executeInstruction: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_executeInstruction`,

        category: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_category`,

        exportAsCsv: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_exportAsCsv`,

        bookmark: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_bookmark`,
        filterByBookmark: `${USER_PERMISSIONS_PREFIXES.business.alarmsList}_filterByBookmark`,

        /**
         * COMMON
         */
        variablesHelp: `${USER_PERMISSIONS_PREFIXES.business.common}_variablesHelp`,

        ...featuresService.get('constants.USER_PERMISSIONS.business.alarmsList.actions'),
      },
    },
    context: {
      actions: {
        createEntity: `${USER_PERMISSIONS_PREFIXES.business.context}_createEntity`,
        editEntity: `${USER_PERMISSIONS_PREFIXES.business.context}_edit`,
        duplicateEntity: `${USER_PERMISSIONS_PREFIXES.business.context}_duplicate`,
        deleteEntity: `${USER_PERMISSIONS_PREFIXES.business.context}_delete`,
        pbehavior: `${USER_PERMISSIONS_PREFIXES.business.context}_pbehavior`,
        massEnable: `${USER_PERMISSIONS_PREFIXES.business.context}_massEnable`,
        massDisable: `${USER_PERMISSIONS_PREFIXES.business.context}_massDisable`,
        listPbehavior: `${USER_PERMISSIONS_PREFIXES.business.context}_listPbehavior`,

        filter: `${USER_PERMISSIONS_PREFIXES.business.context}_filter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.context}_userFilter`,

        category: `${USER_PERMISSIONS_PREFIXES.business.context}_category`,

        exportAsCsv: `${USER_PERMISSIONS_PREFIXES.business.context}_exportAsCsv`,

        /**
         * COMMON
         */
        entityComment: `${USER_PERMISSIONS_PREFIXES.business.common}_entityComment`,
      },
    },
    serviceWeather: {
      actions: {
        entityAck: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_entityAck`,
        entityAssocTicket: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_entityAssocTicket`,
        entityDeclareTicket: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_entityDeclareTicket`,
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

        filter: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_filter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_userFilter`,

        category: `${USER_PERMISSIONS_PREFIXES.business.serviceWeather}_category`,

        /**
         * COMMON
         */
        variablesHelp: `${USER_PERMISSIONS_PREFIXES.business.common}_variablesHelp`,

        entityCommentsList: `${USER_PERMISSIONS_PREFIXES.business.common}_entityCommentsList`,
        createEntityComment: `${USER_PERMISSIONS_PREFIXES.business.common}_createEntityComment`,
        editEntityComment: `${USER_PERMISSIONS_PREFIXES.business.common}_editEntityComment`,
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

        filter: `${USER_PERMISSIONS_PREFIXES.business.map}_filter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.map}_userFilter`,

        category: `${USER_PERMISSIONS_PREFIXES.business.map}_category`,
      },
    },
    barChart: {
      actions: {
        interval: `${USER_PERMISSIONS_PREFIXES.business.barChart}_interval`,

        sampling: `${USER_PERMISSIONS_PREFIXES.business.barChart}_sampling`,

        filter: `${USER_PERMISSIONS_PREFIXES.business.barChart}_filter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.barChart}_userFilter`,
      },
    },
    lineChart: {
      actions: {
        interval: `${USER_PERMISSIONS_PREFIXES.business.lineChart}_interval`,

        sampling: `${USER_PERMISSIONS_PREFIXES.business.lineChart}_sampling`,

        filter: `${USER_PERMISSIONS_PREFIXES.business.lineChart}_filter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.lineChart}_userFilter`,
      },
    },
    pieChart: {
      actions: {
        interval: `${USER_PERMISSIONS_PREFIXES.business.pieChart}_interval`,

        sampling: `${USER_PERMISSIONS_PREFIXES.business.pieChart}_sampling`,

        filter: `${USER_PERMISSIONS_PREFIXES.business.pieChart}_filter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.pieChart}_userFilter`,
      },
    },
    numbers: {
      actions: {
        interval: `${USER_PERMISSIONS_PREFIXES.business.numbers}_interval`,

        sampling: `${USER_PERMISSIONS_PREFIXES.business.numbers}_sampling`,

        filter: `${USER_PERMISSIONS_PREFIXES.business.numbers}_filter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.numbers}_userFilter`,
      },
    },
    userStatistics: {
      actions: {
        interval: `${USER_PERMISSIONS_PREFIXES.business.userStatistics}_interval`,

        filter: `${USER_PERMISSIONS_PREFIXES.business.userStatistics}_filter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.userStatistics}_userFilter`,
      },
    },
    alarmStatistics: {
      actions: {
        interval: `${USER_PERMISSIONS_PREFIXES.business.alarmStatistics}_interval`,

        filter: `${USER_PERMISSIONS_PREFIXES.business.alarmStatistics}_filter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.alarmStatistics}_userFilter`,
      },
    },
    availability: {
      actions: {
        interval: `${USER_PERMISSIONS_PREFIXES.business.availability}_interval`,

        filter: `${USER_PERMISSIONS_PREFIXES.business.availability}_filter`,
        userFilter: `${USER_PERMISSIONS_PREFIXES.business.availability}_userFilter`,

        exportAsCsv: `${USER_PERMISSIONS_PREFIXES.business.availability}_exportAsCsv`,
      },
    },
  },
  api: {
    general: {
      acl: `${USER_PERMISSIONS_PREFIXES.api}_acl`,
      appInfoRead: `${USER_PERMISSIONS_PREFIXES.api}_app_info_read`,
      alarmRead: `${USER_PERMISSIONS_PREFIXES.api}_alarm_read`,
      alarmUpdate: `${USER_PERMISSIONS_PREFIXES.api}_alarm_update`,
      entity: `${USER_PERMISSIONS_PREFIXES.api}_entity`,
      entitycategory: `${USER_PERMISSIONS_PREFIXES.api}_entitycategory`,
      entitycomment: `${USER_PERMISSIONS_PREFIXES.api}_entitycomment`,
      entityservice: `${USER_PERMISSIONS_PREFIXES.api}_entityservice`,
      event: `${USER_PERMISSIONS_PREFIXES.api}_event`,
      view: `${USER_PERMISSIONS_PREFIXES.api}_view`,
      viewgroup: `${USER_PERMISSIONS_PREFIXES.api}_viewgroup`,
      privateViewGroups: `${USER_PERMISSIONS_PREFIXES.api}_private_view_groups`,
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
      metricsSettings: `${USER_PERMISSIONS_PREFIXES.api}_metrics_settings`,
      ratingSettings: `${USER_PERMISSIONS_PREFIXES.api}_rating_settings`,
      corporatePattern: `${USER_PERMISSIONS_PREFIXES.api}_corporate_pattern`,
      exportConfigurations: `${USER_PERMISSIONS_PREFIXES.api}_export_configurations`,
      map: `${USER_PERMISSIONS_PREFIXES.api}_map`,
      shareToken: `${USER_PERMISSIONS_PREFIXES.api}_share_token`,
      widgetTemplate: `${USER_PERMISSIONS_PREFIXES.api}_widgettemplate`,
      maintenance: `${USER_PERMISSIONS_PREFIXES.api}_maintenance`,
      alarmTag: `${USER_PERMISSIONS_PREFIXES.api}_alarm_tag`,
      theme: `${USER_PERMISSIONS_PREFIXES.api}_color_theme`,
      icon: `${USER_PERMISSIONS_PREFIXES.api}_icon`,
      techmetricsSettings: `${USER_PERMISSIONS_PREFIXES.api}_techmetrics_settings`,
      kpiFilter: `${USER_PERMISSIONS_PREFIXES.api}_kpi_filter`,
      messageRateStatsRead: `${USER_PERMISSIONS_PREFIXES.api}_message_rate_stats_read`,
      playlist: `${USER_PERMISSIONS_PREFIXES.api}_playlist`,
      launchEventRecording: `${USER_PERMISSIONS_PREFIXES.api}_launch_event_recording`,
      resendEvents: `${USER_PERMISSIONS_PREFIXES.api}_resend_events`,

      ...featuresService.get('constants.USER_PERMISSIONS.api.general'),
    },
    rules: {
      action: `${USER_PERMISSIONS_PREFIXES.api}_action`,
      dynamicinfos: `${USER_PERMISSIONS_PREFIXES.api}_dynamicinfos`,
      eventFilter: `${USER_PERMISSIONS_PREFIXES.api}_eventfilter`,
      idleRule: `${USER_PERMISSIONS_PREFIXES.api}_idlerule`,
      metaalarmrule: `${USER_PERMISSIONS_PREFIXES.api}_metaalarmrule`,
      flappingRule: `${USER_PERMISSIONS_PREFIXES.api}_flapping_rule`,
      resolveRule: `${USER_PERMISSIONS_PREFIXES.api}_resolve_rule`,
      snmpRule: `${USER_PERMISSIONS_PREFIXES.api}_snmprule`,
      snmpMib: `${USER_PERMISSIONS_PREFIXES.api}_snmpmib`,
      declareTicketRule: `${USER_PERMISSIONS_PREFIXES.api}_declare_ticket_rule`,
      declareTicketExecution: `${USER_PERMISSIONS_PREFIXES.api}_declare_ticket_execution`,
      linkRule: `${USER_PERMISSIONS_PREFIXES.api}_link_rule`,
    },
    remediation: {
      instruction: `${USER_PERMISSIONS_PREFIXES.api}_instruction`,
      jobConfig: `${USER_PERMISSIONS_PREFIXES.api}_job_config`,
      job: `${USER_PERMISSIONS_PREFIXES.api}_job`,
      execution: `${USER_PERMISSIONS_PREFIXES.api}_execution`,
      instructionApprove: `${USER_PERMISSIONS_PREFIXES.api}_instruction_approve`,
    },
    planning: {
      pbehavior: `${USER_PERMISSIONS_PREFIXES.api}_pbehavior`,
      pbehaviorException: `${USER_PERMISSIONS_PREFIXES.api}_pbehaviorexception`,
      pbehaviorReason: `${USER_PERMISSIONS_PREFIXES.api}_pbehaviorreason`,
      pbehaviorType: `${USER_PERMISSIONS_PREFIXES.api}_pbehaviortype`,
    },
  },
};

export const GROUPED_USER_PERMISSIONS = {
  planning: [
    USER_PERMISSIONS.technical.planningType,
    USER_PERMISSIONS.technical.planningReason,
    USER_PERMISSIONS.technical.planningExceptions,
  ],
  remediation: [
    USER_PERMISSIONS.technical.remediationInstruction,
    USER_PERMISSIONS.technical.remediationJob,
    USER_PERMISSIONS.technical.remediationConfiguration,
    USER_PERMISSIONS.technical.remediationStatistic,
  ],
};

export const GROUPED_USER_PERMISSIONS_KEYS = {
  planning: groupedPermissionToPermission(GROUPED_USER_PERMISSIONS.planning),
  remediation: groupedPermissionToPermission(GROUPED_USER_PERMISSIONS.remediation),
};

export const BUSINESS_USER_PERMISSIONS_ACTIONS_MAP = {
  alarmsList: {
    [ALARM_LIST_ACTIONS_TYPES.ack]: USER_PERMISSIONS.business.alarmsList.actions.ack,
    [ALARM_LIST_ACTIONS_TYPES.fastAck]: USER_PERMISSIONS.business.alarmsList.actions.fastAck,
    [ALARM_LIST_ACTIONS_TYPES.ackRemove]: USER_PERMISSIONS.business.alarmsList.actions.ackRemove,
    [ALARM_LIST_ACTIONS_TYPES.pbehaviorAdd]: USER_PERMISSIONS.business.alarmsList.actions.pbehaviorAdd,
    [ALARM_LIST_ACTIONS_TYPES.fastPbehaviorAdd]: USER_PERMISSIONS.business.alarmsList.actions.fastPbehaviorAdd,
    [ALARM_LIST_ACTIONS_TYPES.snooze]: USER_PERMISSIONS.business.alarmsList.actions.snooze,
    [ALARM_LIST_ACTIONS_TYPES.declareTicket]: USER_PERMISSIONS.business.alarmsList.actions.declareTicket,
    [ALARM_LIST_ACTIONS_TYPES.associateTicket]: USER_PERMISSIONS.business.alarmsList.actions.associateTicket,
    [ALARM_LIST_ACTIONS_TYPES.cancel]: USER_PERMISSIONS.business.alarmsList.actions.cancel,
    [ALARM_LIST_ACTIONS_TYPES.unCancel]: USER_PERMISSIONS.business.alarmsList.actions.unCancel,
    [ALARM_LIST_ACTIONS_TYPES.fastCancel]: USER_PERMISSIONS.business.alarmsList.actions.fastCancel,
    [ALARM_LIST_ACTIONS_TYPES.changeState]: USER_PERMISSIONS.business.alarmsList.actions.changeState,
    [ALARM_LIST_ACTIONS_TYPES.history]: USER_PERMISSIONS.business.alarmsList.actions.history,
    [ALARM_LIST_ACTIONS_TYPES.variablesHelp]: USER_PERMISSIONS.business.alarmsList.actions.variablesHelp,
    [ALARM_LIST_ACTIONS_TYPES.comment]: USER_PERMISSIONS.business.alarmsList.actions.comment,
    [ALARM_LIST_ACTIONS_TYPES.exportPdf]: USER_PERMISSIONS.business.alarmsList.actions.exportPdf,
    [ALARM_LIST_ACTIONS_TYPES.linkToMetaAlarm]:
    USER_PERMISSIONS.business.alarmsList.actions.manualMetaAlarmGroup,
    [ALARM_LIST_ACTIONS_TYPES.removeAlarmsFromManualMetaAlarm]:
    USER_PERMISSIONS.business.alarmsList.actions.manualMetaAlarmGroup,
    [ALARM_LIST_ACTIONS_TYPES.removeAlarmsFromAutoMetaAlarm]:
    USER_PERMISSIONS.business.alarmsList.actions.metaAlarmGroup,

    [ALARM_LIST_ACTIONS_TYPES.links]: USER_PERMISSIONS.business.alarmsList.actions.links,
    [ALARM_LIST_ACTIONS_TYPES.correlation]: USER_PERMISSIONS.business.alarmsList.actions.correlation,

    [ALARM_LIST_ACTIONS_TYPES.executeInstruction]:
    USER_PERMISSIONS.business.alarmsList.actions.executeInstruction,

    [ALARM_LIST_ACTIONS_TYPES.addBookmark]: USER_PERMISSIONS.business.alarmsList.actions.bookmark,
    [ALARM_LIST_ACTIONS_TYPES.removeBookmark]: USER_PERMISSIONS.business.alarmsList.actions.bookmark,
  },

  context: {
    [CONTEXT_ACTIONS_TYPES.createEntity]: USER_PERMISSIONS.business.context.actions.createEntity,
    [CONTEXT_ACTIONS_TYPES.editEntity]: USER_PERMISSIONS.business.context.actions.editEntity,
    [CONTEXT_ACTIONS_TYPES.duplicateEntity]: USER_PERMISSIONS.business.context.actions.duplicateEntity,
    [CONTEXT_ACTIONS_TYPES.deleteEntity]: USER_PERMISSIONS.business.context.actions.deleteEntity,
    [CONTEXT_ACTIONS_TYPES.pbehaviorAdd]: USER_PERMISSIONS.business.context.actions.pbehavior,
    [CONTEXT_ACTIONS_TYPES.massEnable]: USER_PERMISSIONS.business.context.actions.massEnable,
    [CONTEXT_ACTIONS_TYPES.massDisable]: USER_PERMISSIONS.business.context.actions.massDisable,
  },

  weather: {
    [WEATHER_ACTIONS_TYPES.entityAck]: USER_PERMISSIONS.business.serviceWeather.actions.entityAck,
    [WEATHER_ACTIONS_TYPES.entityAssocTicket]:
      USER_PERMISSIONS.business.serviceWeather.actions.entityAssocTicket,
    [WEATHER_ACTIONS_TYPES.entityDeclareTicket]:
      USER_PERMISSIONS.business.serviceWeather.actions.entityDeclareTicket,
    [WEATHER_ACTIONS_TYPES.entityValidate]: USER_PERMISSIONS.business.serviceWeather.actions.entityValidate,
    [WEATHER_ACTIONS_TYPES.entityInvalidate]:
      USER_PERMISSIONS.business.serviceWeather.actions.entityInvalidate,
    [WEATHER_ACTIONS_TYPES.entityPause]: USER_PERMISSIONS.business.serviceWeather.actions.entityPause,
    [WEATHER_ACTIONS_TYPES.entityPlay]: USER_PERMISSIONS.business.serviceWeather.actions.entityPlay,
    [WEATHER_ACTIONS_TYPES.entityCancel]: USER_PERMISSIONS.business.serviceWeather.actions.entityCancel,
    [WEATHER_ACTIONS_TYPES.executeInstruction]:
      USER_PERMISSIONS.business.serviceWeather.actions.executeInstruction,

    [WEATHER_ACTIONS_TYPES.entityLinks]: USER_PERMISSIONS.business.serviceWeather.actions.entityLinks,

    [WEATHER_ACTIONS_TYPES.moreInfos]: USER_PERMISSIONS.business.serviceWeather.actions.moreInfos,
    [WEATHER_ACTIONS_TYPES.alarmsList]: USER_PERMISSIONS.business.serviceWeather.actions.alarmsList,
    [WEATHER_ACTIONS_TYPES.pbehaviorList]: USER_PERMISSIONS.business.serviceWeather.actions.pbehaviorList,
    [WEATHER_ACTIONS_TYPES.entityComment]: USER_PERMISSIONS.business.serviceWeather.actions.entityComment,
  },

  counter: {
    [COUNTER_ACTIONS_TYPES.alarmsList]: USER_PERMISSIONS.business.counter.actions.alarmsList,
    [COUNTER_ACTIONS_TYPES.variablesHelp]: USER_PERMISSIONS.business.counter.actions.variablesHelp,
  },
};

export const USER_PERMISSIONS_TO_PAGES_RULES = {
  /**
   * Admin pages
   */
  [USER_PERMISSIONS.technical.healthcheck]: ADMIN_PAGES_RULES.healthcheck,
  [USER_PERMISSIONS.technical.kpi]: ADMIN_PAGES_RULES.kpi,
  [USER_PERMISSIONS.technical.tag]: ADMIN_PAGES_RULES.tag,
  [USER_PERMISSIONS.technical.map]: ADMIN_PAGES_RULES.map,

  /**
   * Grouped
   */
  [GROUPED_USER_PERMISSIONS_KEYS.remediation]: ADMIN_PAGES_RULES.remediation,

  /**
   * Exploitation pages
   */
  [USER_PERMISSIONS.technical.exploitation.eventFilter]: EXPLOITATION_PAGES_RULES.eventFilter,
  [USER_PERMISSIONS.technical.exploitation.snmpRule]: EXPLOITATION_PAGES_RULES.snmpRule,
  [USER_PERMISSIONS.technical.exploitation.dynamicInfo]: EXPLOITATION_PAGES_RULES.dynamicInfo,
  [USER_PERMISSIONS.technical.exploitation.metaAlarmRule]: EXPLOITATION_PAGES_RULES.metaAlarmRule,
  [USER_PERMISSIONS.technical.exploitation.scenario]: EXPLOITATION_PAGES_RULES.scenario,
  [USER_PERMISSIONS.technical.exploitation.declareTicketRule]: EXPLOITATION_PAGES_RULES.declareTicketRule,

  /**
   * Notifications pages
   */
  [USER_PERMISSIONS.technical.notification.instructionStats]: NOTIFICATIONS_PAGES_RULES.instructionStats,
};

export const DOCUMENTATION_LINKS = {
  /**
   * Exploitation
   */
  [USER_PERMISSIONS.technical.exploitation.eventFilter]: 'guide-utilisation/menu-exploitation/filtres-evenements/',
  [USER_PERMISSIONS.technical.exploitation.pbehavior]: 'guide-utilisation/cas-d-usage/comportements_periodiques/',
  [USER_PERMISSIONS.technical.exploitation.snmpRule]: 'interconnexions/Supervision/SNMPtrap/',
  [USER_PERMISSIONS.technical.exploitation.idleRules]: 'guide-utilisation/menu-exploitation/regles-inactivite/',
  [USER_PERMISSIONS.technical.exploitation.resolveRules]: 'guide-utilisation/menu-exploitation/regles-resolution/',
  [USER_PERMISSIONS.technical.exploitation.flappingRules]: 'guide-utilisation/menu-exploitation/regles-bagot/',
  [USER_PERMISSIONS.technical.exploitation.dynamicInfo]: 'guide-utilisation/menu-exploitation/informations-dynamiques/',
  [USER_PERMISSIONS.technical.exploitation.metaAlarmRule]: 'guide-utilisation/menu-exploitation/regles-metaalarme/',
  [USER_PERMISSIONS.technical.exploitation.scenario]: 'guide-utilisation/menu-exploitation/scenarios/',

  /**
   * Admin
   */
  [USER_PERMISSIONS.technical.broadcastMessage]: 'guide-utilisation/interface/broadcast-messages/',
  [USER_PERMISSIONS.technical.playlist]: 'guide-utilisation/interface/playlists/',

  /**
   * Grouped
   */
  [GROUPED_USER_PERMISSIONS_KEYS.planning]: 'guide-administration/moteurs/moteur-pbehavior/#adminitration-de-la-planification',
  [GROUPED_USER_PERMISSIONS_KEYS.remediation]: 'guide-utilisation/remediation/',

  /**
   * Notifications
   */
  // [USER_PERMISSIONS.technical.notification.instructionStats]: '', // TODO: TBD
};
