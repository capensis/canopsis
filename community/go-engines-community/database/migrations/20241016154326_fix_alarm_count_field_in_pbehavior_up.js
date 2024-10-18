db.pbehavior.updateMany({}, {$set: {alarm_count: 0}});
db.periodical_alarm.aggregate([
    {
        $match: {
            "v.resolved": null,
            "v.pbehavior_info.id": {$nin: [null, ""]},
        }
    },
    {
        $group: {
            _id: "$v.pbehavior_info.id",
            alarm_count: {$sum: 1},
        }
    }
]).forEach(function (doc) {
    db.pbehavior.updateOne({_id: doc._id}, {$set: {alarm_count: doc.alarm_count}});
});
