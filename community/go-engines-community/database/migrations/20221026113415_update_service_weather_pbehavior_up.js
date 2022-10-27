db.pbehavior.aggregate([
    {
        $match: {
            name: {$regex: "^downtime-.*"},
            tstop: {$in: [null, 2147483647]},
        }
    },
    {
        $lookup: {
            from: "pbehavior_type",
            localField: "type_",
            foreignField: "_id",
            as: "type"
        }
    },
    {$unwind: "$type"},
    {
        $match: {
            "type.type": "pause",
        }
    },
]).forEach(function (doc) {
    var entity = "";
    if (doc.entity_pattern && doc.entity_pattern.length === 1) {
        if (doc.entity_pattern[0].length === 1) {
            if (doc.entity_pattern[0][0].field === "_id" && doc.entity_pattern[0][0].cond.type === "eq") {
                entity = doc.entity_pattern[0][0].cond.value;
            }
        }
    }

    if (entity !== "") {
        db.pbehavior.updateOne({_id: doc._id}, {
            $set: {
                origin: "ServiceWeather",
                entity: entity
            }
        });
    }
});
