(function () {
    var cursor = db.default_entities.aggregate([
        {$match: {type: "service"}},
        {$lookup: {
                from: "periodical_alarm",
                localField: "_id",
                foreignField: "d",
                as: "alarm"
            }},
        {$unwind: "$alarm"},
        {$match: {"alarm.v.resolved": null}},
    ]);
    while (cursor.hasNext()) {
        var doc = cursor.next();
        db.periodical_alarm.updateOne({_id: doc.alarm._id}, {
            $set: {
                "v.connector": "service",
                "v.connector_name": "service",
            },
        });
    }
})();
