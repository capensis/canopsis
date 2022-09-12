if (!db.default_rights.findOne({_id: "api_corporate_pattern"})) {
    db.default_rights.insertOne({
        _id: "api_corporate_pattern",
        crecord_name: "api_corporate_pattern",
        crecord_type: "action",
        desc: "Corporate patterns"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.api_corporate_pattern": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "models_profile_corporatePattern"})) {
    db.default_rights.insertOne({
        _id: "models_profile_corporatePattern",
        crecord_name: "models_profile_corporatePattern",
        crecord_type: "action",
        desc: "Profile: Corporate patterns",
        type: "CRUD"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.models_profile_corporatePattern": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}
