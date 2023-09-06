db.color_theme.dropIndex("name_1");

db.permission.deleteMany({_id: {$in: ["api_color_theme", "models_color_theme"]}});
db.role.updateMany(
    {},
    {
        $unset: {
            "permissions.api_color_theme": "",
            "permissions.models_color_theme": ""
        }
    },
);
