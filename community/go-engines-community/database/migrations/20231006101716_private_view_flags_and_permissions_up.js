db.viewgroups.updateMany({}, {$set: {"is_private": false}})
db.views.updateMany({}, {$set: {"is_private": false}})
db.viewtabs.updateMany({}, {$set: {"is_private": false}})
db.widgets.updateMany({}, {$set: {"is_private": false}})
db.widget_filters.updateMany({"is_private": {$in: [null, false]}}, {$set: {"widget_private": false, "is_private": false}})
db.widget_filters.updateMany({"is_private": true}, {$set: {"widget_private": true, "is_private": false}})

if (!db.permission.findOne({_id: "api_private_view_groups"})) {
    db.permission.insertOne({
        _id: "api_private_view_groups",
        name: "api_private_view_groups",
        description: "Create private view groups"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.api_private_view_groups": 1}});
}
