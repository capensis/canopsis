if (!db.permission.findOne({_id: "api_entitycomment"})) {
    db.permission.insertOne({
        _id: "api_entitycomment",
        name: "api_entitycomment",
        description: "Entity comments"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_entitycomment": 1
        }
    });
}
