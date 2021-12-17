db.default_rights.insertMany([
    {
        "_id": "api_filter",
        "loader_id": "api_filter",
        "crecord_type": "action",
        "crecord_name": "api_filter",
        "desc": "Filters",
        "type": "CRUD"
    },
]);

db.default_rights.update(
    {
        crecord_name: "admin",
        crecord_type: "role",
    },
    {
        $set: {
            "rights.api_filter": {
                checksum: 15,
                crecord_type: "right",
            },
        },
    }
);
