if (!db.default_rights.findOne({_id: "api_security_read"})) {
    db.default_rights.insertOne({
        _id: "api_security_read",
        crecord_name: "api_security_read",
        crecord_type: "action",
        desc: "API security read"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.api_security_read": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "api_security_update"})) {
    db.default_rights.insertOne({
        _id: "api_security_update",
        crecord_name: "api_security_update",
        crecord_type: "action",
        desc: "API security update"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.api_security_update": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
