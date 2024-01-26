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

if (!db.permission.findOne({_id: "models_stateSetting"})) {
    db.permission.insertOne({
        _id: "models_stateSetting",
        name: "models_stateSetting",
        description: "State settings",
        type: "CRUD"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_stateSetting": 15
        }
    });
}

if (!db.permission.findOne({_id: "models_icon"})) {
    db.permission.insertOne({
        _id: "models_icon",
        name: "models_icon",
        description: "Icons",
        type: "CRUD"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_icon": 15
        }
    });
}

if (!db.permission.findOne({_id: "models_storageSettings"})) {
    db.permission.insertOne({
        _id: "models_storageSettings",
        name: "models_storageSettings",
        description: "Storage settings",
        type: "RW"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_storageSettings": 15
        }
    });
}
