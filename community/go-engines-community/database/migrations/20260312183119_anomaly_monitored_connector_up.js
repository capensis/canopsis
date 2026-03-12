if (!db.permission.findOne({_id: "api_anomaly_monitored_connector"})) {
    db.permission.insertOne({
        _id: "api_anomaly_monitored_connector",
        name: "api_anomaly_monitored_connector",
        type: "CRUD",
        description: "Anomaly monitored connector",
        groups: ["api", "api_general"]
    });

    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_anomaly_monitored_connector": 15
        }
    });
}
