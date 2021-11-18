db.default_rights.insertMany([
    {
        "_id": "api_rating_settings",
        "loader_id": "api_rating_settings",
        "crecord_type": "action",
        "crecord_name": "api_rating_settings",
        "desc": "Rating settings api",
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
            "rights.api_rating_settings": {
                checksum: 15,
                crecord_type: "right",
            },
        },
    }
);
