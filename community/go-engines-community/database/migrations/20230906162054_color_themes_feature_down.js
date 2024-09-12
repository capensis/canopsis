db.color_theme.drop();

db.permission.deleteMany({_id: {$in: ["api_color_theme", "models_profile_color_theme"]}});
db.role.updateMany(
    {},
    {
        $unset: {
            "permissions.api_color_theme": "",
            "permissions.models_profile_color_theme": ""
        }
    },
);

db.user.updateMany({"ui_theme": "canopsis_dark"}, {$set:{"ui_theme": "canopsisDark"}})
db.user.updateMany({"ui_theme": "color_blind"}, {$set:{"ui_theme": "colorBlind"}})
db.user.updateMany({"ui_theme": "color_blind_dark"}, {$set:{"ui_theme": "colorBlindDark"}})
