db.webhook_history.updateMany(
    {
        scenario: {$ne: null},
        status: {$in: [0, 1]}
    },
    [
        {
            $set: {
                rule: "$scenario",
                rule_kind: 0,
            }
        },
        {
            $unset: [
                "scenario"
            ]
        },
    ],
);

db.webhook_history.updateMany(
    {
        declare_ticket_rule: {$ne: null},
        status: {$in: [0, 1]}
    },
    [
        {
            $set: {
                rule: "$declare_ticket_rule",
                rule_kind: 1
            }
        },
        {
            $unset: [
                "declare_ticket_rule"
            ]
        },
    ],
);

