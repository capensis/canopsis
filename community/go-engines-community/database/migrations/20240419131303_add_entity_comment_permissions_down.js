db.permission.deleteMany({
    _id: {
        $in: [
            "api_entitycomment",
            "common_entityCommentsList",
            "common_createEntityComment",
            "common_editEntityComment",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_entitycomment": "",
        "permissions.common_entityCommentsList": "",
        "permissions.common_createEntityComment": "",
        "permissions.common_editEntityComment": "",
    }
});
