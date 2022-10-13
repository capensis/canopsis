if (!db.default_rights.findOne({_id: "api_techmetrics"})) {
    db.default_rights.insertOne({
        _id: "api_techmetrics",
        crecord_name: "api_techmetrics",
        crecord_type: "action",
        desc: "Tech metrics"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.api_techmetrics": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "models_techmetrics"})) {
    db.default_rights.insertOne({
        _id: "models_techmetrics",
        crecord_name: "models_techmetrics",
        crecord_type: "action",
        desc: "Tech metrics"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.models_techmetrics": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
