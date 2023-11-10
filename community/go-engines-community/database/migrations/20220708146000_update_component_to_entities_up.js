db.default_entities.find({type: "component"}).forEach(function (doc) {
    if (!doc.component) {
        db.default_entities.updateOne({_id: doc._id}, {$set: {component: doc._id}});
    }

    db.default_entities.updateMany(
        {
            type: "resource",
            component: null,
            impact: doc._id,
        },
        {
            $set: {
                component: doc._id,
            }
        }
    );
});
