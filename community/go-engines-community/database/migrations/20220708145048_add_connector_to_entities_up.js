db.default_entities.find({type: "connector"}).forEach(function (doc) {
    db.default_entities.updateMany(
        {
            type: "component",
            impact: doc._id,
        },
        {
            $set: {
                connector: doc._id,
            }
        }
    );
    db.default_entities.updateMany(
        {
            type: "resource",
            depends: doc._id,
        },
        {
            $set: {
                connector: doc._id,
            }
        }
    );
});
