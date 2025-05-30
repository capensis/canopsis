db.widgets.updateMany(
    {
        type: "AlarmsList",
        "parameters.clearFilterDisabled": {$ne: null}
    },
    [
        {
            $set: {
                "parameters.clearFilterEnabled": {
                    $cond: {
                        if: {$eq: ["$parameters.clearFilterDisabled", true]},
                        then: false,
                        else: true,
                    }
                }
            }
        },
        {
            $unset: ["parameters.clearFilterDisabled"]
        }
    ]
);
