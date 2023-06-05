db.periodical_alarm.find({"v.pbehavior_info.reason_name": {$exists: true}}).forEach(function (doc) {
    db.periodical_alarm.updateOne({_id: doc._id}, {
        $set: {
            "v.pbehavior_info.reason": doc.v.pbehavior_info.reason_name,
        },
        $unset: {
            "v.pbehavior_info.reason_name": ""
        }
    });
});

db.resolved_alarms.find({"v.pbehavior_info.reason_name": {$exists: true}}).forEach(function (doc) {
    db.resolved_alarms.updateOne({_id: doc._id}, {
        $set: {
            "v.pbehavior_info.reason": doc.v.pbehavior_info.reason_name,
        },
        $unset: {
            "v.pbehavior_info.reason_name": ""
        }
    });
});

db.default_entities.find({"pbehavior_info.reason_name": {$exists: true}}).forEach(function (doc) {
    db.default_entities.updateOne({_id: doc._id}, {
        $set: {
            "pbehavior_info.reason": doc.pbehavior_info.reason_name,
        },
        $unset: {
            "pbehavior_info.reason_name": ""
        }
    });
});
