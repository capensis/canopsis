if (!db.permission.findOne({_id: "crudcontext_pause"})) {
    db.permission.insertOne({
        _id: "crudcontext_pause",
        name: "crudcontext_pause",
        description: "Pause entity",
        groups: ["widgets", "widgets_context", "widgets_context_entityactions"],
        api_permissions: {
            api_entity: 2,
            api_entityservice: 2,
        }
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.crudcontext_pause": 1,
        }
    });
    db.role.updateMany(
        {
            "permissions.crudcontext_pbehavior": 1,
            "permissions.crudcontext_pause": null,
        },
        [
            {
                $set: {
                    "permissions.crudcontext_pause": 1,
                    "permissions.api_entity": {$bitOr: [2, {$toInt: {$ifNull: ["$permissions.api_entity", 0]}}]},
                    "permissions.api_entityservice": {$bitOr: [2, {$toInt: {$ifNull: ["$permissions.api_entityservice", 0]}}]},
                }
            },
        ]
    );
}

db.pbehavior_type.updateOne({name: "Default pause"}, {$set: {icon_name: "$vuetify.icons.pause_action"}});

db.entity_counters.updateMany(
    {inactive_depends: {$exists: false}},
    [
        {
            $set: {
                inactive_depends: {$ifNull: ["$under_pbh", 0]},
                acked_inactive: {$ifNull: ["$acked_under_pbh", 0]},
                paused: 0,
            }
        }
    ]
);

db.default_entities.createIndex(
    {"pause.timestamp": 1},
    {
        name: "pause_timestamp_1",
        partialFilterExpression: {pause: {$exists: true}}
    }
);
