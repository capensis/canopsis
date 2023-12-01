db.icon.drop()

db.permission.deleteMany({
    _id: {
        $in: [
            "api_icon",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_icon": "",
    }
});
