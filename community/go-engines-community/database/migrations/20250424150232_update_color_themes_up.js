db.color_theme.updateOne({_id: "canopsis", "colors.table.shift_row": null}, {
    $set: {
        "colors": {
            "main": {
                "primary": "#2DAC61",
                "secondary": "#2B3E4F",
                "accent": "#678CCA",
                "background": "#FFFFFF",
                "active_color": "#212121",
                "error": "#D22E2E",
                "error_background": "#FFCABE",
                "warning": "#FB8C00",
                "warning_background": "#FFD27A",
                "success": "#22874C",
                "success_background": "#90DFB1",
                "info": "#0A6EBD",
                "info_background": "#A6D9FF"
            },
            "state": {
                "ok": "#22874C",
                "minor": "#FFF176",
                "major": "#FFA800",
                "critical": "#C62828"
            },
            "table": {
                "background": "#FFFFFF",
                "row_color": "#FFFFFF",
                "shift_row": false,
                "shift_row_color": "#F5F5F5",
                "hover_row": true,
                "hover_row_color": "#E5E5E5"
            }
        }
    }
});

db.color_theme.updateOne({_id: "canopsis_dark", "colors.table.shift_row": null}, {
    $set: {
        "colors": {
            "main": {
                "primary": "#2DAC61",
                "secondary": "#2B3E4F",
                "accent": "#678CCA",
                "background": "#424242",
                "active_color": "#FFFFFF",
                "error": "#FB9E9E",
                "error_background": "#995656",
                "warning": "#FFC658",
                "warning_background": "#826742",
                "success": "#90DFB1",
                "success_background": "#4B7D60",
                "info": "#A6D9FF",
                "info_background": "#216B9C"
            },
            "state": {
                "ok": "#4B7D60",
                "minor": "#FFF176",
                "major": "#FFA800",
                "critical": "#995656"
            },
            "table": {
                "background": "#424242",
                "row_color": "#424242",
                "shift_row": false,
                "shift_row_color": "#515151",
                "hover_row": true,
                "hover_row_color": "#616161"
            }
        }
    }
});

db.color_theme.updateOne({_id: "color_blind", "colors.table.shift_row": null}, {
    $set: {
        "colors": {
            "main": {
                "primary": "#2196F3",
                "secondary": "#2B3E4F",
                "accent": "#678CCA",
                "background": "#FFFFFF",
                "active_color": "#212121",
                "error": "#D22E2E",
                "error_background": "#FFCABE",
                "warning": "#FB8C00",
                "warning_background": "#FFDD9A",
                "success": "#368472",
                "success_background": "#A0EDDB",
                "info": "#0A6EBD",
                "info_background": "#C2D0FF"
            },
            "state": {
                "ok": "#22874C",
                "minor": "#FFF176",
                "major": "#FFA800",
                "critical": "#C62828"
            },
            "table": {
                "background": "#FFFFFF",
                "row_color": "#FFFFFF",
                "shift_row": false,
                "shift_row_color": "#F5F5F5",
                "hover_row": true,
                "hover_row_color": "#E5E5E5"
            },
        }
    }
});

db.color_theme.updateOne({_id: "color_blind_dark", "colors.table.shift_row": null}, {
    $set: {
        "colors": {
            "main": {
                "primary": "#2196F3",
                "secondary": "#2B3E4F",
                "accent": "#678CCA",
                "background": "#424242",
                "active_color": "#FFFFFF",
                "error": "#FB9E9E",
                "error_background": "#995656",
                "warning": "#FFC658",
                "warning_background": "#826742",
                "success": "#A0EDDB",
                "success_background": "#368484",
                "info": "#A6D9FF",
                "info_background": "#446AAF"
            },
            "state": {
                "ok": "#368484",
                "minor": "#FFF176",
                "major": "#FFA800",
                "critical": "#995656"
            },
            "table": {
                "background": "#424242",
                "row_color": "#424242",
                "shift_row": false,
                "shift_row_color": "#515151",
                "hover_row": true,
                "hover_row_color": "#616161"
            }
        }
    }
})

db.color_theme.updateMany(
    {
        "colors.table.shift_row": null,
        "colors.table.shift_row_color": {$nin: [null, ""]}
    },
    {
        $set: {
            "colors.table.shift_row": true
        }
    }
);

db.color_theme.updateMany(
    {
        "colors.table.shift_row": null,
        "colors.table.shift_row_color": {$in: [null, ""]}
    },
    {
        $set: {
            "colors.table.shift_row": false
        }
    }
);

db.color_theme.updateMany(
    {
        "colors.table.hover_row": null,
        "colors.table.hover_row_color": {$nin: [null, ""]}
    },
    {
        $set: {
            "colors.table.hover_row": true
        }
    }
);

db.color_theme.updateMany(
    {
        "colors.table.hover_row": null,
        "colors.table.hover_row_color": {$in: [null, ""]}
    },
    {
        $set: {
            "colors.table.hover_row": false
        }
    }
);
