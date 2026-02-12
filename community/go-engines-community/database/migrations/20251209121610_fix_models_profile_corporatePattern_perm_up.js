db.permission.updateOne(
    {
        _id: "models_profile_corporatePattern",
        "api_permissions.api_corporate_pattern": null,
    },
    {
        $set: {
            "api_permissions.api_corporate_pattern": 1,
        },
    },
);


db.role.updateMany(
    {
        "permissions.models_profile_corporatePattern": {$gt: 0},
        "permissions.api_corporate_pattern": null,
    },
    {
        $set: {
            "permissions.api_corporate_pattern": 1,
        },
    },
);
