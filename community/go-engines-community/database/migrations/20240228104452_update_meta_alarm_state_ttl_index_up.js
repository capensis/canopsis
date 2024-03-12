db.meta_alarm_states.dropIndex("meta_alarm_state_created_at_1")
db.meta_alarm_states.createIndex({"expired_at": 1}, {name: "meta_alarm_state_expired_at_1", expireAfterSeconds: 86400})
db.meta_alarm_states.updateMany({expired_at: null}, [
    {
        $set: {
            expired_at: {
                $cond: {
                    if: "$created_at",
                    then: "$created_at",
                    else: {
                        $dateAdd:
                            {
                                startDate: new Date(),
                                unit: "month",
                                amount: 1
                            }
                    },
                }
            },
        }
    },
    {$unset: "created_at"},
]);
