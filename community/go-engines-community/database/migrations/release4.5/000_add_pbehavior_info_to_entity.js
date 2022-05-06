db.periodical_alarm.find({"v.pbehavior_info.id": {$nin: ["", null]}}).forEach(function (doc) {
    var timestamp;
    for (var i = doc.v.steps.length-1; i >= 0; i--) {
        if (doc.v.steps[i]._t === "pbhenter") {
            timestamp = doc.v.steps[i].t;
            break;
        }
    }

    if (timestamp) {
        db.periodical_alarm.updateOne({_id: doc._id}, {$set: {"v.pbehavior_info.timestamp": timestamp}});

        var pbehaviorInfo = doc.v.pbehavior_info;
        pbehaviorInfo.timestamp = timestamp;
        db.default_entities.updateOne({_id: doc.d}, {$set: {"pbehavior_info": pbehaviorInfo}});
    }
});
