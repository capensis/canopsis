db.permission.updateOne(
    {
        _id: "models_profile_corporatePattern",
    },
    {
        $unset: {
            "api_permissions.api_corporate_pattern": "",
        },
    },
);
