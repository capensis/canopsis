db.widgets.updateMany(
    {
        $and: [
            {"parameters.fastPbehaviorReason": {$nin: [null, ""]}},
            {"parameters.fastPbehaviorType": {$nin: [null, ""]}},
            {"parameters.fastPbehaviorNamePrefix": {$nin: [null, ""]}},
        ],
    },
    [
        {
            $set: {
                "parameters.fast_pbehaviors": [
                    {
                        name_prefix: "$parameters.fastPbehaviorNamePrefix",
                        reason: "$parameters.fastPbehaviorReason",
                        type: "$parameters.fastPbehaviorType",
                    },
                ],
            },
        },
        {
            $unset: [
                "parameters.fastPbehaviorReason",
                "parameters.fastPbehaviorType",
                "parameters.fastPbehaviorNamePrefix",
            ],
        },
    ]
)
