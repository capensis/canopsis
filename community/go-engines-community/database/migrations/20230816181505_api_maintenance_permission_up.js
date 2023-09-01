if (!db.permission.findOne({_id: "api_maintenance"})) {
    db.permission.insertOne({
        _id: "api_maintenance",
        crecord_name: "api_maintenance",
        description: "Trigger maintenance mode"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.api_maintenance": 1}});
}

if (!db.permission.findOne({_id: "models_maintenance"})) {
    db.permission.insertOne({
        _id: "models_maintenance",
        crecord_name: "models_maintenance",
        description: "Maintenance"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.models_maintenance": 1}});
}
