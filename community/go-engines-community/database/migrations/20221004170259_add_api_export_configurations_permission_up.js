if (!db.default_rights.findOne({_id: "api_export_configurations"})) {
    db.default_rights.insertOne({
        _id: "api_export_configurations",
        crecord_name: "api_export_configurations",
        crecord_type: "action",
        desc: "Export configurations"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.api_export_configurations": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
