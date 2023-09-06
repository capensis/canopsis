db.color_theme.createIndex({name: 1}, {name: "name_1", unique: true});

if (!db.permission.findOne({_id: "api_color_theme"})) {
    db.permission.insertOne({
        _id: "api_color_theme",
        name: "api_color_theme",
        type: "CRUD",
        description: "Api color themes"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.api_color_theme": 15}});
}

if (!db.permission.findOne({_id: "models_color_theme"})) {
    db.permission.insertOne({
        _id: "models_color_theme",
        name: "models_color_theme",
        type: "CRUD",
        description: "Models color themes"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.models_color_theme": 15}});
}
