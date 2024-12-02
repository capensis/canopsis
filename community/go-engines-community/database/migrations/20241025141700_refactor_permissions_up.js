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
    let unset = {};
    let orConds = [];
    for (const permToReplace of replacePermissions[perm]) {
        unset["permissions." + permToReplace] = "";
        orConds.push({["permissions." + permToReplace]: 1});
    }

    db.role.updateOne(
        {$or: orConds},
        {
            $set: {["permissions." + perm]: 1},
            $unset: unset,
        }
    );
}

const updatedPerms = [
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
        "description": "PBehavior types",
        "name": "api_pbehaviortype",
        "type": "CRUD"
    },
    {
        "_id": "api_pbehaviorreason",
        "description": "PBehavior reasons",
        "name": "api_pbehaviorreason",
        "type": "CRUD"
    },
    {
        "_id": "api_pbehaviorexception",
        "description": "PBehavior exceptions",
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
        "description": "Tech metrics",
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
        "description": "Theme color",
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
    },
    {
        "_id": "models_userview",
        "description": "User views",
        "name": "models_userview",
        "type": "CRUD",
        "api_permissions": {
            "api_viewgroup": 0,
            "api_view": 0
        }
    },
    {
        "_id": "models_role",
        "description": "Roles",
        "name": "models_role",
        "type": "CRUD",
        "api_permissions": {
            "api_acl": 0
        }
    },
    {
        "_id": "models_permission",
        "description": "Rights",
        "name": "models_permission",
        "type": "CRUD",
        "api_permissions": {
            "api_acl": 0
        }
    },
    {
        "_id": "models_user",
        "description": "Users",
        "name": "models_user",
        "type": "CRUD",
        "api_permissions": {
            "api_acl": 0
        }
    },
    {
        "_id": "models_parameters",
        "description": "Parameters - parameters tab",
        "name": "models_parameters",
        "type": "RW",
        "api_permissions_bitmask": {
            "2": {
                "api_user_interface_update": 1
            },
            "1": {
                "api_user_interface_delete": 1
            }
        }
    },
    {
        "_id": "models_broadcastMessage",
        "description": "Broadcast messages",
        "name": "models_broadcastMessage",
        "type": "CRUD",
        "api_permissions": {
            "api_broadcast_message": 0
        }
    },
    {
        "_id": "models_playlist",
        "description": "Playlists",
        "name": "models_playlist",
        "type": "CRUD",
        "api_permissions": {
            "api_playlist": 0
        }
    },
    {
        "_id": "models_planningType",
        "description": "Planning type (Pbehavior)",
        "name": "models_planningType",
        "type": "CRUD",
        "api_permissions": {
            "api_pbehaviortype": 0
        }
    },
    {
        "_id": "models_planningReason",
        "description": "Planning reason (Pbehavior)",
        "name": "models_planningReason",
        "type": "CRUD",
        "api_permissions": {
            "api_pbehaviorreason": 0
        }
    },
    {
        "_id": "models_planningExceptions",
        "description": "Planning exceptions dates (Pbehavior)",
        "name": "models_planningExceptions",
        "type": "CRUD",
        "api_permissions": {
            "api_pbehaviorexception": 0
        }
    },
    {
        "_id": "models_remediationInstruction",
        "description": "Instructions - instructions tab",
        "name": "models_remediationInstruction",
        "type": "CRUD",
        "api_permissions": {
            "api_entitycategory": 4,
            "api_instruction": 0,
            "api_alarm_read": 1,
            "api_entity": 4
        }
    },
    {
        "_id": "models_remediationJob",
        "description": "Instructions - jobs tab",
        "name": "models_remediationJob",
        "type": "CRUD",
        "api_permissions": {
            "api_job": 0
        }
    },
    {
        "_id": "models_remediationConfiguration",
        "description": "Instructions - configurations tab",
        "name": "models_remediationConfiguration",
        "type": "CRUD",
        "api_permissions": {
            "api_job_config": 0
        }
    },
    {
        "_id": "models_exploitation_eventFilter",
        "description": "Event filters",
        "name": "models_exploitation_eventFilter",
        "type": "CRUD",
        "api_permissions": {
            "api_eventfilter": 0,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_exploitation_pbehavior",
        "description": "Pbehaviors",
        "name": "models_exploitation_pbehavior",
        "type": "CRUD",
        "api_permissions": {
            "api_pbehavior": 0,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_exploitation_snmpRule",
        "description": "Snmp rules",
        "name": "models_exploitation_snmpRule",
        "type": "CRUD",
        "api_permissions": {
            "api_snmprule": 0
        }
    },
    {
        "_id": "models_exploitation_dynamicInfo",
        "description": "Dynamic information rules",
        "name": "models_exploitation_dynamicInfo",
        "type": "CRUD",
        "api_permissions": {
            "api_entitycategory": 4,
            "api_dynamicinfos": 0,
            "api_alarm_read": 1,
            "api_entity": 4
        }
    },
    {
        "_id": "models_exploitation_metaAlarmRule",
        "description": "Meta alarm rules and correlation",
        "name": "models_exploitation_metaAlarmRule",
        "type": "CRUD",
        "api_permissions": {
            "api_metaalarmrule": 0,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_exploitation_scenario",
        "description": "Scenarios",
        "name": "models_exploitation_scenario",
        "type": "CRUD",
        "api_permissions": {
            "api_action": 0,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_exploitation_idleRules",
        "description": "Idle rules",
        "name": "models_exploitation_idleRules",
        "type": "CRUD",
        "api_permissions": {
            "api_idlerule": 0,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_exploitation_flappingRules",
        "description": "Flapping rules",
        "name": "models_exploitation_flappingRules",
        "type": "CRUD",
        "api_permissions": {
            "api_flapping_rule": 0,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_exploitation_resolveRules",
        "description": "Resolve rules",
        "name": "models_exploitation_resolveRules",
        "type": "CRUD",
        "api_permissions": {
            "api_resolve_rule": 0,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_exploitation_declareTicketRule",
        "description": "Ticket declaration rules",
        "name": "models_exploitation_declareTicketRule",
        "type": "CRUD",
        "api_permissions": {
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4,
            "api_declare_ticket_rule": 0,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4,
            "api_pbehavior": 4
        }
    },
    {
        "_id": "models_exploitation_linkRule",
        "description": "Link generator",
        "name": "models_exploitation_linkRule",
        "type": "CRUD",
        "api_permissions": {
            "api_link_rule": 0,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_healthcheck",
        "description": "Healthcheck",
        "name": "models_healthcheck",
        "api_permissions": {
            "api_healthcheck": 1,
            "api_message_rate_stats_read": 1
        }
    },
    {
        "_id": "models_healthcheckStatus",
        "description": "Healthcheck status",
        "name": "models_healthcheckStatus",
        "api_permissions": {
            "api_healthcheck": 1
        }
    },
    {
        "_id": "models_kpi",
        "description": "KPI graphs",
        "name": "models_kpi",
        "api_permissions": {
            "api_metrics": 1,
            "api_kpi_filter": 4,
            "api_rating_settings": 4,
            "api_healthcheck": 1
        }
    },
    {
        "_id": "models_kpiFilters",
        "description": "KPI filters",
        "name": "models_kpiFilters",
        "type": "CRUD",
        "api_permissions": {
            "api_kpi_filter": 0,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_kpiRatingSettings",
        "description": "KPI rating settings",
        "name": "models_kpiRatingSettings",
        "type": "CRUD",
        "api_permissions": {
            "api_rating_settings": 0
        }
    },
    {
        "_id": "models_kpiCollectionSettings",
        "description": "KPI collection settings",
        "name": "models_kpiCollectionSettings",
        "type": "CRUD",
        "api_permissions": {
            "api_metrics_settings": 0
        }
    },
    {
        "_id": "models_notification_instructionStats",
        "description": "Instructions statistics",
        "name": "models_notification_instructionStats",
        "type": "CRUD",
        "api_permissions_bitmask": {
            "4": {
                "api_instruction": 4
            }
        }
    },
    {
        "_id": "models_profile_corporatePattern",
        "description": "Corporate patterns",
        "name": "models_profile_corporatePattern",
        "type": "CRUD",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4,
            "api_pbehavior": 4,
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4,
            "api_alarm_read": 1
        }
    },
    {
        "_id": "models_map",
        "description": "Maps",
        "name": "models_map",
        "type": "CRUD",
        "api_permissions": {
            "api_map": 0,
            "api_entity": 4
        }
    },
    {
        "_id": "models_shareToken",
        "description": "Shared token settings",
        "name": "models_shareToken",
        "type": "CRUD",
        "api_permissions_bitmask": {
            "8": {
                "api_share_token": 8
            }
        }
    },
    {
        "_id": "models_techmetrics",
        "description": "Healthcheck - engines' metrics",
        "name": "models_techmetrics",
        "api_permissions": {
            "api_techmetrics": 1
        }
    },
    {
        "_id": "models_widgetTemplate",
        "description": "Parameters - widget templates",
        "name": "models_widgetTemplate",
        "type": "CRUD",
        "api_permissions": {
            "api_widgettemplate": 0,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_dynamicinfos": 4
        }
    },
    {
        "_id": "models_tag",
        "description": "Alarm tags",
        "name": "models_tag",
        "type": "CRUD",
        "api_permissions": {
            "api_alarm_tag": 0,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_eventsRecord",
        "description": "Events recording",
        "name": "models_eventsRecord"
    },
    {
        "_id": "models_privateView",
        "description": "Private views",
        "name": "models_privateView",
        "api_permissions": {
            "api_dynamicinfos": 4,
            "api_state_settings": 4,
            "api_pbehaviorreason": 4,
            "api_metrics": 1,
            "api_pbehaviortype": 4,
            "api_entitycategory": 4,
            "api_entityservice": 4,
            "api_map": 4,
            "api_rating_settings": 4,
            "api_private_view_groups": 1,
            "api_pbehavior": 4,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_junit": 4
        }
    },
    {
        "_id": "listalarm_ack",
        "description": "Ack",
        "name": "listalarm_ack",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_fastAck",
        "description": "Fast ack",
        "name": "listalarm_fastAck",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_cancelAck",
        "description": "Cancel ack",
        "name": "listalarm_cancelAck",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_pbehavior",
        "description": "Edit PBhaviours for alarm",
        "name": "listalarm_pbehavior",
        "api_permissions": {
            "api_pbehavior": 8
        }
    },
    {
        "_id": "listalarm_fastPbehavior",
        "description": "Fast create PBehavior",
        "name": "listalarm_fastPbehavior",
        "api_permissions": {
            "api_pbehavior": 8
        }
    },
    {
        "_id": "listalarm_snoozeAlarm",
        "description": "Snooze",
        "name": "listalarm_snoozeAlarm",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_listPbehavior",
        "description": "Browse pbehavior list",
        "name": "listalarm_listPbehavior",
        "api_permissions": {
            "api_pbehavior": 4
        }
    },
    {
        "_id": "listalarm_declareanIncident",
        "description": "Declare ticket",
        "name": "listalarm_declareanIncident",
        "api_permissions": {
            "api_declare_ticket_execution": 1,
            "api_declare_ticket_rule": 4
        }
    },
    {
        "_id": "listalarm_assignTicketNumber",
        "description": "Associate ticket",
        "name": "listalarm_assignTicketNumber",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_removeAlarm",
        "description": "Cancel",
        "name": "listalarm_removeAlarm",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_fastRemoveAlarm",
        "description": "Fast cancel",
        "name": "listalarm_fastRemoveAlarm",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_unCancel",
        "description": "Uncancel",
        "name": "listalarm_unCancel",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_changeState",
        "description": "Check and lock severity",
        "name": "listalarm_changeState",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_history",
        "description": "View entity history",
        "name": "listalarm_history",
        "api_permissions": {
            "api_alarm_read": 1
        }
    },
    {
        "_id": "listalarm_manualMetaAlarmGroup",
        "description": "Link to / unlink from manual meta alarm",
        "name": "listalarm_manualMetaAlarmGroup",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_metaAlarmGroup",
        "description": "Unlink from meta alarm created by rule",
        "name": "listalarm_metaAlarmGroup",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_comment",
        "description": "Comment",
        "name": "listalarm_comment",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_filter",
        "description": "Set alarm filters",
        "name": "listalarm_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4,
            "api_pbehavior": 4,
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4,
            "api_alarm_read": 1
        }
    },
    {
        "_id": "listalarm_userFilter",
        "description": "Filter alarms",
        "name": "listalarm_userFilter",
        "api_permissions": {
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4,
            "api_view": 4,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4,
            "api_pbehavior": 4
        }
    },
    {
        "_id": "listalarm_bookmark",
        "description": "Add / remove bookmark",
        "name": "listalarm_bookmark",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_filterByBookmark",
        "description": "Filter bookmarked alarms",
        "name": "listalarm_filterByBookmark"
    },
    {
        "_id": "listalarm_remediationInstructionsFilter",
        "description": "Set filters by remediation instructions",
        "name": "listalarm_remediationInstructionsFilter",
        "api_permissions": {
            "api_instruction": 4
        }
    },
    {
        "_id": "listalarm_userRemediationInstructionsFilter",
        "description": "Filter alarms by remediation instructions",
        "name": "listalarm_userRemediationInstructionsFilter",
        "api_permissions": {
            "api_instruction": 4
        }
    },
    {
        "_id": "listalarm_links",
        "description": "Follow links in alarm",
        "name": "listalarm_links"
    },
    {
        "_id": "listalarm_correlation",
        "description": "Group correlated alarms (meta alarms)",
        "name": "listalarm_correlation"
    },
    {
        "_id": "listalarm_executeInstruction",
        "description": "Execute manual instructions",
        "name": "listalarm_executeInstruction",
        "api_permissions": {
            "api_execution": 1
        }
    },
    {
        "_id": "listalarm_category",
        "description": "Filter alarms by category",
        "name": "listalarm_category",
        "api_permissions": {
            "api_entitycategory": 4
        }
    },
    {
        "_id": "listalarm_exportAsCsv",
        "description": "Export alarm list as CSV file",
        "name": "listalarm_exportAsCsv"
    },
    {
        "_id": "listalarm_exportPdf",
        "description": "Export in PDF",
        "name": "listalarm_exportPdf"
    },
    {
        "_id": "serviceweather_entityAck",
        "description": "Ack",
        "name": "serviceweather_entityAck",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "serviceweather_entityAssocTicket",
        "description": "Associate ticket",
        "name": "serviceweather_entityAssocTicket",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "serviceweather_entityComment",
        "description": "Comment alarm",
        "name": "serviceweather_entityComment",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "serviceweather_entityValidate",
        "description": "Validate alarms and change their state to critical",
        "name": "serviceweather_entityValidate",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "serviceweather_entityInvalidate",
        "description": "Invalidate alarms and cancel them",
        "name": "serviceweather_entityInvalidate",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "serviceweather_entityPause",
        "description": "Pause alarms (set the PBehavior type \"Pause\")",
        "name": "serviceweather_entityPause",
        "api_permissions": {
            "api_pbehavior": 8
        }
    },
    {
        "_id": "serviceweather_entityPlay",
        "description": "Activate paused alarms (remove the PBehavior type \"Pause\")",
        "name": "serviceweather_entityPlay",
        "api_permissions": {
            "api_pbehavior": 1
        }
    },
    {
        "_id": "serviceweather_entityCancel",
        "description": "Cancel",
        "name": "serviceweather_entityCancel",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "serviceweather_entityManagePbehaviors",
        "description": "View PBehaviors of services (in the subtab in the services modal windows)",
        "name": "serviceweather_entityManagePbehaviors",
        "api_permissions": {
            "api_pbehavior": 7,
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4,
            "api_pbehaviorexception": 4
        }
    },
    {
        "_id": "serviceweather_entityLinks",
        "description": "Follow links in alarm",
        "name": "serviceweather_entityLinks"
    },
    {
        "_id": "serviceweather_entityDeclareTicket",
        "description": "Declare ticket",
        "name": "serviceweather_entityDeclareTicket",
        "api_permissions": {
            "api_declare_ticket_execution": 1,
            "api_alarm_read": 1,
            "api_declare_ticket_rule": 4
        }
    },
    {
        "_id": "serviceweather_moreInfos",
        "description": "Open \"More infos\" modal",
        "name": "serviceweather_moreInfos",
        "api_permissions": {
            "api_alarm_read": 1
        }
    },
    {
        "_id": "serviceweather_alarmsList",
        "description": "Open the list of alarms available for each service",
        "name": "serviceweather_alarmsList",
        "api_permissions": {
            "api_state_settings": 4,
            "api_entity": 4,
            "api_pbehavior": 4,
            "api_junit": 4,
            "api_alarm_read": 1,
            "api_entityservice": 4
        }
    },
    {
        "_id": "serviceweather_pbehaviorList",
        "description": "View PBehaviors of services (in the subtab in the service entities modal windows)",
        "name": "serviceweather_pbehaviorList",
        "api_permissions": {
            "api_pbehavior": 15,
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4,
            "api_pbehaviorexception": 4
        }
    },
    {
        "_id": "serviceweather_filter",
        "description": "Set entity filters",
        "name": "serviceweather_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4,
            "api_pbehaviortype": 4
        }
    },
    {
        "_id": "serviceweather_userFilter",
        "description": "Filter entities",
        "name": "serviceweather_userFilter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4,
            "api_pbehaviortype": 4,
            "api_view": 4
        }
    },
    {
        "_id": "serviceweather_category",
        "description": "Filter entities by category",
        "name": "serviceweather_category",
        "api_permissions": {
            "api_entitycategory": 4
        }
    },
    {
        "_id": "serviceweather_executeInstruction",
        "description": "Execute manual instructions",
        "name": "serviceweather_executeInstruction",
        "api_permissions": {
            "api_execution": 1
        }
    },
    {
        "_id": "crudcontext_createEntity",
        "description": "Create entity",
        "name": "crudcontext_createEntity",
        "api_permissions": {
            "api_entityservice": 8
        }
    },
    {
        "_id": "crudcontext_edit",
        "description": "Edit entity",
        "name": "crudcontext_edit",
        "api_permissions": {
            "api_entity": 2,
            "api_entityservice": 2
        }
    },
    {
        "_id": "crudcontext_duplicate",
        "description": "Duplicate entity",
        "name": "crudcontext_duplicate",
        "api_permissions": {
            "api_entityservice": 8
        }
    },
    {
        "_id": "crudcontext_delete",
        "description": "Delete entity",
        "name": "crudcontext_delete",
        "api_permissions": {
            "api_entity": 1,
            "api_entityservice": 1
        }
    },
    {
        "_id": "crudcontext_pbehavior",
        "description": "Create / edit pbehavior",
        "name": "crudcontext_pbehavior",
        "api_permissions": {
            "api_pbehavior": 8
        }
    },
    {
        "_id": "crudcontext_listPbehavior",
        "description": "Browse pbehavior list",
        "name": "crudcontext_listPbehavior",
        "api_permissions": {
            "api_pbehavior": 4
        }
    },
    {
        "_id": "crudcontext_filter",
        "description": "Set entity filters",
        "name": "crudcontext_filter",
        "api_permissions": {
            "api_entitycategory": 4,
            "api_pbehavior": 4,
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4,
            "api_alarm_read": 1,
            "api_entity": 4
        }
    },
    {
        "_id": "crudcontext_userFilter",
        "description": "Filter entities",
        "name": "crudcontext_userFilter",
        "api_permissions": {
            "api_pbehavior": 4,
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4,
            "api_view": 4,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "crudcontext_category",
        "description": "Filter entities by category",
        "name": "crudcontext_category",
        "api_permissions": {
            "api_entitycategory": 4
        }
    },
    {
        "_id": "crudcontext_exportAsCsv",
        "description": "Export entities as CSV file",
        "name": "crudcontext_exportAsCsv"
    },
    {
        "_id": "crudcontext_massEnable",
        "description": "Mass action to enable selected entities",
        "name": "crudcontext_massEnable",
        "api_permissions": {
            "api_entity": 2
        }
    },
    {
        "_id": "crudcontext_massDisable",
        "description": "Mass action to disable selected entities",
        "name": "crudcontext_massDisable",
        "api_permissions": {
            "api_entity": 2
        }
    },
    {
        "_id": "common_variablesHelp",
        "description": "See the list of variables (in all widgets)",
        "name": "common_variablesHelp"
    },
    {
        "_id": "counter_alarmsList",
        "description": "View the alarm list associated with counters",
        "name": "counter_alarmsList",
        "api_permissions": {
            "api_pbehavior": 4,
            "api_junit": 4,
            "api_alarm_read": 1,
            "api_entityservice": 4,
            "api_state_settings": 4,
            "api_entity": 4
        }
    },
    {
        "_id": "testingweather_alarmsList",
        "description": "View the alarm list associated with testing weather",
        "name": "testingweather_alarmsList",
        "api_permissions": {
            "api_alarm_read": 1
        }
    },
    {
        "_id": "map_alarmsList",
        "description": "View the alarm list associated with points on maps",
        "name": "map_alarmsList",
        "api_permissions": {
            "api_pbehavior": 4,
            "api_junit": 4,
            "api_alarm_read": 1,
            "api_entityservice": 4,
            "api_state_settings": 4,
            "api_entity": 4
        }
    },
    {
        "_id": "map_filter",
        "description": "Set filters for points on maps",
        "name": "map_filter",
        "api_permissions": {
            "api_pbehaviorreason": 4,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4,
            "api_pbehavior": 4,
            "api_pbehaviortype": 4
        }
    },
    {
        "_id": "map_userFilter",
        "description": "Filter points on maps",
        "name": "map_userFilter",
        "api_permissions": {
            "api_pbehaviorreason": 4,
            "api_view": 4,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4,
            "api_pbehavior": 4,
            "api_pbehaviortype": 4
        }
    },
    {
        "_id": "map_category",
        "description": "Filter points on maps by categories",
        "name": "map_category",
        "api_permissions": {
            "api_entitycategory": 4
        }
    },
    {
        "_id": "userStatistics_interval",
        "description": "Edit time intervals for the data displayed",
        "name": "userStatistics_interval"
    },
    {
        "_id": "userStatistics_filter",
        "description": "Set data filters",
        "name": "userStatistics_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "userStatistics_userFilter",
        "description": "Filter data",
        "name": "userStatistics_userFilter",
        "api_permissions": {
            "api_view": 4,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "alarmStatistics_interval",
        "description": "Edit time intervals for the data displayed",
        "name": "alarmStatistics_interval"
    },
    {
        "_id": "alarmStatistics_filter",
        "description": "Set data filters",
        "name": "alarmStatistics_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "alarmStatistics_userFilter",
        "description": "Filter data",
        "name": "alarmStatistics_userFilter",
        "api_permissions": {
            "api_view": 4,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_remediationStatistic",
        "description": "Instructions - remediation statistics tab",
        "name": "models_remediationStatistic",
        "api_permissions": {
            "api_instruction": 4,
            "api_metrics": 1
        }
    },
    {
        "_id": "barchart_interval",
        "description": "Edit time intervals for the data displayed",
        "name": "barchart_interval"
    },
    {
        "_id": "barchart_sampling",
        "description": "Edit sampling for the data displayed",
        "name": "barchart_sampling"
    },
    {
        "_id": "barchart_filter",
        "description": "Set data filters",
        "name": "barchart_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "barchart_userFilter",
        "description": "Filter data",
        "name": "barchart_userFilter",
        "api_permissions": {
            "api_view": 4,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "linechart_interval",
        "description": "Edit time intervals for the data displayed",
        "name": "linechart_interval"
    },
    {
        "_id": "linechart_sampling",
        "description": "Edit sampling for the data displayed",
        "name": "linechart_sampling"
    },
    {
        "_id": "linechart_filter",
        "description": "Set data filters",
        "name": "linechart_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "linechart_userFilter",
        "description": "Filter data",
        "name": "linechart_userFilter",
        "api_permissions": {
            "api_view": 4,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "piechart_interval",
        "description": "Edit time intervals for the data displayed",
        "name": "piechart_interval"
    },
    {
        "_id": "piechart_sampling",
        "description": "Edit sampling for the data displayed",
        "name": "piechart_sampling"
    },
    {
        "_id": "piechart_filter",
        "description": "Set data filters",
        "name": "piechart_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "piechart_userFilter",
        "description": "Filter data",
        "name": "piechart_userFilter",
        "api_permissions": {
            "api_entitycategory": 4,
            "api_view": 4,
            "api_entity": 4
        }
    },
    {
        "_id": "numbers_interval",
        "description": "Edit time intervals for the data displayed",
        "name": "numbers_interval"
    },
    {
        "_id": "numbers_sampling",
        "description": "Edit sampling for the data displayed",
        "name": "numbers_sampling"
    },
    {
        "_id": "numbers_filter",
        "description": "Set data filters",
        "name": "numbers_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "numbers_userFilter",
        "description": "Filter data",
        "name": "numbers_userFilter",
        "api_permissions": {
            "api_view": 4,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_maintenance",
        "description": "Maintenance",
        "name": "models_maintenance",
        "api_permissions": {
            "api_maintenance": 1
        }
    },
    {
        "_id": "models_profile_color_theme",
        "description": "Theme colors",
        "name": "models_profile_color_theme",
        "type": "CRUD",
        "api_permissions": {
            "api_color_theme": 0
        }
    },
    {
        "_id": "models_stateSetting",
        "description": "State settings",
        "name": "models_stateSetting",
        "type": "CRUD",
        "api_permissions": {
            "api_state_settings": 0
        }
    },
    {
        "_id": "models_icon",
        "description": "Parameters - icons",
        "name": "models_icon",
        "type": "CRUD",
        "api_permissions_bitmask": {
            "8": {
                "api_icon": 1
            },
            "2": {
                "api_icon": 1
            },
            "1": {
                "api_icon": 1
            }
        }
    },
    {
        "_id": "models_view_import_export",
        "description": "Parameters - import / export",
        "name": "models_view_import_export",
        "api_permissions": {
            "api_viewgroup": 3,
            "api_view": 3
        }
    },
    {
        "_id": "models_notification",
        "description": "Parameters - notification settings",
        "name": "models_notification",
        "api_permissions": {
            "api_notification": 1
        }
    },
    {
        "_id": "availability_interval",
        "description": "Edit time intervals for the data displayed",
        "name": "availability_interval"
    },
    {
        "_id": "availability_filter",
        "description": "Set data filters",
        "name": "availability_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "availability_userFilter",
        "description": "Filter data",
        "name": "availability_userFilter",
        "api_permissions": {
            "api_view": 4,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "availability_exportAsCsv",
        "description": "Export availabilities as CSV file",
        "name": "availability_exportAsCsv"
    },
    {
        "_id": "common_entityComment",
        "description": "Manage entity comments (view, create, edit, delete)",
        "name": "common_entityComment",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycomment": 1
        }
    },
    {
        "_id": "AlarmsList",
        "description": "API permissions for AlarmsList widget",
        "name": "AlarmsList",
        "api_permissions_bitmask": {
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_entity": 4,
                "api_state_settings": 4,
                "api_junit": 4,
                "api_metrics": 1,
                "api_alarm_read": 1,
                "api_entityservice": 4,
                "api_pbehavior": 4
            },
            "2": {
                "api_view": 2,
                "api_widgettemplate": 4,
                "api_dynamicinfos": 4
            },
            "1": {
                "api_view": 1
            }
        },
        "hidden": true
    },
    {
        "_id": "Context",
        "description": "API permissions for Context widget",
        "name": "Context",
        "api_permissions_bitmask": {
            "4": {
                "api_alarm_read": 1,
                "api_entityservice": 4,
                "api_state_settings": 4,
                "api_pbehavior": 4,
                "api_metrics": 1,
                "api_view": 4,
                "api_associative_table": 4,
                "api_entity": 4
            },
            "2": {
                "api_view": 2,
                "api_widgettemplate": 4,
                "api_dynamicinfos": 4
            },
            "1": {
                "api_view": 1
            }
        },
        "hidden": true
    },
    {
        "_id": "ServiceWeather",
        "description": "API permissions for ServiceWeather widget",
        "name": "ServiceWeather",
        "api_permissions_bitmask": {
            "1": {
                "api_view": 1
            },
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_entityservice": 4,
                "api_alarm_read": 1,
                "api_state_settings": 4
            },
            "2": {
                "api_dynamicinfos": 4,
                "api_pbehaviortype": 4,
                "api_view": 2,
                "api_widgettemplate": 4,
                "api_entity": 4
            }
        },
        "hidden": true
    },
    {
        "_id": "StatsCalendar",
        "description": "API permissions for StatsCalendar widget",
        "name": "StatsCalendar",
        "api_permissions_bitmask": {
            "1": {
                "api_view": 1
            },
            "4": {
                "api_alarm_read": 1,
                "api_state_settings": 4,
                "api_junit": 4,
                "api_view": 4,
                "api_associative_table": 4,
                "api_entityservice": 4,
                "api_entity": 4,
                "api_pbehavior": 4
            },
            "2": {
                "api_pbehaviorreason": 4,
                "api_widgettemplate": 4,
                "api_pbehavior": 4,
                "api_pbehaviortype": 4,
                "api_entity": 4,
                "api_entitycategory": 4,
                "api_view": 2,
                "api_dynamicinfos": 4,
                "api_alarm_read": 1
            }
        },
        "hidden": true
    },
    {
        "_id": "Text",
        "description": "API permissions for Text widget",
        "name": "Text",
        "api_permissions_bitmask": {
            "4": {
                "api_view": 4,
                "api_associative_table": 4
            },
            "2": {
                "api_view": 2,
                "api_widgettemplate": 4
            },
            "1": {
                "api_view": 1
            }
        },
        "hidden": true
    },
    {
        "_id": "Counter",
        "description": "API permissions for Counter widget",
        "name": "Counter",
        "api_permissions_bitmask": {
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_alarm_read": 1
            },
            "2": {
                "api_entity": 4,
                "api_dynamicinfos": 4,
                "api_alarm_read": 1,
                "api_entitycategory": 4,
                "api_pbehaviortype": 4,
                "api_widgettemplate": 4,
                "api_pbehavior": 4,
                "api_pbehaviorreason": 4,
                "api_view": 2
            },
            "1": {
                "api_view": 1
            }
        },
        "hidden": true
    },
    {
        "_id": "Junit",
        "description": "API permissions for Junit widget",
        "name": "Junit",
        "api_permissions_bitmask": {
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_junit": 4
            },
            "2": {
                "api_view": 2,
                "api_widgettemplate": 4
            },
            "1": {
                "api_view": 1
            }
        },
        "hidden": true
    },
    {
        "_id": "Map",
        "description": "API permissions for Map widget",
        "name": "Map",
        "api_permissions_bitmask": {
            "2": {
                "api_dynamicinfos": 4,
                "api_view": 2,
                "api_widgettemplate": 4,
                "api_alarm_read": 1,
                "api_entity": 4
            },
            "1": {
                "api_view": 1
            },
            "4": {
                "api_entity": 4,
                "api_pbehavior": 4,
                "api_view": 4,
                "api_associative_table": 4,
                "api_map": 4,
                "api_entityservice": 4,
                "api_alarm_read": 1
            }
        },
        "hidden": true
    },
    {
        "_id": "BarChart",
        "description": "API permissions for BarChart widget",
        "name": "BarChart",
        "api_permissions_bitmask": {
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_metrics": 1
            },
            "2": {
                "api_view": 2,
                "api_widgettemplate": 4
            },
            "1": {
                "api_view": 1
            }
        },
        "hidden": true
    },
    {
        "_id": "LineChart",
        "description": "API permissions for LineChart widget",
        "name": "LineChart",
        "api_permissions_bitmask": {
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_metrics": 1
            },
            "2": {
                "api_widgettemplate": 4,
                "api_view": 2
            },
            "1": {
                "api_view": 1
            }
        },
        "hidden": true
    },
    {
        "_id": "PieChart",
        "description": "API permissions for PieChart widget",
        "name": "PieChart",
        "api_permissions_bitmask": {
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_metrics": 1
            },
            "2": {
                "api_view": 2,
                "api_widgettemplate": 4
            },
            "1": {
                "api_view": 1
            }
        },
        "hidden": true
    },
    {
        "_id": "Numbers",
        "description": "API permissions for Numbers widget",
        "name": "Numbers",
        "api_permissions_bitmask": {
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_metrics": 1
            },
            "2": {
                "api_view": 2,
                "api_widgettemplate": 4
            },
            "1": {
                "api_view": 1
            }
        },
        "hidden": true
    },
    {
        "_id": "UserStatistics",
        "description": "API permissions for UserStatistics widget",
        "name": "UserStatistics",
        "api_permissions_bitmask": {
            "4": {
                "api_metrics": 1,
                "api_rating_settings": 4,
                "api_view": 4,
                "api_associative_table": 4
            },
            "2": {
                "api_view": 2,
                "api_widgettemplate": 4
            },
            "1": {
                "api_view": 1
            }
        },
        "hidden": true
    },
    {
        "_id": "AlarmStatistics",
        "description": "API permissions for AlarmStatistics widget",
        "name": "AlarmStatistics",
        "api_permissions_bitmask": {
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_metrics": 1,
                "api_rating_settings": 4
            },
            "2": {
                "api_view": 2,
                "api_widgettemplate": 4
            },
            "1": {
                "api_view": 1
            }
        },
        "hidden": true
    },
    {
        "_id": "Availability",
        "description": "API permissions for Availability widget",
        "name": "Availability",
        "api_permissions_bitmask": {
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_metrics": 1
            },
            "2": {
                "api_widgettemplate": 4,
                "api_view": 2
            },
            "1": {
                "api_view": 1
            }
        },
        "hidden": true
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
        db.permission.updateOne({_id: doc._id}, {$set: {description: doc.view.description}});

        return;
    }

    if (doc.playlist) {
        db.permission.updateOne({_id: doc._id}, {$set: {description: doc.playlist.name}});

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

db.role.aggregate([
    {$match: {type: null}},
    {$project: {permissions: {$objectToArray: "$permissions"}}},
    {$match: {"permissions.k": {$regex: "^(?!api_)"}}}
]).forEach(function (doc) {
    db.role.updateOne({_id: doc._id}, {$set: {type: "ui"}});
});
db.role_template.aggregate([
    {$match: {type: null}},
    {$project: {permissions: {$objectToArray: "$permissions"}}},
    {$match: {"permissions.k": {$regex: "^(?!api_)"}}}
]).forEach(function (doc) {
    db.role_template.updateOne({_id: doc._id}, {$set: {type: "ui"}});
});

db.role.updateMany({type: null}, {$set: {type: "api"}});
db.role_template.updateMany({type: null}, {$set: {type: "api"}});
