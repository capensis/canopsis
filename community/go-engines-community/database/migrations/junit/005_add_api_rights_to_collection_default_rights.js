(function () {
    db.default_rights.insertMany([
        {
            "_id": "api_acl",
            "loader_id": "api_acl",
            "crecord_name": "api_acl",
            "crecord_type": "action",
            "desc": "Roles, permissions, users"
        },
        {
            "_id": "api_state_settings",
            "loader_id": "api_state_settings",
            "crecord_name": "api_state_settings",
            "crecord_type": "action",
            "desc": "State settings"
        },
        {
            "_id": "api_junit",
            "loader_id": "api_junit",
            "crecord_name": "api_junit",
            "crecord_type": "action",
            "desc": "JUnit API",
            "type": "CRUD"
        },
        {
            "_id": "api_datastorage_read",
            "loader_id": "api_datastorage_read",
            "crecord_name": "api_datastorage_read",
            "crecord_type": "action",
            "desc": "Data storage settings read"
        },
        {
            "_id": "api_datastorage_update",
            "loader_id": "api_datastorage_update",
            "crecord_name": "api_datastorage_update",
            "crecord_type": "action",
            "desc": "Data storage settings update"
        }
    ]);

    db.default_rights.updateOne(
        {
            crecord_name: "admin",
            crecord_type: "role",
        },
        {
            $set: {
                "rights.api_acl": {
                    checksum: 1,
                    crecord_type: "right",
                },
                "rights.api_state_settings": {
                    checksum: 1,
                    crecord_type: "right",
                },
                "rights.api_junit": {
                    checksum: 15,
                    crecord_type: "right",
                },
                "rights.api_datastorage_read": {
                    checksum: 1,
                    crecord_type: "right",
                },
                "rights.api_datastorage_update": {
                    checksum: 1,
                    crecord_type: "right",
                },
            },
        }
    );

    db.default_rights.updateMany(
        {
            crecord_name: { "$in": ["Manager", "Support", "Visualisation", "Supervision"] },
            crecord_type: "role",
        },
        {
            $set: {
                "rights.api_acl": {
                    checksum: 1,
                    crecord_type: "right",
                },
            },
        }
    )
})();
