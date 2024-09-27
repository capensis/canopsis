if (!db.permission.findOne({_id: "api_techmetrics_settings"})) {
    db.permission.insertOne({
        _id: "api_techmetrics_settings",
        name: "api_techmetrics_settings",
        description: "Tech metrics settings",
        type: "CRUD"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_techmetrics_settings": 15,
        }
    });
}
