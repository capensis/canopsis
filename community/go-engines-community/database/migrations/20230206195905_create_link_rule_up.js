db.createCollection("link_rule");

if (!db.default_rights.findOne({_id: "api_link_rule"})) {
    db.default_rights.insertOne({
        _id: "api_link_rule",
        crecord_name: "api_link_rule",
        crecord_type: "action",
        description: "Link rule",
        type: "CRUD"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.api_link_rule": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "models_exploitation_linkRule"})) {
    db.default_rights.insertOne({
        _id: "models_exploitation_linkRule",
        crecord_name: "models_exploitation_linkRule",
        crecord_type: "action",
        description: "Exploitation: Link rule",
        type: "CRUD"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.models_exploitation_linkRule": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}
