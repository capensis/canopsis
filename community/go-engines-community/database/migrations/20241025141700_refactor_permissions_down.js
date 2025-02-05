const replacePermissions = {
    common_entityComment: [
        "common_entityCommentsList",
        "common_editEntityComment",
        "common_createEntityComment",
    ],
    listalarm_filter: [
        "listalarm_listFilters",
        "listalarm_editFilter",
        "listalarm_addFilter",
    ],
    listalarm_remediationInstructionsFilter: [
        "listalarm_listRemediationInstructionsFilters",
        "listalarm_editRemediationInstructionsFilter",
        "listalarm_addRemediationInstructionsFilter",
    ],
    listalarm_bookmark: [
        "listalarm_removeBookmark",
        "listalarm_addBookmark",
    ],
    serviceweather_filter: [
        "serviceweather_listFilters",
        "serviceweather_editFilter",
        "serviceweather_addFilter",
    ],
    crudcontext_filter: [
        "crudcontext_listFilters",
        "crudcontext_editFilter",
        "crudcontext_addFilter",
    ],
    map_filter: [
        "map_listFilters",
        "map_editFilter",
        "map_addFilter",
    ],
    userStatistics_filter: [
        "userStatistics_listFilters",
        "userStatistics_editFilter",
        "userStatistics_addFilter",
    ],
    alarmStatistics_filter: [
        "alarmStatistics_listFilters",
        "alarmStatistics_editFilter",
        "alarmStatistics_addFilter",
    ],
    barchart_filter: [
        "barchart_listFilters",
        "barchart_editFilter",
        "barchart_addFilter",
    ],
    linechart_filter: [
        "linechart_listFilters",
        "linechart_editFilter",
        "linechart_addFilter",
    ],
    piechart_filter: [
        "piechart_listFilters",
        "piechart_editFilter",
        "piechart_addFilter",
    ],
    numbers_filter: [
        "numbers_listFilters",
        "numbers_editFilter",
        "numbers_addFilter",
    ],
    availability_filter: [
        "availability_listFilters",
        "availability_editFilter",
        "availability_addFilter",
    ],
};

for (const perm in replacePermissions) {
    let set = {};
    for (const permToReplace of replacePermissions[perm]) {
        set["permissions." + permToReplace] = 1;
    }

    db.role.updateOne(
        {["permissions." + perm]: 1},
        {
            $set: set,
            $unset: {["permissions." + perm]: ""},
        }
    );
}

const updatedPerms = [
    {
        "_id": "models_userview",
        "description": "Userviews",
        "name": "models_userview",
        "type": "CRUD"
    },
    {
        "_id": "models_role",
        "description": "Roles",
        "name": "models_role",
        "type": "CRUD"
    },
    {
        "_id": "models_permission",
        "description": "Rights",
        "name": "models_permission",
        "type": "CRUD"
    },
    {
        "_id": "models_user",
        "description": "Users",
        "name": "models_user",
        "type": "CRUD"
    },
    {
        "_id": "models_parameters",
        "description": "Parameters",
        "name": "models_parameters",
        "type": "RW"
    },
    {
        "_id": "models_broadcastMessage",
        "description": "Broadcast Messages",
        "name": "models_broadcastMessage",
        "type": "CRUD"
    },
    {
        "_id": "models_playlist",
        "description": "Playlists",
        "name": "models_playlist",
        "type": "CRUD"
    },
    {
        "_id": "models_planning",
        "description": "Planning",
        "name": "models_planning",
        "type": "CRUD"
    },
    {
        "_id": "models_planningType",
        "description": "Planning type",
        "name": "models_planningType",
        "type": "CRUD"
    },
    {
        "_id": "models_planningReason",
        "description": "Planning reason",
        "name": "models_planningReason",
        "type": "CRUD"
    },
    {
        "_id": "models_planningExceptions",
        "description": "Planning dates of exceptions",
        "name": "models_planningExceptions",
        "type": "CRUD"
    },
    {
        "_id": "models_remediation",
        "description": "Remediation",
        "name": "models_remediation",
        "type": "CRUD"
    },
    {
        "_id": "models_remediationInstruction",
        "description": "Remediation instruction",
        "name": "models_remediationInstruction",
        "type": "CRUD"
    },
    {
        "_id": "models_remediationJob",
        "description": "Remediation job",
        "name": "models_remediationJob",
        "type": "CRUD"
    },
    {
        "_id": "models_remediationConfiguration",
        "description": "Remediation configuration",
        "name": "models_remediationConfiguration",
        "type": "CRUD"
    },
    {
        "_id": "models_engine",
        "description": "Engines",
        "name": "models_engine"
    },
    {
        "_id": "models_exploitation_eventFilter",
        "description": "Exploitation: Event filter",
        "name": "models_exploitation_eventFilter",
        "type": "CRUD"
    },
    {
        "_id": "models_exploitation_pbehavior",
        "description": "Exploitation: Pbehavior",
        "name": "models_exploitation_pbehavior",
        "type": "CRUD"
    },
    {
        "_id": "models_exploitation_snmpRule",
        "description": "Exploitation: Snmp rule",
        "name": "models_exploitation_snmpRule",
        "type": "CRUD"
    },
    {
        "_id": "models_exploitation_dynamicInfo",
        "description": "Exploitation: Dynamic information rule",
        "name": "models_exploitation_dynamicInfo",
        "type": "CRUD"
    },
    {
        "_id": "models_exploitation_metaAlarmRule",
        "description": "Exploitation: Meta alarm rule",
        "name": "models_exploitation_metaAlarmRule",
        "type": "CRUD"
    },
    {
        "_id": "models_exploitation_scenario",
        "description": "Exploitation: Scenario",
        "name": "models_exploitation_scenario",
        "type": "CRUD"
    },
    {
        "_id": "models_exploitation_idleRules",
        "description": "Exploitation: Idle rule",
        "name": "models_exploitation_idleRules",
        "type": "CRUD"
    },
    {
        "_id": "models_exploitation_flappingRules",
        "description": "Exploitation: Flapping rule",
        "name": "models_exploitation_flappingRules",
        "type": "CRUD"
    },
    {
        "_id": "models_exploitation_resolveRules",
        "description": "Exploitation: Resolve rule",
        "name": "models_exploitation_resolveRules",
        "type": "CRUD"
    },
    {
        "_id": "models_exploitation_declareTicketRule",
        "description": "Exploitation: Declare ticket rule",
        "name": "models_exploitation_declareTicketRule",
        "type": "CRUD"
    },
    {
        "_id": "models_exploitation_linkRule",
        "description": "Exploitation: Link rule",
        "name": "models_exploitation_linkRule",
        "type": "CRUD"
    },
    {
        "_id": "models_healthcheck",
        "description": "Healthcheck",
        "name": "models_healthcheck"
    },
    {
        "_id": "models_healthcheckStatus",
        "description": "Healthcheck status",
        "name": "models_healthcheckStatus"
    },
    {
        "_id": "models_kpi",
        "description": "KPI",
        "name": "models_kpi"
    },
    {
        "_id": "models_kpiFilters",
        "description": "KPI Filters",
        "name": "models_kpiFilters",
        "type": "CRUD"
    },
    {
        "_id": "models_kpiRatingSettings",
        "description": "KPI Rating settings",
        "name": "models_kpiRatingSettings",
        "type": "CRUD"
    },
    {
        "_id": "models_kpiCollectionSettings",
        "description": "KPI collection settings",
        "name": "models_kpiCollectionSettings",
        "type": "CRUD"
    },
    {
        "_id": "models_notification_instructionStats",
        "description": "Notifications: Instructions stats",
        "name": "models_notification_instructionStats",
        "type": "CRUD"
    },
    {
        "_id": "models_profile_corporatePattern",
        "description": "Profile: Corporate patterns",
        "name": "models_profile_corporatePattern",
        "type": "CRUD"
    },
    {
        "_id": "models_map",
        "description": "Map",
        "name": "models_map",
        "type": "CRUD"
    },
    {
        "_id": "models_shareToken",
        "description": "Share token",
        "name": "models_shareToken",
        "type": "CRUD"
    },
    {
        "_id": "models_techmetrics",
        "description": "Tech metrics",
        "name": "models_techmetrics"
    },
    {
        "_id": "models_widgetTemplate",
        "description": "Widget templates",
        "name": "models_widgetTemplate",
        "type": "CRUD"
    },
    {
        "_id": "models_tag",
        "description": "Alarm tags",
        "name": "models_tag",
        "type": "CRUD"
    },
    {
        "_id": "models_eventsRecord",
        "description": "Events record",
        "name": "models_eventsRecord"
    },
    {
        "_id": "listalarm_ack",
        "description": "Rights on listalarm: ack",
        "name": "listalarm_ack"
    },
    {
        "_id": "listalarm_fastAck",
        "description": "Rights on listalarm: fast ack",
        "name": "listalarm_fastAck"
    },
    {
        "_id": "listalarm_cancelAck",
        "description": "Rights on listalarm: cancel ack",
        "name": "listalarm_cancelAck"
    },
    {
        "_id": "listalarm_pbehavior",
        "description": "Rights on listalarm: pbehavior",
        "name": "listalarm_pbehavior"
    },
    {
        "_id": "listalarm_fastPbehavior",
        "description": "Rights on listalarm: fast pbehavior",
        "name": "listalarm_fastPbehavior"
    },
    {
        "_id": "listalarm_snoozeAlarm",
        "description": "Rights on listalarm: snooze alarm",
        "name": "listalarm_snoozeAlarm"
    },
    {
        "_id": "listalarm_listPbehavior",
        "description": "Rights on listalarm: list pbehavior",
        "name": "listalarm_listPbehavior"
    },
    {
        "_id": "listalarm_declareanIncident",
        "description": "Rights on listalarm: declare an incident",
        "name": "listalarm_declareanIncident"
    },
    {
        "_id": "listalarm_assignTicketNumber",
        "description": "Rights on listalarm: assign ticket number",
        "name": "listalarm_assignTicketNumber"
    },
    {
        "_id": "listalarm_removeAlarm",
        "description": "Rights on listalarm: remove alarm",
        "name": "listalarm_removeAlarm"
    },
    {
        "_id": "listalarm_fastRemoveAlarm",
        "description": "Rights on listalarm: fast remove alarm",
        "name": "listalarm_fastRemoveAlarm"
    },
    {
        "_id": "listalarm_unCancel",
        "description": "Rights on listalarm: uncancel alarm",
        "name": "listalarm_unCancel"
    },
    {
        "_id": "listalarm_changeState",
        "description": "Rights on listalarm: change state",
        "name": "listalarm_changeState"
    },
    {
        "_id": "listalarm_history",
        "description": "Rights on listalarm: history",
        "name": "listalarm_history"
    },
    {
        "_id": "listalarm_groupRequest",
        "description": "Rights on listalarm: group request",
        "name": "listalarm_groupRequest"
    },
    {
        "_id": "listalarm_manualMetaAlarmGroup",
        "description": "Rights on listalarm: Manual meta alarm actions",
        "name": "listalarm_manualMetaAlarmGroup"
    },
    {
        "_id": "listalarm_metaAlarmGroup",
        "description": "Rights on listalarm: Meta alarm actions",
        "name": "listalarm_metaAlarmGroup"
    },
    {
        "_id": "listalarm_comment",
        "description": "Rights on listalarm: Access to 'Comment' action",
        "name": "listalarm_comment"
    },
    {
        "_id": "listalarm_listFilters",
        "description": "Rights on listalarm: list filters",
        "name": "listalarm_listFilters"
    },
    {
        "_id": "listalarm_editFilter",
        "description": "Rights on listalarm: edit filters",
        "name": "listalarm_editFilter"
    },
    {
        "_id": "listalarm_addFilter",
        "description": "Rights on listalarm: add filter",
        "name": "listalarm_addFilter"
    },
    {
        "_id": "listalarm_userFilter",
        "description": "Rights on listalarm: User filter",
        "name": "listalarm_userFilter"
    },
    {
        "_id": "listalarm_addBookmark",
        "description": "Rights on listalarm: Access to 'Add bookmark' action",
        "name": "listalarm_addBookmark"
    },
    {
        "_id": "listalarm_filterByBookmark",
        "description": "Rights on listalarm: Access to 'Filter by bookmark' action",
        "name": "listalarm_filterByBookmark"
    },
    {
        "_id": "listalarm_removeBookmark",
        "description": "Rights on listalarm: Access to 'Remove bookmark' action",
        "name": "listalarm_removeBookmark"
    },
    {
        "_id": "listalarm_listRemediationInstructionsFilters",
        "description": "Rights on listalarm: Access to 'Remediation instructions filters: list' action",
        "name": "listalarm_listRemediationInstructionsFilters"
    },
    {
        "_id": "listalarm_editRemediationInstructionsFilter",
        "description": "Rights on listalarm: Access to 'Remediation instructions filters: edit' action",
        "name": "listalarm_editRemediationInstructionsFilter"
    },
    {
        "_id": "listalarm_addRemediationInstructionsFilter",
        "description": "Rights on listalarm: Access to 'Remediation instructions filters: add' action",
        "name": "listalarm_addRemediationInstructionsFilter"
    },
    {
        "_id": "listalarm_userRemediationInstructionsFilter",
        "description": "Rights on listalarm: Access to 'Remediation instructions filters: user' action",
        "name": "listalarm_userRemediationInstructionsFilter"
    },
    {
        "_id": "listalarm_links",
        "description": "Rights on listalarm: Access to Links",
        "name": "listalarm_links"
    },
    {
        "_id": "listalarm_correlation",
        "description": "Rights on listalarm: Access to 'Correlation' action",
        "name": "listalarm_correlation"
    },
    {
        "_id": "listalarm_executeInstruction",
        "description": "Rights on listalarm: Access to 'Execute instruction' action",
        "name": "listalarm_executeInstruction"
    },
    {
        "_id": "listalarm_category",
        "description": "Rights on listalarm: Access to 'Category' action",
        "name": "listalarm_category"
    },
    {
        "_id": "listalarm_exportAsCsv",
        "description": "Rights on listalarm: Access to 'Export as csv' action",
        "name": "listalarm_exportAsCsv"
    },
    {
        "_id": "serviceweather_entityAck",
        "description": "Service weather: Access to 'Ack' action",
        "name": "serviceweather_entityAck"
    },
    {
        "_id": "serviceweather_entityAssocTicket",
        "description": "Service weather: Access to 'Assoc Ticket' action",
        "name": "serviceweather_entityAssocTicket"
    },
    {
        "_id": "serviceweather_entityComment",
        "description": "Service weather: Access to 'Comment' action",
        "name": "serviceweather_entityComment"
    },
    {
        "_id": "serviceweather_entityValidate",
        "description": "Service weather: Access to 'Validate' action",
        "name": "serviceweather_entityValidate"
    },
    {
        "_id": "serviceweather_entityInvalidate",
        "description": "Service weather: Access to 'Invalidate' action",
        "name": "serviceweather_entityInvalidate"
    },
    {
        "_id": "serviceweather_entityPause",
        "description": "Service weather: Access to 'Pause' action",
        "name": "serviceweather_entityPause"
    },
    {
        "_id": "serviceweather_entityPlay",
        "description": "Service weather: Access to 'Play' action",
        "name": "serviceweather_entityPlay"
    },
    {
        "_id": "serviceweather_entityCancel",
        "description": "Service weather: Access to 'Cancel' action",
        "name": "serviceweather_entityCancel"
    },
    {
        "_id": "serviceweather_entityManagePbehaviors",
        "description": "Service weather: Access to pbehaviors management",
        "name": "serviceweather_entityManagePbehaviors"
    },
    {
        "_id": "serviceweather_entityLinks",
        "description": "Service weather: Access to Links",
        "name": "serviceweather_entityLinks"
    },
    {
        "_id": "serviceweather_entityDeclareTicket",
        "description": "Service weather: Access to 'Declare Ticket' action",
        "name": "serviceweather_entityDeclareTicket"
    },
    {
        "_id": "serviceweather_moreInfos",
        "description": "Service weather: Access to 'More infos' modal",
        "name": "serviceweather_moreInfos"
    },
    {
        "_id": "serviceweather_alarmsList",
        "description": "Service weather: Access to 'Alarms list' modal",
        "name": "serviceweather_alarmsList"
    },
    {
        "_id": "serviceweather_pbehaviorList",
        "description": "Service weather: Access to watcher pbehavior list",
        "name": "serviceweather_pbehaviorList"
    },
    {
        "_id": "serviceweather_listFilters",
        "description": "Rights on service weather: List filters",
        "name": "serviceweather_listFilters"
    },
    {
        "_id": "serviceweather_editFilter",
        "description": "Rights on service weather: Edit filters",
        "name": "serviceweather_editFilter"
    },
    {
        "_id": "serviceweather_addFilter",
        "description": "Rights on service weather: Add filter",
        "name": "serviceweather_addFilter"
    },
    {
        "_id": "serviceweather_userFilter",
        "description": "Rights on service weather: User filter",
        "name": "serviceweather_userFilter"
    },
    {
        "_id": "serviceweather_category",
        "description": "Rights on service weather: Access to 'Category' action",
        "name": "serviceweather_category"
    },
    {
        "_id": "serviceweather_executeInstruction",
        "description": "Service weather: Access to execute instruction",
        "name": "serviceweather_executeInstruction"
    },
    {
        "_id": "crudcontext_createEntity",
        "description": "Rights on context: create entity",
        "name": "crudcontext_createEntity"
    },
    {
        "_id": "crudcontext_edit",
        "description": "Rights on context: edit entity",
        "name": "crudcontext_edit"
    },
    {
        "_id": "crudcontext_duplicate",
        "description": "Rights on context: duplicate entity",
        "name": "crudcontext_duplicate"
    },
    {
        "_id": "crudcontext_delete",
        "description": "Rights on context: delete entity",
        "name": "crudcontext_delete"
    },
    {
        "_id": "crudcontext_pbehavior",
        "description": "Rights on context: pbehavior",
        "name": "crudcontext_pbehavior"
    },
    {
        "_id": "crudcontext_listPbehavior",
        "description": "Rights on context: List Pbehaviors",
        "name": "crudcontext_listPbehavior"
    },
    {
        "_id": "crudcontext_deletePbehavior",
        "description": "Rights on context: Delete Pbehaviors",
        "name": "crudcontext_deletePbehavior"
    },
    {
        "_id": "crudcontext_listFilters",
        "description": "Rights on context: List filters",
        "name": "crudcontext_listFilters"
    },
    {
        "_id": "crudcontext_editFilter",
        "description": "Rights on context: Edit filter",
        "name": "crudcontext_editFilter"
    },
    {
        "_id": "crudcontext_addFilter",
        "description": "Rights on context: Add filter",
        "name": "crudcontext_addFilter"
    },
    {
        "_id": "crudcontext_userFilter",
        "description": "Rights on context: User filter",
        "name": "crudcontext_userFilter"
    },
    {
        "_id": "crudcontext_category",
        "description": "Rights on context: Access to 'Category' action",
        "name": "crudcontext_category"
    },
    {
        "_id": "crudcontext_exportAsCsv",
        "description": "Rights on context: Export as csv",
        "name": "crudcontext_exportAsCsv"
    },
    {
        "_id": "crudcontext_massEnable",
        "description": "Rights on context: Mass enable",
        "name": "crudcontext_massEnable"
    },
    {
        "_id": "crudcontext_massDisable",
        "description": "Rights on context: Mass disable",
        "name": "crudcontext_massDisable"
    },
    {
        "_id": "common_variablesHelp",
        "description": "Access to available variables list",
        "name": "common_variablesHelp"
    },
    {
        "_id": "counter_alarmsList",
        "description": "Counter: Access to 'Alarms list' modal",
        "name": "counter_alarmsList"
    },
    {
        "_id": "testingweather_alarmsList",
        "description": "Testing weather: Access to 'Alarms list' modal",
        "name": "testingweather_alarmsList"
    },
    {
        "_id": "map_alarmsList",
        "description": "Map: Access to 'Alarms list' modal",
        "name": "map_alarmsList"
    },
    {
        "_id": "map_listFilters",
        "description": "Rights on map: List filters",
        "name": "map_listFilters"
    },
    {
        "_id": "map_editFilter",
        "description": "Rights on map: Edit filters",
        "name": "map_editFilter"
    },
    {
        "_id": "map_addFilter",
        "description": "Rights on map: Add filter",
        "name": "map_addFilter"
    },
    {
        "_id": "map_userFilter",
        "description": "Rights on map: User filter",
        "name": "map_userFilter"
    },
    {
        "_id": "map_category",
        "description": "Rights on map: Access to 'Category' action",
        "name": "map_category"
    },
    {
        "_id": "userStatistics_interval",
        "description": "Rights on user statistics: Interval",
        "name": "userStatistics_interval"
    },
    {
        "_id": "userStatistics_listFilters",
        "description": "Rights on user statistics: List filters",
        "name": "userStatistics_listFilters"
    },
    {
        "_id": "userStatistics_editFilter",
        "description": "Rights on user statistics: Edit filters",
        "name": "userStatistics_editFilter"
    },
    {
        "_id": "userStatistics_addFilter",
        "description": "Rights on user statistics: Add filter",
        "name": "userStatistics_addFilter"
    },
    {
        "_id": "userStatistics_userFilter",
        "description": "Rights on user statistics: User filter",
        "name": "userStatistics_userFilter"
    },
    {
        "_id": "alarmStatistics_interval",
        "description": "Rights on alarm statistics: Interval",
        "name": "alarmStatistics_interval"
    },
    {
        "_id": "alarmStatistics_listFilters",
        "description": "Rights on alarm statistics: List filters",
        "name": "alarmStatistics_listFilters"
    },
    {
        "_id": "alarmStatistics_editFilter",
        "description": "Rights on alarm statistics: Edit filters",
        "name": "alarmStatistics_editFilter"
    },
    {
        "_id": "alarmStatistics_addFilter",
        "description": "Rights on alarm statistics: Add filter",
        "name": "alarmStatistics_addFilter"
    },
    {
        "_id": "alarmStatistics_userFilter",
        "description": "Rights on alarm statistics: User filter",
        "name": "alarmStatistics_userFilter"
    },
    {
        "_id": "models_remediationStatistic",
        "description": "Remediation statistics",
        "name": "models_remediationStatistic"
    },
    {
        "_id": "barchart_interval",
        "description": "Barchart interval",
        "name": "barchart_interval"
    },
    {
        "_id": "barchart_sampling",
        "description": "Barchart sampling",
        "name": "barchart_sampling"
    },
    {
        "_id": "barchart_listFilters",
        "description": "Barchart list filters",
        "name": "barchart_listFilters"
    },
    {
        "_id": "barchart_editFilter",
        "description": "Barchart edit filter",
        "name": "barchart_editFilter"
    },
    {
        "_id": "barchart_addFilter",
        "description": "Barchart add filter",
        "name": "barchart_addFilter"
    },
    {
        "_id": "barchart_userFilter",
        "description": "Barchart user filter",
        "name": "barchart_userFilter"
    },
    {
        "_id": "linechart_interval",
        "description": "Linechart interval",
        "name": "linechart_interval"
    },
    {
        "_id": "linechart_sampling",
        "description": "Linechart sampling",
        "name": "linechart_sampling"
    },
    {
        "_id": "linechart_listFilters",
        "description": "Linechart list filters",
        "name": "linechart_listFilters"
    },
    {
        "_id": "linechart_editFilter",
        "description": "Linechart edit filter",
        "name": "linechart_editFilter"
    },
    {
        "_id": "linechart_addFilter",
        "description": "Linechart add filter",
        "name": "linechart_addFilter"
    },
    {
        "_id": "linechart_userFilter",
        "description": "Linechart user filter",
        "name": "linechart_userFilter"
    },
    {
        "_id": "piechart_interval",
        "description": "Piechart interval",
        "name": "piechart_interval"
    },
    {
        "_id": "piechart_sampling",
        "description": "Piechart sampling",
        "name": "piechart_sampling"
    },
    {
        "_id": "piechart_listFilters",
        "description": "Piechart list filters",
        "name": "piechart_listFilters"
    },
    {
        "_id": "piechart_editFilter",
        "description": "Piechart edit filter",
        "name": "piechart_editFilter"
    },
    {
        "_id": "piechart_addFilter",
        "description": "Piechart add filter",
        "name": "piechart_addFilter"
    },
    {
        "_id": "piechart_userFilter",
        "description": "Piechart user filter",
        "name": "piechart_userFilter"
    },
    {
        "_id": "numbers_interval",
        "description": "Numbers interval",
        "name": "numbers_interval"
    },
    {
        "_id": "numbers_sampling",
        "description": "Numbers sampling",
        "name": "numbers_sampling"
    },
    {
        "_id": "numbers_listFilters",
        "description": "Numbers list filters",
        "name": "numbers_listFilters"
    },
    {
        "_id": "numbers_editFilter",
        "description": "Numbers edit filter",
        "name": "numbers_editFilter"
    },
    {
        "_id": "numbers_addFilter",
        "description": "Numbers add filter",
        "name": "numbers_addFilter"
    },
    {
        "_id": "numbers_userFilter",
        "description": "Numbers user filter",
        "name": "numbers_userFilter"
    },
    {
        "_id": "models_maintenance",
        "description": "Maintenance",
        "name": "models_maintenance"
    },
    {
        "_id": "models_profile_color_theme",
        "description": "Color theme",
        "name": "models_profile_color_theme",
        "type": "CRUD"
    },
    {
        "_id": "models_stateSetting",
        "description": "State settings",
        "name": "models_stateSetting",
        "type": "CRUD"
    },
    {
        "_id": "models_icon",
        "description": "Icons",
        "name": "models_icon",
        "type": "CRUD"
    },
    {
        "_id": "availability_interval",
        "description": "Availability interval",
        "name": "availability_interval"
    },
    {
        "_id": "availability_listFilters",
        "description": "Availability list filters",
        "name": "availability_listFilters"
    },
    {
        "_id": "availability_editFilter",
        "description": "Availability edit filter",
        "name": "availability_editFilter"
    },
    {
        "_id": "availability_addFilter",
        "description": "Availability add filter",
        "name": "availability_addFilter"
    },
    {
        "_id": "availability_userFilter",
        "description": "Availability user filter",
        "name": "availability_userFilter"
    },
    {
        "_id": "availability_exportAsCsv",
        "description": "Availability export as csv",
        "name": "availability_exportAsCsv"
    },
    {
        "_id": "common_entityCommentsList",
        "description": "Entity comments list",
        "name": "common_entityCommentsList"
    },
    {
        "_id": "common_createEntityComment",
        "description": "Create entity comment",
        "name": "common_createEntityComment"
    },
    {
        "_id": "common_editEntityComment",
        "description": "Edit entity comment",
        "name": "common_editEntityComment"
    },
    {
        "_id": "api_alarm_read",
        "description": "Read alarms",
        "name": "api_alarm_read"
    },
    {
        "_id": "api_alarm_update",
        "description": "Update alarms",
        "name": "api_alarm_update"
    },
    {
        "_id": "api_idlerule",
        "description": "Idle rules",
        "name": "api_idlerule",
        "type": "CRUD"
    },
    {
        "_id": "api_eventfilter",
        "description": "Event filters",
        "name": "api_eventfilter",
        "type": "CRUD"
    },
    {
        "_id": "api_action",
        "description": "Actions",
        "name": "api_action",
        "type": "CRUD"
    },
    {
        "_id": "api_metaalarmrule",
        "description": "Meta-alarm rules",
        "name": "api_metaalarmrule",
        "type": "CRUD"
    },
    {
        "_id": "api_playlist",
        "description": "Playlists",
        "name": "api_playlist",
        "type": "CRUD"
    },
    {
        "_id": "api_dynamicinfos",
        "description": "Dynamic infos",
        "name": "api_dynamicinfos",
        "type": "CRUD"
    },
    {
        "_id": "api_entityservice",
        "description": "Entity service",
        "name": "api_entityservice",
        "type": "CRUD"
    },
    {
        "_id": "api_entitycategory",
        "description": "Entity categories",
        "name": "api_entitycategory",
        "type": "CRUD"
    },
    {
        "_id": "api_entitycomment",
        "description": "Entity comments",
        "name": "api_entitycomment"
    },
    {
        "_id": "api_viewgroup",
        "description": "View groups",
        "name": "api_viewgroup",
        "type": "CRUD"
    },
    {
        "_id": "api_view",
        "description": "Views",
        "name": "api_view",
        "type": "CRUD"
    },
    {
        "_id": "api_widgettemplate",
        "description": "Widget templates",
        "name": "api_widgettemplate",
        "type": "CRUD"
    },
    {
        "_id": "api_pbehavior",
        "description": "PBehaviors",
        "name": "api_pbehavior",
        "type": "CRUD"
    },
    {
        "_id": "api_pbehaviortype",
        "description": "PBehaviorTypes",
        "name": "api_pbehaviortype",
        "type": "CRUD"
    },
    {
        "_id": "api_pbehaviorreason",
        "description": "PBehaviorReasons",
        "name": "api_pbehaviorreason",
        "type": "CRUD"
    },
    {
        "_id": "api_pbehaviorexception",
        "description": "PBehaviorExceptions",
        "name": "api_pbehaviorexception",
        "type": "CRUD"
    },
    {
        "_id": "api_event",
        "description": "Event",
        "name": "api_event"
    },
    {
        "_id": "api_healthcheck",
        "description": "Healthcheck",
        "name": "api_healthcheck"
    },
    {
        "_id": "api_message_rate_stats_read",
        "description": "Message rate statistics",
        "name": "api_message_rate_stats_read"
    },
    {
        "_id": "api_execution",
        "description": "Runs instructions",
        "name": "api_execution"
    },
    {
        "_id": "api_job_config",
        "description": "Job configs",
        "name": "api_job_config",
        "type": "CRUD"
    },
    {
        "_id": "api_job",
        "description": "Jobs",
        "name": "api_job",
        "type": "CRUD"
    },
    {
        "_id": "api_instruction",
        "description": "Instructions",
        "name": "api_instruction",
        "type": "CRUD"
    },
    {
        "_id": "api_file",
        "description": "File",
        "name": "api_file",
        "type": "CRUD"
    },
    {
        "_id": "api_entity",
        "description": "Entity",
        "name": "api_entity",
        "type": "CRUD"
    },
    {
        "_id": "api_contextgraph",
        "description": "Context graph import",
        "name": "api_contextgraph",
        "type": "CRUD"
    },
    {
        "_id": "api_acl",
        "description": "Roles, permissions, users",
        "name": "api_acl",
        "type": "CRUD"
    },
    {
        "_id": "api_state_settings",
        "description": "State settings",
        "name": "api_state_settings",
        "type": "CRUD"
    },
    {
        "_id": "api_broadcast_message",
        "description": "Broadcast Message",
        "name": "api_broadcast_message",
        "type": "CRUD"
    },
    {
        "_id": "api_associative_table",
        "description": "Associative table",
        "name": "api_associative_table",
        "type": "CRUD"
    },
    {
        "_id": "api_user_interface_update",
        "description": "Update user interface",
        "name": "api_user_interface_update"
    },
    {
        "_id": "api_user_interface_delete",
        "description": "Delete user interface",
        "name": "api_user_interface_delete"
    },
    {
        "_id": "api_junit",
        "description": "JUnit API",
        "name": "api_junit",
        "type": "CRUD"
    },
    {
        "_id": "api_datastorage_read",
        "description": "Data storage settings read",
        "name": "api_datastorage_read"
    },
    {
        "_id": "api_datastorage_update",
        "description": "Data storage settings update",
        "name": "api_datastorage_update"
    },
    {
        "_id": "api_notification",
        "description": "Notification settings",
        "name": "api_notification"
    },
    {
        "_id": "api_instruction_approve",
        "description": "Instruction approve",
        "name": "api_instruction_approve"
    },
    {
        "_id": "api_resolve_rule",
        "description": "Resolve rule",
        "name": "api_resolve_rule",
        "type": "CRUD"
    },
    {
        "_id": "api_flapping_rule",
        "description": "Flapping rule",
        "name": "api_flapping_rule",
        "type": "CRUD"
    },
    {
        "_id": "api_metrics",
        "description": "Metrics",
        "name": "api_metrics"
    },
    {
        "_id": "api_rating_settings",
        "description": "Rating settings",
        "name": "api_rating_settings",
        "type": "CRUD"
    },
    {
        "_id": "api_metrics_settings",
        "description": "Metrics settings",
        "name": "api_metrics_settings",
        "type": "CRUD"
    },
    {
        "_id": "api_kpi_filter",
        "description": "KPI filters",
        "name": "api_kpi_filter",
        "type": "CRUD"
    },
    {
        "_id": "api_snmprule",
        "description": "SNMP",
        "name": "api_snmprule",
        "type": "CRUD"
    },
    {
        "_id": "api_snmpmib",
        "description": "SNMP MIB",
        "name": "api_snmpmib",
        "type": "CRUD"
    },
    {
        "_id": "api_corporate_pattern",
        "description": "Corporate patterns",
        "name": "api_corporate_pattern"
    },
    {
        "_id": "api_map",
        "description": "Map",
        "name": "api_map",
        "type": "CRUD"
    },
    {
        "_id": "api_share_token",
        "description": "Share token",
        "name": "api_share_token",
        "type": "CRUD"
    },
    {
        "_id": "api_techmetrics",
        "description": "Tech Metrics",
        "name": "api_techmetrics"
    },
    {
        "_id": "api_export_configurations",
        "description": "Export configurations",
        "name": "api_export_configurations"
    },
    {
        "_id": "api_declare_ticket_rule",
        "description": "Declare ticket rule",
        "name": "api_declare_ticket_rule",
        "type": "CRUD"
    },
    {
        "_id": "api_declare_ticket_execution",
        "description": "Run declare ticket rules",
        "name": "api_declare_ticket_execution"
    },
    {
        "_id": "api_link_rule",
        "description": "Link rule",
        "name": "api_link_rule",
        "type": "CRUD"
    },
    {
        "_id": "api_alarm_tag",
        "description": "Alarm tags",
        "name": "api_alarm_tag",
        "type": "CRUD"
    },
    {
        "_id": "api_maintenance",
        "description": "Trigger maintenance mode",
        "name": "api_maintenance"
    },
    {
        "_id": "api_color_theme",
        "description": "Color theme",
        "name": "api_color_theme",
        "type": "CRUD"
    },
    {
        "_id": "api_private_view_groups",
        "description": "Create private view groups",
        "name": "api_private_view_groups"
    },
    {
        "_id": "api_icon",
        "description": "Create icons",
        "name": "api_icon"
    },
    {
        "_id": "api_techmetrics_settings",
        "description": "Tech metrics settings",
        "name": "api_techmetrics_settings",
        "type": "CRUD"
    },
    {
        "_id": "api_launch_event_recording",
        "description": "Launch event recording",
        "name": "api_launch_event_recording"
    },
    {
        "_id": "api_resend_events",
        "description": "Event recorder resend events",
        "name": "api_resend_events"
    }
];
let updatedPermsByID = {};
for (const perm of updatedPerms) {
    updatedPermsByID[perm._id] = perm;
}

db.permission.aggregate([
    {
        $lookup: {
            from: "views",
            localField: "_id",
            foreignField: "_id",
            as: "view",
        }
    },
    {$unwind: {path: "$view", preserveNullAndEmptyArrays: true}},
    {
        $lookup: {
            from: "view_playlist",
            localField: "_id",
            foreignField: "_id",
            as: "playlist",
        }
    },
    {$unwind: {path: "$playlist", preserveNullAndEmptyArrays: true}},
]).forEach(function (doc) {
    if (doc.view) {
        db.permission.updateOne({_id: doc._id}, {
            $set: {description: "Rights on view : " + doc.view.title},
            $unset: {groups: ""},
        });

        return;
    }

    if (doc.playlist) {
        db.permission.updateOne({_id: doc._id}, {
            $set: {description: "Rights on playlist : " + doc.playlist.name},
            $unset: {groups: ""},
        });

        return;
    }

    if (updatedPermsByID[doc._id]) {
        db.permission.replaceOne({_id: doc._id}, updatedPermsByID[doc._id]);
        delete updatedPermsByID[doc._id];

        return;
    }

    db.permission.deleteOne({_id: doc._id});
    db.role.updateMany({}, {$unset: {["permissions." + doc._id]: ""}})
});

for (const k of Object.keys(updatedPermsByID)) {
    const perm = updatedPermsByID[k];
    db.permission.insertOne(perm);
    let bitmask = 1;
    switch (perm.type) {
        case "CRUD":
            bitmask = 15;

            break;
        case "RW":
            bitmask = 7;

            break;
    }

    db.role.updateOne({name: "admin"}, {$set: {["permissions." + perm._id]: bitmask}})
}

db.role.updateMany({type: {$ne: null}}, {$unset: {type: ""}});
db.role_template.updateMany({type: {$ne: null}}, {$unset: {type: ""}});
