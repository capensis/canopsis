(function () {
    var pipeline = [
        {
            $lookup: {
                from: 'instruction_rating',
                localField: '_id',
                foreignField: 'instruction',
                as: 'rating'
            }
        },
        {
            $unwind: {
                path: "$rating",
                preserveNullAndEmptyArrays: true
            }
        },
        {
            $group: {
                _id: "$_id",
                rating: {
                    "$avg": "$rating.rating"
                }
            }
        },
        {
            $addFields: {
                "rating": {
                    "$cond": {
                        "if": "$rating",
                        "then": "$rating",
                        "else": 0,
                    }
                },
            }
        }
    ];
    var cursor = db.instruction.aggregate(pipeline);

    while (cursor.hasNext()) {
        var doc = cursor.next();
        db.instruction.updateOne({_id: doc._id}, {"$set": {rating: doc.rating}});
    }

    var date = Math.floor(Date.now() / 1000)
    db.instruction_rating.updateMany({}, {"$set": {created: date}});
})();