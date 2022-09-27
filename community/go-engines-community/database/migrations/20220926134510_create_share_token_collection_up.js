db.share_token.createIndex({value: 1}, {name: "value_1", unique: true});

if (!db.default_rights.findOne({_id: "api_share_token"})) {
    db.default_rights.insertOne({
        _id: "api_share_token",
        crecord_name: "api_share_token",
        crecord_type: "action",
        desc: "Share token",
        type: "CRUD"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.api_share_token": {
                checksum: 5,
                crecord_type: "right"
            }
        }
    });
}
