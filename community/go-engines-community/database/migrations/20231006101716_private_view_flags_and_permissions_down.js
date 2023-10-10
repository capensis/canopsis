db.viewgroups.updateMany({}, {$unset: {"is_private": ""}})
db.views.updateMany({}, {$unset: {"is_private": ""}})
db.viewtabs.updateMany({}, {$unset: {"is_private": ""}})
db.widgets.updateMany({}, {$unset: {"is_private": ""}})
db.widget_filters.updateMany({}, {$rename: {"is_user_preference": "is_private"}})

db.permission.deleteMany({_id: {$in: ["api_private_view_groups", "models_privateView"]}});
db.role.updateMany(
    {},
    {
        $unset: {
            "permissions.api_private_view_groups": "",
            "permissions.models_privateView": ""
        }
    },
);
