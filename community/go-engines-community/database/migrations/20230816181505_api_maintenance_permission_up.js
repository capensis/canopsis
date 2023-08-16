if (!db.permission.findOne({_id: "api_maintenance"})) {
    db.permission.insertOne({
        _id: "api_maintenance",
        crecord_name: "api_maintenance",
        description: "Trigger maintenance mode"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.api_maintenance": 1}});
}
