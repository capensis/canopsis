resolveDuplicatesAndCreateIndex("action_scenario");
resolveDuplicatesAndCreateIndex("link_rule");
resolveDuplicatesAndCreateIndex("pbehavior_reason");
resolveDuplicatesAndCreateIndex("pbehavior_type");
resolveDuplicatesAndCreateIndex("pbehavior_exception");
db.entity_category.dropIndex("name_1");
resolveDuplicatesAndCreateIndex("entity_category");
resolveDuplicatesAndCreateIndex("view_playlist");
resolveDuplicatesAndCreateIndex("idle_rule");
resolveDuplicatesAndCreateIndex("resolve_rule");
resolveDuplicatesAndCreateIndex("flapping_rule");
resolveDuplicatesAndCreateIndex("declare_ticket_rule");
resolveDuplicatesAndCreateIndex("job");
resolveDuplicatesAndCreateIndex("job_config");
resolveDuplicatesAndCreateIndex("kpi_filter");

function resolveDuplicatesAndCreateIndex(collectionName) {
    var collection = db.getCollection(collectionName);

    try {
        collection.createIndex({name: 1}, {name: "name_1", unique: true});
    } catch (e) {
        let index = 0;

        collection.aggregate([
            {$group: {_id: "$name", docs: {$push: {_id: "$_id", name: "$name"}}}},
            {$match: {$expr: {$gt: [{$size: "$docs"}, 1]}}},
            {$unwind: "$docs"},
            {$replaceRoot: {newRoot: "$docs"}},
        ]).forEach(function(doc) {
            collection.updateOne(
                {_id: doc._id},
                {$set: {name: doc.name + "_" + (index + 1)}}
            );

            index++;
        });

        collection.createIndex({name: 1}, {name: "name_1", unique: true});
    }
}
