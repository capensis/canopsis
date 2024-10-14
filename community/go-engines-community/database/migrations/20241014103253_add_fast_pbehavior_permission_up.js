if (!db.permission.findOne({_id: "listalarm_fastPbehavior"})) {
    db.permission.insertOne({
        _id: "listalarm_fastPbehavior",
        name: "listalarm_fastPbehavior",
        description: "Rights on listalarm: fast pbehavior",
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.listalarm_fastPbehavior": 1
        }
    });
}
