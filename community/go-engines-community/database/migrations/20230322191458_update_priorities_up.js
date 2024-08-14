var collectionNames = [
    "flapping_rule",
    "idle_rule",
    "pbehavior_type",
    "resolve_rule",
    "action_scenario",
    "instruction",
];

// previous version of index has type "unique"
db.pbehavior_type.dropIndex("priority_1");

if (db.eventfilter.findOne({priority: 0}) && db.eventfilter.findOne({priority: {$gt: 0}})) {
    collectionNames.push("eventfilter");
}

for (var collectionName of collectionNames) {
    var collection = db.getCollection(collectionName);
    collection.find({priority: 0}).forEach(function (doc) {
        var priority = 1;
        collection.updateOne({_id: doc._id}, {
            $set: {priority: priority},
        });

        var cursor = collection.find({
            _id: {$ne: doc._id},
            priority: {$gte: priority},
        }, {}, {sort: {priority: 1}});
        var writeModels = [];
        var seq = priority;
        while (cursor.hasNext()) {
            var nextDoc = cursor.next();

            if (seq == priority && nextDoc.priority != priority) {
                return;
            }
            if (nextDoc.priority != seq) {
                break;
            }

            seq++;
            writeModels.push({
                updateOne: {
                    filter: {_id: nextDoc._id},
                    update: {$set: {priority: seq}}
                }
            });
        }

        if (writeModels.length > 0) {
            collection.bulkWrite(writeModels);
        }
    });
}

db.pbehavior_type.createIndex({priority: 1}, {name: "priority_1"});
