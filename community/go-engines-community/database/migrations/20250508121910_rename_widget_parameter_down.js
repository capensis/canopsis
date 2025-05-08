db.widgets.updateMany(
    {
        type: "AlarmsList",
        "parameters.clearFilterEnabled": {$ne: null}
    },
    [
        {
            $set: {
                "parameters.clearFilterDisabled": {
                    $cond: {
                        if: {$eq: ["$parameters.clearFilterEnabled", true]},
                        then: false,
                        else: true,
                    }
                }
            }
        },
        {
            $unset: ["parameters.clearFilterEnabled"]
        }
    ]
);
