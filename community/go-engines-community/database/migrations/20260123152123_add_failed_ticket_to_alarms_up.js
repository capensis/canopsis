if (!db.permission.findOne({_id: "listalarm_removeAssociatedTicket"})) {
    db.permission.insertOne({
        _id: "listalarm_removeAssociatedTicket",
        name: "listalarm_removeAssociatedTicket",
        description: "Remove associated ticket",
        groups: ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        api_permissions: {api_alarm_update: 1}
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.listalarm_removeAssociatedTicket": 1
        }
    });
    db.role_template.updateOne({name: "Pilotes"}, {
        $set: {
            "permissions.listalarm_removeAssociatedTicket": 1
        }
    });
}

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
