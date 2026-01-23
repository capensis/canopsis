db.periodical_alarm.updateMany(
    {
        "v.tickets._t": "declareticketfail",
    },
    [
        {
            $set: {
                "v.failed_ticket": {
                    $arrayElemAt: ["$v.tickets", -1]
                },
                "v.tickets": {
                    $filter: {
                        input: "$v.tickets",
                        cond: {$ne: ["$$this._t", "declareticketfail"]}
                    }
                }
            }
        },
        {
            $set: {
                "v.failed_ticket": {
                    $cond: {
                        if: {$eq: ["$v.failed_ticket._t", "declareticketfail"]},
                        then: "$v.failed_ticket",
                        else: null
                    }
                },
            }
        }
    ],
);
