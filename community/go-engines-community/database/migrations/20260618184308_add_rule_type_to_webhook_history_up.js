db.webhook_history.updateMany(
    {
        scenario: {$ne: null},
        status: {$in: [0, 1]}
    },
    [
        {
            $set: {
                rule: "$scenario",
                rule_type: 0,
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
                rule_type: 1
            }
        },
        {
            $unset: [
                "declare_ticket_rule"
            ]
        },
    ],
);

db.webhook_history.createIndex({rule: 1, rule_type: 1}, {name: "rule_1_rule_type_1"});
