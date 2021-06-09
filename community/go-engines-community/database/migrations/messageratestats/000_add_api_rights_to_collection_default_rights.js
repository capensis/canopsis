(function () {
    db.default_rights.insertOne({
        "_id": "api_message_rate_stats_read",
        "crecord_type": "action",
        "crecord_name": "api_message_rate_stats_read",
        "desc": "Message rate statistics"
    });
    db.default_rights.updateOne(
        {
            crecord_name: "admin",
            crecord_type: "role",
        },
        {
            $set: {
                "rights.api_message_rate_stats_read": {
                    checksum: 1,
                    crecord_type: "right",
                },
            },
        }
    );
})()