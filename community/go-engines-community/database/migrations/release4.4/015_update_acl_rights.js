(function () {
    db.default_rights.updateOne({_id: "api_acl"}, {$set: {type: "CRUD"}});
    db.default_rights.updateMany(
        {
            crecord_name: "admin",
            crecord_type: "role",
        },
        {
            $set: {
                "rights.api_acl": {
                    checksum: 15,
                    crecord_type: "right",
                },
            },
        }
    );
    db.default_rights.updateMany(
        {
            crecord_name: {$ne: "admin"},
            crecord_type: "role",
            "rights.api_acl.checksum": 1,
        },
        {
            $set: {
                "rights.api_acl": {
                    checksum: 4,
                    crecord_type: "right",
                },
            },
        }
    );
})();
