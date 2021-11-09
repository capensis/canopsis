db.default_rights.insertMany([
    {
        "_id": "api_rating_settings",
        "crecord_type": "action",
        "crecord_name": "api_rating_settings",
        "desc": "Metrics api"
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
                checksum: 1,
                crecord_type: "right",
            },
        },
    }
);
