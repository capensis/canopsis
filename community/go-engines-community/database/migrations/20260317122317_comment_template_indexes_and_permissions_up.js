if (db.getCollectionNames().includes("comment_template")) {
    db.runCommand({collMod: "comment_template", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("comment_template", {changeStreamPreAndPostImages: {enabled: true}})
}

db.comment_template.createIndex({name: 1}, {name: "name_1", unique: true});

if (!db.permission.findOne({_id: "api_comment_template"})) {
    db.permission.insertOne({
        _id: "api_comment_template",
        name: "api_comment_template",
        type: "CRUD",
        description: "Comment template",
        groups: ["api", "api_general"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_comment_template": 15
        }
    });
}

if (!db.permission.findOne({_id: "models_commentTemplate"})) {
    db.permission.insertOne({
        _id: "models_commentTemplate",
        name: "models_commentTemplate",
        type: "CRUD",
        description: "Comment template",
        groups: ["technical", "technical_admin", "technical_admin_settings"],
        api_permissions: {
            api_comment_template: 0
        }
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_commentTemplate": 15
        }
    });
}

