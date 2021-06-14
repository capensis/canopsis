(function () {
    db.default_entities.find({created: null, enable_history: {$nin: [null, []]}}).forEach(function (doc) {
        db.default_entities.updateOne({_id: doc._id}, {
            $set: {
                created: doc.enable_history[0],
            },
        });
    });

    var cursor = db.default_entities.aggregate([
        {$match: {created: null}},
        {$lookup: {
           from: "periodical_alarm",
           localField: "_id",
           foreignField: "d",
           as: "alarm"
        }},
        {$unwind: "$alarm"},
        {$sort: {"alarm.t": 1}},
        {$group: {
            _id: "$_id",
            created: {$first: "$alarm.t"}
        }},
    ]);
    while (cursor.hasNext()) {
        var doc = cursor.next();
        db.default_entities.updateOne({_id: doc._id}, {
            $set: {
                created: doc.created,
            },
        });
    }
})();
