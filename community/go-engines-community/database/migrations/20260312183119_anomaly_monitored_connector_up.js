if (db.configuration.findOne({_id: "data_storage"})) {
    db.configuration.updateOne(
        {
            _id: "data_storage",
            "config.connector_anomalies.delete_after": null
        },
        {
            $set: {
                "config.connector_anomalies.delete_after": {
                    "value": 14,
                    "unit": "d",
                    "enabled": true,
                },
            }
        },
    );
} else {
    db.configuration.insertOne({
        _id: "data_storage",
        "config.connector_anomalies.delete_after": {
            "value": 14,
            "unit": "d",
            "enabled": true,
        },
    });
}


if (!db.permission.findOne({_id: "api_anomaly_monitored_connector"})) {
    db.permission.insertOne({
        _id: "api_anomaly_monitored_connector",
        name: "api_anomaly_monitored_connector",
        type: "CRUD",
        description: "Anomaly monitored connectors",
        groups: ["api", "api_general"]
    });

    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_anomaly_monitored_connector": 15
        }
    });
}


if (!db.permission.findOne({_id: "models_anomalyMonitoredConnector"})) {
    db.permission.insertOne({
        _id: "models_anomalyMonitoredConnector",
        name: "models_anomalyMonitoredConnector",
        type: "CRUD",
        description: "Anomaly monitored connectors",
        groups: ["technical", "technical_admin", "technical_admin_general"],
        api_permissions: {
            api_anomaly_monitored_connector: 0
        }
    });

    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_anomalyMonitoredConnector": 15
        }
    });
}
