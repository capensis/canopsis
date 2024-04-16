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
    [USER_PERMISSIONS_PREFIXES.business.userStatistics]: 'Rights for User Statistics',
    [USER_PERMISSIONS_PREFIXES.business.alarmStatistics]: 'Rights for Alarm Statistics',
    [USER_PERMISSIONS_PREFIXES.business.availability]: 'Rights for Availability',
  },
  api: {
    general: 'General',
    rules: 'Rules',
    remediation: 'Remediation',
    pbehavior: 'PBehavior',
  },
  permissions: {
    /**
     * Business Common Permissions
     */
    [USERS_PERMISSIONS.business.alarmsList.actions.variablesHelp]: {
      name: 'Access to available variables list',
      description: 'Users with this permission can see the list of variables in the alarm list and service weather',
    },

    /**
     * Business Alarms List Permissions
     */
    [USERS_PERMISSIONS.business.alarmsList.actions.ack]: {
      name: 'Rights on alarm list: ack',
      description: 'Users with this permission can acknowledge alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.fastAck]: {
      name: 'Rights on alarm list: fast ack',
      description: 'Users with this permission can do fast acknowledge of alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.ackRemove]: {
      name: 'Rights on alarm list: cancel ack',
      description: 'Users with this permission can cancel ack',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.pbehaviorAdd]: {
      name: 'Rights on alarm list: pbehavior action',
      description: 'Users with this permission can access to the action "Periodical behavior" and edit PBehaviors for alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.snooze]: {
      name: 'Rights on alarm list: snooze alarm',
      description: 'Users with this permission can snooze alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.declareTicket]: {
      name: 'Rights on alarm list: declare ticket',
      description: 'Users with this permission can do the tickets declaration',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.associateTicket]: {
      name: 'Rights on alarm list: associate ticket',
      description: 'Users with this permission can associate a ticket',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.cancel]: {
      name: 'Rights on alarm list: cancel alarm',
      description: 'Users with this permission can cancel alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.fastCancel]: {
      name: 'Rights on alarm list: fast alarm cancelation',
      description: 'Users with this permission can do fast alarms cancelation',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.changeState]: {
      name: 'Rights on alarm list: change state',
      description: 'Users with this permission can change states of alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.history]: {
      name: 'Rights on alarm list: history',
      description: 'Users with this permission can view alarms history',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.manualMetaAlarmGroup]: {
      name: 'Rights on alarm list: Manual meta alarm actions',
      description: 'Users with this permission can apply manual meta alarm rules and group alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.comment]: {
      name: 'Rights on alarm list: Access to \'Comment\' action',
      description: 'Users with this permission can comment alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.listFilters]: {
      name: 'Rights on alarm list: view alarm filters',
      description: 'Users with this permission can view the list of available filters in the alarm list',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.editFilter]: {
      name: 'Rights on alarm list: edit alarm filters',
      description: 'Users with this permission can edit filters for alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.addFilter]: {
      name: 'Rights on alarm list: add alarm filters',
      description: 'Users with this permission can add filters for alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.userFilter]: {
      name: 'Rights on alarm list: show alarm filters',
      description: 'The alarm filter is shown for users with this permission',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.listRemediationInstructionsFilters]: {
      name: 'Rights on alarm list: Access to view filters by remediation instructions',
      description: 'Users with this permission can see and apply the list of created filters by instructions',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.editRemediationInstructionsFilter]: {
      name: 'Rights on alarm list: Access to editing filters by remediation instructions',
      description: 'Users with this permission can edit filters by instructions',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.addRemediationInstructionsFilter]: {
      name: 'Rights on alarm list: Access to adding filters by remediation instructions',
      description: 'Users with this permission can add filters by instructions',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.userRemediationInstructionsFilter]: {
      name: 'Rights on alarm list: Access to filters by remediation instructions',
      description: 'The filter by instructions is shown for users with this permission',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.links]: {
      name: 'Rights on alarm list: Access to Links',
      description: 'Users with this permission can access and follow the links in the alarm list',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.correlation]: {
      name: 'Rights on alarm list: Access to grouping correlated alarms',
      description: 'Users with this permission can enable grouping correlated alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.executeInstruction]: {
      name: 'Rights on alarm list: Access to instructions executions',
      description: 'Users with this permission can execute instructions to remediate alarms',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.category]: {
      name: 'Rights on alarm list: Filter by category',
      description: 'Users with this permission can filter alarm list by category',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.exportAsCsv]: {
      name: 'Rights on alarm list: Access to exporting alarms as CSV',
      description: 'Users with this permission can export alarms to CSV',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.addBookmark]: {
      name: 'Rights on alarm list: Access to adding bookmark to alarms',
      description: 'Users with this permission can add bookmark to alarm',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.removeBookmark]: {
      name: 'Rights on alarm list: Access to removing bookmark from alarm',
      description: 'Users with this permission can remove bookmark from alarm',
    },
    [USERS_PERMISSIONS.business.alarmsList.actions.filterByBookmark]: {
      name: 'Rights on alarm list: Access to filter alarms by only bookmarks',
      description: 'Users with this permission can filter alarms by only bookmarks',
    },

    /**
     * Business Context Explorer Permissions
     */
    [USERS_PERMISSIONS.business.context.actions.createEntity]: {
      name: 'Rights on context explorer: create entity',
      description: 'Users with this permission can create new entities',
    },
    [USERS_PERMISSIONS.business.context.actions.editEntity]: {
      name: 'Rights on context explorer: edit entity',
      description: 'Users with this permission can edit entities',
    },
    [USERS_PERMISSIONS.business.context.actions.duplicateEntity]: {
      name: 'Rights on context explorer: duplicate entity',
      description: 'Users with this permission can duplicate entities',
    },
    [USERS_PERMISSIONS.business.context.actions.deleteEntity]: {
      name: 'Rights on context explorer: delete entity',
      description: 'Users with this permission can delete entities',
    },
    [USERS_PERMISSIONS.business.context.actions.pbehaviorAdd]: {
      name: 'Rights on context explorer: PBehavior action',
      description: 'Users with this permission can access to the action "Periodical behavior" and edit PBehaviors for entities',
    },
    [USERS_PERMISSIONS.business.context.actions.massEnable]: {
      name: 'Rights on context explorer: Mass enable action',
      description: 'Users with this permission can perform mass action to enable selected entities',
    },
    [USERS_PERMISSIONS.business.context.actions.massDisable]: {
      name: 'Rights on context explorer: Mass disable action',
      description: 'Users with this permission can perform mass action to disable selected entities',
    },
    [USERS_PERMISSIONS.business.context.actions.listFilters]: {
      name: 'Rights on context explorer: view filters',
      description: 'Users with this permission can see the list of filters available in the Context explorer',
    },
    [USERS_PERMISSIONS.business.context.actions.editFilter]: {
      name: 'Rights on context explorer: edit filters',
      description: 'Users with this permission can edit entity filters',
    },
    [USERS_PERMISSIONS.business.context.actions.addFilter]: {
      name: 'Rights on context explorer: add filters',
      description: 'Users with this permission can add filters on entities shown in Context explorer',
    },
    [USERS_PERMISSIONS.business.context.actions.userFilter]: {
      name: 'Rights on context explorer: show filters',
      description: 'The entity filter is shown for users with this permission',
    },
    [USERS_PERMISSIONS.business.context.actions.category]: {
      name: 'Rights on context explorer: Filter by category',
      description: 'Users with this permission can filter entities by category',
    },
    [USERS_PERMISSIONS.business.context.actions.exportAsCsv]: {
      name: 'Rights on context explorer: Export as csv',
      description: 'Users with this permission can export entities as CSV file',
    },

    /**
     * Business Service Weather Permissions
     */
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityAck]: {
      name: 'Service weather: Access to \'Ack\' action',
      description: 'Users with this permission can acknowledge alarms',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityAssocTicket]: {
      name: 'Service weather: Access to \'Associate Ticket\' action',
      description: 'Users with this permission can associate tickets for alarms',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityComment]: {
      name: 'Service weather: Access to \'Comment\' action',
      description: 'Users with this permission can add comments',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityValidate]: {
      name: 'Service weather: Access to \'Validate\' action',
      description: 'Users with this permission can validate alarms and change their state to critical',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityInvalidate]: {
      name: 'Service weather: Access to \'Invalidate\' action',
      description: 'Users with this permission can invalidate alarms and cancel them',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityPause]: {
      name: 'Service weather: Access to \'Pause\' action',
      description: 'Users with this permission can pause alarms (set the PBehavior type "Pause")',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityPlay]: {
      name: 'Service weather: Access to \'Play\' action',
      description: 'Users with this permission can activate paused alarms (remove the PBehavior type "Pause")',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityCancel]: {
      name: 'Service weather: Access to \'Cancel\' action',
      description: 'Users with this permission can cancel alarms',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityManagePbehaviors]: {
      name: 'Service weather: Access to pbehaviors management',
      description: 'Users with this permission can access the list of PBehaviors associated to services (in the subtab in the services modal windows)',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.executeInstruction]: {
      name: 'Service weather: Access to execute instruction',
      description: 'Users with this permission can execute instructions for alarms',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.entityLinks]: {
      name: 'Service weather: Access to Links',
      description: 'Users with this permission can see links associated with alarms',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.moreInfos]: {
      name: 'Service weather: Access to \'More infos\' modal',
      description: 'Users with this permission can access to "More infos" modal window',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.alarmsList]: {
      name: 'Service weather: Access to \'Alarms list\' modal',
      description: 'Users with this permission can open the list of alarms available for each service',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.pbehaviorList]: {
      name: 'Service weather: Access to service pbehavior list',
      description: 'Users with this permission can access the list of all PBehaviors of services (in the subtab in the service entities modal windows)',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.listFilters]: {
      name: 'Rights on service weather: View filters',
      description: 'Users with this permission can see the list of filters available',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.editFilter]: {
      name: 'Rights on service weather: Edit filter',
      description: 'Users with this permission can edit applied filters',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.addFilter]: {
      name: 'Rights on service weather: Add filter',
      description: 'Users with this permission can add filters',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.userFilter]: {
      name: 'Rights on service weather: Show filter',
      description: 'The filter is shown for users with this permission',
    },
    [USERS_PERMISSIONS.business.serviceWeather.actions.category]: {
      name: 'Rights on service weather: Filter by category',
      description: 'Users with this permission can filter services by category',
    },

    /**
     * Business Counter Permissions
     */
    [USERS_PERMISSIONS.business.counter.actions.alarmsList]: {
      name: 'Counter: Access to \'Alarms list\' modal',
      description: 'Users with this permission can see the alarm list associated with counters',
    },

    /**
     * Business Testing Weather Permissions
     */
    [USERS_PERMISSIONS.business.testingWeather.actions.alarmsList]: {
      name: 'Testing weather: Access to \'Alarms list\' modal',
      description: 'Users with this permission can see the alarm list associated with testing weather',
    },

    /**
     * Business Testing Weather Permissions
     */
    [USERS_PERMISSIONS.business.map.actions.alarmsList]: {
      name: 'Rights on maps: Access to \'Alarms list\' modal',
      description: 'Users with this permission can see the alarm list associated with points on maps',
    },
    [USERS_PERMISSIONS.business.map.actions.listFilters]: {
      name: 'Rights on maps: View filter',
      description: 'Users with this permission can see the list of filters available',
    },
    [USERS_PERMISSIONS.business.map.actions.editFilter]: {
      name: 'Rights on maps: Edit filter',
      description: 'Users with this permission can edit filters for maps',
    },
    [USERS_PERMISSIONS.business.map.actions.addFilter]: {
      name: 'Rights on maps: Add filter',
      description: 'Users with this permission can add filters for maps',
    },
    [USERS_PERMISSIONS.business.map.actions.userFilter]: {
      name: 'Rights on maps: Show filter',
      description: 'The filter is shown for users with this permission',
    },
    [USERS_PERMISSIONS.business.map.actions.category]: {
      name: 'Rights on maps: Access to \'Category\' action',
      description: 'Users with this permission can filter points by categories',
    },

    /**
     * Business Bar Chart Permissions
     */
    [USERS_PERMISSIONS.business.barChart.actions.interval]: {
      name: 'Barchart: interval',
      description: 'Users with this permission can edit time intervals for the data displayed',
    },
    [USERS_PERMISSIONS.business.barChart.actions.sampling]: {
      name: 'Barchart: sampling',
      description: 'Users with this permission can change sampling for the data displayed',
    },
    [USERS_PERMISSIONS.business.barChart.actions.listFilters]: {
      name: 'Barchart: view filters',
      description: 'Users with this permission can see the list of filters available',
    },
    [USERS_PERMISSIONS.business.barChart.actions.editFilter]: {
      name: 'Barchart: edit filters',
      description: 'Users with this permission can edit filters',
    },
    [USERS_PERMISSIONS.business.barChart.actions.addFilter]: {
      name: 'Barchart: add filters',
      description: 'Users with this permission can add filters',
    },
    [USERS_PERMISSIONS.business.barChart.actions.userFilter]: {
      name: 'Barchart: show filters',
      description: 'The filter is shown for users with this permission',
    },

    /**
     * Business Line Chart Permissions
     */
    [USERS_PERMISSIONS.business.lineChart.actions.interval]: {
      name: 'Linechart: interval',
      description: 'Users with this permission can change',
    },
    [USERS_PERMISSIONS.business.lineChart.actions.sampling]: {
      name: 'Linechart: sampling',
      description: 'Users with this permission can change sampling for the data displayed',
    },
    [USERS_PERMISSIONS.business.lineChart.actions.listFilters]: {
      name: 'Linechart: view filters',
      description: 'Users with this permission can see the list of filters available',
    },
    [USERS_PERMISSIONS.business.lineChart.actions.editFilter]: {
      name: 'Linechart: edit filters',
      description: 'Users with this permission can edit filters',
    },
    [USERS_PERMISSIONS.business.lineChart.actions.addFilter]: {
      name: 'Linechart: add filters',
      description: 'Users with this permission can add filters',
    },
    [USERS_PERMISSIONS.business.lineChart.actions.userFilter]: {
      name: 'Linechart: show filters',
      description: 'The filter is shown for users with this permission',
    },

    /**
     * Business Pie Chart Permissions
     */
    [USERS_PERMISSIONS.business.pieChart.actions.interval]: {
      name: 'Piechart: interval',
      description: 'Users with this permission can change time interval for the data displayed',
    },
    [USERS_PERMISSIONS.business.pieChart.actions.sampling]: {
      name: 'Piechart: sampling',
      description: 'Users with this permission can change sampling for the data displayed',
    },
    [USERS_PERMISSIONS.business.pieChart.actions.listFilters]: {
      name: 'Piechart: view filters',
      description: 'Users with this permission can see the list of filters available',
    },
    [USERS_PERMISSIONS.business.pieChart.actions.editFilter]: {
      name: 'Piechart: edit filters',
      description: 'Users with this permission can edit filters',
    },
    [USERS_PERMISSIONS.business.pieChart.actions.addFilter]: {
      name: 'Piechart: add filters',
      description: 'Users with this permission can add filters',
    },
    [USERS_PERMISSIONS.business.pieChart.actions.userFilter]: {
      name: 'Piechart: show filters',
      description: 'The filter is shown for users with this permission',
    },

    /**
     * Business Numbers Permissions
     */
    [USERS_PERMISSIONS.business.numbers.actions.interval]: {
      name: 'Numbers: interval',
      description: 'Users with this permission can change time interval for the data displayed',
    },
    [USERS_PERMISSIONS.business.numbers.actions.sampling]: {
      name: 'Numbers: sampling',
      description: 'Users with this permission can change sampling for the data displayed',
    },
    [USERS_PERMISSIONS.business.numbers.actions.listFilters]: {
      name: 'Numbers: view filters',
      description: 'Users with this permission can see the list of filters available',
    },
    [USERS_PERMISSIONS.business.numbers.actions.editFilter]: {
      name: 'Numbers: edit filters',
      description: 'Users with this permission can edit filters',
    },
    [USERS_PERMISSIONS.business.numbers.actions.addFilter]: {
      name: 'Numbers: add filters',
      description: 'Users with this permission can add filters',
    },
    [USERS_PERMISSIONS.business.numbers.actions.userFilter]: {
      name: 'Numbers: show filters',
      description: 'The filter is shown for users with this permission',
    },

    /**
     * Business User Statistics
     */
    [USERS_PERMISSIONS.business.userStatistics.actions.interval]: {
      name: 'User Statistics: interval',
      description: 'Users with this permission can change time interval for the data displayed',
    },
    [USERS_PERMISSIONS.business.userStatistics.actions.listFilters]: {
      name: 'User Statistics: view filters',
      description: 'Users with this permission can see the list of filters available',
    },
    [USERS_PERMISSIONS.business.userStatistics.actions.editFilter]: {
      name: 'User Statistics: edit filters',
      description: 'Users with this permission can edit filters',
    },
    [USERS_PERMISSIONS.business.userStatistics.actions.addFilter]: {
      name: 'User Statistics: add filters',
      description: 'Users with this permission can add filters',
    },
    [USERS_PERMISSIONS.business.userStatistics.actions.userFilter]: {
      name: 'User Statistics: show filters',
      description: 'The filter is shown for users with this permission',
    },

    /**
     * Business Alarm Statistics
     */
    [USERS_PERMISSIONS.business.alarmStatistics.actions.interval]: {
      name: 'Alarm Statistics: interval',
      description: 'Users with this permission can change time interval for the data displayed',
    },
    [USERS_PERMISSIONS.business.alarmStatistics.actions.listFilters]: {
      name: 'Alarm Statistics: view filters',
      description: 'Users with this permission can see the list of filters available',
    },
    [USERS_PERMISSIONS.business.alarmStatistics.actions.editFilter]: {
      name: 'Alarm Statistics: edit filters',
      description: 'Users with this permission can edit filters',
    },
    [USERS_PERMISSIONS.business.alarmStatistics.actions.addFilter]: {
      name: 'Alarm Statistics: add filters',
      description: 'Users with this permission can add filters',
    },
    [USERS_PERMISSIONS.business.alarmStatistics.actions.userFilter]: {
      name: 'Alarm Statistics: show filters',
      description: 'The filter is shown for users with this permission',
    },

    /**
     * Business Availability
     */
    [USERS_PERMISSIONS.business.availability.actions.interval]: {
      name: 'Availability: interval',
      description: 'Users with this permission can change time interval for the data displayed',
    },
    [USERS_PERMISSIONS.business.availability.actions.listFilters]: {
      name: 'Availability: view filters',
      description: 'Users with this permission can see the list of filters available',
    },
    [USERS_PERMISSIONS.business.availability.actions.editFilter]: {
      name: 'Availability: edit filters',
      description: 'Users with this permission can edit filters',
    },
    [USERS_PERMISSIONS.business.availability.actions.addFilter]: {
      name: 'Availability: add filters',
      description: 'Users with this permission can add filters',
    },
    [USERS_PERMISSIONS.business.availability.actions.userFilter]: {
      name: 'Availability: show filters',
      description: 'The filter is shown for users with this permission',
    },
    [USERS_PERMISSIONS.business.availability.actions.exportAsCsv]: {
      name: 'Availability: Export as csv',
      description: 'Users with this permission can export availabilities as CSV file',
    },

    /**
     * Technical General Permissions
     */
    [USERS_PERMISSIONS.technical.view]: {
      name: 'Views',
      description: 'This permission defines the access to the list of views',
    },
    [USERS_PERMISSIONS.technical.privateView]: {
      name: 'Private views',
      description: 'This permission defines the access to the list of private views',
    },
    [USERS_PERMISSIONS.technical.role]: {
      name: 'Roles',
      description: 'This permission defines the access to the list of roles',
    },
    [USERS_PERMISSIONS.technical.permission]: {
      name: 'Rights',
      description: 'This permission defines the access to the list of Rights',
    },
    [USERS_PERMISSIONS.technical.user]: {
      name: 'Users',
      description: 'This permission defines the access to the list of users',
    },
    [USERS_PERMISSIONS.technical.parameters]: {
      name: 'Parameters',
      description: 'This permission defines the access to the Canopsis settings and parameters',
    },
    [USERS_PERMISSIONS.technical.broadcastMessage]: {
      name: 'Broadcast Messages',
      description: 'This permission defines the access to the Broadcast messages admin panel',
    },
    [USERS_PERMISSIONS.technical.playlist]: {
      name: 'Playlists',
      description: 'This permission defines the access to the Playlists settings',
    },
    [USERS_PERMISSIONS.technical.planning]: {
      name: 'Planning',
      description: 'This permission defines the access to the Planning and PBehavior settings',
    },
    [USERS_PERMISSIONS.technical.planningType]: {
      name: 'Planning type',
      description: 'This permission defines the access to the PBehavior types',
    },
    [USERS_PERMISSIONS.technical.planningReason]: {
      name: 'Planning reason',
      description: 'This permission defines the access to the PBehavior reasons',
    },
    [USERS_PERMISSIONS.technical.planningExceptions]: {
      name: 'Planning dates of exceptions',
      description: 'This permission defines the access to exception dates for PBehaviors',
    },
    [USERS_PERMISSIONS.technical.remediation]: {
      name: 'Remediation',
      description: 'This permission defines the access to the Remediation admin panel',
    },
    [USERS_PERMISSIONS.technical.remediationInstruction]: {
      name: 'Remediation instruction',
      description: 'This permission defines the access to the list of Instructions',
    },
    [USERS_PERMISSIONS.technical.remediationJob]: {
      name: 'Remediation job',
      description: 'This permission defines the access to the list of Jobs',
    },
    [USERS_PERMISSIONS.technical.remediationConfiguration]: {
      name: 'Remediation configuration',
      description: 'This permission defines the access to the Remediation configuration',
    },
    [USERS_PERMISSIONS.technical.remediationStatistic]: {
      name: 'Remediation statistics',
      description: 'This permission defines the access to the Remediation statistics',
    },
    [USERS_PERMISSIONS.technical.healthcheck]: {
      name: 'Healthcheck',
      description: 'This permission defines the access to the Healthcheck functionality',
    },
    [USERS_PERMISSIONS.technical.techmetrics]: {
      name: 'Tech metrics',
      description: 'This permission defines the access to the Tech metrics',
    },
    [USERS_PERMISSIONS.technical.engine]: {
      name: 'Engines',
      description: 'This permission defines the access to the Engines configuration',
    },
    [USERS_PERMISSIONS.technical.healthcheckStatus]: {
      name: 'Healthcheck status',
      description: 'The system healthcheck status is shown in the header for users with this permission',
    },
    [USERS_PERMISSIONS.technical.kpi]: {
      name: 'KPI',
      description: 'This permission defines the access to KPI metrics',
    },
    [USERS_PERMISSIONS.technical.kpiFilters]: {
      name: 'KPI Filters',
      description: 'This permission defines the access to filters for the KPI metrics',
    },
    [USERS_PERMISSIONS.technical.kpiRatingSettings]: {
      name: 'KPI Rating settings',
      description: 'This permission defines the access to the KPI Rating settings',
    },
    [USERS_PERMISSIONS.technical.kpiCollectionSettings]: {
      name: 'KPI Collection settings',
      description: 'This permission defines the access to the KPI Collection settings',
    },
    [USERS_PERMISSIONS.technical.map]: {
      name: 'Map editor',
      description: 'This permission defines the access to the map editor',
    },
    [USERS_PERMISSIONS.technical.shareToken]: {
      name: 'Share token',
      description: 'This permission defines the access to the Shared tokens settings',
    },
    [USERS_PERMISSIONS.technical.widgetTemplate]: {
      name: 'Widget templates',
      description: 'This permission defines the access to the Widget templates',
    },
    [USERS_PERMISSIONS.technical.maintenance]: {
      name: 'Maintenance mode',
      description: 'This permission defines the access to the Maintenance mode',
    },
    [USERS_PERMISSIONS.technical.tag]: {
      name: 'Tags management',
      description: 'This permission defines the access to the Tags management',
    },

    /**
     * Technical Exploitation Permissions
     */
    [USERS_PERMISSIONS.technical.exploitation.eventFilter]: {
      name: 'Exploitation: Event filters',
      description: 'This permission defines the access to the event filters',
    },
    [USERS_PERMISSIONS.technical.exploitation.pbehavior]: {
      name: 'Exploitation: Pbehaviors',
      description: 'This permission defines the access to the PBehavior events',
    },
    [USERS_PERMISSIONS.technical.exploitation.snmpRule]: {
      name: 'Exploitation: Snmp rules',
      description: 'This permission defines the access to the SNMP rules',
    },
    [USERS_PERMISSIONS.technical.exploitation.dynamicInfo]: {
      name: 'Exploitation: Dynamic information rules',
      description: 'This permission defines the access to the dynamic infos functionality',
    },
    [USERS_PERMISSIONS.technical.exploitation.metaAlarmRule]: {
      name: 'Exploitation: Meta alarm rules',
      description: 'This permission defines the access to the meta alarm rules and correlation',
    },
    [USERS_PERMISSIONS.technical.exploitation.scenario]: {
      name: 'Exploitation: Scenarios',
      description: 'This permission defines the access to the scenarios functionalitiy',
    },
    [USERS_PERMISSIONS.technical.exploitation.idleRules]: {
      name: 'Exploitation: Idle rules',
      description: 'This permission defines the access to the idle rules',
    },
    [USERS_PERMISSIONS.technical.exploitation.flappingRules]: {
      name: 'Exploitation: Flapping rules',
      description: 'This permission defines the access to the flapping rules',
    },
    [USERS_PERMISSIONS.technical.exploitation.resolveRules]: {
      name: 'Exploitation: Resolve rules',
      description: 'This permission defines the access to the resolve rules',
    },
    [USERS_PERMISSIONS.technical.exploitation.declareTicketRule]: {
      name: 'Exploitation: Declare ticket rules',
      description: 'This permission defines the access to the ticket declaration functionality',
    },
    [USERS_PERMISSIONS.technical.exploitation.linkRule]: {
      name: 'Exploitation: Link rules',
      description: 'This permission defines the access to the links and link rules',
    },

    /**
     * Technical Notification Permissions
     */
    [USERS_PERMISSIONS.technical.notification.instructionStats]: {
      name: 'Notifications: Instructions stats',
      description: 'This permission defines the access to the notifications associated with instructions statistics',
    },

    /**
     * Technical Profile Permissions
     */
    [USERS_PERMISSIONS.technical.profile.corporatePattern]: {
      name: 'Profile: Corporate patterns',
      description: 'This permission defines the access to the corporate patterns functionality',
    },
    [USERS_PERMISSIONS.technical.profile.theme]: {
      name: 'Themes',
      description: 'This permission defines the access to the theme colors',
    },

    /**
     * API Permissions
     */
    [USERS_PERMISSIONS.api.general.acl]: {
      name: 'Roles, permissions, users',
      description: 'Access to API route to CRUD roles, permissions and users',
    },
    [USERS_PERMISSIONS.api.general.appInfoRead]: {
      name: 'Read app information',
      description: 'Access to API route to read app information',
    },
    [USERS_PERMISSIONS.api.general.alarmRead]: {
      name: 'Read alarms',
      description: 'Access to API route to read alarms',
    },
    [USERS_PERMISSIONS.api.general.alarmUpdate]: {
      name: 'Update alarms',
      description: 'Access to API route to update alarms',
    },
    [USERS_PERMISSIONS.api.general.entity]: {
      name: 'Entity',
      description: 'Access to API route to CRUD entities',
    },
    [USERS_PERMISSIONS.api.general.entityservice]: {
      name: 'Entity service',
      description: 'Access to API route to CRUD services',
    },
    [USERS_PERMISSIONS.api.general.entitycategory]: {
      name: 'Entity categories',
      description: 'Access to API route to CRUD entity categories',
    },
    [USERS_PERMISSIONS.api.general.event]: {
      name: 'Event',
      description: 'Access to API route for events',
    },
    [USERS_PERMISSIONS.api.general.view]: {
      name: 'Views',
      description: 'Access to API route to CRUD views',
    },
    [USERS_PERMISSIONS.api.general.viewgroup]: {
      name: 'View groups',
      description: 'Access to API route to CRUD view groups',
    },
    [USERS_PERMISSIONS.api.general.privateViewGroups]: {
      name: 'Private view groups',
      description: 'Access to API route to CRUD private view groups',
    },
    [USERS_PERMISSIONS.api.general.userInterfaceUpdate]: {
      name: 'Update user interface',
      description: 'Access to API route to update user interface',
    },
    [USERS_PERMISSIONS.api.general.userInterfaceDelete]: {
      name: 'Delete user interface',
      description: 'Access to API route to delete user interface',
    },
    [USERS_PERMISSIONS.api.general.datastorageRead]: {
      name: 'Data storage settings read',
      description: 'Access to API route to read data storage settings',
    },
    [USERS_PERMISSIONS.api.general.datastorageUpdate]: {
      name: 'Data storage settings update',
      description: 'Access to API route to edit data storage settings',
    },
    [USERS_PERMISSIONS.api.general.associativeTable]: {
      name: 'Associative table',
      description: 'Access to API route with some associated data storage (dynamic infos templates, etc)',
    },
    [USERS_PERMISSIONS.api.general.stateSettings]: {
      name: 'State settings',
      description: 'Access to API route to state settings',
    },
    [USERS_PERMISSIONS.api.general.files]: {
      name: 'File',
      description: 'Access to API route to CRUD files',
    },
    [USERS_PERMISSIONS.api.general.healthcheck]: {
      name: 'Healthcheck',
      description: 'Access to API route for healthcheck',
    },
    [USERS_PERMISSIONS.api.general.techmetrics]: {
      name: 'Tech Metrics',
      description: 'Access to API route to tech metrics',
    },
    [USERS_PERMISSIONS.api.general.contextgraph]: {
      name: 'Context graph import',
      description: 'Access to API route for the context graph import',
    },
    [USERS_PERMISSIONS.api.general.broadcastMessage]: {
      name: 'Broadcast Message',
      description: 'Access to API route for broadcast messages',
    },
    [USERS_PERMISSIONS.api.general.junit]: {
      name: 'JUnit',
      description: 'Access to API route to JUnit API',
    },
    [USERS_PERMISSIONS.api.general.notifications]: {
      name: 'Notification settings',
      description: 'Access to API route for notification settings',
    },
    [USERS_PERMISSIONS.api.general.metrics]: {
      name: 'Metrics',
      description: 'Access to API route for metrics',
    },
    [USERS_PERMISSIONS.api.general.metricsSettings]: {
      name: 'Metrics settings',
      description: 'Access to API route for metric settings',
    },
    [USERS_PERMISSIONS.api.general.ratingSettings]: {
      name: 'Rating settings',
      description: 'Access to API route for rating settings',
    },
    [USERS_PERMISSIONS.api.general.filter]: {
      name: 'KPI filters',
      description: 'Access to API route to KPI filters',
    },
    [USERS_PERMISSIONS.api.general.corporatePattern]: {
      name: 'Corporate patterns',
      description: 'Access to API route for corporate patterns',
    },
    [USERS_PERMISSIONS.api.general.exportConfigurations]: {
      name: 'Export configurations',
      description: 'Access to API route to export configuration',
    },
    [USERS_PERMISSIONS.api.general.map]: {
      name: 'Map',
      description: 'Access to API route to CRUD maps',
    },
    [USERS_PERMISSIONS.api.general.shareToken]: {
      name: 'Share token',
      description: 'Access to API route to CRUD shared tokens',
    },
    [USERS_PERMISSIONS.api.general.declareTicketExecution]: {
      name: 'Run declare ticket rules',
      description: 'Access to API route to run declare ticket rules',
    },
    [USERS_PERMISSIONS.api.general.widgetTemplate]: {
      name: 'Widget templates',
      description: 'Access to API route to CRUD widget templates',
    },
    [USERS_PERMISSIONS.api.general.maintenance]: {
      name: 'Maintenance mode',
      description: 'Access to API route to the maintenance mode',
    },
    [USERS_PERMISSIONS.api.general.theme]: {
      name: 'Themes',
      description: 'Access to API route to the themes',
    },

    [USERS_PERMISSIONS.api.rules.action]: {
      name: 'Actions',
      description: 'Users with this permission can CRUD actions by API',
    },
    [USERS_PERMISSIONS.api.rules.dynamicinfos]: {
      name: 'Dynamic infos',
      description: 'Users with this permission can CRUD dynamic infos by API',
    },
    [USERS_PERMISSIONS.api.rules.eventFilter]: {
      name: 'Event filter',
      description: 'Users with this permission can CRUD event filters by API',
    },
    [USERS_PERMISSIONS.api.rules.idleRule]: {
      name: 'Idle rule',
      description: 'Users with this permission can CRUD idle rules by API',
    },
    [USERS_PERMISSIONS.api.rules.metaalarmrule]: {
      name: 'Meta alarm rule',
      description: 'Users with this permission can CRUD meta alarm rules by API',
    },
    [USERS_PERMISSIONS.api.rules.playlist]: {
      name: 'Playlists',
      description: 'Users with this permission can CRUD playlists by API',
    },
    [USERS_PERMISSIONS.api.rules.flappingRule]: {
      name: 'Flapping rule',
      description: 'Users with this permission can CRUD flapping rules by API',
    },
    [USERS_PERMISSIONS.api.rules.resolveRule]: {
      name: 'Resolve rule',
      description: 'Users with this permission can CRUD resolve rules by API',
    },
    [USERS_PERMISSIONS.api.rules.snmpRule]: {
      name: 'SNMP rule',
      description: 'Users with this permission can CRUD SNMP rules by API',
    },
    [USERS_PERMISSIONS.api.rules.snmpMib]: {
      name: 'SNMP MIB',
      description: 'Users with this permission can CRUD SNMP MIB by API',
    },
    [USERS_PERMISSIONS.api.rules.declareTicketRule]: {
      name: 'Declare ticket rule',
      description: 'Users with this permission can CRUD declare ticket rules by API',
    },
    [USERS_PERMISSIONS.api.rules.linkRule]: {
      name: 'Link rule',
      description: 'Users with this permission can CRUD links and link rules by API',
    },

    [USERS_PERMISSIONS.api.remediation.instruction]: {
      name: 'Instructions',
      description: 'Users with this permission can CRUD instructions by API',
    },
    [USERS_PERMISSIONS.api.remediation.jobConfig]: {
      name: 'Job configs',
      description: 'Users with this permission can CRUD job configurations by API',
    },
    [USERS_PERMISSIONS.api.remediation.job]: {
      name: 'Jobs',
      description: 'Users with this permission can CRUD jobs by API',
    },
    [USERS_PERMISSIONS.api.remediation.execution]: {
      name: 'Runs instructions',
      description: 'Users with this permission can run instructions by API',
    },
    [USERS_PERMISSIONS.api.remediation.instructionApprove]: {
      name: 'Instruction approve',
      description: 'Users with this permission can approve instructions by API',
    },
    [USERS_PERMISSIONS.api.remediation.messageRateStatsRead]: {
      name: 'Message rate statistics',
      description: 'Users with this permission can access message rate statistics by API',
    },

    [USERS_PERMISSIONS.api.pbehavior.pbehavior]: {
      name: 'PBehaviors',
      description: 'Users with this permission can CRUD PBehavior events dates by API',
    },
    [USERS_PERMISSIONS.api.pbehavior.pbehaviorException]: {
      name: 'PBehavior exceptions',
      description: 'Users with this permission can CRUD PBehavior exceptions dates by API',
    },
    [USERS_PERMISSIONS.api.pbehavior.pbehaviorReason]: {
      name: 'PBehavior reasons',
      description: 'Users with this permission can CRUD PBehavior reasons dates by API',
    },
    [USERS_PERMISSIONS.api.pbehavior.pbehaviorType]: {
      name: 'PBehavior types',
      description: 'Users with this permission can CRUD PBehavior types dates by API',
    },
  },
};
