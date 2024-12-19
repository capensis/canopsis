import { USER_PERMISSIONS_GROUPS, USER_VIEWS_PERMISSIONS, USER_PERMISSIONS } from '@/constants';

export default {
  title: {
    /**
     * VIEWS PERMISSIONS
     */
    [USER_VIEWS_PERMISSIONS.viewActions]: 'View: actions',
    [USER_VIEWS_PERMISSIONS.viewGeneral]: 'View: general',

    /**
     * GROUPS
     */
    [USER_PERMISSIONS_GROUPS.commonviews]: 'Views',
    [USER_PERMISSIONS_GROUPS.commonviewsPlaylist]: 'Playlists',
    [USER_PERMISSIONS_GROUPS.widgets]: 'Widgets',
    [USER_PERMISSIONS_GROUPS.widgetsCommon]: 'Common',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmstatistics]: 'Alarm statistics',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmstatisticsWidgetsettings]: 'Alarm statistics widget settings',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmstatisticsViewsettings]: 'Alarm statistics view settings',
    [USER_PERMISSIONS_GROUPS.widgetsAvailability]: 'Availability',
    [USER_PERMISSIONS_GROUPS.widgetsAvailabilityActions]: 'Availability actions',
    [USER_PERMISSIONS_GROUPS.widgetsAvailabilityWidgetsettings]: 'Availability widget settings',
    [USER_PERMISSIONS_GROUPS.widgetsAvailabilityViewsettings]: 'Availability view settings',
    [USER_PERMISSIONS_GROUPS.widgetsBarchart]: 'Bar chart',
    [USER_PERMISSIONS_GROUPS.widgetsBarchartWidgetsettings]: 'Bar chart widget settings',
    [USER_PERMISSIONS_GROUPS.widgetsBarchartViewsettings]: 'Bar chart view settings',
    [USER_PERMISSIONS_GROUPS.widgetsCounter]: 'Counter',
    [USER_PERMISSIONS_GROUPS.widgetsContext]: 'Context explorer',
    [USER_PERMISSIONS_GROUPS.widgetsContextViewsettings]: 'Context explorer view settings',
    [USER_PERMISSIONS_GROUPS.widgetsContextEntityactions]: 'Context explorer entity actions',
    [USER_PERMISSIONS_GROUPS.widgetsContextActions]: 'Context explorer actions',
    [USER_PERMISSIONS_GROUPS.widgetsContextWidgetsettings]: 'Context explorer widget settings',
    [USER_PERMISSIONS_GROUPS.widgetsLinechart]: 'Line chart',
    [USER_PERMISSIONS_GROUPS.widgetsLinechartWidgetsettings]: 'Line chart widget settings',
    [USER_PERMISSIONS_GROUPS.widgetsLinechartViewsettings]: 'Line chart view settings',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmslist]: 'Alarms list',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmslistAlarmactions]: 'Alarm actions',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmslistViewsettings]: 'Alarm list view settings',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmslistActions]: 'Alarm list actions',
    [USER_PERMISSIONS_GROUPS.widgetsAlarmslistWidgetsettings]: 'Alarm list widget settings',
    [USER_PERMISSIONS_GROUPS.widgetsMap]: 'Mapping',
    [USER_PERMISSIONS_GROUPS.widgetsMapViewsettings]: 'Maps view settings',
    [USER_PERMISSIONS_GROUPS.widgetsMapWidgetsettings]: 'Maps widget settings',
    [USER_PERMISSIONS_GROUPS.widgetsNumbers]: 'Numbers',
    [USER_PERMISSIONS_GROUPS.widgetsNumbersWidgetsettings]: 'Numbers widget settings',
    [USER_PERMISSIONS_GROUPS.widgetsNumbersViewsettings]: 'Number view settings',
    [USER_PERMISSIONS_GROUPS.widgetsPiechart]: 'Pie chart',
    [USER_PERMISSIONS_GROUPS.widgetsPiechartWidgetsettings]: 'Pie chart widget settings',
    [USER_PERMISSIONS_GROUPS.widgetsPiechartViewsettings]: 'Pie chart view settings',
    [USER_PERMISSIONS_GROUPS.widgetsServiceweather]: 'Service weather',
    [USER_PERMISSIONS_GROUPS.widgetsServiceweatherAlarmactions]: 'Service weather alarm actions',
    [USER_PERMISSIONS_GROUPS.widgetsServiceweatherViewsettings]: 'Service weather view settings',
    [USER_PERMISSIONS_GROUPS.widgetsServiceweatherWidgetsettings]: 'Service weather widget settings',
    [USER_PERMISSIONS_GROUPS.widgetsTestingweather]: 'Testing weather',
    [USER_PERMISSIONS_GROUPS.widgetsUserstatistics]: 'User/Alarm statistics',
    [USER_PERMISSIONS_GROUPS.widgetsUserstatisticsWidgetsettings]: 'User/Alarm statistics widget settings',
    [USER_PERMISSIONS_GROUPS.widgetsUserstatisticsViewsettings]: 'User/Alarm statistics view settings',
    [USER_PERMISSIONS_GROUPS.api]: 'API',
    [USER_PERMISSIONS_GROUPS.apiGeneral]: 'General',
    [USER_PERMISSIONS_GROUPS.apiRules]: 'Rules',
    [USER_PERMISSIONS_GROUPS.apiRemediation]: 'Remediation',
    [USER_PERMISSIONS_GROUPS.apiPlanning]: 'Planning',
    [USER_PERMISSIONS_GROUPS.technical]: 'Technical',
    [USER_PERMISSIONS_GROUPS.technicalAdmin]: 'Admin',
    [USER_PERMISSIONS_GROUPS.technicalAdminCommunication]: 'Communication',
    [USER_PERMISSIONS_GROUPS.technicalAdminGeneral]: 'General',
    [USER_PERMISSIONS_GROUPS.technicalAdminAccess]: 'Access',
    [USER_PERMISSIONS_GROUPS.technicalExploitation]: 'Exploitation',
    [USER_PERMISSIONS_GROUPS.technicalNotification]: 'Notifications',
    [USER_PERMISSIONS_GROUPS.technicalViewsandwidgets]: 'Views and widgets',
    [USER_PERMISSIONS_GROUPS.technicalProfile]: 'Profile',
    [USER_PERMISSIONS_GROUPS.technicalToken]: 'Token',

    /**
     * Business Common Permissions
     */
    [USER_PERMISSIONS.business.alarmsList.actions.variablesHelp]: 'See the list of variables (in all widgets)',
    [USER_PERMISSIONS.business.context.actions.entityComment]: 'Manage entity comments (view, create, edit, delete)',

    /**
     * Business Alarms List Permissions
     */
    [USER_PERMISSIONS.business.alarmsList.actions.ack]: 'Ack',
    [USER_PERMISSIONS.business.alarmsList.actions.fastAck]: 'Fast ack',
    [USER_PERMISSIONS.business.alarmsList.actions.ackRemove]: 'Cancel ack',
    [USER_PERMISSIONS.business.alarmsList.actions.pbehaviorAdd]: 'Edit PBhaviours for alarm',
    [USER_PERMISSIONS.business.alarmsList.actions.fastPbehaviorAdd]: 'Fast pbehavior',
    [USER_PERMISSIONS.business.alarmsList.actions.snooze]: 'Snooze alarm',
    [USER_PERMISSIONS.business.alarmsList.actions.declareTicket]: 'Declare ticket',
    [USER_PERMISSIONS.business.alarmsList.actions.associateTicket]: 'Associate ticket',
    [USER_PERMISSIONS.business.alarmsList.actions.cancel]: 'Remove alarm',
    [USER_PERMISSIONS.business.alarmsList.actions.unCancel]: 'Uncalcel alarm',
    [USER_PERMISSIONS.business.alarmsList.actions.fastCancel]: 'Fast remove alarm',
    [USER_PERMISSIONS.business.alarmsList.actions.changeState]: 'Check and lock severity',
    [USER_PERMISSIONS.business.alarmsList.actions.history]: 'View alarm history',
    [USER_PERMISSIONS.business.alarmsList.actions.manualMetaAlarmGroup]: 'Link to manual meta alarm rule / Unlink',
    [USER_PERMISSIONS.business.alarmsList.actions.comment]: 'Comment alarm',
    [USER_PERMISSIONS.business.alarmsList.actions.filter]: 'Set alarm filters',
    [USER_PERMISSIONS.business.alarmsList.actions.userFilter]: 'Filter alarms',
    [USER_PERMISSIONS.business.alarmsList.actions.remediationInstructionsFilter]: 'Set filters by remediation instructions',
    [USER_PERMISSIONS.business.alarmsList.actions.userRemediationInstructionsFilter]: 'Filter alarms by remediation instructions',
    [USER_PERMISSIONS.business.alarmsList.actions.links]: 'Follow link in alarm',
    [USER_PERMISSIONS.business.alarmsList.actions.correlation]: 'Group correlated alarms (meta alarms)',
    [USER_PERMISSIONS.business.alarmsList.actions.executeInstruction]: 'Execute manual instructions',
    [USER_PERMISSIONS.business.alarmsList.actions.category]: 'Filter alarms by category',
    [USER_PERMISSIONS.business.alarmsList.actions.exportPdf]: 'Export in PDF',
    [USER_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: 'Export alarm list as CSV',
    [USER_PERMISSIONS.business.alarmsList.actions.metaAlarmGroup]: 'Unlink alarm from auto meta alarm',
    [USER_PERMISSIONS.business.alarmsList.actions.bookmark]: 'Add / remove bookmark',
    [USER_PERMISSIONS.business.alarmsList.actions.filterByBookmark]: 'Filter bookmarked alarms',

    /**
     * Business Context Explorer Permissions
     */
    [USER_PERMISSIONS.business.context.actions.createEntity]: 'Create entity',
    [USER_PERMISSIONS.business.context.actions.editEntity]: 'Edit entity',
    [USER_PERMISSIONS.business.context.actions.duplicateEntity]: 'Duplicate entity',
    [USER_PERMISSIONS.business.context.actions.deleteEntity]: 'Delete entity',
    [USER_PERMISSIONS.business.context.actions.massEnable]: 'Mass action to enable selected entities',
    [USER_PERMISSIONS.business.context.actions.massDisable]: 'Mass action to disable selected entities',
    [USER_PERMISSIONS.business.context.actions.pbehavior]: 'Set PBehavior',
    [USER_PERMISSIONS.business.context.actions.filter]: 'Set entities filters',
    [USER_PERMISSIONS.business.context.actions.userFilter]: 'Filter entities',
    [USER_PERMISSIONS.business.context.actions.category]: 'Filter entities by category',
    [USER_PERMISSIONS.business.context.actions.exportAsCsv]: 'Export entities as CSV file',

    /**
     * Business Service Weather Permissions
     */
    [USER_PERMISSIONS.business.serviceWeather.actions.entityAck]: 'Ack',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityAssocTicket]: 'Associate ticket',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityDeclareTicket]: 'Declare ticket',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityComment]: 'Comment alarm',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityValidate]: 'Validate alarms and change their state to critical',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityInvalidate]: 'Invalidate alarms and cancel them',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityPause]: 'Pause alarms (set the PBehavior type "Pause")',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityPlay]: 'Activate paused alarms (remove the PBehavior type "Pause")',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityCancel]: 'Access the list of PBehaviors associated to services (in the subtab in the services modal windows)',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityManagePbehaviors]: 'View PBehaviors of services (in the subtab in the services modal windows)',
    [USER_PERMISSIONS.business.serviceWeather.actions.executeInstruction]: 'Execute manual instructions',
    [USER_PERMISSIONS.business.serviceWeather.actions.entityLinks]: 'Follow link in alarm',
    [USER_PERMISSIONS.business.serviceWeather.actions.moreInfos]: 'Open "More infos" modal',
    [USER_PERMISSIONS.business.serviceWeather.actions.alarmsList]: 'Open the list of alarms available for each service',
    [USER_PERMISSIONS.business.serviceWeather.actions.pbehaviorList]: 'View PBehaviors of services (in the subtab in the service entities modal windows)',
    [USER_PERMISSIONS.business.serviceWeather.actions.filter]: 'Set alarm filters',
    [USER_PERMISSIONS.business.serviceWeather.actions.userFilter]: 'Filter alarms',
    [USER_PERMISSIONS.business.serviceWeather.actions.category]: 'Filter alarms by category',

    /**
     * Business Counter Permissions
     */
    [USER_PERMISSIONS.business.counter.actions.alarmsList]: 'View the alarm list associated with counters',

    /**
     * Business Testing Weather Permissions
     */
    [USER_PERMISSIONS.business.testingWeather.actions.alarmsList]: 'Open the list of alarms available',

    /**
     * Business Map Permissions
     */
    [USER_PERMISSIONS.business.map.actions.alarmsList]: 'View the alarm list associated with points on maps',
    [USER_PERMISSIONS.business.map.actions.filter]: 'Set filters for points on maps',
    [USER_PERMISSIONS.business.map.actions.userFilter]: 'Filter points on maps',
    [USER_PERMISSIONS.business.map.actions.category]: 'Filter points on maps by categories',

    /**
     * Business Bar Chart Permissions
     */
    [USER_PERMISSIONS.business.barChart.actions.interval]: 'Edit time intervals for the data displayed',
    [USER_PERMISSIONS.business.barChart.actions.sampling]: 'Edit sampling for the data displayed',
    [USER_PERMISSIONS.business.barChart.actions.filter]: 'Set data filters',
    [USER_PERMISSIONS.business.barChart.actions.userFilter]: 'Filter data',

    /**
     * Business Line Chart Permissions
     */
    [USER_PERMISSIONS.business.lineChart.actions.interval]: 'Edit time intervals for the data displayed',
    [USER_PERMISSIONS.business.lineChart.actions.sampling]: 'Edit sampling for the data displayed',
    [USER_PERMISSIONS.business.lineChart.actions.filter]: 'Set data filters',
    [USER_PERMISSIONS.business.lineChart.actions.userFilter]: 'Filter data',

    /**
     * Business Pie Chart Permissions
     */
    [USER_PERMISSIONS.business.pieChart.actions.interval]: 'Edit time intervals for the data displayed',
    [USER_PERMISSIONS.business.pieChart.actions.sampling]: 'Edit sampling for the data displayed',
    [USER_PERMISSIONS.business.pieChart.actions.filter]: 'Set data filters',
    [USER_PERMISSIONS.business.pieChart.actions.userFilter]: 'Filter data',

    /**
     * Business Numbers Permissions
     */
    [USER_PERMISSIONS.business.numbers.actions.interval]: 'Edit time intervals for the data displayed',
    [USER_PERMISSIONS.business.numbers.actions.sampling]: 'Edit sampling for the data displayed',
    [USER_PERMISSIONS.business.numbers.actions.filter]: 'Set data filters',
    [USER_PERMISSIONS.business.numbers.actions.userFilter]: 'Filter data',

    /**
     * Business User Statistics
     */
    [USER_PERMISSIONS.business.userStatistics.actions.interval]: 'Edit time intervals for the data displayed',
    [USER_PERMISSIONS.business.userStatistics.actions.filter]: 'Set data filters',
    [USER_PERMISSIONS.business.userStatistics.actions.userFilter]: 'Filter data',

    /**
     * Business Alarm Statistics
     */
    [USER_PERMISSIONS.business.alarmStatistics.actions.interval]: 'Edit time intervals for the data displayed',
    [USER_PERMISSIONS.business.alarmStatistics.actions.filter]: 'Set data filters',
    [USER_PERMISSIONS.business.alarmStatistics.actions.userFilter]: 'Filter data',

    /**
     * Business Availability
     */
    [USER_PERMISSIONS.business.availability.actions.interval]: 'Edit time intervals for the data displayed',
    [USER_PERMISSIONS.business.availability.actions.filter]: 'Set data filters',
    [USER_PERMISSIONS.business.availability.actions.userFilter]: 'Filter data',
    [USER_PERMISSIONS.business.availability.actions.exportAsCsv]: 'Export availabilities as CSV file',

    /**
     * Technical Admin Communication
     */
    [USER_PERMISSIONS.technical.broadcastMessage]: 'Broadcast Messages',
    [USER_PERMISSIONS.technical.playlist]: 'Playlists',

    /**
     * Technical Admin General
     */
    [USER_PERMISSIONS.technical.eventsRecord]: 'Events records',
    [USER_PERMISSIONS.technical.healthcheck]: 'Healthcheck',
    [USER_PERMISSIONS.technical.healthcheckStatus]: 'Healthcheck status',
    [USER_PERMISSIONS.technical.icon]: 'Parameters - icons',
    [USER_PERMISSIONS.technical.kpi]: 'KPI Graphs',
    [USER_PERMISSIONS.technical.kpiCollectionSettings]: 'KPI Collection settings',
    [USER_PERMISSIONS.technical.kpiFilters]: 'KPI Filters',
    [USER_PERMISSIONS.technical.kpiRatingSettings]: 'KPI Rating settings',
    [USER_PERMISSIONS.technical.maintenance]: 'Maintenance mode',
    [USER_PERMISSIONS.technical.map]: 'Maps',
    [USER_PERMISSIONS.technical.parameters]: 'Parameters - parameters tab',
    [USER_PERMISSIONS.technical.planningExceptions]: 'Planning exceptions dates (Pbehavior)',
    [USER_PERMISSIONS.technical.planningReason]: 'Planning reason (Pbehavior)',
    [USER_PERMISSIONS.technical.planningType]: 'Planning type (Pbehavior)',
    [USER_PERMISSIONS.technical.remediationConfiguration]: 'Instructions - configurations tab',
    [USER_PERMISSIONS.technical.remediationInstruction]: 'Instructions - instructions tab',
    [USER_PERMISSIONS.technical.remediationJob]: 'Instructions - jobs tab',
    [USER_PERMISSIONS.technical.remediationStatistic]: 'Instructions - remediation statistics tab',
    [USER_PERMISSIONS.technical.stateSetting]: 'State settings',
    [USER_PERMISSIONS.technical.storageSettings]: 'Storage settings',
    [USER_PERMISSIONS.technical.tag]: 'Tags management',
    [USER_PERMISSIONS.technical.techmetrics]: 'Healthcheck - engines\' metrics',
    [USER_PERMISSIONS.technical.widgetTemplate]: 'Parameters - widget templates',
    [USER_PERMISSIONS.technical.viewImportExport]: 'Parameters - import / export',

    /**
     * Technical Admin Access
     */
    [USER_PERMISSIONS.technical.permission]: 'Rights',
    [USER_PERMISSIONS.technical.role]: 'Roles',
    [USER_PERMISSIONS.technical.user]: 'Users',

    /**
     * Technical Admin Exploitation
     */
    [USER_PERMISSIONS.technical.exploitation.declareTicketRule]: 'Ticket declaration rules',
    [USER_PERMISSIONS.technical.exploitation.dynamicInfo]: 'Dynamic information rules',
    [USER_PERMISSIONS.technical.exploitation.eventFilter]: 'Event filters',
    [USER_PERMISSIONS.technical.exploitation.flappingRules]: 'Flapping rules',
    [USER_PERMISSIONS.technical.exploitation.idleRules]: 'Idle rules',
    [USER_PERMISSIONS.technical.exploitation.linkRule]: 'Link generator',
    [USER_PERMISSIONS.technical.exploitation.metaAlarmRule]: 'Meta alarm rules and correlation',
    [USER_PERMISSIONS.technical.exploitation.pbehavior]: 'Pbehaviors',
    [USER_PERMISSIONS.technical.exploitation.resolveRules]: 'Resolve rules',
    [USER_PERMISSIONS.technical.exploitation.scenario]: 'Scenarios',
    [USER_PERMISSIONS.technical.exploitation.snmpRule]: 'Snmp rules',

    /**
     * Technical Admin Notification
     */
    [USER_PERMISSIONS.technical.notification.common]: 'Parameters - notification settings ',
    [USER_PERMISSIONS.technical.notification.instructionStats]: 'Instructions stats',

    /**
     * Technical Admin Views and widgets
     */
    [USER_PERMISSIONS.technical.privateView]: 'Private views',
    [USER_PERMISSIONS.technical.view]: 'Views',

    /**
     * Technical Admin Profile
     */
    [USER_PERMISSIONS.technical.profile.theme]: 'Theme colors',
    [USER_PERMISSIONS.technical.profile.corporatePattern]: 'Corporate patterns',

    /**
     * Technical Admin Token
     */
    [USER_PERMISSIONS.technical.shareToken]: 'Shared token settings',

    /**
     * API Permissions General
     */
    [USER_PERMISSIONS.api.general.acl]: 'Roles, permissions, users',
    [USER_PERMISSIONS.api.general.alarmRead]: 'Read alarms',
    [USER_PERMISSIONS.api.general.alarmTag]: 'Alarm tags',
    [USER_PERMISSIONS.api.general.alarmUpdate]: 'Update alarms',
    [USER_PERMISSIONS.api.general.associativeTable]: 'Associative tables',
    [USER_PERMISSIONS.api.general.broadcastMessage]: 'Broadcast Message',
    [USER_PERMISSIONS.api.general.theme]: 'Theme colors',
    [USER_PERMISSIONS.api.general.contextgraph]: 'Context graph import',
    [USER_PERMISSIONS.api.general.corporatePattern]: 'Corporate patterns',
    [USER_PERMISSIONS.api.general.datastorageRead]: 'Data storage settings read',
    [USER_PERMISSIONS.api.general.datastorageUpdate]: 'Data storage settings update',
    [USER_PERMISSIONS.api.general.entity]: 'Entities',
    [USER_PERMISSIONS.api.general.entitycategory]: 'Entity categories',
    [USER_PERMISSIONS.api.general.entitycomment]: 'Entity comments',
    [USER_PERMISSIONS.api.general.entityservice]: 'Entity services',
    [USER_PERMISSIONS.api.general.event]: 'Events',
    [USER_PERMISSIONS.api.general.exportConfigurations]: 'Export configurations',
    [USER_PERMISSIONS.api.general.files]: 'Files',
    [USER_PERMISSIONS.api.general.healthcheck]: 'Healthcheck',
    [USER_PERMISSIONS.api.general.icon]: 'Icons',
    [USER_PERMISSIONS.api.general.junit]: 'JUnit',
    [USER_PERMISSIONS.api.general.kpiFilter]: 'KPI Filters',
    [USER_PERMISSIONS.api.general.launchEventRecording]: 'Launch events recording',
    [USER_PERMISSIONS.api.general.maintenance]: 'Maintenance mode',
    [USER_PERMISSIONS.api.general.map]: 'Maps',
    [USER_PERMISSIONS.api.general.messageRateStatsRead]: 'Message rate statistics',
    [USER_PERMISSIONS.api.general.metrics]: 'Metrics',
    [USER_PERMISSIONS.api.general.metricsSettings]: 'Metrics settings',
    [USER_PERMISSIONS.api.general.notifications]: 'Notification settings',
    [USER_PERMISSIONS.api.general.playlist]: 'Playlists',
    [USER_PERMISSIONS.api.general.privateViewGroups]: 'Private view groups',
    [USER_PERMISSIONS.api.general.ratingSettings]: 'Rating settings',
    [USER_PERMISSIONS.api.general.resendEvents]: 'Resend events',
    [USER_PERMISSIONS.api.general.shareToken]: 'Share tokens',
    [USER_PERMISSIONS.api.general.stateSettings]: 'State settings',
    [USER_PERMISSIONS.api.general.techmetrics]: 'Tech metrics',
    [USER_PERMISSIONS.api.general.techmetricsSettings]: 'Tech metrics settings',
    [USER_PERMISSIONS.api.general.userInterfaceDelete]: 'Delete user interface',
    [USER_PERMISSIONS.api.general.userInterfaceUpdate]: 'Update user interface',
    [USER_PERMISSIONS.api.general.view]: 'Views',
    [USER_PERMISSIONS.api.general.viewgroup]: 'View groups',
    [USER_PERMISSIONS.api.general.widgetTemplate]: 'Widget templates',

    /**
     * API Permissions Rules
     */
    [USER_PERMISSIONS.api.rules.action]: 'Actions',
    [USER_PERMISSIONS.api.rules.declareTicketExecution]: 'Run declare ticket rules',
    [USER_PERMISSIONS.api.rules.declareTicketRule]: 'Declare ticket rules',
    [USER_PERMISSIONS.api.rules.dynamicinfos]: 'Dynamic infos',
    [USER_PERMISSIONS.api.rules.eventFilter]: 'Event filters',
    [USER_PERMISSIONS.api.rules.flappingRule]: 'Flapping rules',
    [USER_PERMISSIONS.api.rules.idleRule]: 'Idle rule',
    [USER_PERMISSIONS.api.rules.linkRule]: 'Link generator',
    [USER_PERMISSIONS.api.rules.metaalarmrule]: 'Meta alarm rules',
    [USER_PERMISSIONS.api.rules.resolveRule]: 'Resolve rules',
    [USER_PERMISSIONS.api.rules.snmpRule]: 'SNMP rules',
    [USER_PERMISSIONS.api.rules.snmpMib]: 'SNMP MIB',

    /**
     * API Permissions Remediation
     */
    [USER_PERMISSIONS.api.remediation.execution]: 'Runs instructions',
    [USER_PERMISSIONS.api.remediation.instruction]: 'Instructions',
    [USER_PERMISSIONS.api.remediation.instructionApprove]: 'Approve instructions',
    [USER_PERMISSIONS.api.remediation.job]: 'Jobs',
    [USER_PERMISSIONS.api.remediation.jobConfig]: 'Job configs',

    /**
     * API Permissions Planning
     */
    [USER_PERMISSIONS.api.planning.pbehavior]: 'PBehaviors',
    [USER_PERMISSIONS.api.planning.pbehaviorException]: 'PBehavior exceptions',
    [USER_PERMISSIONS.api.planning.pbehaviorReason]: 'PBehavior reasons',
    [USER_PERMISSIONS.api.planning.pbehaviorType]: 'PBehavior types',
  },
};
