db.permission.deleteMany({
    _id: {
        $in: [
            "api_comment_template",
            "models_commentTemplate"
        ]
    }
});

db.role.updateMany({}, {
    $unset: {
        "permissions.api_comment_template": "",
        "permissions.models_commentTemplate": ""
    }
});
