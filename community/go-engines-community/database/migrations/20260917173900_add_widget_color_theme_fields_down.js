db.color_theme.updateMany(
    {},
    {
        $unset: {
            "colors.widget": ""
        }
    }
);
