db.default_rights.insertMany([
    {
        _id: "api_entity_read",
        loader_id: "api_entity_read",
        crecord_name: "api_entity_read",
        crecord_type: "action",
        desc: "Read entities",
    },
    {
        _id: "api_entity_update",
        loader_id: "api_entity_update",
        crecord_name: "api_entity_update",
        crecord_type: "action",
        desc: "Update entities",
    },
    {
        _id: "api_entity_delete",
        loader_id: "api_entity_delete",
        crecord_name: "api_entity_delete",
        crecord_type: "action",
        desc: "Delete entities",
    },
    {
        _id: "api_alarm_read",
        loader_id: "api_alarm_read",
        crecord_name: "api_alarm_read",
        crecord_type: "action",
        desc: "Read alarms",
    },
    {
        _id: "api_alarm_update",
        loader_id: "api_alarm_update",
        crecord_name: "api_alarm_update",
        crecord_type: "action",
        desc: "Update alarms",
    },
    {
        _id: "api_alarm_delete",
        loader_id: "api_alarm_delete",
        crecord_name: "api_alarm_delete",
        crecord_type: "action",
        desc: "Delete alarms",
    },
    {
        _id: "api_alarmfilter",
        loader_id: "api_alarmfilter",
        crecord_name: "api_alarmfilter",
        crecord_type: "action",
        desc: "Alarm filters",
        type: "CRUD",
    },
    {
        _id: "api_idlerule",
        loader_id: "api_idlerule",
        crecord_name: "api_idlerule",
        crecord_type: "action",
        desc: "Idle rules",
        type: "CRUD",
    },
    {
        _id: "api_eventfilter",
        loader_id: "api_eventfilter",
        crecord_name: "api_eventfilter",
        crecord_type: "action",
        desc: "Event filters",
        type: "CRUD",
    },
    {
        _id: "api_action",
        loader_id: "api_action",
        crecord_name: "api_action",
        crecord_type: "action",
        desc: "Actions",
        type: "CRUD",
    },
    {
        _id: "api_metaalarmrule",
        loader_id: "api_metaalarmrule",
        crecord_name: "api_metaalarmrule",
        crecord_type: "action",
        desc: "Meta-alarm rules",
        type: "CRUD",
    },
    {
        _id: "api_playlist",
        loader_id: "api_playlist",
        crecord_name: "api_playlist",
        crecord_type: "action",
        desc: "Playlists",
        type: "CRUD",
    },
    {
        _id: "api_dynamicinfos",
        loader_id: "api_dynamicinfos",
        crecord_name: "api_dynamicinfos",
        crecord_type: "action",
        desc: "Dynamic infos",
        type: "CRUD",
    },
    {
        _id: "api_watcher",
        loader_id: "api_watcher",
        crecord_name: "api_watcher",
        crecord_type: "action",
        desc: "Watchers",
        type: "CRUD",
    },
    {
        _id: "api_heartbeat",
        loader_id: "api_heartbeat",
        crecord_name: "api_heartbeat",
        crecord_type: "action",
        desc: "Heartbeats",
        type: "CRUD",
    },
    {
        _id: "api_viewgroup",
        loader_id: "api_viewgroup",
        crecord_name: "api_viewgroup",
        crecord_type: "action",
        desc: "View groups",
        type: "CRUD",
    },
    {
        _id: "api_view",
        loader_id: "api_view",
        crecord_name: "api_view",
        crecord_type: "action",
        desc: "Views",
        type: "CRUD",
    },
    {
        _id: "api_pbehavior",
        loader_id: "api_pbehavior",
        crecord_name: "api_pbehavior",
        crecord_type: "action",
        desc: "PBehaviors",
        type: "CRUD",
    },
    {
        _id: "api_pbehaviortype",
        loader_id: "api_pbehaviortype",
        crecord_name: "api_pbehaviortype",
        crecord_type: "action",
        desc: "PBehaviorTypes",
        type: "CRUD",
    },
    {
        _id: "api_pbehaviorreason",
        loader_id: "api_pbehaviorreason",
        crecord_name: "api_pbehaviorreason",
        crecord_type: "action",
        desc: "PBehaviorReasons",
        type: "CRUD",
    },
    {
        _id: "api_pbehaviorexception",
        loader_id: "api_pbehaviorexception",
        crecord_name: "api_pbehaviorexception",
        crecord_type: "action",
        desc: "PBehaviorExceptions",
        type: "CRUD",
    },
    {
        _id: "api_event",
        loader_id: "api_event",
        crecord_name: "api_event",
        crecord_type: "action",
        desc: "Event",
    },
    {
        _id: "api_engine",
        loader_id: "api_engine",
        crecord_name: "api_engine",
        crecord_type: "action",
        desc: "Engine Info",
    },
    /* Remediation rights */
    {
        _id: "api_execution",
        loader_id: "api_execution",
        crecord_name: "api_execution",
        crecord_type: "action",
        desc: "Runs instructions",
    },
    {
        _id: "api_job_config",
        loader_id: "api_job_config",
        crecord_name: "api_job_config",
        crecord_type: "action",
        desc: "Job configs",
        type: "CRUD",
    },
    {
        _id: "api_job",
        loader_id: "api_job",
        crecord_name: "api_job",
        crecord_type: "action",
        desc: "Jobs",
        type: "CRUD",
    },
    {
        _id: "api_instruction",
        loader_id: "api_instruction",
        crecord_name: "api_instruction",
        crecord_type: "action",
        desc: "Instructions",
        type: "CRUD",
    },
    {
        _id: "api_file",
        loader_id: "api_file",
        crecord_name: "api_file",
        crecord_type: "action",
        desc: "File",
        type: "CRUD",
    }
]);
db.default_rights.find({
    "crecord_name": {
        "$in": [
            "api_entity_read",
            "api_entity_update",
            "api_entity_delete",
            "api_alarm_read",
            "api_alarm_update",
            "api_alarm_delete",
            "api_event",
            "api_engine",
            "api_execution",
        ],
    }
}).forEach(function (doc) {
    db.default_rights.update(
        {
            crecord_name: "admin",
            crecord_type: "role",
        },
        {
            $set: {
                ['rights.' + doc._id]: {
                    checksum: 1,
                    crecord_type: "right",
                },
            },
        }
    )
});
db.default_rights.find({
    "crecord_name": {
        "$in": [
            "api_alarmfilter",
            "api_idlerule",
            "api_eventfilter",
            "api_action",
            "api_webhook",
            "api_metaalarmrule",
            "api_playlist",
            "api_dynamicinfos",
            "api_heartbeat",
            "api_watcher",
            "api_viewgroup",
            "api_view",
            "api_pbehavior",
            "api_pbehaviortype",
            "api_pbehaviorreason",
            "api_pbehaviorexception",
            /* Remediation rights */
            "api_job_config",
            "api_job",
            "api_instruction",
            "api_file"
        ],
    }
}).forEach(function (doc) {
    db.default_rights.update(
        {
            crecord_name: "admin",
            crecord_type: "role",
        },
        {
            $set: {
                ['rights.' + doc._id]: {
                    checksum: 15,
                    crecord_type: "right",
                },
            },
        }
    )
});
db.default_rights.find({
    "crecord_name": {
        "$in": [
            "api_alarmfilter",
            "api_idlerule",
            "api_eventfilter",
            "api_action",
            "api_webhook",
            "api_metaalarmrule",
            "api_playlist",
            "api_dynamicinfos",
            "api_heartbeat",
            "api_watcher",
            "api_viewgroup",
            "api_view",
            "api_pbehavior",
            "api_pbehaviortype",
            "api_pbehaviorreason",
            "api_pbehaviorexception",
            /* Remediation rights */
            "api_job_config",
            "api_job",
            "api_instruction",
            "api_file"
        ],
    }
}).forEach(function (doc) {
    db.default_rights.update(
        {
            crecord_name: "Visualisation",
            crecord_type: "role",
        },
        {
            $set: {
                ['rights.' + doc._id]: {
                    checksum: 4,
                    crecord_type: "right",
                },
            },
        }
    )
});
db.default_rights.find({
    "crecord_name": {
        "$in": [
            "api_playlist",
            "api_viewgroup",
        ],
    }
}).forEach(function (doc) {
    db.default_rights.updateMany(
        {
            crecord_name: { "$in": ["Manager", "Support", "Supervision"] },
            crecord_type: "role",
        },
        {
            $set: {
                ['rights.' + doc._id]: {
                    checksum: 4,
                    crecord_type: "right",
                },
            },
        }
    )
});
db.default_rights.find({
    "crecord_name": {
        "$in": [
            "api_alarm_read"
        ],
    }
}).forEach(function (doc) {
    db.default_rights.updateMany(
        {
            crecord_name: { "$in": ["Manager", "Support", "Visualisation", "Supervision"] },
            crecord_type: "role",
        },
        {
            $set: {
                ['rights.' + doc._id]: {
                    checksum: 1,
                    crecord_type: "right",
                },
            },
        }
    )
});
