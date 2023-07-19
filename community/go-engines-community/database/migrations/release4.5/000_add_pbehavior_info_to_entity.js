db.periodical_alarm.aggregate([
    {
        $match: {
            "v.pbehavior_info.id": {$nin: ["", null]},
            "v.pbehavior_info.timestamp": null,
        }
    },
    {
        $lookup: {
            from: "pbehavior",
            localField: "v.pbehavior_info.id",
            foreignField: "_id",
            as: "pbehavior"
        }
    },
    {
        $unwind: {path: "$pbehavior", preserveNullAndEmptyArrays: true}
    }
]).forEach(function (doc) {
    var timestamp = -1;
    for (var i = doc.v.steps.length - 1; i >= 0; i--) {
        if (doc.v.steps[i]._t === "pbhenter") {
            timestamp = doc.v.steps[i].t;
            break;
        }
    }

    if (timestamp <= 0) {
        if (doc.pbehavior && doc.pbehavior.created > doc.t) {
            timestamp = doc.pbehavior.created;
        } else {
            timestamp = doc.t;
        }
    }

    db.periodical_alarm.updateOne({_id: doc._id}, {
        $set: {
            "v.pbehavior_info.timestamp": timestamp,
        }
    });

    var pbehaviorInfo = doc.v.pbehavior_info;
    pbehaviorInfo.timestamp = timestamp;
    db.default_entities.updateOne({_id: doc.d}, {
        $set: {
            "pbehavior_info": pbehaviorInfo,
        }
    });
});
