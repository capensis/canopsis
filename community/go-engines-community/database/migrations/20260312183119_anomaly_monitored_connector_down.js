db.configuration.updateOne({_id: "data_storage"}, {
    $unset: {
        "config.connector_anomalies": "",
    }
});


db.permission.deleteMany(
    {
        _id: {
            $in: [
                "api_anomaly_monitored_connector",
                "models_anomalyMonitoredConnector"
            ]
        }
    }
);

db.role.updateMany({}, {
    $unset: {
        "permissions.api_anomaly_monitored_connector": "",
        "permissions.models_anomalyMonitoredConnector": "",
    }
});
