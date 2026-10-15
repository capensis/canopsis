db.permission.deleteOne({
    _id: "crudcontext_pause"
});

db.role.updateMany({}, {
    $unset: {
        "permissions.crudcontext_pause": ""
    }
});

db.pbehavior_type.updateOne({name: "Default pause"}, {$set: {icon_name: "pause"}});

db.entity_counters.updateMany({}, {$unset: {inactive_depends: "", acked_inactive: "", paused: ""}});

db.default_entities.dropIndex("pause_timestamp_1");
