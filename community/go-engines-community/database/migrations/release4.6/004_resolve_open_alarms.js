(function () {
    db.periodical_alarm.aggregate([
        {$match: {"v.resolved": null}},
        {$sort: {t: -1}},
        {
            $group: {
                _id: "$d",
                count: {$sum: 1},
                ids: {$push: "$_id"}
            }
        },
        {$match: {count: {$gt: 1}}},
    ], {allowDiskUse: true}).forEach(function (doc) {
        var now = Math.ceil((new Date()).getTime() / 1000);
        var ids = doc.ids;
        ids.shift();
        db.periodical_alarm.updateMany({_id: {$in: ids}}, {
            $set: {
                "v.resolved": now,
            }
        })
    });
})();
