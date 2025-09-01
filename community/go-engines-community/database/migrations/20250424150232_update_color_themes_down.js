db.color_theme.updateOne({_id: "canopsis", "colors.table.shift_row": {$ne: null}}, {
    $set: {
        "colors": {
            "main": {
                "primary": "#2fab63",
                "secondary": "#2b3e4f",
                "accent": "#82b1ff",
                "background": "#ffffff",
                "active_color": "#000",
                "error": "#ff5252",
                "warning": "#fb8c00",
                "success": "#4caf50",
                "info": "#2196f3"
            },
            "state": {
                "ok": "#00a65a",
                "minor": "#fcdc00",
                "major": "#ff9900",
                "critical": "#f56954"
            },
            "table": {
                "background": "#fff",
                "row_color": "#fff",
                "hover_row_color": "#eee"
            }
        }
    }
});

db.color_theme.updateOne({_id: "canopsis_dark", "colors.table.shift_row": {$ne: null}}, {
    $set: {
        "colors": {
            "main": {
                "primary": "#2fab63",
                "secondary": "#2b3e4f",
                "accent": "#82b1ff",
                "background": "#424242",
                "active_color": "#fff",
                "error": "#ff8b8b",
                "warning": "#fb8c00",
                "success": "#4caf50",
                "info": "#2196f3"
            },
            "state": {
                "ok": "#00a65a",
                "minor": "#fcdc00",
                "major": "#ff9900",
                "critical": "#f56954"
            },
            "table": {
                "background": "#424242",
                "row_color": "#424242",
                "hover_row_color": "#616161"
            }
        }
    }
});

db.color_theme.updateOne({_id: "color_blind", "colors.table.shift_row": {$ne: null}}, {
    $set: {
        "colors": {
            "main": {
                "primary": "#2196f3",
                "secondary": "#2b3e4f",
                "accent": "#82b1ff",
                "background": "#ffffff",
                "active_color": "#000",
                "error": "#ff5252",
                "warning": "#fb8c00",
                "success": "#4caf50",
                "info": "#2196f3"
            },
            "state": {
                "ok": "#00a65a",
                "minor": "#fcdc00",
                "major": "#ff9900",
                "critical": "#f56954"
            },
            "table": {
                "background": "#fff",
                "row_color": "#fff",
                "hover_row_color": "#eee"
            }
        }
    }
});

db.color_theme.updateOne({_id: "color_blind_dark", "colors.table.shift_row": {$ne: null}}, {
    $set: {
        "colors": {
            "main": {
                "primary": "#2196f3",
                "secondary": "#2b3e4f",
                "accent": "#82b1ff",
                "background": "#424242",
                "active_color": "#fff",
                "error": "#ff8b8b",
                "warning": "#fb8c00",
                "success": "#4caf50",
                "info": "#2196f3"
            },
            "state": {
                "ok": "#00a65a",
                "minor": "#fcdc00",
                "major": "#ff9900",
                "critical": "#f56954"
            },
            "table": {
                "background": "#424242",
                "row_color": "#424242",
                "hover_row_color": "#616161"
            }
        }
    }
});

db.color_theme.updateMany(
    {
        "colors.table.shift_row": true,
    },
    {
        $unset: {
            "colors.table.shift_row": ""
        }
    }
);

db.color_theme.updateMany(
    {
        "colors.table.shift_row": false,
    },
    {
        $unset: {
            "colors.table.shift_row": "",
            "colors.table.shift_row_color": ""
        }
    }
);

db.color_theme.updateMany(
    {
        "colors.table.hover_row": true
    },
    {
        $unset: {
            "colors.table.hover_row": ""
        }
    }
);

db.color_theme.updateMany(
    {
        "colors.table.hover_row": false
    },
    {
        $unset: {
            "colors.table.hover_row": "",
            "colors.table.hover_row_color": ""
        }
    }
);
