db.color_theme.updateMany(
    {
        _id: {$in: ["canopsis", "color_blind"]},
        "colors.widget": null
    },
    {
        $set: {
            "colors.widget": {
                title_background: "#FFFFFF",
                content_background: "#FFFFFF",
                border: false,
                border_color: "#C7C7C7",
                border_radius: 0
            }
        }
    }
);

db.color_theme.updateMany(
    {
        _id: {$in: ["canopsis_dark", "color_blind_dark"]},
        "colors.widget": null
    },
    {
        $set: {
            "colors.widget": {
                title_background: "#424242",
                content_background: "#424242",
                border: false,
                border_color: "#7D7D7D",
                border_radius: 0
            }
        }
    }
);
