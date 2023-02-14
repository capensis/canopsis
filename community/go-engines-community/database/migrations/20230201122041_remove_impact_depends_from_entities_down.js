db.default_entities.find({
    type: "resource",
    impact: null,
    depends: null,
}).forEach(function (doc) {
    var impact = [];
    var depends = [];

    if (doc.component !== "") {
        impact.push(doc.component);
    }

    if (doc.services && doc.services.length > 0) {
        impact = impact.concat(doc.services);
    }

    if (doc.connector !== "") {
        depends.push(doc.connector);
    }

    db.default_entities.updateOne({_id: doc._id}, {
        $set: {
            impact: impact,
            depends: depends,
        }
    });
});
db.default_entities.aggregate([
    {
        $match: {
            type: "component",
            impact: null,
            depends: null,
        }
    },
    {
        $graphLookup: {
            from: "default_entities",
            startWith: "$_id",
            connectFromField: "_id",
            connectToField: "component",
            restrictSearchWithMatch: {type: "resource"},
            as: "resources",
            maxDepth: 0,
        }
    },
    {
        $project: {
            connector: 1,
            services: 1,
            resources: {$map: {input: "$resources", in: "$$this._id"}}
        }
    }
]).forEach(function (doc) {
    var impact = [];
    var depends = [];

    if (doc.connector !== "") {
        impact.push(doc.connector);
    }

    if (doc.services && doc.services.length > 0) {
        impact = impact.concat(doc.services);
    }

    if (doc.resources && doc.resources.length > 0) {
        depends = doc.resources;
    }

    db.default_entities.updateOne({_id: doc._id}, {
        $set: {
            impact: impact,
            depends: depends,
        }
    });
});
db.default_entities.aggregate([
    {
        $match: {
            type: "connector",
            impact: null,
            depends: null,
        }
    },
    {
        $graphLookup: {
            from: "default_entities",
            startWith: "$_id",
            connectFromField: "_id",
            connectToField: "connector",
            restrictSearchWithMatch: {type: "resource"},
            as: "resources",
            maxDepth: 0,
        }
    },
    {
        $graphLookup: {
            from: "default_entities",
            startWith: "$_id",
            connectFromField: "_id",
            connectToField: "connector",
            restrictSearchWithMatch: {type: "component"},
            as: "components",
            maxDepth: 0,
        }
    },
    {
        $project: {
            services: 1,
            resources: {$map: {input: "$resources", in: "$$this._id"}},
            components: {$map: {input: "$components", in: "$$this._id"}},
        }
    }
]).forEach(function (doc) {
    var impact = [];
    var depends = [];

    if (doc.resources && doc.resources.length > 0) {
        impact = impact.concat(doc.resources);
    }

    if (doc.services && doc.services.length > 0) {
        impact = impact.concat(doc.services);
    }

    if (doc.components && doc.components.length > 0) {
        depends = doc.components;
    }

    db.default_entities.updateOne({_id: doc._id}, {
        $set: {
            impact: impact,
            depends: depends,
        }
    });
});
db.default_entities.aggregate([
    {
        $match: {
            type: "service",
            impact: null,
            depends: null,
        }
    },
    {
        $graphLookup: {
            from: "default_entities",
            startWith: "$_id",
            connectFromField: "_id",
            connectToField: "services",
            as: "depends",
            maxDepth: 0,
        }
    },
    {
        $project: {
            services: 1,
            depends: {$map: {input: "$depends", in: "$$this._id"}},
        }
    }
]).forEach(function (doc) {
    db.default_entities.updateOne({_id: doc._id}, {
        $set: {
            impact: doc.services,
            depends: doc.depends,
        }
    });
});

db.default_entities.updateMany({}, {$unset: {services: ""}});

db.default_entities.dropIndex("connector_1");
db.default_entities.dropIndex("component_1");
db.default_entities.dropIndex("services_1");
