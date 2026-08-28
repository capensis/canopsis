db.role.updateMany({
        "$and": [
            {
                "$or": [
                    {
                        "permissions.models_exploitation_linkRule": {
                            "$ne": null
                        }
                    },
                    {
                        "permissions.models_exploitation_eventFilter": {
                            "$ne": null
                        }
                    }
                ]
            },
            {
                "$or": [
                    {
                        "permissions.api_external_data_table": null
                    },
                    {
                        "permissions.api_external_data_table": {
                            "$bitsAllClear": [
                                2
                            ]
                        }
                    }
                ]
            }
        ]
    },
    [
        {
            $set: {
                "permissions.api_external_data_table": {$bitOr: [4, {$toInt: {$ifNull: ["$permissions.api_external_data_table", 0]}}]},

            }
        }
    ]
);

db.permission.updateMany({_id: {$in: ["models_exploitation_linkRule", "models_exploitation_eventFilter"]}}, {
    $set: {
        "api_permissions.api_external_data_table": 4
    }
});
