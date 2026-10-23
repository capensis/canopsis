db.color_theme.updateMany(
    {"colors.icons": {$exists: false}},
    {$set: {"colors.icons": {
        ack: "#9C27B0",
        canceled: "#607D8B",
        children: "#6D4C41",
        parent: "#6D4C41",
        comment: "#BA68C8",
        snooze: "#E91E63",
        tickets: "#2196F3",
        cancel_delay: "#808080"
    }}}
);
