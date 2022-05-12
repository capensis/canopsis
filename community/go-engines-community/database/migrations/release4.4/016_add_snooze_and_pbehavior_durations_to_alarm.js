var alarmCollections = ["periodical_alarm", "resolved_alarms"];
for (var i in alarmCollections) {
    var collection = alarmCollections[i];
    db.getCollection(collection).find({ "v.steps._t": { $in: ["pbhenter", "snooze"] } }).forEach(function (doc) {
        var snoozeDuration = 0;
        var pbhDuration = 0;
        var pbhEnterTs = 0;
        var now = Math.ceil((new Date()).getTime() / 1000);
        doc.v.steps.forEach(function (step) {
            switch (step._t) {
                case "snooze":
                    if (step.val < now) {
                        snoozeDuration += step.val - step.t;
                    }
                    break;
                case "pbhenter":
                    if (step.pbehavior_canonical_type !== "active") {
                        pbhEnterTs = step.t;
                    }
                    break;
                case "pbhleave":
                    if (pbhEnterTs > 0) {
                        pbhDuration += step.t - pbhEnterTs;
                        pbhEnterTs = 0;
                    }
                    break;
            }
        });

        var set = {};
        if (snoozeDuration > 0) {
            set["v.snooze_duration"] = snoozeDuration;
        }
        if (pbhDuration > 0) {
            set["v.pbh_inactive_duration"] = pbhDuration;
        }
        if (Object.keys(set).length > 0) {
            db.getCollection(collection).updateOne({ _id: doc._id }, { $set: set });
        }
    });
}