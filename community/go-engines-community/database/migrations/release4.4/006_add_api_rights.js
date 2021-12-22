(function () {
    db.default_rights.insertMany([
        {
            "_id": "api_notification",
            "loader_id": "api_notification",
            "crecord_type": "action",
            "crecord_name": "api_notification",
            "desc": "Notification settings"
        },
        {
            "_id": "api_instruction_approve",
            "loader_id": "api_instruction_approve",
            "crecord_type": "action",
            "crecord_name": "api_instruction_approve",
            "desc": "Instruction approve"
        },
    ]);

    db.default_rights.update(
        {
            crecord_name: "admin",
            crecord_type: "role",
        },
        {
            $set: {
                "rights.api_instruction_approve": {
                    checksum: 1,
                    crecord_type: "right",
                },
                "rights.api_notification": {
                    checksum: 1,
                    crecord_type: "right",
                },
            },
        }
    );
})();
