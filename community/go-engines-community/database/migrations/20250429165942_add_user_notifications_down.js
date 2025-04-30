if (db.user_notification_settings.findOne()) {
    db.user_notification_settings.renameCollection("notification");
}

db.user_notification.drop();
