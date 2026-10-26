db.eventfilter.updateMany(
    {name: {$exists: false}},
    [
        {
            $set: {
                name: {
                    $cond: [
                        {$eq: ["$description", ""]},
                        "$_id",
                        "$description"
                    ]
                }
            }
        }
    ]
);

resolveDuplicatesAndCreateIndex("eventfilter", {name: "name_1", unique: true});
resolveDuplicatesAndCreateIndex("dynamic_infos", {name: "name_1", unique: true});
resolveDuplicatesAndCreateIndex("meta_alarm_rules", {
    name: "name_1",
    unique: true,
    partialFilterExpression: {soft_deleted: null},
});

function resolveDuplicatesAndCreateIndex(collectionName, indexOptions) {
    var collection = db.getCollection(collectionName);
    var indexExists = false;

    collection.getIndexes().forEach(function (index) {
        if (index.name === indexOptions.name) {
            indexExists = true;
        }
    });

    if (indexExists) {
        return;
    }

    var duplicateMatch = {name: {$exists: true, $ne: null}};
    if (indexOptions.partialFilterExpression) {
        for (var key in indexOptions.partialFilterExpression) {
            duplicateMatch[key] = indexOptions.partialFilterExpression[key];
        }
    }

    try {
        collection.createIndex({name: 1}, indexOptions);
    } catch (e) {
        collection.aggregate([
            {$match: duplicateMatch},
            {$group: {_id: "$name", docs: {$push: {_id: "$_id", name: "$name"}}}},
            {$match: {$expr: {$gt: [{$size: "$docs"}, 1]}}},
        ]).forEach(function (group) {
            for (var i = 1; i < group.docs.length; i++) {
                collection.updateOne(
                    {_id: group.docs[i]._id},
                    {$set: {name: group.docs[i].name + "-" + i}}
                );
            }
        });

        collection.createIndex({name: 1}, indexOptions);
    }
}
