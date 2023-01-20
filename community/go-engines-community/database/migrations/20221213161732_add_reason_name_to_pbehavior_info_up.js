db.pbehavior_reason.find({}).forEach(function (doc) {
    db.periodical_alarm.updateMany({"v.pbehavior_info.reason": doc.name}, {
        $set: {
            "v.pbehavior_info.reason": doc._id,
            "v.pbehavior_info.reason_name": doc.name,
        }
    });

    db.resolved_alarms.updateMany({"v.pbehavior_info.reason": doc.name}, {
        $set: {
            "v.pbehavior_info.reason": doc._id,
            "v.pbehavior_info.reason_name": doc.name,
        }
    });

    db.default_entities.updateMany({"pbehavior_info.reason": doc.name}, {
        $set: {
            "pbehavior_info.reason": doc._id,
            "pbehavior_info.reason_name": doc.name,
        }
    });
});
