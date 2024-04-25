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

if (!db.permission.findOne({_id: "common_entityCommentsList"})) {
    db.permission.insertOne({
        _id: "common_entityCommentsList",
        name: "common_entityCommentsList",
        description: "Entity comments list"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.common_entityCommentsList": 1
        }
    });
}

if (!db.permission.findOne({_id: "common_createEntityComment"})) {
    db.permission.insertOne({
        _id: "common_createEntityComment",
        name: "common_createEntityComment",
        description: "Create entity comment"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.common_createEntityComment": 1
        }
    });
}

if (!db.permission.findOne({_id: "common_editEntityComment"})) {
    db.permission.insertOne({
        _id: "common_editEntityComment",
        name: "common_editEntityComment",
        description: "Edit entity comment"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.common_editEntityComment": 1
        }
    });
}
