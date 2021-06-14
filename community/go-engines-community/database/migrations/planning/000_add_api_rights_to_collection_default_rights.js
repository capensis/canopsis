db.default_rights.insertMany([
    {
        _id: "api_entity_read",
        crecord_name: "api_entity_read",
        crecord_type: "action",
        desc: "Read entities",
    },
    {
        _id: "api_entity_update",
        crecord_name: "api_entity_update",
        crecord_type: "action",
        desc: "Update entities",
    },
    {
        _id: "api_entity_delete",
        crecord_name: "api_entity_delete",
        crecord_type: "action",
        desc: "Delete entities",
    },
    {
        _id: "api_alarm_read",
        crecord_name: "api_alarm_read",
        crecord_type: "action",
        desc: "Read alarms",
    },
    {
        _id: "api_alarm_update",
        crecord_name: "api_alarm_update",
        crecord_type: "action",
        desc: "Update alarms",
    },
    {
        _id: "api_alarm_delete",
        crecord_name: "api_alarm_delete",
        crecord_type: "action",
        desc: "Delete alarms",
    },
    {
        _id: "api_alarmfilter",
        crecord_name: "api_alarmfilter",
        crecord_type: "action",
        desc: "Alarm filters",
        type: "CRUD",
    },
    {
        _id: "api_idlerule",
        crecord_name: "api_idlerule",
        crecord_type: "action",
        desc: "Idle rules",
        type: "CRUD",
    },
    {
        _id: "api_eventfilter",
        crecord_name: "api_eventfilter",
        crecord_type: "action",
        desc: "Event filters",
        type: "CRUD",
    },
    {
        _id: "api_action",
        crecord_name: "api_action",
        crecord_type: "action",
        desc: "Actions",
        type: "CRUD",
    },
    {
        _id: "api_metaalarmrule",
        crecord_name: "api_metaalarmrule",
        crecord_type: "action",
        desc: "Meta-alarm rules",
        type: "CRUD",
    },
    {
        _id: "api_playlist",
        crecord_name: "api_playlist",
        crecord_type: "action",
        desc: "Playlists",
        type: "CRUD",
    },
    {
        _id: "api_dynamicinfos",
        crecord_name: "api_dynamicinfos",
        crecord_type: "action",
        desc: "Dynamic infos",
        type: "CRUD",
    },
    {
        _id: "api_heartbeat",
        crecord_name: "api_heartbeat",
        crecord_type: "action",
        desc: "Heartbeats",
        type: "CRUD",
    },
    {
        _id: "api_watcher",
        crecord_name: "api_watcher",
        crecord_type: "action",
        desc: "Watchers",
        type: "CRUD",
    },
    {
        _id: "api_viewgroup",
        crecord_name: "api_viewgroup",
        crecord_type: "action",
        desc: "View groups",
        type: "CRUD",
    },
    {
        _id: "api_view",
        crecord_name: "api_view",
        crecord_type: "action",
        desc: "Views",
        type: "CRUD",
    },
    {
        _id: "api_pbehavior",
        crecord_name: "api_pbehavior",
        crecord_type: "action",
        desc: "PBehaviors",
        type: "CRUD",
    },
    {
        _id: "api_pbehaviortype",
        crecord_name: "api_pbehaviortype",
        crecord_type: "action",
        desc: "PBehaviorTypes",
        type: "CRUD",
    },
    {
        _id: "api_pbehaviorreason",
        crecord_name: "api_pbehaviorreason",
        crecord_type: "action",
        desc: "PBehaviorReasons",
        type: "CRUD",
    },
    {
        _id: "api_pbehaviorexception",
        crecord_name: "api_pbehaviorexception",
        crecord_type: "action",
        desc: "PBehaviorExceptions",
        type: "CRUD",
    },
    {
        _id: "api_event",
        crecord_name: "api_event",
        crecord_type: "action",
        desc: "Event",
    },
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
            "api_pbehaviorexception"
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
