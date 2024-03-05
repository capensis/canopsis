db.viewgroups.updateMany({is_private: null}, {$set: {"is_private": false}});
db.views.updateMany({is_private: null}, {$set: {"is_private": false}});
db.viewtabs.updateMany({is_private: null}, {$set: {"is_private": false}});
db.widgets.updateMany({is_private: null}, {$set: {"is_private": false}});
db.widget_filters.updateMany({is_user_preference: null, "is_private": {$in: [null, false]}}, {$set: {"is_user_preference": false, "is_private": false}});
db.widget_filters.updateMany({is_user_preference: null, "is_private": true}, {$set: {"is_user_preference": true, "is_private": false}});

if (!db.permission.findOne({_id: "api_private_view_groups"})) {
    db.permission.insertOne({
        _id: "api_private_view_groups",
        name: "api_private_view_groups",
        description: "Create private view groups"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.api_private_view_groups": 1}});
}

if (!db.permission.findOne({_id: "models_privateView"})) {
    db.permission.insertOne({
        _id: "models_privateView",
        name: "models_privateView",
        description: "Private view"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.models_privateView": 1}});
}
