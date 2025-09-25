if (db.user_notification_settings.findOne()) {
    db.user_notification_settings.renameCollection("notification");
}

db.user_notification.drop();

if (!db.permission.findOne({_id: "models_notification"})) {
    db.permission.insertOne({
        "_id": "models_notification",
        "description": "Parameters - notification settings",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_notification",
        "api_permissions": {
            "api_notification": 1
        }
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_notification": 1
        }
    });
    db.role.updateMany({"permissions.models_notification_common": 1}, {
        $set: {
            "permissions.models_notification": 1
        }
    });
}

if (!db.permission_group.findOne({_id: "technical_notification"})) {
    db.permission_group.insertOne({
        "_id": "technical_notification",
        "name": "technical_notification",
        "position": 54
    });
}

if (!db.permission.findOne({_id: "models_notification_instructionStats"})) {
    db.permission.insertOne({
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
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_notification_instructionStats": 15
        }
    });
    db.role.updateMany({"permissions.models_instructionStats": {$gt: 0}}, [
        {
            $set: {
                "permissions.models_notification_instructionStats": "$permissions.models_instructionStats"
            }
        }
    ]);
}

db.permission.deleteMany({
    _id: {
        $in: [
            "models_notification_common",
            "models_instructionStats",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.models_notification_common": "",
        "permissions.models_instructionStats": "",
    }
});
