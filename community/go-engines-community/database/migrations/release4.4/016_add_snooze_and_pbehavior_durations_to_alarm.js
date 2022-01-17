db.periodical_alarm.find({"v.steps._t": {$in: ["pbhenter", "snooze"]}}).forEach(function (doc) {
    var snoozeDuration = 0;
    var pbhDuration = 0;
    var pbhEnterTs = 0;
    doc.v.steps.forEach(function (step) {
        switch (step._t) {
            case "snooze":
                snoozeDuration += step.val - step.t;
                break;
            case "pbhenter":
                if (step.cannonical_type !== "active") {
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
        db.periodical_alarm.updateOne({_id: doc._id}, {$set: set});
    }
});
