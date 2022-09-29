db.createCollection("map");

if (!db.default_rights.findOne({_id: "api_map"})) {
    db.default_rights.insertOne({
        _id: "api_map",
        crecord_name: "api_map",
        crecord_type: "action",
        desc: "Map",
        type: "CRUD"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.api_map": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "models_map"})) {
    db.default_rights.insertOne({
        _id: "models_map",
        crecord_name: "models_map",
        crecord_type: "action",
        desc: "Map",
        type: "CRUD"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.models_map": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}
