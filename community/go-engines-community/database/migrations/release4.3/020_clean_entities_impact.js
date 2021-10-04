(function () {
    var cursor = db.default_entities.aggregate([
        {$match: {type: "service"}},
        {$lookup: {
                from: "default_entities",
                localField: "_id",
                foreignField: "impact",
                as: "depends_by_impact"
            }},
        {$addFields: {
                depends_by_impact: {$map: {
                        input: "$depends_by_impact",
                        in: "$$this._id"
                    }}
            }},
        {$project: {
                diff: { $setDifference: [ "$depends_by_impact", "$depends" ] }
            }},
        {$match: {$expr: {$gt: [{$size: "$diff"}, 0]}}}
    ]);
    while (cursor.hasNext()) {
        var doc = cursor.next();
        db.default_entities.updateMany({_id: {$in: doc.diff}}, {
            $pull: {
                "impact": doc._id,
            },
        });
    }
})();
