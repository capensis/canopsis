if (!db.default_rights.findOne({_id: "api_metrics_settings"})) {
    db.default_rights.insertOne({
        _id: "api_metrics_settings",
        crecord_name: "api_metrics_settings",
        crecord_type: "action",
        description: "Metrics settings",
        type: "CRUD"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.api_metrics_settings": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}
