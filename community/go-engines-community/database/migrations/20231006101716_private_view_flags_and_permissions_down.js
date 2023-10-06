db.viewgroups.updateMany({}, {$unset: {"is_private": ""}})
db.views.updateMany({}, {$unset: {"is_private": ""}})
db.viewtabs.updateMany({}, {$unset: {"is_private": ""}})
db.widgets.updateMany({}, {$unset: {"is_private": ""}})
db.widget_filters.updateMany({}, {$rename: {"widget_private": "is_private"}})

db.permission.deleteOne({_id: "api_private_view_groups"});
db.role.updateMany(
    {},
    {
        $unset: {
            "permissions.api_private_view_groups": ""
        }
    },
);
