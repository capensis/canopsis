db.permission.deleteMany({
    _id: {
        $in: [
            "listalarm_fastPbehavior",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.listalarm_fastPbehavior": "",
    }
});
