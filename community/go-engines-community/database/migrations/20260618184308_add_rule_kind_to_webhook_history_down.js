db.webhook_history.dropIndex("rule_1_rule_kind_1");

db.webhook_history.updateMany(
    {
        rule_kind: 0,
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
                "rule_kind"
            ]
        },
    ],
);

db.webhook_history.updateMany(
    {
        rule_kind: 1,
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
                "rule_kind"
            ]
        },
    ],
);
