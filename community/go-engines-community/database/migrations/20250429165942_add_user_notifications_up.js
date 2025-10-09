if (db.notification.findOne()) {
    db.notification.renameCollection("user_notification_settings");
}

db.user_notification.createIndex({user: 1, time: 1}, {name: "user_1_time_1"});
db.user_notification.createIndex({roles: 1, time: 1}, {name: "roles_1_time_1"});
db.user_notification.createIndex({"rule._id": 1, type: 1}, {name: "rule._id_1_type_1"});

if (!db.permission.findOne({_id: "models_notification_common"})) {
    db.permission.insertOne({
        "_id": "models_notification_common",
        "description": "Parameters - notification settings",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_notification_common",
        "api_permissions": {
            "api_notification": 1
        }
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_notification_common": 1
        }
    });
    db.role.updateMany({"permissions.models_notification": 1}, {
        $set: {
            "permissions.models_notification_common": 1
        }
    });
}

if (!db.permission.findOne({_id: "models_instructionStats"})) {
    db.permission.insertOne({
        "_id": "models_instructionStats",
        "description": "Instructions - instructions stats tab",
        "groups": ["technical", "technical_admin", "technical_admin_general"],
        "name": "models_instructionStats",
        "type": "CRUD",
        "api_permissions_bitmask": {
            "4": {
                "api_instruction": 4
            }
        }
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_instructionStats": 15
        }
    });
    db.role.updateMany({"permissions.models_notification_instructionStats": {$gt: 0}}, [
        {
            $set: {
                "permissions.models_instructionStats": "$permissions.models_notification_instructionStats"
            }
        }
    ]);
}

db.permission.deleteMany({
    _id: {
        $in: [
            "models_notification",
            "models_notification_instructionStats",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.models_notification": "",
        "permissions.models_notification_instructionStats": "",
    }
});
db.permission_group.deleteOne({_id: "technical_notification"});
