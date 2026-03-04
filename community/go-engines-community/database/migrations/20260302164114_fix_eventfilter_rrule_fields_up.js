db.eventfilter.updateMany(
    {
        start: {$ne: null},
        stop: {$ne: null}
    },
    [
        {
            $set: {
                resolved_start: "$start",
                resolved_stop: "$stop",
            }
        },
        {
            $unset: [
                "next_resolved_start",
                "next_resolved_stop",
            ]
        }
    ]
);
