db.periodical_alarm.updateMany(
    {
        "v.last_comment": {$ne: null},
    },
    [
        {
            $set: {
                "v.comments": {
                    $filter: {
                        input: "$v.steps",
                        as: "step",
                        cond: {$eq: ["$$step._t", "comment"]}
                    }
                }
            }
        }
    ],
)

db.resolved_alarms.updateMany(
    {
        "v.last_comment": {$ne: null},
    },
    [
        {
            $set: {
                "v.comments": {
                    $filter: {
                        input: "$v.steps",
                        as: "step",
                        cond: {$eq: ["$$step._t", "comment"]}
                    }
                }
            }
        }
    ],
)

