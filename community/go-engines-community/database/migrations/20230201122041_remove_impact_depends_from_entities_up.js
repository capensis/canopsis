db.default_entities.updateMany({services: null}, {$set: {services: []}});

db.default_entities.aggregate([
    {
        $graphLookup: {
            from: "default_entities",
            startWith: "$impact",
            connectFromField: "impact",
            connectToField: "_id",
            restrictSearchWithMatch: {type: "service"},
            as: "services",
            maxDepth: 0,
        }
    },
    {
        $match: {services: {$ne: []}}
    },
    {
        $project: {
            services: {$map: {input: "$services", in: "$$this._id"}}
        }
    }
]).forEach(function (doc) {
    db.default_entities.updateOne({_id: doc._id}, {$set: {services: doc.services}});
});

db.default_entities.updateMany({}, {$unset: {impact: "", depends: ""}});

db.default_entities.getIndexes().forEach(function (index) {
    Object.keys(index.key).forEach(function (field) {
        if (field === "impact" || field === "depends") {
            db.default_entities.dropIndex(index.name);
        }
    });
});

db.default_entities.createIndex({connector: 1}, {name: "connector_1"});
db.default_entities.createIndex({component: 1}, {name: "component_1"});
db.default_entities.createIndex({services: 1}, {name: "services_1"});
