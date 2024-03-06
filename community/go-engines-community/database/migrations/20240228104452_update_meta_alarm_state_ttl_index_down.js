db.meta_alarm_states.dropIndex("meta_alarm_state_expired_at_1");
db.meta_alarm_states.createIndex({"created_at": 1}, {name: "meta_alarm_state_created_at_1", expireAfterSeconds: 86400});
db.meta_alarm_states.updateMany({expired_at: {$ne: null}}, [
    {
        $set: {
            created_at: "$expired_at",
        }
    },
    {$unset: "expired_at"},
]);
