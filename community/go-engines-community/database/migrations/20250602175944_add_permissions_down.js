db.permission.deleteMany({
    _id: {
        $in: [
            "api_pbehavior_patterns",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_pbehavior_patterns": "",
    }
});
