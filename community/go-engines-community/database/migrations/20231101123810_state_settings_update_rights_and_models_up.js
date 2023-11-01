db.state_settings.createIndex({title: 1}, {name: "title_1", unique: true});

if (!db.state_settings.findOne({_id: "service"})) {
    db.state_settings.insertOne({
        _id: "service",
        title: "Service",
        method: "worst",
        on_top: 2
    });
}

db.state_settings.updateOne({_id: "junit"}, {$set:{title: "Junit", on_top: 1}, $unset:{type: ""}});

db.permission.updateOne({_id:"api_state_settings"}, {$set: {type: "CRUD"}});
db.role.updateMany({"permissions.api_state_settings": {$exists: true}}, {$set: {"permissions.api_state_settings": 15}});
