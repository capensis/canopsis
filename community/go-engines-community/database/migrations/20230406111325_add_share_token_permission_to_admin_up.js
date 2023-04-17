db.default_rights.updateOne({_id: "admin"}, {
    $set: {
        "rights.api_share_token": {
            checksum: 15,
            crecord_type: "right"
        }
    }
});
