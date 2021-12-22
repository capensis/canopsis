(function () {
    db.default_rights.insertMany([
        {
            "_id": "api_resolve_rule",
            "loader_id": "api_resolve_rule",
            "crecord_name": "api_resolve_rule",
            "crecord_type": "action",
            "desc": "Resolve rule",
            "type": "CRUD"
        },
        {
            "_id": "api_flapping_rule",
            "loader_id": "api_flapping_rule",
            "crecord_name": "api_flapping_rule",
            "crecord_type": "action",
            "desc": "Flapping rule",
            "type": "CRUD"
        },
    ]);
    db.default_rights.updateMany(
        {
            crecord_name: "admin",
            crecord_type: "role",
        },
        {
            $set: {
                "rights.api_resolve_rule": {
                    checksum: 15,
                    crecord_type: "right",
                },
                "rights.api_flapping_rule": {
                    checksum: 15,
                    crecord_type: "right",
                },
            },
        }
    );
})();
