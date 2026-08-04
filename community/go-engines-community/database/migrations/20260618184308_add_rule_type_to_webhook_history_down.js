db.webhook_history.dropIndex("rule_1_rule_type_1");

db.webhook_history.updateMany(
    {
        rule_type: 0,
        status: {$in: [0, 1]}
    },
    [
        {
            $set: {
                scenario: "$rule",
            }
        },
        {
            $unset: [
                "rule",
                "rule_type"
            ]
        },
    ],
);

db.webhook_history.updateMany(
    {
        rule_type: 1,
        status: {$in: [0, 1]}
    },
    [
        {
            $set: {
                declare_ticket_rule: "$rule",
            }
        },
        {
            $unset: [
                "rule",
                "rule_type"
            ]
        },
    ],
);
