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

const permGroups = [
    {
        "_id": "widgets",
        "position": 1
    },
    {
        "_id": "widgets_common",
        "position": 2
    },
    {
        "_id": "widgets_alarmslist",
        "position": 3
    },
    {
        "_id": "widgets_alarmslist_widgetsettings",
        "position": 4
    },
    {
        "_id": "widgets_alarmslist_viewsettings",
        "position": 5
    },
    {
        "_id": "widgets_alarmslist_alarmactions",
        "position": 6
    },
    {
        "_id": "widgets_alarmslist_actions",
        "position": 7
    },
    {
        "_id": "widgets_context",
        "position": 8
    },
    {
        "_id": "widgets_context_widgetsettings",
        "position": 9
    },
    {
        "_id": "widgets_context_viewsettings",
        "position": 10
    },
    {
        "_id": "widgets_context_entityactions",
        "position": 11
    },
    {
        "_id": "widgets_context_actions",
        "position": 12
    },
    {
        "_id": "widgets_serviceweather",
        "position": 13
    },
    {
        "_id": "widgets_serviceweather_widgetsettings",
        "position": 14
    },
    {
        "_id": "widgets_serviceweather_viewsettings",
        "position": 15
    },
    {
        "_id": "widgets_serviceweather_alarmactions",
        "position": 16
    },
    {
        "_id": "widgets_counter",
        "position": 17
    },
    {
        "_id": "widgets_testingweather",
        "position": 18
    },
    {
        "_id": "widgets_map",
        "position": 19
    },
    {
        "_id": "widgets_map_widgetsettings",
        "position": 20
    },
    {
        "_id": "widgets_map_viewsettings",
        "position": 21
    },
    {
        "_id": "widgets_barchart",
        "position": 22
    },
    {
        "_id": "widgets_barchart_widgetsettings",
        "position": 23
    },
    {
        "_id": "widgets_barchart_viewsettings",
        "position": 24
    },
    {
        "_id": "widgets_linechart",
        "position": 25
    },
    {
        "_id": "widgets_linechart_widgetsettings",
        "position": 26
    },
    {
        "_id": "widgets_linechart_viewsettings",
        "position": 27
    },
    {
        "_id": "widgets_piechart",
        "position": 28
    },
    {
        "_id": "widgets_piechart_widgetsettings",
        "position": 29
    },
    {
        "_id": "widgets_piechart_viewsettings",
        "position": 30
    },
    {
        "_id": "widgets_numbers",
        "position": 31
    },
    {
        "_id": "widgets_numbers_widgetsettings",
        "position": 32
    },
    {
        "_id": "widgets_numbers_viewsettings",
        "position": 33
    },
    {
        "_id": "widgets_userstatistics",
        "position": 34
    },
    {
        "_id": "widgets_userstatistics_widgetsettings",
        "position": 35
    },
    {
        "_id": "widgets_userstatistics_viewsettings",
        "position": 36
    },
    {
        "_id": "widgets_alarmstatistics",
        "position": 37
    },
    {
        "_id": "widgets_alarmstatistics_widgetsettings",
        "position": 38
    },
    {
        "_id": "widgets_alarmstatistics_viewsettings",
        "position": 39
    },
    {
        "_id": "widgets_availability",
        "position": 40
    },
    {
        "_id": "widgets_availability_widgetsettings",
        "position": 41
    },
    {
        "_id": "widgets_availability_viewsettings",
        "position": 42
    },
    {
        "_id": "widgets_availability_actions",
        "position": 43
    },
    {
        "_id": "commonviews",
        "position": 44
    },
    {
        "_id": "commonviews_playlist",
        "position": -1
    },
    {
        "_id": "technical",
        "position": 46
    },
    {
        "_id": "technical_admin",
        "position": 47
    },
    {
        "_id": "technical_admin_access",
        "position": 48
    },
    {
        "_id": "technical_admin_communication",
        "position": 49
    },
    {
        "_id": "technical_admin_general",
        "position": 50
    },
    {
        "_id": "technical_viewsandwidgets",
        "position": 51
    },
    {
        "_id": "technical_token",
        "position": 52
    },
    {
        "_id": "technical_exploitation",
        "position": 53
    },
    {
        "_id": "technical_notification",
        "position": 54
    },
    {
        "_id": "technical_profile",
        "position": 55
    },
    {
        "_id": "api",
        "position": 56
    },
    {
        "_id": "api_general",
        "position": 57
    },
    {
        "_id": "api_rules",
        "position": 58
    },
    {
        "_id": "api_remediation",
        "position": 59
    },
    {
        "_id": "api_planning",
        "position": 60
    }
];
if (!db.permission_group.findOne({})) {
    db.permission_group.insertMany(permGroups);
}

const updatedPerms = [
    {
        "_id": "api_alarm_read",
        "description": "Read alarms",
        "groups": ["api", "api_general"],
        "name": "api_alarm_read"
    },
    {
        "_id": "api_alarm_update",
        "description": "Update alarms",
        "groups": ["api", "api_general"],
        "name": "api_alarm_update"
    },
    {
        "_id": "api_playlist",
        "description": "Playlists",
        "groups": ["api", "api_general"],
        "name": "api_playlist",
        "type": "CRUD"
    },
    {
        "_id": "api_entityservice",
        "description": "Entity service",
        "groups": ["api", "api_general"],
        "name": "api_entityservice",
        "type": "CRUD"
    },
    {
        "_id": "api_entitycategory",
        "description": "Entity categories",
        "groups": ["api", "api_general"],
        "name": "api_entitycategory",
        "type": "CRUD"
    },
    {
        "_id": "api_entitycomment",
        "description": "Entity comments",
        "groups": ["api", "api_general"],
        "name": "api_entitycomment"
    },
    {
        "_id": "api_viewgroup",
        "description": "View groups",
        "groups": ["api", "api_general"],
        "name": "api_viewgroup",
        "type": "CRUD"
    },
    {
        "_id": "api_view",
        "description": "Views",
        "groups": ["api", "api_general"],
        "name": "api_view",
        "type": "CRUD"
    },
    {
        "_id": "api_widgettemplate",
        "description": "Widget templates",
        "groups": ["api", "api_general"],
        "name": "api_widgettemplate",
        "type": "CRUD"
    },
    {
        "_id": "api_event",
        "description": "Event",
        "groups": ["api", "api_general"],
        "name": "api_event"
    },
    {
        "_id": "api_healthcheck",
        "description": "Healthcheck",
        "groups": ["api", "api_general"],
        "name": "api_healthcheck"
    },
    {
        "_id": "api_message_rate_stats_read",
        "description": "Message rate statistics",
        "groups": ["api", "api_general"],
        "name": "api_message_rate_stats_read"
    },
    {
        "_id": "api_file",
        "description": "File",
        "groups": ["api", "api_general"],
        "name": "api_file",
        "type": "CRUD"
    },
    {
        "_id": "api_entity",
        "description": "Entity",
        "groups": ["api", "api_general"],
        "name": "api_entity",
        "type": "CRUD"
    },
    {
        "_id": "api_contextgraph",
        "description": "Context graph import",
        "groups": ["api", "api_general"],
        "name": "api_contextgraph",
        "type": "CRUD"
    },
    {
        "_id": "api_acl",
        "description": "Roles, permissions, users",
        "groups": ["api", "api_general"],
        "name": "api_acl",
        "type": "CRUD"
    },
    {
        "_id": "api_state_settings",
        "description": "State settings",
        "groups": ["api", "api_general"],
        "name": "api_state_settings",
        "type": "CRUD"
    },
    {
        "_id": "api_broadcast_message",
        "description": "Broadcast Message",
        "groups": ["api", "api_general"],
        "name": "api_broadcast_message",
        "type": "CRUD"
    },
    {
        "_id": "api_associative_table",
        "description": "Associative table",
        "groups": ["api", "api_general"],
        "name": "api_associative_table",
        "type": "CRUD"
    },
    {
        "_id": "api_user_interface_update",
        "description": "Update user interface",
        "groups": ["api", "api_general"],
        "name": "api_user_interface_update"
    },
    {
        "_id": "api_user_interface_delete",
        "description": "Delete user interface",
        "groups": ["api", "api_general"],
        "name": "api_user_interface_delete"
    },
    {
        "_id": "api_junit",
        "description": "JUnit API",
        "groups": ["api", "api_general"],
        "name": "api_junit",
        "type": "CRUD"
    },
    {
        "_id": "api_datastorage_read",
        "description": "Data storage settings read",
        "groups": ["api", "api_general"],
        "name": "api_datastorage_read"
    },
    {
        "_id": "api_datastorage_update",
        "description": "Data storage settings update",
        "groups": ["api", "api_general"],
        "name": "api_datastorage_update"
    },
    {
        "_id": "api_notification",
        "description": "Notification settings",
        "groups": ["api", "api_general"],
        "name": "api_notification"
    },
    {
        "_id": "api_metrics",
        "description": "Metrics",
        "groups": ["api", "api_general"],
        "name": "api_metrics"
    },
    {
        "_id": "api_rating_settings",
        "description": "Rating settings",
        "groups": ["api", "api_general"],
        "name": "api_rating_settings",
        "type": "CRUD"
    },
    {
        "_id": "api_metrics_settings",
        "description": "Metrics settings",
        "groups": ["api", "api_general"],
        "name": "api_metrics_settings",
        "type": "CRUD"
    },
    {
        "_id": "api_kpi_filter",
        "description": "KPI filters",
        "groups": ["api", "api_general"],
        "name": "api_kpi_filter",
        "type": "CRUD"
    },
    {
        "_id": "api_corporate_pattern",
        "description": "Corporate patterns",
        "groups": ["api", "api_general"],
        "name": "api_corporate_pattern"
    },
    {
        "_id": "api_map",
        "description": "Map",
        "groups": ["api", "api_general"],
        "name": "api_map",
        "type": "CRUD"
    },
    {
        "_id": "api_share_token",
        "description": "Share token",
        "groups": ["api", "api_general"],
        "name": "api_share_token",
        "type": "CRUD"
    },
    {
        "_id": "api_techmetrics",
        "description": "Tech metrics",
        "groups": ["api", "api_general"],
        "name": "api_techmetrics"
    },
    {
        "_id": "api_export_configurations",
        "description": "Export configurations",
        "groups": ["api", "api_general"],
        "name": "api_export_configurations"
    },
    {
        "_id": "api_alarm_tag",
        "description": "Alarm tags",
        "groups": ["api", "api_general"],
        "name": "api_alarm_tag",
        "type": "CRUD"
    },
    {
        "_id": "api_maintenance",
        "description": "Trigger maintenance mode",
        "groups": ["api", "api_general"],
        "name": "api_maintenance"
    },
    {
        "_id": "api_color_theme",
        "description": "Theme color",
        "groups": ["api", "api_general"],
        "name": "api_color_theme",
        "type": "CRUD"
    },
    {
        "_id": "api_private_view_groups",
        "description": "Create private view groups",
        "groups": ["api", "api_general"],
        "name": "api_private_view_groups"
    },
    {
        "_id": "api_icon",
        "description": "Create icons",
        "groups": ["api", "api_general"],
        "name": "api_icon"
    },
    {
        "_id": "api_techmetrics_settings",
        "description": "Tech metrics settings",
        "groups": ["api", "api_general"],
        "name": "api_techmetrics_settings",
        "type": "CRUD"
    },
    {
        "_id": "api_launch_event_recording",
        "description": "Launch event recording",
        "groups": ["api", "api_general"],
        "name": "api_launch_event_recording"
    },
    {
        "_id": "api_resend_events",
        "description": "Event recorder resend events",
        "groups": ["api", "api_general"],
        "name": "api_resend_events"
    },
    {
        "_id": "api_dynamicinfos",
        "description": "Dynamic infos",
        "groups": ["api", "api_rules"],
        "name": "api_dynamicinfos",
        "type": "CRUD"
    },
    {
        "_id": "api_idlerule",
        "description": "Idle rules",
        "groups": ["api", "api_rules"],
        "name": "api_idlerule",
        "type": "CRUD"
    },
    {
        "_id": "api_eventfilter",
        "description": "Event filters",
        "groups": ["api", "api_rules"],
        "name": "api_eventfilter",
        "type": "CRUD"
    },
    {
        "_id": "api_action",
        "description": "Actions",
        "groups": ["api", "api_rules"],
        "name": "api_action",
        "type": "CRUD"
    },
    {
        "_id": "api_resolve_rule",
        "description": "Resolve rule",
        "groups": ["api", "api_rules"],
        "name": "api_resolve_rule",
        "type": "CRUD"
    },
    {
        "_id": "api_flapping_rule",
        "description": "Flapping rule",
        "groups": ["api", "api_rules"],
        "name": "api_flapping_rule",
        "type": "CRUD"
    },
    {
        "_id": "api_link_rule",
        "description": "Link rule",
        "groups": ["api", "api_rules"],
        "name": "api_link_rule",
        "type": "CRUD"
    },
    {
        "_id": "api_metaalarmrule",
        "description": "Meta-alarm rules",
        "groups": ["api", "api_rules"],
        "name": "api_metaalarmrule",
        "type": "CRUD"
    },
    {
        "_id": "api_declare_ticket_rule",
        "description": "Declare ticket rule",
        "groups": ["api", "api_rules"],
        "name": "api_declare_ticket_rule",
        "type": "CRUD"
    },
    {
        "_id": "api_declare_ticket_execution",
        "description": "Run declare ticket rules",
        "groups": ["api", "api_rules"],
        "name": "api_declare_ticket_execution"
    },
    {
        "_id": "api_snmprule",
        "description": "SNMP",
        "groups": ["api", "api_rules"],
        "name": "api_snmprule",
        "type": "CRUD"
    },
    {
        "_id": "api_snmpmib",
        "description": "SNMP MIB",
        "groups": ["api", "api_rules"],
        "name": "api_snmpmib",
        "type": "CRUD"
    },
    {
        "_id": "api_job_config",
        "description": "Job configs",
        "groups": ["api", "api_remediation"],
        "name": "api_job_config",
        "type": "CRUD"
    },
    {
        "_id": "api_job",
        "description": "Jobs",
        "groups": ["api", "api_remediation"],
        "name": "api_job",
        "type": "CRUD"
    },
    {
        "_id": "api_instruction",
        "description": "Instructions",
        "groups": ["api", "api_remediation"],
        "name": "api_instruction",
        "type": "CRUD"
    },
    {
        "_id": "api_instruction_approve",
        "description": "Instruction approve",
        "groups": ["api", "api_remediation"],
        "name": "api_instruction_approve"
    },
    {
        "_id": "api_execution",
        "description": "Runs instructions",
        "groups": ["api", "api_remediation"],
        "name": "api_execution"
    },
    {
        "_id": "api_pbehavior",
        "description": "PBehaviors",
        "groups": ["api", "api_planning"],
        "name": "api_pbehavior",
        "type": "CRUD"
    },
    {
        "_id": "api_pbehaviortype",
        "description": "PBehavior types",
        "groups": ["api", "api_planning"],
        "name": "api_pbehaviortype",
        "type": "CRUD"
    },
    {
        "_id": "api_pbehaviorreason",
        "description": "PBehavior reasons",
        "groups": ["api", "api_planning"],
        "name": "api_pbehaviorreason",
        "type": "CRUD"
    },
    {
        "_id": "api_pbehaviorexception",
        "description": "PBehavior exceptions",
        "groups": ["api", "api_planning"],
        "name": "api_pbehaviorexception",
        "type": "CRUD"
    },
    {
        "_id": "models_permission",
        "description": "Rights",
        "groups": ["technical", "technical_admin", "technical_admin_access"],
        "name": "models_permission",
        "type": "CRUD",
        "api_permissions": {
            "api_acl": 0
        }
    },
    {
        "_id": "models_role",
        "description": "Roles",
        "groups": ["technical", "technical_admin", "technical_admin_access"],
        "name": "models_role",
        "type": "CRUD",
        "api_permissions": {
            "api_acl": 0
        }
    },
    {
        "_id": "models_user",
        "description": "Users",
        "groups": ["technical", "technical_admin", "technical_admin_access"],
        "name": "models_user",
        "type": "CRUD",
        "api_permissions": {
            "api_acl": 0
        }
    },
    {
        "_id": "models_broadcastMessage",
        "description": "Broadcast messages",
        "groups": ["technical", "technical_admin", "technical_admin_communication"],
        "name": "models_broadcastMessage",
        "type": "CRUD",
        "api_permissions": {
            "api_broadcast_message": 0
        }
    },
    {
        "_id": "models_playlist",
        "description": "Playlists",
        "groups": ["technical", "technical_admin", "technical_admin_communication"],
        "name": "models_playlist",
        "type": "CRUD",
        "api_permissions": {
            "api_playlist": 0
        }
    },
    {
        "_id": "models_healthcheck",
        "description": "Healthcheck",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_healthcheck",
        "api_permissions": {
            "api_healthcheck": 1,
            "api_message_rate_stats_read": 1
        }
    },
    {
        "_id": "models_healthcheckStatus",
        "description": "Healthcheck status",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_healthcheckStatus",
        "api_permissions": {
            "api_healthcheck": 1
        }
    },
    {
        "_id": "models_techmetrics",
        "description": "Healthcheck - engines' metrics",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_techmetrics",
        "api_permissions": {
            "api_techmetrics": 1
        }
    },
    {
        "_id": "models_remediationInstruction",
        "description": "Instructions - instructions tab",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_remediationInstruction",
        "type": "CRUD",
        "api_permissions": {
            "api_instruction": 0,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_remediationJob",
        "description": "Instructions - jobs tab",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_remediationJob",
        "type": "CRUD",
        "api_permissions": {
            "api_job": 0
        }
    },
    {
        "_id": "models_remediationConfiguration",
        "description": "Instructions - configurations tab",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_remediationConfiguration",
        "type": "CRUD",
        "api_permissions": {
            "api_job_config": 0
        }
    },
    {
        "_id": "models_remediationStatistic",
        "description": "Instructions - remediation statistics tab",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_remediationStatistic",
        "api_permissions": {
            "api_metrics": 1,
            "api_instruction": 4
        }
    },
    {
        "_id": "models_kpi",
        "description": "KPI graphs",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_kpi",
        "api_permissions": {
            "api_kpi_filter": 4,
            "api_rating_settings": 4,
            "api_healthcheck": 1,
            "api_metrics": 1
        }
    },
    {
        "_id": "models_kpiFilters",
        "description": "KPI filters",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
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
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_kpiRatingSettings",
        "type": "CRUD",
        "api_permissions": {
            "api_rating_settings": 0
        }
    },
    {
        "_id": "models_kpiCollectionSettings",
        "description": "KPI collection settings",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_kpiCollectionSettings",
        "type": "CRUD",
        "api_permissions": {
            "api_metrics_settings": 0
        }
    },
    {
        "_id": "models_maintenance",
        "description": "Maintenance",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_maintenance",
        "api_permissions": {
            "api_maintenance": 1
        }
    },
    {
        "_id": "models_map",
        "description": "Maps",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_map",
        "type": "CRUD",
        "api_permissions": {
            "api_map": 0,
            "api_entity": 4
        }
    },
    {
        "_id": "models_parameters",
        "description": "Parameters - parameters tab",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
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
        "_id": "models_widgetTemplate",
        "description": "Parameters - widget templates",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_widgetTemplate",
        "type": "CRUD",
        "api_permissions": {
            "api_dynamicinfos": 4,
            "api_widgettemplate": 0,
            "api_alarm_read": 1,
            "api_entity": 4
        }
    },
    {
        "_id": "models_icon",
        "description": "Parameters - icons",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
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
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_view_import_export",
        "api_permissions": {
            "api_viewgroup": 3,
            "api_view": 3
        }
    },
    {
        "_id": "models_notification",
        "description": "Parameters - notification settings",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_notification",
        "api_permissions": {
            "api_notification": 1
        }
    },
    {
        "_id": "models_planningType",
        "description": "Planning type (Pbehavior)",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_planningType",
        "type": "CRUD",
        "api_permissions": {
            "api_pbehaviortype": 0
        }
    },
    {
        "_id": "models_planningReason",
        "description": "Planning reason (Pbehavior)",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_planningReason",
        "type": "CRUD",
        "api_permissions": {
            "api_pbehaviorreason": 0
        }
    },
    {
        "_id": "models_planningExceptions",
        "description": "Planning exceptions dates (Pbehavior)",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_planningExceptions",
        "type": "CRUD",
        "api_permissions": {
            "api_pbehaviorexception": 0
        }
    },
    {
        "_id": "models_stateSetting",
        "description": "State settings",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_stateSetting",
        "type": "CRUD",
        "api_permissions": {
            "api_state_settings": 0
        }
    },
    {
        "_id": "models_storageSettings",
        "description": "Storage settings",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_storageSettings",
        "type": "RW",
        "api_permissions_bitmask": {
            "4": {
                "api_datastorage_read": 1
            },
            "2": {
                "api_datastorage_update": 1
            }
        }
    },
    {
        "_id": "models_tag",
        "description": "Tags management",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
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
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_eventsRecord",
        "api_permissions": {
            "api_launch_event_recording": 1,
            "api_resend_events": 1
        }
    },
    {
        "_id": "models_privateView",
        "description": "Private views",
        "groups": ["technical", "technical_viewsandwidgets"],
        "name": "models_privateView",
        "api_permissions": {
            "api_rating_settings": 4,
            "api_junit": 4,
            "api_alarm_read": 1,
            "api_state_settings": 4,
            "api_metrics": 1,
            "api_pbehaviorreason": 4,
            "api_pbehaviortype": 4,
            "api_entity": 4,
            "api_private_view_groups": 1,
            "api_entityservice": 4,
            "api_dynamicinfos": 4,
            "api_map": 4,
            "api_pbehavior": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "models_userview",
        "description": "Views",
        "groups": ["technical", "technical_viewsandwidgets"],
        "name": "models_userview",
        "type": "CRUD",
        "api_permissions": {
            "api_view": 0,
            "api_viewgroup": 0
        }
    },
    {
        "_id": "models_shareToken",
        "description": "Shared token settings",
        "groups": ["technical", "technical_token"],
        "name": "models_shareToken",
        "type": "CRUD",
        "api_permissions_bitmask": {
            "8": {
                "api_share_token": 8
            }
        }
    },
    {
        "_id": "models_exploitation_eventFilter",
        "description": "Event filters",
        "groups": ["technical", "technical_exploitation"],
        "name": "models_exploitation_eventFilter",
        "type": "CRUD",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4,
            "api_eventfilter": 0
        }
    },
    {
        "_id": "models_exploitation_pbehavior",
        "description": "Pbehaviors",
        "groups": ["technical", "technical_exploitation"],
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
        "groups": ["technical", "technical_exploitation"],
        "name": "models_exploitation_snmpRule",
        "type": "CRUD",
        "api_permissions": {
            "api_snmprule": 0
        }
    },
    {
        "_id": "models_exploitation_dynamicInfo",
        "description": "Dynamic information rules",
        "groups": ["technical", "technical_exploitation"],
        "name": "models_exploitation_dynamicInfo",
        "type": "CRUD",
        "api_permissions": {
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4,
            "api_dynamicinfos": 0
        }
    },
    {
        "_id": "models_exploitation_metaAlarmRule",
        "description": "Meta alarm rules and correlation",
        "groups": ["technical", "technical_exploitation"],
        "name": "models_exploitation_metaAlarmRule",
        "type": "CRUD",
        "api_permissions": {
            "api_entitycategory": 4,
            "api_metaalarmrule": 0,
            "api_alarm_read": 1,
            "api_entity": 4
        }
    },
    {
        "_id": "models_exploitation_scenario",
        "description": "Scenarios",
        "groups": ["technical", "technical_exploitation"],
        "name": "models_exploitation_scenario",
        "type": "CRUD",
        "api_permissions": {
            "api_entitycategory": 4,
            "api_action": 0,
            "api_alarm_read": 1,
            "api_entity": 4
        }
    },
    {
        "_id": "models_exploitation_idleRules",
        "description": "Idle rules",
        "groups": ["technical", "technical_exploitation"],
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
        "groups": ["technical", "technical_exploitation"],
        "name": "models_exploitation_flappingRules",
        "type": "CRUD",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4,
            "api_flapping_rule": 0,
            "api_alarm_read": 1
        }
    },
    {
        "_id": "models_exploitation_resolveRules",
        "description": "Resolve rules",
        "groups": ["technical", "technical_exploitation"],
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
        "groups": ["technical", "technical_exploitation"],
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
        "groups": ["technical", "technical_exploitation"],
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
        "_id": "models_notification_instructionStats",
        "description": "Instructions statistics",
        "groups": ["technical", "technical_notification"],
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
        "groups": ["technical", "technical_profile"],
        "name": "models_profile_corporatePattern",
        "type": "CRUD",
        "api_permissions": {
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4,
            "api_pbehavior": 4,
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4
        }
    },
    {
        "_id": "models_profile_color_theme",
        "description": "Theme colors",
        "groups": ["technical", "technical_profile"],
        "name": "models_profile_color_theme",
        "type": "CRUD",
        "api_permissions": {
            "api_color_theme": 0
        }
    },
    {
        "_id": "common_variablesHelp",
        "description": "See the list of variables (in all widgets)",
        "groups": ["widgets", "widgets_common"],
        "name": "common_variablesHelp"
    },
    {
        "_id": "common_entityComment",
        "description": "Manage entity comments (view, create, edit, delete)",
        "groups": ["widgets", "widgets_common"],
        "name": "common_entityComment",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycomment": 1
        }
    },
    {
        "_id": "listalarm_filter",
        "description": "Set alarm filters",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_widgetsettings"],
        "name": "listalarm_filter",
        "api_permissions": {
            "api_pbehavior": 4,
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "listalarm_remediationInstructionsFilter",
        "description": "Set filters by remediation instructions",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_widgetsettings"],
        "name": "listalarm_remediationInstructionsFilter",
        "api_permissions": {
            "api_instruction": 4
        }
    },
    {
        "_id": "listalarm_userFilter",
        "description": "Filter alarms",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_viewsettings"],
        "name": "listalarm_userFilter",
        "api_permissions": {
            "api_view": 4,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4,
            "api_pbehavior": 4,
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4
        }
    },
    {
        "_id": "listalarm_userRemediationInstructionsFilter",
        "description": "Filter alarms by remediation instructions",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_viewsettings"],
        "name": "listalarm_userRemediationInstructionsFilter",
        "api_permissions": {
            "api_instruction": 4
        }
    },
    {
        "_id": "listalarm_category",
        "description": "Filter alarms by category",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_viewsettings"],
        "name": "listalarm_category",
        "api_permissions": {
            "api_entitycategory": 4
        }
    },
    {
        "_id": "listalarm_filterByBookmark",
        "description": "Filter bookmarked alarms",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_viewsettings"],
        "name": "listalarm_filterByBookmark"
    },
    {
        "_id": "listalarm_correlation",
        "description": "Group correlated alarms (meta alarms)",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_viewsettings"],
        "name": "listalarm_correlation"
    },
    {
        "_id": "listalarm_ack",
        "description": "Ack",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_ack",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_fastAck",
        "description": "Fast ack",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_fastAck",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_cancelAck",
        "description": "Cancel ack",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_cancelAck",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_pbehavior",
        "description": "Edit PBhaviours for alarm",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_pbehavior",
        "api_permissions": {
            "api_pbehavior": 8
        }
    },
    {
        "_id": "listalarm_fastPbehavior",
        "description": "Fast create PBehavior",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_fastPbehavior",
        "api_permissions": {
            "api_pbehavior": 8
        }
    },
    {
        "_id": "listalarm_snoozeAlarm",
        "description": "Snooze",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_snoozeAlarm",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_declareanIncident",
        "description": "Declare ticket",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_declareanIncident",
        "api_permissions": {
            "api_declare_ticket_execution": 1,
            "api_declare_ticket_rule": 4
        }
    },
    {
        "_id": "listalarm_assignTicketNumber",
        "description": "Associate ticket",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_assignTicketNumber",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_removeAlarm",
        "description": "Cancel",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_removeAlarm",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_fastRemoveAlarm",
        "description": "Fast cancel",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_fastRemoveAlarm",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_unCancel",
        "description": "Uncancel",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_unCancel",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_changeState",
        "description": "Check and lock severity",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_changeState",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_history",
        "description": "View entity history",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_history",
        "api_permissions": {
            "api_alarm_read": 1
        }
    },
    {
        "_id": "listalarm_manualMetaAlarmGroup",
        "description": "Link to / unlink from manual meta alarm",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_manualMetaAlarmGroup",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_metaAlarmGroup",
        "description": "Unlink from meta alarm created by rule",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_metaAlarmGroup",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_comment",
        "description": "Comment",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_comment",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_bookmark",
        "description": "Add / remove bookmark",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_bookmark",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "listalarm_links",
        "description": "Follow links in alarm",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_links"
    },
    {
        "_id": "listalarm_executeInstruction",
        "description": "Execute manual instructions",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_executeInstruction",
        "api_permissions": {
            "api_execution": 1
        }
    },
    {
        "_id": "listalarm_exportPdf",
        "description": "Export in PDF",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        "name": "listalarm_exportPdf"
    },
    {
        "_id": "listalarm_exportAsCsv",
        "description": "Export alarm list as CSV file",
        "groups": ["widgets", "widgets_alarmslist", "widgets_alarmslist_actions"],
        "name": "listalarm_exportAsCsv"
    },
    {
        "_id": "crudcontext_filter",
        "description": "Set entity filters",
        "groups": ["widgets", "widgets_context", "widgets_context_widgetsettings"],
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
        "groups": ["widgets", "widgets_context", "widgets_context_viewsettings"],
        "name": "crudcontext_userFilter",
        "api_permissions": {
            "api_entitycategory": 4,
            "api_pbehavior": 4,
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4,
            "api_view": 4,
            "api_alarm_read": 1,
            "api_entity": 4
        }
    },
    {
        "_id": "crudcontext_category",
        "description": "Filter entities by category",
        "groups": ["widgets", "widgets_context", "widgets_context_viewsettings"],
        "name": "crudcontext_category",
        "api_permissions": {
            "api_entitycategory": 4
        }
    },
    {
        "_id": "crudcontext_createEntity",
        "description": "Create entity",
        "groups": ["widgets", "widgets_context", "widgets_context_entityactions"],
        "name": "crudcontext_createEntity",
        "api_permissions": {
            "api_entityservice": 8
        }
    },
    {
        "_id": "crudcontext_edit",
        "description": "Edit entity",
        "groups": ["widgets", "widgets_context", "widgets_context_entityactions"],
        "name": "crudcontext_edit",
        "api_permissions": {
            "api_entity": 2,
            "api_entityservice": 2
        }
    },
    {
        "_id": "crudcontext_duplicate",
        "description": "Duplicate entity",
        "groups": ["widgets", "widgets_context", "widgets_context_entityactions"],
        "name": "crudcontext_duplicate",
        "api_permissions": {
            "api_entityservice": 8
        }
    },
    {
        "_id": "crudcontext_delete",
        "description": "Delete entity",
        "groups": ["widgets", "widgets_context", "widgets_context_entityactions"],
        "name": "crudcontext_delete",
        "api_permissions": {
            "api_entity": 1,
            "api_entityservice": 1
        }
    },
    {
        "_id": "crudcontext_pbehavior",
        "description": "Create / edit pbehavior",
        "groups": ["widgets", "widgets_context", "widgets_context_entityactions"],
        "name": "crudcontext_pbehavior",
        "api_permissions": {
            "api_pbehavior": 8
        }
    },
    {
        "_id": "crudcontext_massEnable",
        "description": "Mass action to enable selected entities",
        "groups": ["widgets", "widgets_context", "widgets_context_entityactions"],
        "name": "crudcontext_massEnable",
        "api_permissions": {
            "api_entity": 2
        }
    },
    {
        "_id": "crudcontext_massDisable",
        "description": "Mass action to disable selected entities",
        "groups": ["widgets", "widgets_context", "widgets_context_entityactions"],
        "name": "crudcontext_massDisable",
        "api_permissions": {
            "api_entity": 2
        }
    },
    {
        "_id": "crudcontext_exportAsCsv",
        "description": "Export entities as CSV file",
        "groups": ["widgets", "widgets_context", "widgets_context_actions"],
        "name": "crudcontext_exportAsCsv"
    },
    {
        "_id": "serviceweather_filter",
        "description": "Set entity filters",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_widgetsettings"],
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
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_viewsettings"],
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
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_viewsettings"],
        "name": "serviceweather_category",
        "api_permissions": {
            "api_entitycategory": 4
        }
    },
    {
        "_id": "serviceweather_entityAck",
        "description": "Ack",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_entityAck",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "serviceweather_entityAssocTicket",
        "description": "Associate ticket",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_entityAssocTicket",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "serviceweather_entityComment",
        "description": "Comment alarm",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_entityComment",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "serviceweather_entityValidate",
        "description": "Validate alarms and change their state to critical",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_entityValidate",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "serviceweather_entityInvalidate",
        "description": "Invalidate alarms and cancel them",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_entityInvalidate",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "serviceweather_entityPause",
        "description": "Pause alarms (set the PBehavior type \"Pause\")",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_entityPause",
        "api_permissions": {
            "api_pbehavior": 8
        }
    },
    {
        "_id": "serviceweather_entityPlay",
        "description": "Activate paused alarms (remove the PBehavior type \"Pause\")",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_entityPlay",
        "api_permissions": {
            "api_pbehavior": 1
        }
    },
    {
        "_id": "serviceweather_entityCancel",
        "description": "Cancel",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_entityCancel",
        "api_permissions": {
            "api_alarm_update": 1
        }
    },
    {
        "_id": "serviceweather_entityManagePbehaviors",
        "description": "View PBehaviors of services (in the subtab in the services modal windows)",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_entityManagePbehaviors",
        "api_permissions": {
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4,
            "api_pbehaviorexception": 4,
            "api_pbehavior": 7
        }
    },
    {
        "_id": "serviceweather_entityLinks",
        "description": "Follow links in alarm",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_entityLinks"
    },
    {
        "_id": "serviceweather_entityDeclareTicket",
        "description": "Declare ticket",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
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
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_moreInfos",
        "api_permissions": {
            "api_alarm_read": 1
        }
    },
    {
        "_id": "serviceweather_alarmsList",
        "description": "Open the list of alarms available for each service",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_alarmsList",
        "api_permissions": {
            "api_alarm_read": 1,
            "api_entityservice": 4,
            "api_state_settings": 4,
            "api_entity": 4,
            "api_pbehavior": 4,
            "api_junit": 4
        }
    },
    {
        "_id": "serviceweather_pbehaviorList",
        "description": "View PBehaviors of services (in the subtab in the service entities modal windows)",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_pbehaviorList",
        "api_permissions": {
            "api_pbehavior": 15,
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4,
            "api_pbehaviorexception": 4
        }
    },
    {
        "_id": "serviceweather_executeInstruction",
        "description": "Execute manual instructions",
        "groups": ["widgets", "widgets_serviceweather", "widgets_serviceweather_alarmactions"],
        "name": "serviceweather_executeInstruction",
        "api_permissions": {
            "api_execution": 1
        }
    },
    {
        "_id": "counter_alarmsList",
        "description": "View the alarm list associated with counters",
        "groups": ["widgets", "widgets_counter"],
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
        "groups": ["widgets", "widgets_testingweather"],
        "name": "testingweather_alarmsList",
        "api_permissions": {
            "api_alarm_read": 1
        }
    },
    {
        "_id": "map_filter",
        "description": "Set filters for points on maps",
        "groups": ["widgets", "widgets_map", "widgets_map_widgetsettings"],
        "name": "map_filter",
        "api_permissions": {
            "api_pbehavior": 4,
            "api_pbehaviortype": 4,
            "api_pbehaviorreason": 4,
            "api_alarm_read": 1,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "map_alarmsList",
        "description": "View the alarm list associated with points on maps",
        "groups": ["widgets", "widgets_map", "widgets_map_viewsettings"],
        "name": "map_alarmsList",
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
        "_id": "map_userFilter",
        "description": "Filter points on maps",
        "groups": ["widgets", "widgets_map", "widgets_map_viewsettings"],
        "name": "map_userFilter",
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
        "_id": "map_category",
        "description": "Filter points on maps by categories",
        "groups": ["widgets", "widgets_map", "widgets_map_viewsettings"],
        "name": "map_category",
        "api_permissions": {
            "api_entitycategory": 4
        }
    },
    {
        "_id": "barchart_filter",
        "description": "Set data filters",
        "groups": ["widgets", "widgets_barchart", "widgets_barchart_widgetsettings"],
        "name": "barchart_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "barchart_interval",
        "description": "Edit time intervals for the data displayed",
        "groups": ["widgets", "widgets_barchart", "widgets_barchart_viewsettings"],
        "name": "barchart_interval"
    },
    {
        "_id": "barchart_sampling",
        "description": "Edit sampling for the data displayed",
        "groups": ["widgets", "widgets_barchart", "widgets_barchart_viewsettings"],
        "name": "barchart_sampling"
    },
    {
        "_id": "barchart_userFilter",
        "description": "Filter data",
        "groups": ["widgets", "widgets_barchart", "widgets_barchart_viewsettings"],
        "name": "barchart_userFilter",
        "api_permissions": {
            "api_view": 4,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "linechart_filter",
        "description": "Set data filters",
        "groups": ["widgets", "widgets_linechart", "widgets_linechart_widgetsettings"],
        "name": "linechart_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "linechart_interval",
        "description": "Edit time intervals for the data displayed",
        "groups": ["widgets", "widgets_linechart", "widgets_linechart_viewsettings"],
        "name": "linechart_interval"
    },
    {
        "_id": "linechart_sampling",
        "description": "Edit sampling for the data displayed",
        "groups": ["widgets", "widgets_linechart", "widgets_linechart_viewsettings"],
        "name": "linechart_sampling"
    },
    {
        "_id": "linechart_userFilter",
        "description": "Filter data",
        "groups": ["widgets", "widgets_linechart", "widgets_linechart_viewsettings"],
        "name": "linechart_userFilter",
        "api_permissions": {
            "api_view": 4,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "piechart_filter",
        "description": "Set data filters",
        "groups": ["widgets", "widgets_piechart", "widgets_piechart_widgetsettings"],
        "name": "piechart_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "piechart_interval",
        "description": "Edit time intervals for the data displayed",
        "groups": ["widgets", "widgets_piechart", "widgets_piechart_viewsettings"],
        "name": "piechart_interval"
    },
    {
        "_id": "piechart_sampling",
        "description": "Edit sampling for the data displayed",
        "groups": ["widgets", "widgets_piechart", "widgets_piechart_viewsettings"],
        "name": "piechart_sampling"
    },
    {
        "_id": "piechart_userFilter",
        "description": "Filter data",
        "groups": ["widgets", "widgets_piechart", "widgets_piechart_viewsettings"],
        "name": "piechart_userFilter",
        "api_permissions": {
            "api_view": 4,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "numbers_filter",
        "description": "Set data filters",
        "groups": ["widgets", "widgets_numbers", "widgets_numbers_widgetsettings"],
        "name": "numbers_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "numbers_interval",
        "description": "Edit time intervals for the data displayed",
        "groups": ["widgets", "widgets_numbers", "widgets_numbers_viewsettings"],
        "name": "numbers_interval"
    },
    {
        "_id": "numbers_sampling",
        "description": "Edit sampling for the data displayed",
        "groups": ["widgets", "widgets_numbers", "widgets_numbers_viewsettings"],
        "name": "numbers_sampling"
    },
    {
        "_id": "numbers_userFilter",
        "description": "Filter data",
        "groups": ["widgets", "widgets_numbers", "widgets_numbers_viewsettings"],
        "name": "numbers_userFilter",
        "api_permissions": {
            "api_entitycategory": 4,
            "api_view": 4,
            "api_entity": 4
        }
    },
    {
        "_id": "userStatistics_filter",
        "description": "Set data filters",
        "groups": ["widgets", "widgets_userstatistics", "widgets_userstatistics_widgetsettings"],
        "name": "userStatistics_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "userStatistics_interval",
        "description": "Edit time intervals for the data displayed",
        "groups": ["widgets", "widgets_userstatistics", "widgets_userstatistics_viewsettings"],
        "name": "userStatistics_interval"
    },
    {
        "_id": "userStatistics_userFilter",
        "description": "Filter data",
        "groups": ["widgets", "widgets_userstatistics", "widgets_userstatistics_viewsettings"],
        "name": "userStatistics_userFilter",
        "api_permissions": {
            "api_view": 4,
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "alarmStatistics_filter",
        "description": "Set data filters",
        "groups": ["widgets", "widgets_alarmstatistics", "widgets_alarmstatistics_widgetsettings"],
        "name": "alarmStatistics_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "alarmStatistics_interval",
        "description": "Edit time intervals for the data displayed",
        "groups": ["widgets", "widgets_alarmstatistics", "widgets_alarmstatistics_viewsettings"],
        "name": "alarmStatistics_interval"
    },
    {
        "_id": "alarmStatistics_userFilter",
        "description": "Filter data",
        "groups": ["widgets", "widgets_alarmstatistics", "widgets_alarmstatistics_viewsettings"],
        "name": "alarmStatistics_userFilter",
        "api_permissions": {
            "api_entitycategory": 4,
            "api_view": 4,
            "api_entity": 4
        }
    },
    {
        "_id": "availability_filter",
        "description": "Set data filters",
        "groups": ["widgets", "widgets_availability", "widgets_availability_widgetsettings"],
        "name": "availability_filter",
        "api_permissions": {
            "api_entity": 4,
            "api_entitycategory": 4
        }
    },
    {
        "_id": "availability_interval",
        "description": "Edit time intervals for the data displayed",
        "groups": ["widgets", "widgets_availability", "widgets_availability_viewsettings"],
        "name": "availability_interval"
    },
    {
        "_id": "availability_userFilter",
        "description": "Filter data",
        "groups": ["widgets", "widgets_availability", "widgets_availability_viewsettings"],
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
        "groups": ["widgets", "widgets_availability", "widgets_availability_actions"],
        "name": "availability_exportAsCsv"
    },
    {
        "_id": "AlarmsList",
        "description": "API permissions for AlarmsList widget",
        "name": "AlarmsList",
        "api_permissions_bitmask": {
            "4": {
                "api_state_settings": 4,
                "api_metrics": 1,
                "api_alarm_read": 1,
                "api_entityservice": 4,
                "api_entity": 4,
                "api_pbehavior": 4,
                "api_junit": 4,
                "api_view": 4,
                "api_associative_table": 4
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
                "api_pbehavior": 4,
                "api_metrics": 1,
                "api_view": 4,
                "api_associative_table": 4,
                "api_entity": 4,
                "api_alarm_read": 1,
                "api_entityservice": 4,
                "api_state_settings": 4
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
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_entityservice": 4,
                "api_alarm_read": 1,
                "api_state_settings": 4
            },
            "2": {
                "api_pbehaviortype": 4,
                "api_view": 2,
                "api_widgettemplate": 4,
                "api_entity": 4,
                "api_dynamicinfos": 4
            },
            "1": {
                "api_view": 1
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
                "api_entityservice": 4,
                "api_state_settings": 4,
                "api_pbehavior": 4,
                "api_junit": 4,
                "api_associative_table": 4,
                "api_alarm_read": 1,
                "api_view": 4,
                "api_entity": 4
            },
            "2": {
                "api_view": 2,
                "api_widgettemplate": 4,
                "api_alarm_read": 1,
                "api_entitycategory": 4,
                "api_pbehaviortype": 4,
                "api_pbehaviorreason": 4,
                "api_dynamicinfos": 4,
                "api_entity": 4,
                "api_pbehavior": 4
            }
        },
        "hidden": true
    },
    {
        "_id": "Text",
        "description": "API permissions for Text widget",
        "name": "Text",
        "api_permissions_bitmask": {
            "1": {
                "api_view": 1
            },
            "4": {
                "api_associative_table": 4,
                "api_view": 4
            },
            "2": {
                "api_view": 2,
                "api_widgettemplate": 4
            }
        },
        "hidden": true
    },
    {
        "_id": "Counter",
        "description": "API permissions for Counter widget",
        "name": "Counter",
        "api_permissions_bitmask": {
            "1": {
                "api_view": 1
            },
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_alarm_read": 1
            },
            "2": {
                "api_entitycategory": 4,
                "api_dynamicinfos": 4,
                "api_alarm_read": 1,
                "api_pbehavior": 4,
                "api_pbehaviortype": 4,
                "api_pbehaviorreason": 4,
                "api_view": 2,
                "api_widgettemplate": 4,
                "api_entity": 4
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
        "_id": "Map",
        "description": "API permissions for Map widget",
        "name": "Map",
        "api_permissions_bitmask": {
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_map": 4,
                "api_entityservice": 4,
                "api_alarm_read": 1,
                "api_entity": 4,
                "api_pbehavior": 4
            },
            "2": {
                "api_dynamicinfos": 4,
                "api_view": 2,
                "api_widgettemplate": 4,
                "api_alarm_read": 1,
                "api_entity": 4
            },
            "1": {
                "api_view": 1
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
                "api_associative_table": 4,
                "api_metrics": 1,
                "api_view": 4
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
        "_id": "PieChart",
        "description": "API permissions for PieChart widget",
        "name": "PieChart",
        "api_permissions_bitmask": {
            "2": {
                "api_view": 2,
                "api_widgettemplate": 4
            },
            "1": {
                "api_view": 1
            },
            "4": {
                "api_metrics": 1,
                "api_view": 4,
                "api_associative_table": 4
            }
        },
        "hidden": true
    },
    {
        "_id": "Numbers",
        "description": "API permissions for Numbers widget",
        "name": "Numbers",
        "api_permissions_bitmask": {
            "1": {
                "api_view": 1
            },
            "4": {
                "api_view": 4,
                "api_associative_table": 4,
                "api_metrics": 1
            },
            "2": {
                "api_view": 2,
                "api_widgettemplate": 4
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
                "api_metrics": 1,
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
    if (doc.groups !== undefined || doc.hidden !== undefined) {
        delete updatedPermsByID[doc._id];

        return;
    }

    if (doc.view) {
        db.permission.updateOne({_id: doc._id}, {$set: {
            description: doc.view.title,
            groups: ["commonviews", doc.view.group_id],
        }});

        return;
    }

    if (doc.playlist) {
        db.permission.updateOne({_id: doc._id}, {$set: {
            description: doc.playlist.name,
            groups: ["commonviews", "commonviews_playlist"],
        }});

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
