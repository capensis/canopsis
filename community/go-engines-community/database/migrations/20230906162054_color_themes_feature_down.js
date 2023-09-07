db.color_theme.drop();

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
