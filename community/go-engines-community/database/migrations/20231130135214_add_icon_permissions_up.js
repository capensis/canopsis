db.createCollection("icon")

if (!db.permission.findOne({_id: "api_icon"})) {
    db.permission.insertOne({
        _id: "api_icon",
        name: "api_icon",
        description: "Create icons"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_icon": 1
        }
    });
}
