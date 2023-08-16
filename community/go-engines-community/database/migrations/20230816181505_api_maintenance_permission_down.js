db.permission.deleteOne({_id: "api_maintenance"});
db.role.updateMany(
    {
        _id: "admin"
    },
    {
        $unset: {
            "rights.api_maintenance": 1
        }
    },
);
