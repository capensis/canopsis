if (db.notification.findOne()) {
    db.notification.renameCollection("user_notification_settings");
}

db.user_notification.createIndex({user: 1, time: 1}, {name: "user_1_time_1"});
db.user_notification.createIndex({roles: 1, time: 1}, {name: "roles_1_time_1"});
db.user_notification.createIndex({"rule._id": 1, type: 1}, {name: "rule._id_1_type_1"});
