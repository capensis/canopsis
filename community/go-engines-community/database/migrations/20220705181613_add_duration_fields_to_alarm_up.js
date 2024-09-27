var updateResolvedPipeline = [
    {
        $set: {
            "v.duration": {
                "$subtract": [
                    "$v.resolved",
                    "$v.creation_date"
                ],
            },
            "v.current_state_duration": {
                "$subtract": [
                    "$v.resolved",
                    "$v.state.t"
                ],
            },
        }
    },
];

function updateResolvedAlarms(collectionName) {
    return function (doc) {
        if (!doc.v || !doc.v.steps) {
            return;
        }
        var snoozeDuration = 0;
        var pbhInactiveDuration = 0;
        var inactivePeriods = [];
        var pbhIndex = null;

        for (var step of doc.v.steps) {
            var stepTs = step.t;

            if (isNaN(stepTs)) {
                continue
            }

            if (stepTs < doc.v.creation_date) {
                stepTs = doc.v.creation_date;
            }

            switch (step._t) {
                case "snooze":
                    var snoozeStartTs = stepTs;
                    var snoozeEndTs = step.val;

                    if (isNaN(snoozeEndTs)) {
                        break
                    }

                    if (snoozeEndTs > doc.v.resolved) {
                        snoozeEndTs = doc.v.resolved;
                    }

                    snoozeDuration += snoozeEndTs - snoozeStartTs;
                    inactivePeriods.push({
                        start: snoozeStartTs,
                        end: snoozeEndTs,
                    });

                    break;
                case "pbhenter":
                    if (pbhIndex !== null) {
                        break;
                    }

                    if (step.pbehavior_canonical_type !== "active") {
                        pbhIndex = inactivePeriods.length;
                        inactivePeriods.push({
                            start: stepTs,
                        });
                    }
                    break;
                case "pbhleave":
                    if (pbhIndex === null) {
                        break;
                    }

                    var pbhEnterTs = inactivePeriods[pbhIndex].start;
                    var pbhLeaveTs = stepTs;
                    pbhInactiveDuration += pbhLeaveTs - pbhEnterTs;
                    inactivePeriods[pbhIndex].end = pbhLeaveTs;

                    pbhIndex = null;
                    break;
            }
        }

        if (pbhIndex !== null) {
            var pbhEnterTs = inactivePeriods[pbhIndex].start;
            var pbhLeaveTs = doc.v.resolved;
            pbhInactiveDuration += pbhLeaveTs - pbhEnterTs;
            inactivePeriods[pbhIndex].end = pbhLeaveTs;
        }

        var inactiveDuration = 0;
        for (var i = 0; i < inactivePeriods.length;) {
            var start = inactivePeriods[i].start;
            var end = inactivePeriods[i].end;

            i++;
            var endPeriod = end;
            for (var j = i; j < inactivePeriods.length; j++) {
                if (inactivePeriods[j].start > end) {
                    break;
                }

                if (inactivePeriods[j].end > end) {
                    endPeriod = inactivePeriods[j].start;
                    i = j;
                } else {
                    i = j + 1;
                }
            }

            inactiveDuration += endPeriod - start;
        }

        db.getCollection(collectionName).updateOne({_id: doc._id}, {
            $set: {
                "v.inactive_duration": toInt(inactiveDuration),
                "v.active_duration": toInt(doc.v.duration - inactiveDuration),
                "v.snooze_duration": toInt(snoozeDuration),
                "v.pbh_inactive_duration": toInt(pbhInactiveDuration),
            }
        });
    };
}

function updateOpenedAlarms(collectionName) {
    return function (doc) {
        if (!doc.v || !doc.v.steps) {
            return;
        }
        var snoozeDuration = 0;
        var pbhInactiveDuration = 0;
        var inactivePeriods = [];
        var pbhIndex = null;
        var now = Math.ceil((new Date()).getTime() / 1000);

        for (var step of doc.v.steps) {
            var stepTs = step.t;

            if (isNaN(stepTs)) {
                continue
            }

            if (stepTs < doc.v.creation_date) {
                stepTs = doc.v.creation_date;
            }

            switch (step._t) {
                case "snooze":
                    var snoozeStartTs = stepTs;
                    var snoozeEndTs = step.val;

                    if (isNaN(snoozeEndTs)) {
                        break
                    }

                    if (snoozeEndTs < now) {
                        snoozeDuration += snoozeEndTs - snoozeStartTs;
                        inactivePeriods.push({
                            start: snoozeStartTs,
                            end: snoozeEndTs,
                        });
                    } else {
                        inactivePeriods.push({
                            start: snoozeStartTs,
                            end: null,
                        });
                    }
                    break;
                case "pbhenter":
                    if (pbhIndex !== null) {
                        break;
                    }

                    if (step.pbehavior_canonical_type !== "active") {
                        pbhIndex = inactivePeriods.length;
                        inactivePeriods.push({
                            start: stepTs,
                            end: null,
                        });
                    }
                    break;
                case "pbhleave":
                    if (pbhIndex === null) {
                        break;
                    }

                    var pbhEnterTs = inactivePeriods[pbhIndex].start;
                    var pbhLeaveTs = stepTs;
                    pbhInactiveDuration += pbhLeaveTs - pbhEnterTs;
                    inactivePeriods[pbhIndex].end = pbhLeaveTs;

                    pbhIndex = null;
                    break;
            }
        }

        var inactiveDuration = 0;
        var inactiveStart = null;
        for (var i = 0; i < inactivePeriods.length;) {
            var start = inactivePeriods[i].start;
            var end = inactivePeriods[i].end;

            if (end === null) {
                if (i === inactivePeriods.length - 1) {
                    inactiveStart = start;
                    if (i > 0 && inactivePeriods[i-1].end !== null && inactivePeriods[i-1].end > start) {
                        inactiveStart = inactivePeriods[i-1].end;
                        inactiveDuration += inactiveStart - start
                    }
                } else if (inactivePeriods[inactivePeriods.length - 1].end === null) {
                    inactiveStart = inactivePeriods[inactivePeriods.length - 1].start;
                    inactiveDuration += inactiveStart - start;
                } else {
                    inactiveStart = inactivePeriods[inactivePeriods.length - 1].end;
                    inactiveDuration += inactiveStart - start;
                }

                break;
            }

            i++;
            var endPeriod = end;
            for (var j = i; j < inactivePeriods.length; j++) {
                if (inactivePeriods[j].start > end) {
                    break;
                }

                if (inactivePeriods[j].end === null || inactivePeriods[j].end > end) {
                    endPeriod = inactivePeriods[j].start;
                    i = j;
                } else {
                    i = j + 1;
                }
            }

            inactiveDuration += endPeriod - start;
        }

        var set = {
            "v.inactive_duration": toInt(inactiveDuration),
            "v.snooze_duration": toInt(snoozeDuration),
            "v.pbh_inactive_duration": toInt(pbhInactiveDuration),
        }
        if (inactiveStart !== null) {
            set["v.inactive_start"] = toInt(inactiveStart);
        }

        db.getCollection(collectionName).updateOne({_id: doc._id}, {
            $set: set,
        });
    };
}

db.periodical_alarm.updateMany({
    "v.resolved": {$ne: null},
    "v.duration": null,
}, updateResolvedPipeline);
db.resolved_alarms.updateMany({
    "v.duration": null
}, updateResolvedPipeline);

db.periodical_alarm.find({
    "v.resolved": {$ne: null},
    "v.active_duration": null,
}).forEach(updateResolvedAlarms("periodical_alarm"));
db.resolved_alarms.find({
    "v.active_duration": null
}).forEach(updateResolvedAlarms("resolved_alarms"));

db.periodical_alarm.find({
    "v.resolved": null,
    "v.inactive_duration": null,
}).forEach(updateOpenedAlarms("periodical_alarm"));
