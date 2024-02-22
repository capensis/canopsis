db.state_settings.createIndex({title: 1}, {name: "title_1", unique: true});

if (!db.state_settings.findOne({_id: "service"})) {
    db.state_settings.insertOne({
        _id: "service",
        title: "Service",
        method: "worst",
        on_top: 2,
        enabled: true
    });
}

db.state_settings.updateOne({_id: "junit"}, {
    $set: {
        title: "Junit",
        on_top: 1,
        enabled: true
    },
    $unset: {
        type: ""
    }
});

if (!db.permission.findOne({_id: "models_stateSetting"})) {
    db.permission.insertOne({
        _id: "models_stateSetting",
        name: "models_stateSetting",
        type: "CRUD",
        description: "State settings"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.models_stateSetting": 15}});
}

db.permission.updateOne({_id:"api_state_settings"}, {$set: {type: "CRUD"}});
db.role.updateMany({"permissions.api_state_settings": 1}, {$set: {"permissions.api_state_settings": 15}});
