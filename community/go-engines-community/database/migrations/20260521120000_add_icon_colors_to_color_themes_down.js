db.color_theme.updateMany(
    {"colors.icons": {$exists: true}},
    {$unset: {"colors.icons": ""}}
);
