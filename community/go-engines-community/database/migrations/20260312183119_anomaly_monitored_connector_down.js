db.permission.deleteOne({_id: "api_anomaly_monitored_connector"});

db.role.updateMany({}, {
    $unset: {
        "permissions.api_anomaly_monitored_connector": "",
    }
});
