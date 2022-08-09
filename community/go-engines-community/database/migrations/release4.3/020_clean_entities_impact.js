(function () {
    var offset = 0;
    var limit = 100;

    while (true) {
        var cursor = db.default_entities.aggregate([
            {$match: {type: "service"}},
            {$sort: {_id: 1}},
            {$skip: offset},
            {$limit: limit},
            {
                $project: {
                    depends: 1,
                }
            },
            {
                $lookup: {
                    from: "default_entities",
                    localField: "_id",
                    foreignField: "impact",
                    as: "depends_by_impact"
                }
            },
            {
                $addFields: {
                    depends_by_impact: {
                        $map: {
                            input: "$depends_by_impact",
                            in: "$$this._id"
                        }
                    }
                }
            },
            {
                $project: {
                    diff: {$setDifference: ["$depends_by_impact", "$depends"]}
                }
            },
            {$match: {$expr: {$gt: [{$size: "$diff"}, 0]}}}
        ]);
        if (!cursor.hasNext()) {
            return;
        }

        while (cursor.hasNext()) {
            var doc = cursor.next();

            if (doc.diff && doc.diff.length > 0) {
                db.default_entities.updateMany({_id: {$in: doc.diff}}, {
                    $pull: {
                        "impact": doc._id,
                    },
                });
            }
        }

        offset += limit;
    }
})();
