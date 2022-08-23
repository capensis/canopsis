if (!db.default_rights.findOne({_id: "models_remediationStatistic"})) {
    db.default_rights.insertOne({
        _id: "models_remediationStatistic",
        crecord_name: "models_remediationStatistic",
        crecord_type: "action",
        desc: "Remediation statistics"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.models_remediationStatistic": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
