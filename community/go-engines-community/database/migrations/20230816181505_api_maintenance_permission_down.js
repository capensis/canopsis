db.permission.deleteMany({_id: {$in: ["api_maintenance", "models_maintenance"]}});
db.role.updateMany(
    {},
    {
        $unset: {
            "permissions.api_maintenance": "",
            "permissions.models_maintenance": ""
        }
    },
);
