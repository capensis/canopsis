db.default_rights.insertMany([
    {
        "_id": "api_metrics",
        "loader_id": "api_metrics",
        "crecord_type": "action",
        "crecord_name": "api_metrics",
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
            "rights.api_metrics": {
                checksum: 1,
                crecord_type: "right",
            },
        },
    }
);
